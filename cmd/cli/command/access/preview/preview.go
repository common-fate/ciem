package preview

import (
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name:  "preview",
	Usage: "Preview available entitlements for a principal",
	Subcommands: []*cli.Command{
		&userAccess,
		&entitlementCommand,
	},
}
