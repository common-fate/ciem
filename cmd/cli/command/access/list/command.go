package list

import "github.com/urfave/cli/v2"

var Command = cli.Command{
	Name: "list",
	Subcommands: []*cli.Command{
		&availableCommand,
		&requestsCommand,
		&approversCommand,
	},
}
