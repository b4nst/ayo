package main

import (
	"context"
	"strings"

	ayoai "github.com/banst/ayo/pkg/ai"
	"github.com/banst/ayo/pkg/log"
	"github.com/banst/ayo/pkg/tool"
)

type Chat struct {
	Args []string `arg:"" help:"The message to send to the assistant"`
}

func (c *Chat) Run(ctx context.Context, cli *CLI) error {
	log := log.FromContext(ctx).With("cmd", "chat")

	log.Debug("loading tools", "dir", cli.ToolDir)
	tools, err := tool.LoadTools(cli.ToolDir)
	if err != nil {
		return err
	}
	log.Debug("tools loaded", "total", len(tools))

	ai := ayoai.NewAI(tools)
	log.Debug("ai loaded")

	return ai.Chat(ctx, strings.Join(c.Args, " "))
}
