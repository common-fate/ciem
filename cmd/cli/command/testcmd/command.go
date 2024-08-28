package testcmd

import "github.com/urfave/cli/v2"

var Command = cli.Command{
	Name:    "tests",
	Aliases: []string{"test"},
	Subcommands: []*cli.Command{
		&runCommand,
		&createCommand,
	},
}
