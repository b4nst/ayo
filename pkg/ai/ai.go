package ai

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ollama/ollama/api"
	"github.com/ollama/ollama/envconfig"

	"github.com/banst/ayo/pkg/log"
	"github.com/banst/ayo/pkg/tool"
)

const logPackagePrefix = "ai::"

type AI struct {
	model    string
	client   *api.Client
	stream   bool
	tools    map[string]*tool.Tool
	apiTools []api.Tool
}

func NewAI(tools map[string]*tool.Tool, model string) *AI {
	client := api.NewClient(envconfig.Host(), http.DefaultClient)
	apiTools := make([]api.Tool, len(tools))
	for _, t := range tools {
		apiTools = append(apiTools, t.Tool)
	}

	return &AI{
		client:   client,
		tools:    tools,
		stream:   false,
		apiTools: apiTools,
		model:    model,
	}
}

func (ai *AI) Chat(ctx context.Context, message string) error {
	log := log.FromContext(ctx).With("caller", logPackagePrefix+"AI:Chat")

	for history := []api.Message{Message(message, RoleUser)}; ; {
		answer, err := ai.query(ctx, history)
		if err != nil {
			return err
		}
		history = append(history, answer)

		if len(answer.ToolCalls) <= 0 {
			fmt.Println(answer.Content)
			return nil
		}

		for _, tc := range answer.ToolCalls {
			f := tc.Function
			t, ok := ai.tools[f.Name]
			if !ok {
				return fmt.Errorf("tool %s not found", f.Name)
			}
			log.Debug("running tool", "tool", f.Name, "args", f.Arguments)
			toolresp, err := t.Run(f.Arguments)
			if err != nil {
				return err
			}
			history = append(history, Message(string(toolresp), RoleTool))
		}
	}
}

func (ai *AI) query(ctx context.Context, messages []api.Message) (api.Message, error) {
	log := log.FromContext(ctx).With("caller", logPackagePrefix+"AI:query")

	var answer api.Message
	handler := func(resp api.ChatResponse) error {
		log.Debug("received response", "metrics", resp.Metrics, "model", resp.Model)
		answer = resp.Message
		return nil
	}

	req := &api.ChatRequest{
		Model:    ai.model,
		Messages: messages,
		Tools:    ai.apiTools,
		Stream:   &ai.stream,
	}
	err := ai.client.Chat(ctx, req, handler)
	if err != nil {
		return answer, err
	}

	return answer, nil
}

func chatResponse(ch chan<- api.Message) func(api.ChatResponse) error {
	return func(resp api.ChatResponse) error {
		ch <- resp.Message
		return nil
	}
}
