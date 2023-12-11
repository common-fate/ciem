package authz

import (
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name: "authz",
	Subcommands: []*cli.Command{
		&evaluateCommand,
	},
}
