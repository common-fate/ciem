package policyset

import (
	"errors"

	"github.com/urfave/cli/v2"
)

var validateCommand = cli.Command{
	Name: "validate",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "file-pattern", Usage: "Policy set file(s) to validate. May be a glob.", Value: "**/*.cedar"},
		&cli.StringFlag{Name: "file-pattern-dir", Usage: "Directory to use as the base for the --file-pattern argument", Value: "."},
	},
	Action: func(c *cli.Context) error {
		return errors.New("this command is deprecated, please use the Cedar CLI for policy set validation: https://docs.commonfate.io/authz/validation")
	},
}
