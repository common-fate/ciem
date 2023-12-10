package policy

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/common-fate/clio"
	"github.com/common-fate/sdk/service/authz/policyset"
	"github.com/urfave/cli/v2"
	"golang.org/x/net/http2"
)

type policyFile struct {
	input    policyset.Input
	filePath string
}

var applyCommand = cli.Command{
	Name: "apply",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "glob", Value: "*.cedar", Usage: "File pattern to match"},
		&cli.BoolFlag{Name: "dry-run", Usage: "Print policies to apply but don't make any changes"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		client := policyset.NewClient(policyset.Opts{
			HTTPClient: newInsecureClient(),
			BaseURL:    "http://127.0.0.1:5050",
		})

		policysets := map[string]policyFile{}

		pattern := c.String("glob")

		dryRun := c.Bool("dry-run")

		// Find all the files with the specified glob pattern in the current directory
		files, err := filepath.Glob(pattern)
		if err != nil {
			return err
		}

		for _, filePath := range files {
			// Process the file here or add it to a list
			clio.Debugf("Found .cedar file:", filePath)

			// use the basename of the file as the ID of the policy
			baseName := path.Base(filePath)
			id := strings.TrimSuffix(baseName, path.Ext(baseName))

			if existing, ok := policysets[id]; ok {
				return fmt.Errorf("found two policies with conflicting filenames (%s and %s): please rename one of the policy files", existing.filePath, filePath)
			}

			f, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}

			policysets[id] = policyFile{
				filePath: filePath,
				input: policyset.Input{
					ID:   id,
					Text: string(f),
				},
			}

			if dryRun {
				clio.Infof("would apply %s: %s", id, filePath)
			} else {
				clio.Infof("applying %s: %s", id, filePath)
			}
		}

		if len(policysets) == 0 {
			return fmt.Errorf("no policies found matching pattern %s", pattern)
		}

		if dryRun {
			// return early and don't apply anything
			return nil
		}

		var input policyset.BatchPutInput

		for _, ps := range policysets {
			input.PolicySets = append(input.PolicySets, ps.input)
		}

		_, err = client.BatchPut(ctx, input)
		if err != nil {
			return err
		}

		clio.Success("applied Policy Sets")
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
