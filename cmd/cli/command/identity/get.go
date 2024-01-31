package identity

import (
	"errors"
	"fmt"

	"connectrpc.com/connect"
	"github.com/common-fate/cli/treeprint"
	"github.com/common-fate/sdk/config"
	"github.com/common-fate/sdk/eid"
	accessv1alpha1 "github.com/common-fate/sdk/gen/commonfate/access/v1alpha1"
	identitysvc "github.com/common-fate/sdk/service/identity"
	"github.com/urfave/cli/v2"
	"google.golang.org/protobuf/encoding/protojson"
)

var getCommand = cli.Command{
	Name:  "get",
	Usage: "Get the current authenticated identity",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "output", Value: "text", Usage: "output format ('text', 'tree', or 'json')"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		client := identitysvc.NewFromConfig(cfg)

		res, err := client.GetCallerIdentity(ctx, connect.NewRequest(&accessv1alpha1.GetCallerIdentityRequest{}))
		if err != nil {
			return err
		}

		output := c.String("output")
		switch output {
		case "text":
			pid := eid.FromAPI(res.Msg.Principal.Eid)
			fmt.Printf("%s (%s)\n", res.Msg.Principal.Display(), pid)

		case "tree":
			tree := treeprint.New()
			lastNode := tree
			for _, link := range res.Msg.Chain {
				id := eid.FromAPI(link.Id)
				if link.Label != "" {
					lastNode = lastNode.AddMetaBranch(link.Label, id)
				} else {
					lastNode = lastNode.AddBranch(id)
				}
			}

			fmt.Println(tree.String())

		case "json":
			resJSON, err := protojson.Marshal(res.Msg)
			if err != nil {
				return err
			}
			fmt.Println(string(resJSON))

		default:
			return errors.New("invalid --output flag, valid values are [json, text, tree]")
		}

		return nil
	},
}
