package authz

import (
	"github.com/common-fate/cli/cmd/cli/command/authz/schema"
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name: "authz",
	Subcommands: []*cli.Command{
		&evaluateCommand,
		&schema.Command,
	},
}
