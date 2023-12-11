package entity

import "github.com/urfave/cli/v2"

var Command = cli.Command{
	Name:    "entity",
	Aliases: []string{"entities"},
	Subcommands: []*cli.Command{
		&putCommand,
		&deleteCommand,
		&listCommand,
	},
}
