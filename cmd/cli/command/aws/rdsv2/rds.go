package rdsv2

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"connectrpc.com/connect"
	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/briandowns/spinner"
	accessCmd "github.com/common-fate/cli/cmd/cli/command/access"
	"github.com/common-fate/cli/printdiags"
	"github.com/common-fate/clio"
	"github.com/common-fate/clio/clierr"
	"github.com/common-fate/grab"
	"github.com/common-fate/granted/pkg/assume"
	"github.com/common-fate/sdk/config"
	"github.com/common-fate/sdk/eid"
	databaseproxyv1alpha1 "github.com/common-fate/sdk/gen/commonfate/access/databaseproxy/v1alpha1"
	"github.com/common-fate/sdk/gen/commonfate/access/databaseproxy/v1alpha1/databaseproxyv1alpha1connect"
	accessv1alpha1 "github.com/common-fate/sdk/gen/commonfate/access/v1alpha1"
	entityv1alpha1 "github.com/common-fate/sdk/gen/commonfate/entity/v1alpha1"
	"github.com/common-fate/sdk/service/access"
	"github.com/common-fate/sdk/service/access/grants"
	"github.com/common-fate/sdk/service/entity"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name:  "rds",
	Usage: "Perform RDS Operations",
	Subcommands: []*cli.Command{
		&proxyCommand,
	},
}

