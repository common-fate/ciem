package main

import (
	"errors"
	"net/url"
	"os"
	"strings"

	"go.uber.org/zap"

	"github.com/AlecAivazis/survey/v2"
	"github.com/common-fate/ciem/cmd/cli/command"
	"github.com/common-fate/ciem/cmd/cli/command/access"
	"github.com/common-fate/ciem/cmd/cli/command/auditlog"
	"github.com/common-fate/ciem/cmd/cli/command/authz"
	"github.com/common-fate/ciem/cmd/cli/command/aws"
	"github.com/common-fate/ciem/cmd/cli/command/entity"
	"github.com/common-fate/ciem/cmd/cli/command/identity"
	"github.com/common-fate/ciem/cmd/cli/command/policyset"
	build "github.com/common-fate/ciem/internal"
	"github.com/common-fate/clio"
	"github.com/common-fate/clio/clierr"
	"github.com/common-fate/sdk/config"
	"github.com/common-fate/sdk/loginflow"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:      "cf",
		Writer:    os.Stderr,
		Usage:     "https://commonfate.io",
		UsageText: "cf [options] [command]",
		Version:   build.Version,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "api-url", Usage: "override the Common Fate API URL"},
			&cli.BoolFlag{Name: "verbose", Usage: "Enable verbose logging, effectively sets environment variable CF_LOG=DEBUG"},
		},
		Before: func(c *cli.Context) error {
			if c.Bool("verbose") {
				clio.SetLevelFromString("debug")
			}

			if c.Args().First() == "configure" {
				return nil
			}
			_, err := config.LoadDefault(c.Context)
			if err == nil {
				return nil
			}

			clio.Debugw("failed to load config from file, will prompt user to configure", zap.Error(err))

			// prompt for an App URL to load initial config
			clio.Info("It looks like this is your first time using the Common Fate CLI")
			clio.Info("To get started, you need to configure the CLI to connect to your team's Common Fate deployment")
			clio.Info("This is simple, just enter the URL of your deployment, e.g https://commonfate.example.com")

			var u string
			err = survey.AskOne(&survey.Input{
				Message: "Enter the URL of your teams Common Fate deployment:",
			}, &u, survey.WithValidator(func(ans interface{}) error {
				a := EnsureURLScheme(ans.(string))
				url, err := url.Parse(a)
				if err != nil {
					return err
				}
				if url.Path != "" {
					return errors.New("URL should not include a path")
				}
				return nil
			}))
			if err != nil {
				return err
			}
			baseUrl, err := url.Parse(EnsureURLScheme(u))
			if err != nil {
				return err
			}
			url := baseUrl.JoinPath("/config.json")

			err = command.ConfigureFromURL(url.String())
			if err != nil {
				return err
			}
			clio.Success("Your CLI has been configured successfully!")
			clio.Infof("Opening your browser to try and log you in to %s...", baseUrl.String())

			// try and log the user in immediately
			cfg, err := config.LoadDefault(c.Context)
			if err != nil {
				return err
			}

			lf := loginflow.NewFromConfig(cfg)

			return lf.Login(c.Context)
		},
		Commands: []*cli.Command{
			&command.Login,
			&command.Logout,
			&identity.Command,
			&command.Auth,
			&policyset.Command,
			&entity.Command,
			&authz.Command,
			&access.Command,
			&auditlog.Command,
			&command.Configure,
			&command.Context,
			&aws.Command,
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

func EnsureURLScheme(u string) string {
	if !(strings.HasPrefix(u, "https://") || (strings.HasPrefix(u, "http://"))) {
		clio.Debugf("URL did not have a scheme so the default https:// will be added")
		return "https://" + u
	}
	return u
}
