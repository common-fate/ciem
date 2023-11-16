package use

import (
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/common-fate/clio"
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name:        "use",
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
		var result string
		err := survey.AskOne(&survey.Select{
			Message: "What Jira ticket are you needing this access for?",
			Options: []string{
				"CF-117: Deploy new static website",
				"CF-118: Optimise API POST handling",
				"CF-119: Redrive message queue",
				"CF-120: Add audit trail to user logins",
			},
		}, &result)
		if err != nil {
			return err
		}

		time.Sleep(time.Second)

		clio.Successf("access is approved for the next 2h")
		clio.Infof("to use the permissions with the AWS CLI, run: 'export AWS_PROFILE=prod' or 'assume prod'")
		return nil
	},
}
