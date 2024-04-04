package policyset

import (
	"context"
	"fmt"
	"os"

	"connectrpc.com/connect"
	"github.com/common-fate/cli/doubleglob"
	"github.com/common-fate/sdk/config"
	authzv1alpha1 "github.com/common-fate/sdk/gen/commonfate/authz/v1alpha1"
	"github.com/common-fate/sdk/gen/commonfate/authz/v1alpha1/authzv1alpha1connect"
	"github.com/common-fate/sdk/service/authz/validation"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

var validateCommand = cli.Command{
	Name: "validate",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "file-pattern", Usage: "Policy set file(s) to validate. May be a glob.", Value: "**/*.cedar"},
		&cli.StringFlag{Name: "file-pattern-dir", Usage: "Directory to use as the base for the --file-pattern argument", Value: "."},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		client := validation.NewFromConfig(cfg)

		fs := os.DirFS(c.String("file-pattern-dir"))
		files, err := doubleglob.Glob(fs, c.String("file-pattern"))
		if err != nil {
			return err
		}

		if len(files) == 1 {
			fmt.Printf("validating 1 policy\n")
		} else {
			fmt.Printf("validating %v policies\n", len(files))
		}

		results := map[string]*authzv1alpha1.ValidationResult{}

		for _, fp := range files {
			results[fp] = validatePolicy(ctx, client, fp)
		}

		var hasIssues bool

		// print a summary for each policy
		for _, fp := range files {
			result := results[fp]

			status := color.GreenString("ok")
			if len(result.Warnings) > 0 {
				status = color.YellowString("warn")
				hasIssues = true
			}
			if len(result.Errors) > 0 {
				status = color.RedString("FAILED")
				hasIssues = true
			}

			fmt.Printf("%s ... %s\n", fp, status)
		}

		if !hasIssues {
			os.Exit(0)
		}

		// print any warnings and failures
		fmt.Printf("\nissues:\n\n")

		var hasErrors bool

		for _, fp := range files {
			result := results[fp]

			if len(result.Errors) == 0 && len(result.Warnings) == 0 {
				// don't print the filename out if there are no warnings or errors
				continue
			}

			fmt.Printf("---- %s ----\n\n", fp)

			for _, w := range result.Warnings {
				if w.PolicyId != "" {
					fmt.Printf("\t[WARNING] %s: %s\n\n", w.PolicyId, w.Message)
				} else {
					fmt.Printf("\t[WARNING] %s\n\n", w.Message)
				}
			}

			for _, w := range result.Errors {
				hasErrors = true

				if w.PolicyId != "" {
					fmt.Printf("\t[ERROR] %s: %s\n\n", w.PolicyId, w.Message)
				} else {
					fmt.Printf("\t[ERROR] %s\n\n", w.Message)
				}
			}
		}

		if hasErrors {
			os.Exit(1)
		}

		return nil
	},
}

func validatePolicy(ctx context.Context, client authzv1alpha1connect.ValidationServiceClient, filename string) *authzv1alpha1.ValidationResult {
	f, err := os.ReadFile(filename)
	if err != nil {
		return &authzv1alpha1.ValidationResult{
			Errors: []*authzv1alpha1.ValidationError{
				{
					Message: fmt.Sprintf("error loading policy file: %s", err.Error()),
				},
			},
		}
	}

	res, err := client.ValidatePolicySetText(ctx, connect.NewRequest(&authzv1alpha1.ValidatePolicySetTextRequest{
		PolicySetText: string(f),
	}))
	if err != nil {
		return &authzv1alpha1.ValidationResult{
			Errors: []*authzv1alpha1.ValidationError{
				{
					Message: fmt.Sprintf("error performing policy validation: %s", err.Error()),
				},
			},
		}
	}

	return res.Msg.Result
}
