package entities

import "github.com/urfave/cli/v2"

var Command = cli.Command{
	Name: "entities",
	Subcommands: []*cli.Command{
		&putCommand,
		&deleteCommand,
	},
}
