package access

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/bufbuild/connect-go"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/lipgloss"
	"github.com/common-fate/ciem/multiselecttable"
	"github.com/common-fate/ciem/printdiags"
	"github.com/common-fate/clio"
	"github.com/common-fate/grab"
	"github.com/common-fate/sdk/config"
	"github.com/common-fate/sdk/eid"
	accessv1alpha1 "github.com/common-fate/sdk/gen/commonfate/access/v1alpha1"
	"github.com/common-fate/sdk/service/access"
	"github.com/fatih/color"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/urfave/cli/v2"
	"google.golang.org/protobuf/encoding/protojson"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("#ED77C0"))

type model struct {
	filter textinput.Model
	table  multiselecttable.Model

	rows         []multiselecttable.Row
	filteredRows []multiselecttable.Row
}

type resetSelectedMsg struct {
}

func (m model) Init() tea.Cmd {
	return tea.Batch(func() tea.Msg {
		return resetSelectedMsg{}
	}, textinput.Blink)

}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var tableCmd tea.Cmd
	var filterCmd tea.Cmd
	switch msg := msg.(type) {
	case resetSelectedMsg:
		m.table.Selected = make(map[int]bool)
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			m.table.Selected[m.table.Cursor()] = !m.table.Selected[m.table.Cursor()]
		case "ctrl+e":
			return m, tea.Quit
		case "down", "up":

			// prevents the filter input responding to up down
			m.table, tableCmd = m.table.Update(msg)
			return m, tableCmd
		}
	}
	trimmed := strings.Trim(m.filter.Value(), " ")
	m.filteredRows = []multiselecttable.Row{}
	if trimmed == "" {
		m.table.SetRows(m.rows)
	} else {
		filters := slices.Compact(strings.Split(trimmed, " "))

		// fuzzy match on any colum in the row
	Row:
		for i, r := range m.rows {
			for _, v := range r {
				for _, f := range filters {
					if fuzzy.Match(f, v) {
						m.filteredRows = append(m.filteredRows, m.rows[i])
						continue Row
					}
				}

			}
		}
		m.table.SetRows(m.filteredRows)
	}

	m.table, tableCmd = m.table.Update(msg)
	m.filter, filterCmd = m.filter.Update(msg)
	return m, tea.Batch(tableCmd, filterCmd)
}

func (m model) View() string {
	return baseStyle.Render(m.filter.View(), "\n\n", m.table.View(), "\n"+"(ctrl+e to confirm selection. ctrl+c to quit)") + "\n"
}
func selector(ctx context.Context) ([]*accessv1alpha1.Availability, error) {

	cfg, err := config.LoadDefault(ctx)
	if err != nil {
		return nil, err
	}
	client := access.NewFromConfig(cfg)

	availabilities, err := grab.AllPages(ctx, func(ctx context.Context, nextToken *string) ([]*accessv1alpha1.Availability, *string, error) {
		res, err := client.QueryAvailabilities(ctx, connect.NewRequest(&accessv1alpha1.QueryAvailabilitiesRequest{
			PageToken: grab.Value(nextToken),
		}))
		if err != nil {
			return nil, nil, err
		}
		return res.Msg.Availabilities, &res.Msg.NextPageToken, nil
	})
	if err != nil {
		return nil, err
	}
	_ = availabilities

	columns := []multiselecttable.Column{
		{Title: "TARGET", Width: 30},
		{Title: "NAME", Width: 45},
		{Title: "ROLE", Width: 30},
		{Title: "DURATION", Width: 10},
		{Title: "SELECTOR", Width: 20},
		{Title: "PRIORITY", Width: 10},
	}
	rows := []multiselecttable.Row{}

	var suggestions []string
	for _, a := range availabilities {
		rows = append(rows, []string{a.Target.Eid.Display(), a.Target.Name, a.Role.Name, a.Duration.AsDuration().String(), a.TargetSelector.Id, strconv.FormatUint(uint64(a.Priority), 10)})
		suggestions = append(suggestions, a.Target.Eid.Display(), a.Target.Name, a.Role.Name)
	}
	t := multiselecttable.New(
		multiselecttable.WithColumns(columns),
		multiselecttable.WithRows(rows),
		multiselecttable.WithFocused(true),
		multiselecttable.WithHeight(7),
	)

	s := multiselecttable.DefaultStyles()

	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#ED77C0")).
		BorderBottom(true).
		Bold(false).Foreground(lipgloss.Color("#619EFF"))
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color("#30D15D")).
		Bold(false)
	s.Highlighted = s.Highlighted.
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color("#30D15D")).
		Bold(true)

	t.SetStyles(s)
	ti := textinput.New()
	ti.Placeholder = "Filter by target or role, use spaces to search by multiple terms"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 100
	ti.SetSuggestions(suggestions)

	ti.ShowSuggestions = true
	m := model{table: t, filter: ti, rows: rows}
	final, err := tea.NewProgram(m).Run()
	if err != nil {
		return nil, err
	}

	m = final.(model)
	var result []*accessv1alpha1.Availability
	for i := range m.table.Selected {
		result = append(result, availabilities[i])
	}
	return result, nil
}