var proxyCommand = cli.Command{
	Name:  "proxy",
	Usage: "Run a database proxy",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "target", Required: true},
		&cli.StringFlag{Name: "role", Required: true},
		&cli.BoolFlag{Name: "confirm", Aliases: []string{"y"}, Usage: "skip the confirmation prompt"},
		&cli.IntFlag{Name: "port", Value: 3306, Usage: "The local port to forward the database connection to"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		// ensure required CLI tools are installed
		err = CheckDependencies()
		if err != nil {
			return err
		}

		target := c.String("target")
		role := c.String("role")
		client := access.NewFromConfig(cfg)
		apiURL, err := url.Parse(cfg.APIURL)
		if err != nil {
			return err
		}
		req := accessv1alpha1.BatchEnsureRequest{
			Entitlements: []*accessv1alpha1.EntitlementInput{
				{
					Target: &accessv1alpha1.Specifier{
						Specify: &accessv1alpha1.Specifier_Lookup{
							Lookup: target,
						},
					},
					Role: &accessv1alpha1.Specifier{
						Specify: &accessv1alpha1.Specifier_Lookup{
							Lookup: role,
						},
					},
				},
			},
			DryRun: !c.Bool("confirm"),
		}
		var ensuredGrant *accessv1alpha1.GrantState
		for {
			var hasChanges bool
			si := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
			si.Suffix = grab.If(req.DryRun, " planning access changes...", " ensuring access...")
			si.Writer = os.Stderr
			si.Start()

			res, err := client.BatchEnsure(ctx, connect.NewRequest(&req))
			if err != nil {
				si.Stop()
				return err
			}

			si.Stop()

			clio.Debugw("BatchEnsure response", "response", res)

			names := map[eid.EID]string{}
			for _, g := range res.Msg.Grants {
				names[eid.New("Access::Grant", g.Grant.Id)] = g.Grant.Name

				exp := "<invalid expiry>"

				if g.Grant.ExpiresAt != nil {
					exp = accessCmd.ShortDur(time.Until(g.Grant.ExpiresAt.AsTime()))
				}
				if g.Change > 0 {
					hasChanges = true
				}

				switch g.Change {

				case accessv1alpha1.GrantChange_GRANT_CHANGE_ACTIVATED:
					if req.DryRun {
						color.New(color.BgHiGreen).Printf("[WILL ACTIVATE]")
						color.New(color.FgGreen).Printf(" %s will be activated for %s: %s\n", g.Grant.Name, exp, accessCmd.RequestURL(apiURL, g.Grant))
					} else {
						ensuredGrant = g
						color.New(color.BgHiGreen).Printf("[ACTIVATED]")
						color.New(color.FgGreen).Printf(" %s was activated for %s: %s\n", g.Grant.Name, exp, accessCmd.RequestURL(apiURL, g.Grant))
					}
					continue
				case accessv1alpha1.GrantChange_GRANT_CHANGE_EXTENDED:
					if req.DryRun {
						color.New(color.BgBlue).Printf("[WILL EXTEND]")
						color.New(color.FgBlue).Printf(" %s will be extended for another %s: %s\n", g.Grant.Name, exp, accessCmd.RequestURL(apiURL, g.Grant))
					} else {
						ensuredGrant = g
						color.New(color.BgBlue).Printf("[EXTENDED]")
						color.New(color.FgBlue).Printf(" %s was extended for another %s: %s\n", g.Grant.Name, exp, accessCmd.RequestURL(apiURL, g.Grant))
					}
					continue
				case accessv1alpha1.GrantChange_GRANT_CHANGE_REQUESTED:
					if req.DryRun {
						color.New(color.BgHiYellow, color.FgBlack).Printf("[WILL REQUEST]")
						color.New(color.FgYellow).Printf(" %s will require approval\n", g.Grant.Name)
					} else {
						color.New(color.BgHiYellow, color.FgBlack).Printf("[REQUESTED]")
						color.New(color.FgYellow).Printf(" %s requires approval: %s\n", g.Grant.Name, accessCmd.RequestURL(apiURL, g.Grant))
					}
					continue
				case accessv1alpha1.GrantChange_GRANT_CHANGE_PROVISIONING_FAILED:
					if req.DryRun {
						// shouldn't happen in the dry-run request but handle anyway
						color.New(color.FgRed).Printf("[ERROR] %s will fail provisioning\n", g.Grant.Name)
					} else {
						// shouldn't happen in the dry-run request but handle anyway
						color.New(color.FgRed).Printf("[ERROR] %s failed provisioning: %s\n", g.Grant.Name, accessCmd.RequestURL(apiURL, g.Grant))
					}
					continue
				}

				switch g.Grant.Status {
				case accessv1alpha1.GrantStatus_GRANT_STATUS_ACTIVE:
					ensuredGrant = g
					color.New(color.FgGreen).Printf("[ACTIVE] %s is already active for the next %s: %s\n", g.Grant.Name, exp, accessCmd.RequestURL(apiURL, g.Grant))
					continue
				case accessv1alpha1.GrantStatus_GRANT_STATUS_PENDING:
					color.New(color.FgWhite).Printf("[PENDING] %s is already pending: %s\n", g.Grant.Name, accessCmd.RequestURL(apiURL, g.Grant))
					continue
				case accessv1alpha1.GrantStatus_GRANT_STATUS_CLOSED:
					color.New(color.FgWhite).Printf("[CLOSED] %s is closed but was still returned: %s\n. This is most likely due to an error in Common Fate and should be reported to our team: support@commonfate.io.", g.Grant.Name, accessCmd.RequestURL(apiURL, g.Grant))
					continue
				}

				color.New(color.FgWhite).Printf("[UNSPECIFIED] %s is in an unspecified status: %s\n. This is most likely due to an error in Common Fate and should be reported to our team: support@commonfate.io.", g.Grant.Name, accessCmd.RequestURL(apiURL, g.Grant))
			}

			printdiags.Print(res.Msg.Diagnostics, names)

			if req.DryRun && hasChanges {
				if !accessCmd.IsTerminal(os.Stdin.Fd()) {
					return errors.New("detected a noninteractive terminal: to apply the planned changes please re-run 'cf access ensure' with the --confirm flag")
				}

				confirm := survey.Confirm{
					Message: "Apply proposed access changes",
				}
				var proceed bool
				err = survey.AskOne(&confirm, &proceed)
				if err != nil {
					return err
				}
				if !proceed {
					clio.Info("no access changes")
				}
				req.DryRun = false
				continue
			} else {
				break
			}
		}

		// if its not yet active, we can just exit the process
		if ensuredGrant == nil {
			clio.Debug("exiting because grant status is not active, or a grant was not found")
			return nil
		}

		grantsClient := grants.NewFromConfig(cfg)

		children, err := grab.AllPages(ctx, func(ctx context.Context, nextToken *string) ([]*entityv1alpha1.Entity, *string, error) {
			res, err := grantsClient.QueryGrantChildren(ctx, connect.NewRequest(&accessv1alpha1.QueryGrantChildrenRequest{
				Id:        ensuredGrant.Grant.Id,
				PageToken: grab.Value(nextToken),
			}))
			if err != nil {
				return nil, nil, err
			}
			return res.Msg.Entities, &res.Msg.NextPageToken, nil
		})
		if err != nil {
			return err
		}

		commandData := CommandData{
			LocalPort: strconv.Itoa((c.Int("port"))),
			ProxyPort: "3307",
		}

		for _, child := range children {
			if child.Eid.Type == GrantOutputType {
				err = entity.Unmarshal(child, &commandData.GrantOutput)
				if err != nil {
					return err
				}
			}
		}

		if commandData.GrantOutput.Grant.ID == "" {
			return errors.New("did not find a grant output entity in query grant children response")
		}
		creds, err := GrantedCredentialProcess(commandData)
		if err != nil {
			return err
		}
		awsCfg, err := awsConfig.LoadDefaultConfig(ctx, awsConfig.WithCredentialsProvider(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return creds, nil
		})), awsConfig.WithRegion(commandData.GrantOutput.ProxyRegion))
		if err != nil {
			return err
		}

		ecsClient := ecs.NewFromConfig(awsCfg)

		listTasksResp, err := ecsClient.ListTasks(ctx, &ecs.ListTasksInput{
			Cluster:     aws.String("common-fate-demo-cluster"),
			ServiceName: aws.String("common-fate-demo-demo-rds-proxy"),
		})
		if err != nil {
			return fmt.Errorf("failed to list tasks, %w", err)
		}

		// Describe tasks
		describeTasksResp, err := ecsClient.DescribeTasks(ctx, &ecs.DescribeTasksInput{
			Cluster: aws.String("common-fate-demo-cluster"),
			Tasks:   listTasksResp.TaskArns,
		})
		if err != nil {
			return fmt.Errorf("failed to describe tasks, %w", err)
		}

		if len(describeTasksResp.Tasks) > 1 {
			return fmt.Errorf("expected only one task to be returned")
		}
		task := describeTasksResp.Tasks[0]

		for _, container := range task.Containers {
			if grab.Value(container.Name) == "aws-rds-proxy-container" {
				runtimeID := grab.Value(container.RuntimeId)
				taskID := strings.Split(runtimeID, "-")[0]
				commandData.SSMSessionTarget = fmt.Sprintf("ecs:common-fate-demo-cluster_%s_%s", taskID, runtimeID)
			}
		}

		passwordExchangeData := commandData
		passwordExchangeData.ProxyPort = "9999"
		// @todo mayby find another open port automatically
		passwordExchangeData.LocalPort = "9999"

		notifyCh := make(chan struct{}, 1) // Buffer of 1 to prevent blocking
		cmd := exec.Command("aws", formatSSMCommandArgs(passwordExchangeData)...)
		clio.Info("running aws ssm command", "command", "aws "+strings.Join(formatSSMCommandArgs(passwordExchangeData), " "))
		// might need to handle errors better here, the idea is to wait till we can exchange the auth token for the password
		cmd.Stderr = io.MultiWriter(DebugWriter{}, NewNotifyingWriter(io.Discard, "Waiting for connections", notifyCh))
		cmd.Stdout = io.MultiWriter(DebugWriter{}, NewNotifyingWriter(io.Discard, "Waiting for connections", notifyCh))
		cmd.Stdin = os.Stdin
		cmd.Env = PrepareAWSCLIEnv(creds, passwordExchangeData)

		// Start the command in a separate goroutine
		err = cmd.Start()
		if err != nil {
			log.Fatal(err)
		}

		select {
		case <-notifyCh:
		case <-time.After(time.Second * 15):
			return errors.New("timed out waiting for password exchange from proxy server, you can try running again with --verbose flag to see debugging logs")
		}

		defer func() {
			err = cmd.Process.Signal(os.Interrupt)
			if err != nil {
				clio.Error(err)
			}
		}()

		// once we see the signal, we can start the token exchange
		dbClient := databaseproxyv1alpha1connect.NewDatabaseProxyServiceClient(cfg.HTTPClient, "http://localhost:9999")
		exchange, err := dbClient.Exchange(ctx, connect.NewRequest(&databaseproxyv1alpha1.ExchangeRequest{
			GrantId: ensuredGrant.Grant.Id,
		}))
		if err != nil {
			return err
		}
		// end the exchange proxy
		err = cmd.Process.Signal(os.Interrupt)
		if err != nil {
			return err
		}

		clio.Infof("starting database proxy on port %v", commandData.LocalPort)
		clio.Infof("You can connect to the database using this connection string '%s:%s@tcp(127.0.0.1:%s)/%s?allowCleartextPasswords=1'", exchange.Msg.DatabaseUser, exchange.Msg.DatabasePassword, commandData.LocalPort, exchange.Msg.DatabaseName)
		cmd = exec.Command("aws", formatSSMCommandArgs(commandData)...)
		clio.Debugw("running aws ssm command", "command", "aws "+strings.Join(formatSSMCommandArgs(commandData), " "))
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Env = PrepareAWSCLIEnv(creds, commandData)

		// Start the command in a separate goroutine
		err = cmd.Start()
		if err != nil {
			return err
		}

		// Set up a channel to receive OS signals
		sigs := make(chan os.Signal, 1)
		// Notify sigs on os.Interrupt (Ctrl+C)
		signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

		// Wait for a termination signal in a separate goroutine
		go func() {
			<-sigs
			clio.Info("Received interrupt signal, shutting down...")
			if err := cmd.Process.Signal(os.Interrupt); err != nil {
				clio.Error("Error sending SIGTERM to process:", err)
			}
		}()

		// Wait for the command to finish
		err = cmd.Wait()
		if err != nil {
			clio.Error("Proxy connection failed with error:", err)
		} else {
			clio.Info("Proxy connection closed successfully")
		}
		return nil
	},
}

