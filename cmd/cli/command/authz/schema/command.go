package schema

import "github.com/urfave/cli/v2"

var Command = cli.Command{
	Name:  "schema",
	Usage: "Manage Cedar schemas used for authorization",
	Subcommands: []*cli.Command{
		&getCommand,
	},
}
