package request

import (
	"github.com/common-fate/ciem/cmd/cli/command/jit/request/access"
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name:  "request",
	Usage: "Request access to an entitlement",
	Subcommands: []*cli.Command{
		&access.Command,
		&approveCommand,
		&closeCommand,
		&list,
	},
}