// DebugWriter is an io.Writer that writes messages using clio.Debug.
type DebugWriter struct{}

// Write implements the io.Writer interface for DebugWriter.
func (dw DebugWriter) Write(p []byte) (n int, err error) {
	message := string(p)
	clio.Debug(message)
	return len(p), nil
}

type NotifyingWriter struct {
	writer   io.Writer
	phrase   string
	notifyCh chan struct{}
	buffer   bytes.Buffer
}

func NewNotifyingWriter(writer io.Writer, phrase string, notifyCh chan struct{}) *NotifyingWriter {
	return &NotifyingWriter{
		writer:   writer,
		phrase:   phrase,
		notifyCh: notifyCh,
	}
}

func (nw *NotifyingWriter) Write(p []byte) (n int, err error) {
	// Write to the buffer first
	nw.buffer.Write(p)
	// Check if the phrase is in the buffer
	if strings.Contains(nw.buffer.String(), nw.phrase) {
		// Notify the channel in a non-blocking way
		select {
		case nw.notifyCh <- struct{}{}:
		default:
		}
		// Clear the buffer up to the phrase
		nw.buffer.Reset()
	}
	// Write to the underlying writer
	return nw.writer.Write(p)
}

func PrepareAWSCLIEnv(creds aws.Credentials, commandData CommandData) []string {
	return append(SanitisedEnv(), assume.EnvKeys(creds, commandData.GrantOutput.ProxyRegion)...)
}

