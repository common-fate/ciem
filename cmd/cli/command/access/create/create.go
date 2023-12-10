package create

import "github.com/urfave/cli/v2"

var Command = cli.Command{
	Name: "create",
	Subcommands: []*cli.Command{
		&selectorCommand,
	},
}
