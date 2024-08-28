package testcmd

import (
	"fmt"
	"os"

	"github.com/common-fate/clio"
	"github.com/urfave/cli/v2"
)

var createCommand = cli.Command{
	Name:  "create",
	Usage: "Create an example YAML file for defining tests",
	Flags: []cli.Flag{
		&cli.PathFlag{Name: "file", Aliases: []string{"f"}, Required: true, Usage: "the path to the YAML file specifying the tests to run"},
	},
	Action: func(c *cli.Context) error {
		filePath := c.Path("file")

		fileContent := []byte(`# group-tests tests whether a user is a member of a group or not.
group-tests:
  - user: alice@example.com
    group: 2df6c5c4-0e09-477a-b796-9ad9bd756d83
    is-member: true

# access-tests tests whether a user is allowed to access a particular entitlement or not.
access-tests:
  - user: alice@example.com
    target: development-aws-account
    role: AWSReadOnlyAccess
    expected-result: auto-approved

  - user: alice@example.com
    # you can also use Cedar entity IDs here for the target or role
    target: AWS::Account::"123456789012"
    role: AWS::IDC::PermissionSet::"3d06e8e5-93d3-4bfd-ae4f-e5caf43c99ad"
    expected-result: requires-approval

  - user: alice@example.com
    target: very-sensitive-aws-account
    role: AWSReadOnlyAccess
    expected-result: no-access
`)

		err := os.WriteFile(filePath, fileContent, 0644)
		if err != nil {
			return fmt.Errorf("error writing to file %q: %w", filePath, err)
		}

		clio.Successf("created %s", filePath)

		return nil
	},
}
