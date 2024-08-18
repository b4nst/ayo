package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ollama/ollama/api"
)

func main() {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	messages := []api.Message{
		api.Message{
			Role:    "user",
			Content: "Pick Jira task to tackle for me",
		},
	}

	tools := []api.Tool{
		api.Tool{
			Type: "function",
			Function: api.ToolFunction{
				Name:        "git_commit",
				Description: "Commit changes",
				Parameters: Parameters{
					Type:     "object",
					Required: []string{"message"},
					Properties: map[string]struct {
						Type        string   `json:"type"`
						Description string   `json:"description"`
						Enum        []string `json:"enum,omitempty"`
					}{
						"message": {
							Type:        "string",
							Description: "Commit message, conventional commit format",
						},
					},
				},
			},
		},
	}

	stream := false
	fmt.Println("Stream: ", stream)
	ctx := context.Background()
	req := &api.ChatRequest{
		Model:    "llama3.1",
		Messages: messages,
		Tools:    tools,
		Stream:   &stream,
	}

	respFunc := func(resp api.ChatResponse) error {
		fmt.Printf("%+v\n", resp)
		return nil
	}

	err = client.Chat(ctx, req, respFunc)
	if err != nil {
		log.Fatal(err)
	}
}

type Parameters struct {
	Type       string   `json:"type"`
	Required   []string `json:"required"`
	Properties map[string]struct {
		Type        string   `json:"type"`
		Description string   `json:"description"`
		Enum        []string `json:"enum,omitempty"`
	} `json:"properties"`
}
