package command

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/common-fate/clio"
	"github.com/common-fate/sdk/config"
	"github.com/urfave/cli/v2"
)

var Configure = cli.Command{
	Name:      "configure",
	Usage:     "Configure CLI",
	ArgsUsage: "The frontend url for your deployment",
	Flags:     []cli.Flag{&cli.StringFlag{Name: "context-name", Value: "default"}},
	Action: func(c *cli.Context) error {
		u := c.Args().First()
		if u == "" {
			return errors.New("please provide a url argument")
		}
		url, err := url.Parse(u)
		if err != nil {
			return err
		}
		url = url.JoinPath("/config.json")

		err = ConfigureFromURL(url.String(), WithContextName(c.String("context-name")))
		if err != nil {
			return err
		}
		clio.Success("Successfully updated config")
		return nil
	},
}

var Context = cli.Command{
	Name:  "context",
	Usage: "Manage your current CLI context",
	Subcommands: []*cli.Command{
		{
			Name: "switch",
			Action: func(c *cli.Context) error {
				contexts, err := config.ListContexts()
				if err != nil {
					return err
				}
				var context string
				err = survey.AskOne(&survey.Select{
					Options: contexts,
					Message: "select a context",
				}, &context)
				if err != nil {
					return err
				}

				err = config.SwitchContext(context)
				if err != nil {
					return err
				}

				clio.Successf("Successfully switch context to %s", context)
				return nil
			},
		},
	},
}

type Config struct {
	OauthClientId    string `json:"oauthClientId"`
	CliOAuthClientId string `json:"cliOAuthClientId"`
	OauthAuthority   string `json:"oauthAuthority"`
	APIURL           string `json:"apiUrl"`
	AccessAPIURL     string `json:"accessApiUrl"`
	AuthzGraphAPIURL string `json:"authzGraphApiUrl"`
	TeamName         string `json:"teamName"`
	FaviconUrl       string `json:"faviconUrl"`
	IconUrl          string `json:"iconUrl"`
}

type ConfigureFromURLOpts struct {
	ContextName string
}

func WithContextName(name string) func(o *ConfigureFromURLOpts) {
	return func(o *ConfigureFromURLOpts) {
		o.ContextName = name
	}
}

func ConfigureFromURL(u string, opts ...func(o *ConfigureFromURLOpts)) error {
	confOpts := ConfigureFromURLOpts{
		ContextName: "default",
	}
	for _, f := range opts {
		f(&confOpts)
	}

	res, err := http.DefaultClient.Get(u)
	if err != nil {
		return err
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var cfg Config
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		return err
	}

	newConfig := config.Default()

	newConfig.Contexts[confOpts.ContextName] = config.Context{
		APIURL:       cfg.APIURL,
		AccessURL:    cfg.AccessAPIURL,
		OIDCIssuer:   strings.TrimSuffix(cfg.OauthAuthority, "/"),
		OIDCClientID: cfg.CliOAuthClientId,
	}
	newConfig.CurrentContext = confOpts.ContextName
	err = config.Save(newConfig)
	if err != nil {
		return err
	}

	return nil
}
