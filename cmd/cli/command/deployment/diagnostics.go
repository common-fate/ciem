package deployment

import (
	"errors"
	"fmt"
	"os"
	"time"

	"connectrpc.com/connect"
	"github.com/common-fate/cli/table"
	"github.com/common-fate/sdk/config"
	diagnosticv1alpha1 "github.com/common-fate/sdk/gen/commonfate/control/diagnostic/v1alpha1"
	"github.com/common-fate/sdk/service/control/diagnostic"
	"github.com/urfave/cli/v2"
	"google.golang.org/protobuf/encoding/protojson"
)

var diagnosticsCommand = cli.Command{
	Name:  "diagnostics",
	Usage: "Retrieve diagnostics about your deployment",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "output", Value: "text", Usage: "output format ('text' or 'json')"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		outputFormat := c.String("output")

		if outputFormat != "text" && outputFormat != "json" {
			return errors.New("--output flag must be either 'text' or 'json'")
		}

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		// fetch each diagnostic through separate API calls, and then combine them here
		// into the main 'full' set of diagnostics.
		//
		// This has been implemented to avoid having a single massive 'AllDiagnostics' endpoint
		// which may be expensive to call - plus, if something is wrong the entire API may return with an error.
		var all diagnosticv1alpha1.AllDiagnostics

		client := diagnostic.NewFromConfig(cfg)

		tokenMetadata, err := client.GetOAuthTokenMetadata(ctx, connect.NewRequest(&diagnosticv1alpha1.GetOAuthTokenMetadataRequest{}))
		if err != nil {
			return err
		}

		all.OauthTokenMetadata = tokenMetadata.Msg

		switch outputFormat {
		case "text":
			fmt.Println("OAUTH TOKEN METADATA")
			tbl := table.New(os.Stdout)
			tbl.Columns("ID", "APPNAME", "EXPIRES")

			for _, t := range all.OauthTokenMetadata.Tokens {
				exp := "-"

				if !t.ExpiresAt.AsTime().IsZero() {
					exp = t.ExpiresAt.AsTime().Format(time.RFC3339)
				}

				tbl.Row(t.Id, t.AppName, exp)
			}

			err = tbl.Flush()
			if err != nil {
				return err
			}

		case "json":
			resJSON, err := protojson.Marshal(&all)
			if err != nil {
				return err
			}
			fmt.Println(string(resJSON))
		}

		return nil
	},
}
