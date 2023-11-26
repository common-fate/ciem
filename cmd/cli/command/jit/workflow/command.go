package workflow

import "github.com/urfave/cli/v2"

var Command = cli.Command{
	Name: "workflow",
	Subcommands: []*cli.Command{
		&createCommand,
	},
}
