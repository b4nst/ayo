package main

import (
	"context"
	"strings"

	ayoai "github.com/banst/ayo/pkg/ai"
	"github.com/banst/ayo/pkg/log"
	"github.com/banst/ayo/pkg/tool"
)

type Chat struct {
	Model string   `help:"The model to use for the chatbot" default:"llama3.1"`
	Args  []string `arg:"" help:"The message to send to the assistant"`
}

func (c *Chat) Run(ctx context.Context, cli *CLI) error {
	log := log.FromContext(ctx).With("caller", logPackagePrefix+"Chat:Run")

	log.Debug("loading tools", "dir", cli.Toolbox)
	tools, err := tool.LoadAll(cli.Toolbox)
	if err != nil {
		return err
	}
	log.Debug("tools loaded", "total", len(tools))

	ai := ayoai.NewAI(tools, c.Model)
	log.Debug("ai loaded", "model", c.Model)

	return ai.Chat(ctx, strings.Join(c.Args, " "))
}
