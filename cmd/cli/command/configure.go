package command

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/common-fate/clio"
	"github.com/common-fate/sdk/config"
	"github.com/urfave/cli/v2"
)

var Configure = cli.Command{
	Name:      "configure",
	Usage:     "Configure CLI",
	ArgsUsage: "The frontend url for your deployment",
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

		err = ConfigureFromURL(url.String())
		if err != nil {
			return err
		}
		clio.Success("Successfully updated config")
		return nil
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

func ConfigureFromURL(u string) error {
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

	newConfig.Contexts["default"] = config.Context{
		APIURL:       cfg.APIURL,
		AccessURL:    cfg.AccessAPIURL,
		OIDCIssuer:   strings.TrimSuffix(cfg.OauthAuthority, "/"),
		OIDCClientID: cfg.CliOAuthClientId,
	}
	newConfig.CurrentContext = "default"
	err = config.Save(newConfig)
	if err != nil {
		return err
	}

	return nil
}
