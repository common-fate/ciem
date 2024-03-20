package log

import (
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name: "log",
	Subcommands: []*cli.Command{
		&queryCommand,
	},
}
