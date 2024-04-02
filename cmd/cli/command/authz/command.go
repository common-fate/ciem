package authz

import (
	"github.com/common-fate/cli/cmd/cli/command/authz/policyset"
	"github.com/common-fate/cli/cmd/cli/command/authz/schema"
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name:  "authz",
	Usage: "Manage Common Fate authorization",
	Subcommands: []*cli.Command{
		&evaluateCommand,
		&schema.Command,
		&policyset.Command,
	},
}
