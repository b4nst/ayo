package tool

import (
	"strings"
	"text/template"
)

func RenderArgs(args []string, values map[string]any) ([]string, error) {
	tmpl := template.New("cmd")
	tmpl.Option("missingkey=error")
	sb := new(strings.Builder)

	cmdargs := make([]string, len(args))

	// Apply templating to all arguments
	for i, arg := range args {
		targ, err := tmpl.Parse(arg)
		if err != nil {
			return nil, err
		}
		err = targ.Execute(sb, values)
		if err != nil {
			return nil, err
		}
		cmdargs[i] = sb.String()
		sb.Reset()
	}

	return cmdargs, nil
}
