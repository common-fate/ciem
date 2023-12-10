package access

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"
	"github.com/bufbuild/connect-go"
	"github.com/common-fate/ciem/printdiags"
	"github.com/common-fate/clio"
	accessv1alpha1 "github.com/common-fate/sdk/gen/commonfate/access/v1alpha1"
	"github.com/common-fate/sdk/gen/commonfate/access/v1alpha1/accessv1alpha1connect"
	"github.com/common-fate/sdk/uid"
	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
	"google.golang.org/protobuf/encoding/protojson"
)

func dryRun(ctx context.Context, apiURL *url.URL, client accessv1alpha1connect.AccessServiceClient, req *accessv1alpha1.BatchEnsureRequest, jsonOutput bool) (bool, error) {
	req.DryRun = true

	si := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	si.Suffix = " planning access changes..."
	si.Writer = os.Stderr
	si.Start()

	res, err := client.BatchEnsure(ctx, connect.NewRequest(req))
	if err != nil {
		si.Stop()
		return false, err
	}

	si.Stop()

	clio.Debugw("BatchEnsure response", "response", res)

	if jsonOutput {
		resJSON, err := protojson.Marshal(res.Msg)
		if err != nil {
			return false, err
		}
		fmt.Println(string(resJSON))

		return false, errors.New("exiting because --output=json was specified: use --output=text to show an interactive prompt, or use --confirm to proceed with the changes")
	}

	names := map[uid.UID]string{}

	var hasChanges bool

	for _, g := range res.Msg.Grants {
		names[uid.New("Access::Grant", g.Grant.Id)] = g.Grant.Name

		exp := "<invalid expiry>"

		if g.Grant.ExpiresAt != nil {
			exp = shortDur(time.Until(g.Grant.ExpiresAt.AsTime()))
		}

		if g.Change > 0 {
			hasChanges = true
		}

		switch g.Change {
		case accessv1alpha1.GrantChange_GRANT_CHANGE_ACTIVATED:
			color.New(color.BgHiGreen).Printf("[WILL ACTIVATE]")
			color.New(color.FgGreen).Printf(" %s will be activated for %s: %s\n", g.Grant.Name, exp, requestURL(apiURL, g.Grant))
			continue

		case accessv1alpha1.GrantChange_GRANT_CHANGE_EXTENDED:
			color.New(color.BgBlue).Printf("[WILL EXTEND]")
			color.New(color.FgBlue).Printf(" %s will be extended for another %s: %s\n", g.Grant.Name, exp, requestURL(apiURL, g.Grant))
			continue

		case accessv1alpha1.GrantChange_GRANT_CHANGE_REQUESTED:
			color.New(color.BgHiYellow, color.FgBlack).Printf("[WILL REQUEST]")
			color.New(color.FgYellow).Printf(" %s will require approval\n", g.Grant.Name)
			continue

		case accessv1alpha1.GrantChange_GRANT_CHANGE_PROVISIONING_FAILED:
			// shouldn't happen in the dry-run request but handle anyway
			color.New(color.FgRed).Printf("[ERROR] %s will fail provisioning\n", g.Grant.Name)
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

	if !hasChanges {
		return false, nil
	}

	if !isTerminal(os.Stdin.Fd()) {
		return false, errors.New("detected a noninteractive terminal: to apply the planned changes please re-run 'cf access ensure' with the --confirm flag")
	}

	confirm := survey.Confirm{
		Message: "Apply proposed access changes",
	}
	var proceed bool
	err = survey.AskOne(&confirm, &proceed)
	if err != nil {
		return false, err
	}
	return true, nil
}

func isTerminal(fd uintptr) bool {
	return isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd)
}

func requestURL(apiURL *url.URL, grant *accessv1alpha1.Grant) string {
	p := apiURL.JoinPath("access", "requests", grant.AccessRequestId)
	return p.String()
}