// SanitisedEnv returns the environment variables excluding specific AWS keys.
// used so that existing aws creds in the terminal are not passed through to downstream programs like the AWS cli
func SanitisedEnv() []string {
	// List of AWS keys to remove from the environment.
	awsKeys := map[string]struct{}{
		"AWS_ACCESS_KEY_ID":         {},
		"AWS_SECRET_ACCESS_KEY":     {},
		"AWS_SESSION_TOKEN":         {},
		"AWS_PROFILE":               {},
		"AWS_REGION":                {},
		"AWS_DEFAULT_REGION":        {},
		"AWS_SESSION_EXPIRATION":    {},
		"AWS_CREDENTIAL_EXPIRATION": {},
	}

	var cleanedEnv []string
	for _, env := range os.Environ() {
		// Split the environment variable into key and value
		parts := strings.SplitN(env, "=", 2)
		key := parts[0]

		// If the key is not one of the AWS keys, include it in the cleaned environment
		if _, found := awsKeys[key]; !found {
			cleanedEnv = append(cleanedEnv, env)
		}
	}
	return cleanedEnv
}

type CommandData struct {
	GrantOutput      GrantOutput
	LocalPort        string
	ProxyPort        string
	SSMSessionTarget string
}

func formatSSMCommandArgs(data CommandData) []string {
	out := []string{
		"ssm",
		"start-session",
		fmt.Sprintf("--target=%s", data.SSMSessionTarget),
		"--document-name=AWS-StartPortForwardingSession",
		"--parameters",
		fmt.Sprintf(`{"portNumber":["%s"], "localPortNumber":["%s"]}`, data.ProxyPort, data.LocalPort),
	}

	return out
}

