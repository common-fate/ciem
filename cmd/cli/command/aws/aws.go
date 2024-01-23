package aws

import (
	"github.com/common-fate/ciem/cmd/cli/command/aws/rds"
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name:  "aws",
	Usage: "Perform AWS Operations",
	Subcommands: []*cli.Command{
		&rds.Command,
	},
}
