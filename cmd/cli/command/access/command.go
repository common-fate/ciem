package access

import "github.com/urfave/cli/v2"

var Command = cli.Command{
	Name: "access",
	Subcommands: []*cli.Command{
		&showCommand,
	},
}
