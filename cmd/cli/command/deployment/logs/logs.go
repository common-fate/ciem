package logs

import "github.com/urfave/cli/v2"

var Command = cli.Command{
	Name:        "logs",
	Description: "View recent application logs from Cloudwatch or stream them in real time",
	Usage:       "View recent application logs from Cloudwatch or stream them in real time",
	Action:      cli.ShowSubcommandHelp,
	Subcommands: []*cli.Command{&getCommand, &watchCommand},
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "namespace", Value: "common-fate", Usage: "Your Common Fate infrastructure deployment namespace"},
		&cli.StringFlag{Name: "stage", Value: "prod", Usage: "Your Common Fate infrastructure deployment stage"},
	},
}

// the services names are defined here for this CLI command, and may be different in other usages
var ServiceNames = []string{
	"access-handler",
	"authz",
	"control-plane",
	"otel-collector",
	"provisioner",
	"builtin-provisioner",
	"web",
	"worker",
}
