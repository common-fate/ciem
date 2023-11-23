package command

import (
	"github.com/common-fate/ciem/config"
	"github.com/common-fate/ciem/loginflow"
	"github.com/urfave/cli/v2"
)

var Login = cli.Command{
	Name:  "login",
	Usage: "Log in to Common Fate",
	Action: func(c *cli.Context) error {
		cfg, err := config.LoadDefault(c.Context)
		if err != nil {
			return err
		}

		lf := loginflow.NewFromConfig(cfg)

		return lf.Login(c.Context)
	},
}

var Logout = cli.Command{
	Name:  "logout",
	Usage: "Log out of Common Fate",
	Action: func(c *cli.Context) error {
		cfg, err := config.LoadDefault(c.Context)
		if err != nil {
			return err
		}

		return cfg.TokenStore.Clear()
	},
}
