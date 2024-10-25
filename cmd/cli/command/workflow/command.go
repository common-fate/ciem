package workflow

import (
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name:  "workflow",
	Usage: "Manage Common Fate Access Workflows",
	Subcommands: []*cli.Command{
		&listCommand,
		&deleteCommand,
	},
}
