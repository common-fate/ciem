package entity

import (
	"crypto/tls"
	"encoding/json"
	"net"
	"net/http"
	"os"

	"github.com/common-fate/clio"
	"github.com/common-fate/sdk/service/entity"
	"github.com/urfave/cli/v2"
	"golang.org/x/net/http2"
)

var putCommand = cli.Command{
	Name: "put",
	Flags: []cli.Flag{
		&cli.PathFlag{Name: "file", Required: true},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		client := entity.NewClient(entity.Opts{
			HTTPClient: newInsecureClient(),
			BaseURL:    "http://127.0.0.1:5050",
		})

		f, err := os.ReadFile(c.Path("file"))
		if err != nil {
			return err
		}

		var entities []entity.EntityJSON

		err = json.Unmarshal(f, &entities)
		if err != nil {
			return err
		}

		_, err = client.BatchPutJSON(ctx, entity.BatchPutJSONInput{
			Entities: entities,
		})
		if err != nil {
			return err
		}

		clio.Successf("put %v entities", len(entities))
		return nil
	},
}

func newInsecureClient() *http.Client {
	return &http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, _ *tls.Config) (net.Conn, error) {
				// If you're also using this client for non-h2c traffic, you may want
				// to delegate to tls.Dial if the network isn't TCP or the addr isn't
				// in an allowlist.
				return net.Dial(network, addr)
			},
			// Don't forget timeouts!
		},
	}
}
