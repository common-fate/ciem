package auditlog

import "github.com/urfave/cli/v2"

var Command = cli.Command{
	Name: "auditlog",
	Subcommands: []*cli.Command{
		&listCommand,
	},
}
