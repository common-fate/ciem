package integration

import "github.com/urfave/cli/v2"

var Command = cli.Command{
	Name: "integration",
	Subcommands: []*cli.Command{
		&deleteOauthTokenCommand,
	},
}
