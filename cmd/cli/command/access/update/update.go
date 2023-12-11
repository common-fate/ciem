package update

import "github.com/urfave/cli/v2"

var Command = cli.Command{
	Name: "update",
	Subcommands: []*cli.Command{
		&selectorCommand,
	},
}
