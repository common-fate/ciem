package policy

import (
	"fmt"
	"os"

	"github.com/bufbuild/connect-go"
	authzv1alpha1 "github.com/common-fate/ciem/gen/commonfate/authz/v1alpha1"
	"github.com/common-fate/ciem/service/index"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

var simulateCommand = cli.Command{
	Name: "simulate",
	Flags: []cli.Flag{
		&cli.PathFlag{Name: "file", Required: true},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		// cfg, err := config.LoadDefault(ctx)
		// if err != nil {
		// 	return err
		// }

		client := index.NewClient(newInsecureClient(), "http://127.0.0.1:5050")

		f, err := os.ReadFile(c.Path("file"))
		if err != nil {
			return err
		}

		res, err := client.Simulate(ctx, connect.NewRequest(&authzv1alpha1.SimulateRequest{
			Policies: []*authzv1alpha1.Policy{
				{
					Id:    "policy.0",
					Cedar: string(f),
				},
			},
		}))
		if err != nil {
			return err
		}

		red := color.New(color.FgRed).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()

		if len(res.Msg.Diff.Added) == 0 && len(res.Msg.Diff.Removed) == 0 {
			fmt.Println("no changes in access relationships")
		}

		if len(res.Msg.Diff.Added) > 0 {
			fmt.Println(green("The following access relationships would be added by this change:\n"))

			for _, rel := range res.Msg.Diff.Added {
				fmt.Printf(green("+  %s\t>\t%s\t>\t%s\n"), formatUID(rel.Principal), formatUID(rel.Action), formatUID(rel.Resource))
			}

			if len(res.Msg.Diff.Added) > 0 {
				fmt.Printf("\n\n\n")
			}

		}

		if len(res.Msg.Diff.Removed) > 0 {
			fmt.Println(red("The following access relationships would be removed by this change:\n"))

			for _, rel := range res.Msg.Diff.Removed {
				fmt.Printf(red("-  %s\t>\t%s\t>\t%s\n"), formatUID(rel.Principal), formatUID(rel.Action), formatUID(rel.Resource))
			}
		}

		return nil
	},
}

func formatUID(uid *authzv1alpha1.UID) string {
	return fmt.Sprintf(`%s::%s`, uid.Type, uid.Id)
}
