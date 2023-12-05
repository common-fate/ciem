package policy

import "github.com/urfave/cli/v2"

var Command = cli.Command{
	Name: "policy",
	Subcommands: []*cli.Command{
		&applyCommand,
		&listCommand,
	},
}
