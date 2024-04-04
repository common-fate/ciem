package policyset

import "github.com/urfave/cli/v2"

var Command = cli.Command{
	Name: "policyset",
	Subcommands: []*cli.Command{
		&createCommand,
		&listCommand,
		&updateCommand,
		&deleteCommand,
		&validateCommand,
	},
}
