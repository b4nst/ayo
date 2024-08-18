package main

import (
	"fmt"
	"text/tabwriter"

	"github.com/alecthomas/kong"

	"github.com/banst/ayo/cmd/ayo/config"
)

// Version represents the version command.
type Version struct {
	Short bool `help:"Print just the semver version." short:"s"`
}

// Run runs the version command.
func (v *Version) Run(ctx *kong.Context) error {
	if v.Short {
		_, err := fmt.Fprintln(ctx.Stdout, config.Version)
		if err != nil {
			return err
		}
		return nil
	}

	w := tabwriter.NewWriter(ctx.Stdout, 0, 0, 3, ' ', 0)
	format := "Version:\t%s\nCommit:\t%s\nDate:\t%s\nBuilt by:\t%s\n"
	_, err := fmt.Fprintf(w, format, config.Version, config.Commit, config.Date, config.BuiltBy)
	if err != nil {
		return err
	}
	return w.Flush()
}
