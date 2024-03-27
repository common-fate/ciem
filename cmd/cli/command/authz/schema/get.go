package schema

import (
	"fmt"

	"connectrpc.com/connect"
	"github.com/common-fate/sdk/config"
	authzv1alpha1 "github.com/common-fate/sdk/gen/commonfate/authz/v1alpha1"
	schemasvc "github.com/common-fate/sdk/service/authz/schema"
	"github.com/urfave/cli/v2"
)

var getCommand = cli.Command{
	Name:  "get",
	Usage: "Retrieve the Cedar schema in JSON format",
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		client := schemasvc.NewFromConfig(cfg)

		res, err := client.GetSchemaJSONString(ctx, connect.NewRequest(&authzv1alpha1.GetSchemaJSONStringRequest{}))
		if err != nil {
			return err
		}

		fmt.Println(res.Msg.Schema)

		return nil
	},
}
