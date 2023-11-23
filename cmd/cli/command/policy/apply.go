package policy

import (
	"crypto/tls"
	"net"
	"net/http"
	"os"

	"github.com/common-fate/ciem/service/authz"
	"github.com/common-fate/clio"
	"github.com/urfave/cli/v2"
	"golang.org/x/net/http2"
)

var applyCommand = cli.Command{
	Name: "apply",
	Flags: []cli.Flag{
		&cli.PathFlag{Name: "file", Required: true},
		&cli.StringFlag{Name: "id", Required: true},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		// cfg, err := config.LoadDefault(ctx)
		// if err != nil {
		// 	return err
		// }

		client := authz.NewClient(newInsecureClient(), "http://127.0.0.1:5050")

		f, err := os.ReadFile(c.Path("file"))
		if err != nil {
			return err
		}

		_, err = client.BatchPutPolicy(ctx, authz.BatchPutPolicyInput{
			Policies: []authz.Policy{
				{
					ID:    c.String("id"),
					Cedar: string(f),
				},
			},
		})
		if err != nil {
			return err
		}

		clio.Success("applied policy")
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