// CredentialProcessOutput represents the JSON output format of the credential process.
type CredentialProcessOutput struct {
	Version         int       `json:"Version"`
	AccessKeyId     string    `json:"AccessKeyId"`
	SecretAccessKey string    `json:"SecretAccessKey"`
	SessionToken    string    `json:"SessionToken,omitempty"`
	Expiration      time.Time `json:"Expiration,omitempty"`
}

// ParseCredentialProcessOutput parses the JSON output of a credential process and returns aws.Credentials.
func ParseCredentialProcessOutput(credentialProcessOutput string) (aws.Credentials, error) {
	var output CredentialProcessOutput
	err := json.Unmarshal([]byte(credentialProcessOutput), &output)
	if err != nil {
		return aws.Credentials{}, fmt.Errorf("error parsing credential process output: %w", err)
	}

	return aws.Credentials{
		AccessKeyID:     output.AccessKeyId,
		SecretAccessKey: output.SecretAccessKey,
		SessionToken:    output.SessionToken,
		CanExpire:       !output.Expiration.IsZero(),
		Expires:         output.Expiration,
	}, nil
}

func CheckDependencies() error {
	_, err := exec.LookPath("granted")
	if err != nil {
		// The executable was not found in the PATH
		if _, ok := err.(*exec.Error); ok {
			return clierr.New("the required cli 'granted' was not found on your path", clierr.Info("Granted is required to access AWS via SSO, please follow the instructions here to install it https://docs.commonfate.io/granted/getting-started/"))
		}
		return err
	}
	_, err = exec.LookPath("aws")
	if err != nil {
		// The executable was not found in the PATH
		if _, ok := err.(*exec.Error); ok {
			return clierr.New("the required cli 'aws' was not found on your path", clierr.Info("The AWS cli is required to access dastabases via SSM Session Manager, please follow the instructions here to install it https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-welcome.html"))
		}
		return err
	}
	return nil
}

func GrantedCredentialProcess(commandData CommandData) (aws.Credentials, error) {
	// the grant id is used for teh profile to avoid issues with the credential cache in granted credential-process, it also gets the benefit of this cache per grant
	configFile := fmt.Sprintf(`[profile %s]
sso_account_id = %s
sso_role_name = %s
sso_start_url = %s
sso_region = %s
region = %s
`, commandData.GrantOutput.Grant.ID, commandData.GrantOutput.ProxyAccountID, commandData.GrantOutput.SSORoleName, commandData.GrantOutput.SSOStartURL, commandData.GrantOutput.SSORegion, commandData.GrantOutput.ProxyRegion)

	file, err := os.CreateTemp(os.TempDir(), "")
	if err != nil {
		return aws.Credentials{}, err
	}
	defer file.Close()
	defer os.Remove(file.Name())
	clio.Debugf("temporary config file generated at %s\n\n%s", file.Name(), configFile)
	_, err = file.Write([]byte(configFile))
	if err != nil {
		return aws.Credentials{}, err
	}
	err = file.Close()
	if err != nil {
		return aws.Credentials{}, err
	}

	cmd := exec.Command("granted", "credential-process", "--auto-login", "--profile", commandData.GrantOutput.Grant.ID)
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env, "AWS_CONFIG_FILE="+file.Name())

	var stdout strings.Builder
	var stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	clio.Debugw("granted credentials process stderr output", "stderr", stderr.String())
	if err != nil {
		return aws.Credentials{}, err
	}
	return ParseCredentialProcessOutput(stdout.String())
}