var ensureCommand = cli.Command{
	Name:  "ensure",
	Usage: "Ensure access to some entitlements (will request, active, or extend access as necessary)",
	Flags: []cli.Flag{
		&cli.StringSliceFlag{Name: "target", Required: false},
		&cli.StringSliceFlag{Name: "role", Required: false},
		&cli.StringFlag{Name: "output", Value: "text", Usage: "output format ('text' or 'json')"},
		&cli.BoolFlag{Name: "confirm", Aliases: []string{"y"}, Usage: "skip the confirmation prompt"},
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

		targets := c.StringSlice("target")
		roles := c.StringSlice("role")

		if len(targets) != len(roles) {
			return errors.New("you need to provide --role flag for each --target flag. For example:\n'cf access ensure --target AWS::Account::123456789012 --role AdministratorAccess --target OtherAccount --role Developer")
		}

		req := accessv1alpha1.BatchEnsureRequest{}

		if len(targets) == 0 {
			selections, err := selector(ctx)
			if err != nil {
				return err
			}

			if len(selections) == 0 {
				clio.Info("No availabilities selected")
				return nil
			}
			for _, selection := range selections {
				req.Entitlements = append(req.Entitlements, &accessv1alpha1.EntitlementInput{
					Target: &accessv1alpha1.Specifier{
						Specify: &accessv1alpha1.Specifier_Lookup{
							Lookup: selection.Target.Name,
						},
					},
					Role: &accessv1alpha1.Specifier{
						Specify: &accessv1alpha1.Specifier_Lookup{
							Lookup: selection.Role.Name,
						},
					},
				})
			}
		} else {
			for i, target := range targets {
				req.Entitlements = append(req.Entitlements, &accessv1alpha1.EntitlementInput{
					Target: &accessv1alpha1.Specifier{
						Specify: &accessv1alpha1.Specifier_Lookup{
							Lookup: target,
						},
					},
					Role: &accessv1alpha1.Specifier{
						Specify: &accessv1alpha1.Specifier_Lookup{
							Lookup: roles[i],
						},
					},
				})
			}
		}

		apiURL, err := url.Parse(cfg.APIURL)
		if err != nil {
			return err
		}

		client := access.NewFromConfig(cfg)

		if !c.Bool("confirm") {
			jsonOutput := c.String("output") == "json"

			// run the dry-run first
			hasChanges, err := dryRun(ctx, apiURL, client, &req, jsonOutput)
			if err != nil {
				return err
			}
			if !hasChanges {
				fmt.Println("no access changes")
				return nil
			}
		}

		// if we get here, dry-run has passed the user has confirmed they want to proceed.
		req.DryRun = false

		si := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		si.Suffix = " ensuring access..."
		si.Writer = os.Stderr
		si.Start()

		res, err := client.BatchEnsure(ctx, connect.NewRequest(&req))
		if err != nil {
			si.Stop()
			return err
		}

		si.Stop()

		clio.Debugw("BatchEnsure response", "response", res)

		if outputFormat == "text" {

			// tree := treeprint.New()

			names := map[eid.EID]string{}

			for _, g := range res.Msg.Grants {
				names[eid.New("Access::Grant", g.Grant.Id)] = g.Grant.Name

				exp := "<invalid expiry>"

				if g.Grant.ExpiresAt != nil {
					exp = shortDur(time.Until(g.Grant.ExpiresAt.AsTime()))
				}

				switch g.Change {
				case accessv1alpha1.GrantChange_GRANT_CHANGE_ACTIVATED:
					color.New(color.BgHiGreen).Printf("[ACTIVATED]")
					color.New(color.FgGreen).Printf(" %s was activated for %s: %s\n", g.Grant.Name, exp, requestURL(apiURL, g.Grant))
					continue

				case accessv1alpha1.GrantChange_GRANT_CHANGE_EXTENDED:
					color.New(color.BgBlue).Printf("[EXTENDED]")
					color.New(color.FgBlue).Printf(" %s was extended for another %s: %s\n", g.Grant.Name, exp, requestURL(apiURL, g.Grant))
					continue

				case accessv1alpha1.GrantChange_GRANT_CHANGE_REQUESTED:
					color.New(color.BgHiYellow, color.FgBlack).Printf("[REQUESTED]")
					color.New(color.FgYellow).Printf(" %s requires approval: %s\n", g.Grant.Name, requestURL(apiURL, g.Grant))
					continue

				case accessv1alpha1.GrantChange_GRANT_CHANGE_PROVISIONING_FAILED:
					// shouldn't happen in the dry-run request but handle anyway
					color.New(color.FgRed).Printf("[ERROR] %s failed provisioning: %s\n", g.Grant.Name, requestURL(apiURL, g.Grant))
					continue
				}

				switch g.Grant.Status {
				case accessv1alpha1.GrantStatus_GRANT_STATUS_ACTIVE:
					color.New(color.FgGreen).Printf("[ACTIVE] %s is already active for the next %s: %s\n", g.Grant.Name, exp, requestURL(apiURL, g.Grant))
					continue
				case accessv1alpha1.GrantStatus_GRANT_STATUS_PENDING:
					color.New(color.FgWhite).Printf("[PENDING] %s is already pending: %s\n", g.Grant.Name, requestURL(apiURL, g.Grant))
					continue
				case accessv1alpha1.GrantStatus_GRANT_STATUS_CLOSED:
					color.New(color.FgWhite).Printf("[CLOSED] %s is closed but was still returned: %s\n. This is most likely due to an error in Common Fate and should be reported to our team: support@commonfate.io.", g.Grant.Name, requestURL(apiURL, g.Grant))
					continue
				}

				color.New(color.FgWhite).Printf("[UNSPECIFIED] %s is in an unspecified status: %s\n. This is most likely due to an error in Common Fate and should be reported to our team: support@commonfate.io.", g.Grant.Name, requestURL(apiURL, g.Grant))
			}

			printdiags.Print(res.Msg.Diagnostics, names)

		}

		if outputFormat == "json" {
			resJSON, err := protojson.Marshal(res.Msg)
			if err != nil {
				return err
			}
			fmt.Println(string(resJSON))
		}

		return nil
	},
}

func shortDur(d time.Duration) string {
	if d > time.Minute {
		d = d.Round(time.Minute)
	} else {
		d = d.Round(time.Second)
	}

	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = s[:len(s)-2]
	}
	if strings.HasSuffix(s, "h0m") {
		s = s[:len(s)-2]
	}
	return s
}
