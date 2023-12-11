package access

import (
	"github.com/common-fate/ciem/cmd/cli/command/access/approve"
	"github.com/common-fate/ciem/cmd/cli/command/access/close"
	"github.com/common-fate/ciem/cmd/cli/command/access/list"
	"github.com/common-fate/ciem/cmd/cli/command/access/update"
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name:  "access",
	Usage: "Manage and request access to entitlements",
	Subcommands: []*cli.Command{
		&ensureCommand,
		&list.Command,
		&update.Command,
		&close.Command,
		&approve.Command,
	},
}
