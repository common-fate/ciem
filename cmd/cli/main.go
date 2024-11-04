package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	availabilityspec "github.com/common-fate/cli/cmd/cli/command/availability_spec"
	"github.com/common-fate/cli/cmd/cli/command/testcmd"
	"github.com/common-fate/cli/cmd/cli/command/workflow"

	"go.uber.org/zap"

	"github.com/AlecAivazis/survey/v2"
	"github.com/common-fate/cli/cmd/cli/command"
	"github.com/common-fate/cli/cmd/cli/command/access"
	"github.com/common-fate/cli/cmd/cli/command/auditlog"
	"github.com/common-fate/cli/cmd/cli/command/authz"
	"github.com/common-fate/cli/cmd/cli/command/deployment"
	"github.com/common-fate/cli/cmd/cli/command/entity"
	"github.com/common-fate/cli/cmd/cli/command/identity"
	"github.com/common-fate/cli/cmd/cli/command/integration"
	"github.com/common-fate/cli/internal/build"
	glidecli "github.com/common-fate/glide-cli"

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

			if c.Args().First() == "oss" {
				return nil
			}

			if c.Args().First() == "configure" {
				return nil
			}
			_, err := config.LoadDefault(c.Context)
			if err == nil {
				return nil
			}
			if err != config.ErrConfigFileNotFound {
				return fmt.Errorf("error loading Common Fate config: %w", err)
			}

			if os.Getenv("CI") == "true" {
				return fmt.Errorf("failed to load Common Fate config from file (returning early because the 'CI' environment variable is 'true', which usually indicates you're running this in a CI environment): %w", err)
			}

			clio.Debugw("failed to load config from file, will prompt user to configure", zap.Error(err))

			// prompt for an App URL to load initial config
			clio.Info("It looks like this is your first time using the Common Fate CLI")
			clio.Info("To get started, you need to configure the CLI to connect to your team's Common Fate deployment")
			clio.Info("Enter the URL of your deployment, e.g https://commonfate.example.com")

			var u string
			err = survey.AskOne(&survey.Input{
				Message: "The URL of your team's Common Fate deployment:",
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
			&entity.Command,
			&authz.Command,
			&access.Command,
			&glidecli.OSSSubCommand,
			&auditlog.Command,
			&command.Configure,
			&command.Context,
			&deployment.Command,
			&integration.Command,
			&testcmd.Command,
			&workflow.Command,
			&availabilityspec.Command,
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
