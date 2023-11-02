package main

import (
	"os"

	"go.uber.org/zap"

	"github.com/common-fate/ciem/cmd/cli/command"
	"github.com/common-fate/clio"
	"github.com/common-fate/clio/clierr"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:      "accessclient",
		Writer:    os.Stderr,
		Usage:     "https://commonfate.io",
		UsageText: "accessclient [options] [command]",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "api-url", Usage: "override the Common Fate API URL"},
			&cli.BoolFlag{Name: "verbose", Usage: "Enable verbose logging, effectively sets environment variable CF_LOG=DEBUG"},
		},
		Before: func(ctx *cli.Context) error {
			if ctx.Bool("verbose") {
				clio.SetLevelFromString("debug")
			}

			return nil
		},
		Commands: []*cli.Command{
			&command.Login,
			&command.Auth,
			&command.List,
		},
	}
	clio.SetLevelFromEnv("CF_LOG")
	zap.ReplaceGlobals(clio.G())

	err := app.Run(os.Args)
	if err != nil {
		// if the error is an instance of clierr.PrintCLIErrorer then print the error accordingly
		if cliError, ok := err.(clierr.PrintCLIErrorer); ok {
			cliError.PrintCLIError()
		} else {
			clio.Error(err.Error())
		}
		os.Exit(1)
	}
}
