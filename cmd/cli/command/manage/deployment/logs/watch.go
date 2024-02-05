package logs

import (
	"fmt"
	"strings"
	"sync"

	"github.com/TylerBrock/saw/blade"
	sawconfig "github.com/TylerBrock/saw/config"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/common-fate/clio"
	"github.com/urfave/cli/v2"
)

var watchCommand = cli.Command{
	Name: "watch",
	Flags: []cli.Flag{
		&cli.StringSliceFlag{Name: "service", Aliases: []string{"s"}, Usage: "The service to watch logs for. Services: " + strings.Join(ServiceNames, ", "), Required: false},
		&cli.StringFlag{Name: "filter", Usage: "Filter logs using a keyword, see the AWS documentation for details and syntax https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html"},
	},
	Description: "Stream logs from CloudWatch",
	Action: func(c *cli.Context) error {
		services := c.StringSlice("service")
		err := validateServices(services)
		if err != nil {
			return err
		}
		ctx := c.Context
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			return err
		}

		wg := sync.WaitGroup{}
		// if no services supplied, watch all
		if len(services) == 0 {
			services = ServiceNames
		}
		for _, service := range services {

			wg.Add(1)
			go func(lg, s string) {
				clio.Infof("Starting to watch logs for %s, log group id: %s", s, lg)
				watchEvents(lg, cfg.Region, c.String("filter"))
				wg.Done()
			}(fmt.Sprintf("%s-%s-%s", c.String("namespace"), c.String("stage"), service), service)
		}

		wg.Wait()

		return nil
	},
}

func watchEvents(group string, region string, filter string) {
	sawcfg := sawconfig.Configuration{
		Group:  group,
		Filter: filter,
	}

	outputcfg := sawconfig.OutputConfiguration{
		Pretty: true,
	}
	// The Blade api from saw is not very configurable
	// The most we can do is pass in a Region
	b := blade.NewBlade(&sawcfg, &sawconfig.AWSConfiguration{Region: region}, &outputcfg)
	b.StreamEvents()
}
