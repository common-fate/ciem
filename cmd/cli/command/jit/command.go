package jit

import (
	"github.com/common-fate/ciem/cmd/cli/command/jit/available"
	"github.com/common-fate/ciem/cmd/cli/command/jit/request"
	"github.com/common-fate/ciem/cmd/cli/command/jit/workflow"
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name: "jit",
	Subcommands: []*cli.Command{
		&workflow.Command,
		&available.Command,
		&request.Command,
	},
}
