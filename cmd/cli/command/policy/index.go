package policy

import (
	"github.com/bufbuild/connect-go"
	"github.com/common-fate/clio"
	authzv1alpha1 "github.com/common-fate/sdk/gen/commonfate/authz/v1alpha1"
	"github.com/common-fate/sdk/service/authz/index"
	"github.com/urfave/cli/v2"
)

var indexCommand = cli.Command{
	Name: "index",
	Action: func(c *cli.Context) error {
		ctx := c.Context
		client := index.NewClient(newInsecureClient(), "http://127.0.0.1:5050")

		res, err := client.StartIndexJob(ctx, connect.NewRequest(&authzv1alpha1.StartIndexJobRequest{}))
		if err != nil {
			return err
		}

		clio.Successf("started indexing: %s", res.Msg.JobId)
		return nil
	},
}
