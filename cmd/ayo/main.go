package main

import (
	"context"

	"github.com/alecthomas/kong"

	"github.com/banst/ayo/pkg/log"
)

const logPackagePrefix = "main::"

type CLI struct {
	Verbose   bool   `help:"Enable verbose mode." short:"v" default:"false"`
	LogFormat string `help:"Log format to use." enum:"json,pretty,logfmt" default:"pretty"`
	Toolbox   string `help:"The directory to load tools from" default:"~/.config/ayo/toolbox" type:"existingdir" env:"AYO_TOOLBOX"` //nolint: lll

	Chat    Chat    `cmd:"" help:"Send a message to the chatbot" default:"withargs"`
	Version Version `cmd:"" help:"Show the version information"`
	Tool    Tool    `cmd:"" help:"Get information about a tool"`
}

func main() {
	var cli CLI
	cmd := kong.Parse(&cli, kong.Name("ayo"), kong.Description("A cli AI assistant"))

	// Create logger
	l := log.New(cmd.Stdout, cli.LogFormat, cli.Verbose)

	// Bind context
	ctx := log.ContextWithLogger(context.Background(), l)
	cmd.BindTo(ctx, (*context.Context)(nil)) // cannot bind to context directly, since it's an interface

	cmd.FatalIfErrorf(cmd.Run())
}
