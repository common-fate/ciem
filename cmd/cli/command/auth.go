package command

import (
	"encoding/json"
	"fmt"

	"github.com/common-fate/ciem/config"
	"github.com/common-fate/ciem/loginflow"
	"github.com/urfave/cli/v2"
)

var Auth = cli.Command{
	Name:  "auth",
	Usage: "Manage Common Fate authentication",
	Subcommands: []*cli.Command{
		&tokenCommand,
		&refreshCommand,
	},
}

var tokenCommand = cli.Command{
	Name:  "token",
	Usage: "Print details about the Common Fate authentication token",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "show-sensitive-values", Usage: "Show sensitive values"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		tok, err := cfg.TokenStore.Token()
		if err != nil {
			return err
		}

		show := c.Bool("show-sensitive-values")

		if !show && tok.AccessToken != "" {
			tok.AccessToken = `redacted (use --show-sensitive-values to show)`
		}
		if !show && tok.RefreshToken != "" {
			tok.RefreshToken = `redacted (use --show-sensitive-values to show)`
		}

		tokenStr, err := json.Marshal(tok)
		if err != nil {
			return err
		}

		fmt.Println(string(tokenStr))

		return nil
	},
}

var refreshCommand = cli.Command{
	Name:  "refresh",
	Usage: "Force a refresh of the access token",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "show-sensitive-values", Usage: "Show sensitive values"},
	},
	Action: func(c *cli.Context) error {
		cfg, err := config.LoadDefault(c.Context)
		if err != nil {
			return err
		}

		lf := loginflow.NewFromConfig(cfg)

		tok, err := lf.RefreshToken(c.Context)
		if err != nil {
			return err
		}

		show := c.Bool("show-sensitive-values")

		if !show && tok.AccessToken != "" {
			tok.AccessToken = `redacted (use --show-sensitive-values to show)`
		}
		if !show && tok.RefreshToken != "" {
			tok.RefreshToken = `redacted (use --show-sensitive-values to show)`
		}

		tokenStr, err := json.Marshal(tok)
		if err != nil {
			return err
		}

		fmt.Println(string(tokenStr))

		return nil
	},
}
