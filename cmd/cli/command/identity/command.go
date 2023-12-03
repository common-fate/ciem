package identity

import "github.com/urfave/cli/v2"

var Command = cli.Command{
	Name:    "identity",
	Aliases: []string{"id"},
	Subcommands: []*cli.Command{
		&getCommand,
	},
}
