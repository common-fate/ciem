package available

import "github.com/urfave/cli/v2"

var Command = cli.Command{
	Name:    "available",
	Aliases: []string{"av"},
	Usage:   "Query available Just-In-Time entitlements",
	Subcommands: []*cli.Command{
		&listCommand,
	},
}
