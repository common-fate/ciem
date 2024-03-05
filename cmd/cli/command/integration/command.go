package integration

import "github.com/urfave/cli/v2"

var Command = cli.Command{
	Name:    "integration",
	Aliases: []string{"in"},
	Subcommands: []*cli.Command{
		&resetCommand,
	},
}
