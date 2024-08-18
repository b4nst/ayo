package ai

import (
	"github.com/ollama/ollama/api"
)

type MessageRole string

const (
	RoleUser   MessageRole = "user"
	RoleSystem MessageRole = "system"
	RoleTool   MessageRole = "tool"
)

func Message(content string, role MessageRole) api.Message {
	return api.Message{
		Role:    string(role),
		Content: content,
	}
}
