package policy

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/common-fate/sdk/service/authz/policyset"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

var listCommand = cli.Command{
	Name: "list",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "output", Value: "text", Usage: "output format ('text', or 'json')"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		client := policyset.NewClient(policyset.Opts{
			HTTPClient: newInsecureClient(),
			BaseURL:    "http://127.0.0.1:5050",
		})

		var out policyset.ListOutput

		call := client.ListPolicySetsRequest(policyset.ListInput{})

		err := call.Pages(ctx, func(lpso policyset.ListOutput) error {
			out.PolicySets = append(out.PolicySets, lpso.PolicySets...)
			return nil
		})
		if err != nil {
			return err
		}

		output := c.String("output")
		switch output {
		case "json":
			outJSON, err := json.Marshal(out.PolicySets)
			if err != nil {
				return err
			}

			fmt.Println(string(outJSON))

		case "text":
			boldGreen := color.New(color.Bold, color.FgBlue)
			for _, ps := range out.PolicySets {
				for _, p := range ps.Policies {
					boldGreen.Println(p.ID)
					fmt.Printf("%s\n\n", p.Text)
				}
			}

		default:
			return errors.New("invalid --output flag, valid values are [json, text]")
		}

		return nil
	},
}
