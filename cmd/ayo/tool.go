package main

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/alecthomas/kong"

	"github.com/banst/ayo/pkg/tool"
)

type Tools []string

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
			fmt.Fprintln(tw)
		}
		fmt.Fprintf(tw, "File:\t%s\n", tf)

		t, err := tool.Load(tf)
		if err != nil {
			fmt.Fprintf(tw, "Error:\t%s\n", err)
		} else {
			printToolInfo(tw, t)
		}

		if err := tw.Flush(); err != nil {
			return err
		}
	}

	return nil
}

func printToolInfo(w io.Writer, t *tool.Tool) {
	params := []string{}
	for name, p := range t.Function.Parameters.Properties {
		params = append(params, fmt.Sprintf("%s %s", name, p.Type))
	}
	function := fmt.Sprintf("%s(%s)", t.Function.Name, strings.Join(params, ", "))

	cmd := strings.Join(append([]string{t.Cmd}, t.Args...), " ")

	format := "Function:\t%s\nDescription:\t%s\nCommand:\t%s\n"
	fmt.Fprintf(w, format, function, t.Function.Description, cmd)
}
