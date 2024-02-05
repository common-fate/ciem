package manage

import (
	"github.com/common-fate/cli/cmd/cli/command/manage/deployment"
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name:    "manage",
	Aliases: []string{"mng"},
	Usage:   "Manage Common Fate",
	Subcommands: []*cli.Command{
		&deployment.Command,
	},
}
