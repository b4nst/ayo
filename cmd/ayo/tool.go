package main

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/alecthomas/kong"

	"github.com/banst/ayo/pkg/tool"
)

type Tool struct {
	Tools []string `arg:"" optional:"" help:"Filename of the tool to get information about. If empty, all tools from the toolbox are listed." type:"existingfile"` //nolint:lll // Struct tags are long.
}

func (t *Tool) Run(context *kong.Context, cli *CLI) error {
	tw := tabwriter.NewWriter(context.Stdout, 0, 0, 3, ' ', 0)

	if len(t.Tools) == 0 {
		var err error
		t.Tools, err = tool.ListFiles(cli.Toolbox)
		if err != nil {
			return err
		}
	}

	for i, tf := range t.Tools {
		if i > 0 {
			if _, err := fmt.Fprintln(tw); err != nil {
				return err
			}
		}

		if _, err := fmt.Fprintf(tw, "File:\t%s\n", tf); err != nil {
			return err
		}

		if t, err := tool.Load(tf); err != nil {
			if _, err := fmt.Fprintf(tw, "Error:\t%s\n", err); err != nil {
				return err
			}
		} else {
			params := []string{}
			for name, p := range t.Function.Parameters.Properties {
				params = append(params, fmt.Sprintf("%s %s", name, p.Type))
			}
			function := fmt.Sprintf("%s(%s)", t.Function.Name, strings.Join(params, ", "))

			cmd := strings.Join(append([]string{t.Cmd}, t.Args...), " ")

			format := "Function:\t%s\nDescription:\t%s\nCommand:\t%s\n"
			if _, err := fmt.Fprintf(tw, format, function, t.Function.Description, cmd); err != nil {
				return err
			}
		}

		if err := tw.Flush(); err != nil {
			return err
		}
	}

	return nil
}
