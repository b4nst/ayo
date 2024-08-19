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
	apiTools := make([]api.Tool, 0, len(tools))
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
	init := Message("You are a CLI assistant with tool calling capabilities. Answer the question the best you can. Only use tool when you need to. Don't use a tool you don't have access to.", RoleSystem)

	for history := []api.Message{
		init,
		Message(message, RoleUser),
	}; ; {
		answer, err := ai.query(ctx, history)
		if err != nil {
			return err
		}
		history = append(history, answer)

		if len(answer.ToolCalls) == 0 {
			//nolint:forbidigo // TODO: refactor this print.
			fmt.Println(answer.Content)
			return nil
		}

		for _, tc := range answer.ToolCalls {
			log.Debug("tool call requested", "function", tc.Function.Name, "args", tc.Function.Arguments)
			f := tc.Function
			t, ok := ai.tools[f.Name]
			if !ok {
				history = append(history, Message("You don't have access to that tool.", RoleSystem))
				continue
			}
			toolresp, err := t.Run(f.Arguments)
			if err != nil {
				history = append(history, Message(fmt.Sprintf("Error running tool: %v", err), RoleTool))
				continue
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
