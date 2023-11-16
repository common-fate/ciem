package open

import (
	"os"
	"os/exec"

	"github.com/common-fate/clio"
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name:        "open",
	Subcommands: []*cli.Command{&aws},
}

var aws = cli.Command{
	Name:        "aws",
	Subcommands: []*cli.Command{&s3},
}

var s3 = cli.Command{
	Name:        "s3",
	Subcommands: []*cli.Command{&bucket},
}

var bucket = cli.Command{
	Name: "bucket",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "account"},
		&cli.StringFlag{Name: "name"},
		&cli.StringFlag{Name: "action"},
	},
	Action: func(c *cli.Context) error {
		cmd := exec.Command("assumego", "-c", "Sandbox-1/GrantedAdministratorAccess", "--cd", "https://s3.console.aws.amazon.com/s3/buckets/example-bucket-to-access", "--browser-profile", "prod/AdaptiveAccess")
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "GRANTED_ALIAS_CONFIGURED=true")
		err := cmd.Run()
		if err != nil {
			return err
		}

		clio.Successf("opened AWS console")
		return nil
	},
}
