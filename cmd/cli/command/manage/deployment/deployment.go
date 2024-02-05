package deployment

import (
	"github.com/common-fate/cli/cmd/cli/command/manage/deployment/logs"
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name:    "deployment",
	Aliases: []string{"dep"},
	Usage:   "Manage your Common Fate Deployment",
	Subcommands: []*cli.Command{
		&logs.Command,
	},
}
