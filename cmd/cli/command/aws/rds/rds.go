package rds

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"connectrpc.com/connect"
	"go.uber.org/zap"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/briandowns/spinner"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	accessCmd "github.com/common-fate/cli/cmd/cli/command/access"
	"github.com/common-fate/cli/printdiags"
	"github.com/common-fate/clio"
	"github.com/common-fate/clio/clierr"
	"github.com/common-fate/grab"
	"github.com/common-fate/granted/pkg/assume"

	"github.com/common-fate/sdk/config"
	"github.com/common-fate/sdk/eid"
	accessv1alpha1 "github.com/common-fate/sdk/gen/commonfate/access/v1alpha1"
	entityv1alpha1 "github.com/common-fate/sdk/gen/commonfate/entity/v1alpha1"
	"github.com/common-fate/sdk/handshake"
	"github.com/common-fate/sdk/service/access"
	"github.com/common-fate/sdk/service/access/grants"
	"github.com/common-fate/sdk/service/entity"
	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
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
		&cli.StringFlag{Name: "target"},
		&cli.StringFlag{Name: "role"},
		&cli.BoolFlag{Name: "confirm", Aliases: []string{"y"}, Usage: "skip the confirmation prompt"},
		&cli.IntFlag{Name: "mysql-port", Value: 3306, Usage: "The local port to forward the mysql database connection to"},
		&cli.IntFlag{Name: "postgres-port", Value: 5432, Usage: "The local port to forward the postgres database connection to"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		err = cfg.Initialize(ctx, config.InitializeOpts{})
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

		if target == "" && role == "" {
			entitlements, err := grab.AllPages(ctx, func(ctx context.Context, nextToken *string) ([]*accessv1alpha1.Entitlement, *string, error) {
				res, err := client.QueryEntitlements(ctx, connect.NewRequest(&accessv1alpha1.QueryEntitlementsRequest{
					PageToken:  grab.Value(nextToken),
					TargetType: grab.Ptr("AWS::RDS::Database"),
				}))
				if err != nil {
					return nil, nil, err
				}
				return res.Msg.Entitlements, &res.Msg.NextPageToken, nil
			})
			if err != nil {
				return err
			}

			type Column struct {
				Title string
				Width int
			}
			cols := []Column{{Title: "Database", Width: 40}, {Title: "Role", Width: 40}}
			var s = make([]string, 0, len(cols))
			for _, col := range cols {
				style := lipgloss.NewStyle().Width(col.Width).MaxWidth(col.Width).Inline(true)
				renderedCell := style.Render(runewidth.Truncate(col.Title, col.Width, "…"))
				s = append(s, lipgloss.NewStyle().Bold(true).Padding(0).Render(renderedCell))
			}
			header := lipgloss.NewStyle().PaddingLeft(2).Render(lipgloss.JoinHorizontal(lipgloss.Left, s...))
			var options []huh.Option[*accessv1alpha1.Entitlement]

			for _, entitlement := range entitlements {
				style := lipgloss.NewStyle().Width(cols[0].Width).MaxWidth(cols[0].Width).Inline(true)
				target := lipgloss.NewStyle().Bold(true).Padding(0).Render(style.Render(runewidth.Truncate(entitlement.Target.Display(), cols[0].Width, "…")))

				style = lipgloss.NewStyle().Width(cols[1].Width).MaxWidth(cols[1].Width).Inline(true)
				role := lipgloss.NewStyle().Bold(true).Padding(0).Render(style.Render(runewidth.Truncate(entitlement.Role.Display(), cols[1].Width, "…")))

				options = append(options, huh.Option[*accessv1alpha1.Entitlement]{
					Key:   lipgloss.JoinHorizontal(lipgloss.Left, target, role),
					Value: entitlement,
				})
			}

			selector := huh.NewSelect[*accessv1alpha1.Entitlement]().
				Options(options...).
				Title("Select a database to connect to").
				Description(header).WithTheme(huh.ThemeBase())
			err = selector.Run()
			if err != nil {
				return err
			}

			entitlement := selector.GetValue().(*accessv1alpha1.Entitlement)

			target = entitlement.Target.Eid.Display()
			role = entitlement.Role.Eid.Display()

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

		// find an unused local port to use for the ssm server
		// the user doesn't directly connect to this, they connect through our local proxy
		// which adds authentication
		ssmPortforwardLocalPort, err := GrabUnusedPort()
		if err != nil {
			return err
		}

		clio.Debugf("starting SSM portforward on local port: %s", ssmPortforwardLocalPort)

		commandData := CommandData{
			// the proxy server always runs on port 8080
			SSMPortForwardServerPort: "8080",
			SSMPortForwardLocalPort:  ssmPortforwardLocalPort,
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

		// @TODO consider embedding granted here instead of having the external dependency
		creds, err := GrantedCredentialProcess(commandData)
		if err != nil {
			return err
		}

		// the port that the user connects to
		mysqlPort := strconv.Itoa((c.Int("mysql-port")))
		postgresPort := strconv.Itoa((c.Int("postgres-port")))

		notifyCh := make(chan struct{})

		var cmd *exec.Cmd

		// in local dev you can skip using ssm and just use a local port forward instead
		if os.Getenv("CF_DEV_PROXY") == "true" {
			cmd = exec.Command("socat", fmt.Sprintf("TCP-LISTEN:%s,fork", ssmPortforwardLocalPort), "TCP:127.0.0.1:8081")
			go func() { notifyCh <- struct{}{} }()
		} else {
			cmd = exec.Command("aws", formatSSMCommandArgs(commandData)...)
		}

		clio.Debugw("running aws ssm command", "command", "aws "+strings.Join(formatSSMCommandArgs(commandData), " "))

		si := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		si.Suffix = " Starting database proxy..."
		si.Writer = os.Stderr
		si.Start()
		defer si.Stop()

		cmd.Stderr = io.MultiWriter(NewNotifyingWriter(io.Discard, "Waiting for connections...", notifyCh), DebugWriter{})
		cmd.Stdout = io.MultiWriter(NewNotifyingWriter(io.Discard, "Waiting for connections...", notifyCh), DebugWriter{})
		cmd.Stdin = os.Stdin
		cmd.Env = PrepareAWSCLIEnv(creds, commandData)

		// Start the command in a separate goroutine
		err = cmd.Start()
		if err != nil {
			return err
		}

		// listen for interrupt signals and forward them on
		// listen for a context cancellation

		// Set up a channel to receive OS signals
		sigs := make(chan os.Signal, 1)
		// Notify sigs on os.Interrupt (Ctrl+C)
		signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

		ctx, cancel := context.WithCancel(ctx)
		eg, ctx := errgroup.WithContext(ctx)

		eg.Go(func() error {
			select {
			case <-notifyCh:
				si.Stop()
			case <-ctx.Done():
				si.Stop()
				return nil
			}

			var connectionString, cliString, port string
			yellow := color.New(color.FgYellow)
			switch commandData.GrantOutput.Database.Engine {
			case "postgres":
				port = postgresPort
				connectionString = yellow.Sprintf("postgresql://%s:password@127.0.0.1:%s/%s?sslmode=disable", commandData.GrantOutput.User.Username, postgresPort, commandData.GrantOutput.Database.Database)
				cliString = yellow.Sprintf(`psql "postgresql://%s:password@127.0.0.1:%s/%s?sslmode=disable"`, commandData.GrantOutput.User.Username, postgresPort, commandData.GrantOutput.Database.Database)
			case "mysql":
				port = mysqlPort
				connectionString = yellow.Sprintf("%s:password@tcp(127.0.0.1:%s)/%s", commandData.GrantOutput.User.Username, mysqlPort, commandData.GrantOutput.Database.Database)
				cliString = yellow.Sprintf(`mysql -u %s -p'password' -h 127.0.0.1 -P %s %s`, commandData.GrantOutput.User.Username, mysqlPort, commandData.GrantOutput.Database.Database)
			default:
				return fmt.Errorf("unsupported database engine: %s, maybe you need to update your `cf` cli", commandData.GrantOutput.Database.Engine)
			}

			clio.NewLine()
			clio.Infof("Database proxy ready for connections on 127.0.0.1:%s", port)
			clio.NewLine()

			clio.Infof("You can connect now using this connection string: %s", connectionString)
			clio.NewLine()

			clio.Infof("Or using the %s cli: %s", commandData.GrantOutput.Database.Engine, cliString)
			clio.NewLine()

			eg.Go(func() error {
				defer cancel()

				ln, err := net.Listen("tcp", "localhost:"+port)
				if err != nil {
					clio.Errorf("failed to start listening for connections on port: %s error: %w", port, err)
				}
				defer ln.Close()

				for {
					connChan := make(chan net.Conn)
					errChan := make(chan error)

					go func() {
						conn, err := ln.Accept()
						if err != nil {
							errChan <- err
							return
						}

						clio.Debug("accepted connection")
						connChan <- conn
					}()

					select {
					case <-ctx.Done():
						clio.Debug("context cancelled shutting down port forward")
						return nil // Context cancelled, exit the loop
					case err := <-errChan:
						clio.Errorf("Failed to accept new connection: %w", err)
						return err
					case conn := <-connChan:
						go func() {
							serverConn, err := net.Dial("tcp", "localhost:"+commandData.SSMPortForwardLocalPort)
							if err != nil {
								_ = conn.Close()
								clio.Errorf("Failed to establish a connection to the remote proxy server while accepting local connection: %w", err)
								return
							}

							handshaker := handshake.NewHandshakeClient(serverConn, ensuredGrant.Grant.Id, cfg.TokenSource)
							handshakeResult, err := handshaker.Handshake()
							// if the handshake fails, we bail because we won't be able to make any connections to this database
							if err != nil {
								_ = conn.Close()
								clio.Errorf("Failed to authenticate connection to the remove proxy server while accepting local connection: %w", err)
								return
							}

							clio.Debugw("handshakeResult", "result", handshakeResult)

							clio.Infof("Connection accepted for session [%v]", handshakeResult.ConnectionID)

							// when the handshake is successful for a connection
							// begin proxying traffic
							go func() {
								defer conn.Close()
								_, err := io.Copy(serverConn, conn)
								if err != nil {
									clio.Debugw("Error writing data from client to server usually this is just because the database proxy session ended.", "connectionId", handshakeResult.ConnectionID, zap.Error(err))
								}
								clio.Infof("Connection ended for session [%v]", handshakeResult.ConnectionID)
							}()
							go func() {
								defer conn.Close()
								_, err := io.Copy(conn, serverConn)
								if err != nil {
									clio.Debugw("Error writing data from server to client usually this is just because the database proxy session ended.", "connectionId", handshakeResult.ConnectionID, zap.Error(err))
								}
							}()
						}()
					}
				}
			})
			return nil
		})

		eg.Go(func() error {
			select {
			case <-sigs:
				clio.Info("Received interrupt signal, shutting down database proxy...")
				cancel()
			case <-ctx.Done():
				clio.Info("Shutting down database proxy...")
			}
			if err := cmd.Process.Signal(os.Interrupt); err != nil {
				clio.Errorf("Error sending SIGTERM to AWS SSM process: %w", err)
			}
			return nil
		})

		// Wait for the command to finish
		eg.Go(func() error {
			defer cancel()
			err = cmd.Wait()
			if err != nil {
				return fmt.Errorf("AWS SSM port forward session closed with an error: %s", err)
			}
			return nil
		})

		return eg.Wait()
	},
}

func GrabUnusedPort() (string, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return "", err
	}

	port := listener.Addr().(*net.TCPAddr).Port
	err = listener.Close()
	if err != nil {
		return "", err
	}
	return strconv.Itoa(port), nil
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
	return append(SanitisedEnv(), assume.EnvKeys(creds, commandData.GrantOutput.Database.Region)...)
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
	GrantOutput              AWSRDS
	SSMPortForwardLocalPort  string
	SSMPortForwardServerPort string
}

func formatSSMCommandArgs(data CommandData) []string {
	out := []string{
		"ssm",
		"start-session",
		fmt.Sprintf("--target=%s", data.GrantOutput.SSMSessionTarget),
		"--document-name=AWS-StartPortForwardingSession",
		"--parameters",
		fmt.Sprintf(`{"portNumber":["%s"], "localPortNumber":["%s"]}`, data.SSMPortForwardServerPort, data.SSMPortForwardLocalPort),
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
`, commandData.GrantOutput.Grant.ID, commandData.GrantOutput.Database.Account.ID, commandData.GrantOutput.SSORoleName, commandData.GrantOutput.SSOStartURL, commandData.GrantOutput.SSORegion, commandData.GrantOutput.Database.Region)

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
	if err != nil {
		clio.Error("granted credentials process failed")
		clio.Log(stderr.String())
		return aws.Credentials{}, err
	} else {
		clio.Debugw("granted credential process completed without os error", "stderr", stderr.String())
	}
	return ParseCredentialProcessOutput(stdout.String())
}
