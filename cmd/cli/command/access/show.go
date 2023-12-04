package access

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/bufbuild/connect-go"
	"github.com/common-fate/ciem/table"
	authzv1alpha1 "github.com/common-fate/sdk/gen/commonfate/authz/v1alpha1"
	"github.com/common-fate/sdk/service/authz/index"
	"github.com/common-fate/sdk/uid"
	"github.com/urfave/cli/v2"
	"golang.org/x/net/http2"
)

var showCommand = cli.Command{
	Name: "show",
	Subcommands: []*cli.Command{
		&showUserCommand,
	},
}

var showUserCommand = cli.Command{
	Name: "user",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "id", Required: true},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		client := index.NewClient(newInsecureClient(), "http://127.0.0.1:5050")

		id := uid.New("User", c.String("id"))

		res, err := client.LookupResources(ctx, connect.NewRequest(&authzv1alpha1.LookupResourcesRequest{
			Universe:    "default",
			Environment: "production",
			Principal:   id.ToAPI(),
		}))
		if err != nil {
			return err
		}
		w := table.New(os.Stdout)
		w.Columns("TYPE", "ID", "ACTIONS")

		for _, resource := range res.Msg.Resources {
			var actions []string
			for _, act := range resource.Actions {
				actions = append(actions, fmt.Sprintf("%s::%s", act.Action.Type, act.Action.Id))
			}

			w.Row(resource.Resource.Uid.Type, resource.Resource.Uid.Id, strings.Join(actions, ", "))
		}

		err = w.Flush()
		if err != nil {
			return err
		}

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
