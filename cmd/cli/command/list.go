package command

import (
	"github.com/urfave/cli/v2"
)

var List = cli.Command{
	Name:  "list",
	Usage: "List available entitlements",
	Subcommands: []*cli.Command{
		&gcpList,
	},
}

var gcpList = cli.Command{
	Name:  "gcp",
	Usage: "List available GCP entitlements",
	Action: func(c *cli.Context) error {
		// ctx := c.Context

		// cfg, err := config.LoadDefault(ctx)
		// if err != nil {
		// 	return err
		// }

		// client := access.NewFromConfig(cfg)

		// res, err := client.ListAccessRequests(ctx, connect.NewRequest(&accessv1alpha1.ListAccessRequestsRequest{}))
		// if err != nil {
		// 	return err
		// }

		// w := table.New(os.Stdout)
		// w.Columns("PROJECT", "ROLE")

		// for _, e := range res.Msg.AccessRequests {
		// 	switch v := e.Entitlement.Resource.Resource.(type) {
		// 	case *accessv1alpha1.Resource_AwsAccount:
		// 		w.Row(v.AwsAccount.AccountId, v.AwsAccount.Role)
		// 	case *accessv1alpha1.Resource_GcpProject:
		// 		w.Row(v.GcpProject.Project, v.GcpProject.Role)
		// 	}
		// }

		// w.Flush()

		return nil
	},
}
