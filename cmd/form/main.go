package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/bufbuild/connect-go"
	huh "github.com/charmbracelet/huh"
	"github.com/common-fate/clio"
	"github.com/common-fate/grab"
	"github.com/common-fate/sdk/config"
	accessv1alpha1 "github.com/common-fate/sdk/gen/commonfate/access/v1alpha1"
	"github.com/common-fate/sdk/service/access"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type EnsureForm struct {
	selection []string
}

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ensureForm := EnsureForm{}
	ctx := context.Background()
	// Should we run in accessible mode?
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))
	cfg, err := config.LoadDefault(ctx)
	if err != nil {
		return err
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
		return err
	}
	_ = availabilities

	columns := []huh.Column{
		{Title: "TARGET", Width: 30},
		{Title: "NAME", Width: 45},
		{Title: "ROLE", Width: 30},
		{Title: "DURATION", Width: 10},
		{Title: "SELECTOR", Width: 20},
		{Title: "PRIORITY", Width: 10},
	}

	rows := []huh.Row{}

	// var suggestions []string
	for _, a := range availabilities {
		rows = append(rows, huh.Row{
			Key:    a.Id,
			Values: []string{a.Target.Eid.Display(), a.Target.Name, a.Role.Name, a.Duration.AsDuration().String(), a.TargetSelector.Id, strconv.FormatUint(uint64(a.Priority), 10)},
		})
		// suggestions = append(suggestions, a.Target.Eid.Display(), a.Target.Name, a.Role.Name)
	}

	form := huh.NewForm(

		huh.NewGroup(

			huh.NewMultiSelectTable().
				Title(fmt.Sprintf("Availabilities (total: %v)", len(rows))).
				Description("Select one or more.").
				Columns(columns).
				Options(rows...).
				Validate(func(t []string) error {
					if len(t) <= 0 {
						return fmt.Errorf("at least one topping is required")
					}
					return nil
				}).
				Value(&ensureForm.selection).
				Filterable(true).
				Height(7).WithFilterFunc(func(filter string, option huh.Row) bool {
				for _, v := range option.Values {
					if fuzzy.Match(strings.ToLower(filter), v) {
						return true
					}
				}
				return false
			}),
		),
	).WithAccessible(accessible)

	err = form.Run()

	if err != nil {
		return err
	}

	// done := make(chan bool)
	// action := func() {

	// 	<-done
	// }

	// _ = spinner.New().Title("Planning access changes...").Accessible(accessible).Action(action).Run()

	req := accessv1alpha1.BatchEnsureRequest{}
Outer:
	for _, selection := range ensureForm.selection {
		for _, availability := range availabilities {
			if availability.Id == selection {
				req.Entitlements = append(req.Entitlements, &accessv1alpha1.EntitlementInput{
					Target: &accessv1alpha1.Specifier{
						Specify: &accessv1alpha1.Specifier_Lookup{
							Lookup: availability.Target.Name,
						},
					},
					Role: &accessv1alpha1.Specifier{
						Specify: &accessv1alpha1.Specifier_Lookup{
							Lookup: availability.Role.Name,
						},
					},
				})
				continue Outer
			}
		}
	}

	res, err := client.BatchEnsure(ctx, connect.NewRequest(&req))
	if err != nil {
		// done <- true
		return err
	}
	// done <- true

	var confirm bool
	err = huh.NewConfirm().
		Title("Would you like to request access now?").
		Value(&confirm).
		Affirmative("Yes!").
		Negative("No.").Run()
	if err != nil {
		return err
	}

	if !confirm {
		fmt.Println("no access changes")
		return nil
	}

	// _ = spinner.New().Title("Planning access changes...").Accessible(accessible).Action(action).Run()
	req.DryRun = false
	res, err = client.BatchEnsure(ctx, connect.NewRequest(&req))
	if err != nil {
		// done <- true
		return err
	}
	// done <- true

	clio.Infow("BatchEnsure response", "response", res)

	// // Print order summary.
	// {
	// 	var sb strings.Builder
	// 	keyword := func(s string) string {
	// 		return lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render(s)
	// 	}
	// 	fmt.Fprintf(&sb,
	// 		"%s\n\nOne %s%s, topped with %s with %s on the side.",
	// 		lipgloss.NewStyle().Bold(true).Render("BURGER RECEIPT"),
	// 		keyword(order.Burger.Spice.String()),
	// 		keyword(order.Burger.Type),
	// 		keyword(xstrings.EnglishJoin(order.Burger.Toppings, true)),
	// 		keyword(order.Side),
	// 	)

	// 	name := order.Name
	// 	if name != "" {
	// 		name = ", " + name
	// 	}
	// 	fmt.Fprintf(&sb, "\n\nThanks for your order%s!", name)

	// 	if order.Discount {
	// 		fmt.Fprint(&sb, "\n\nEnjoy 15% off.")
	// 	}

	// 	fmt.Println(
	// 		lipgloss.NewStyle().
	// 			Width(40).
	// 			BorderStyle(lipgloss.RoundedBorder()).
	// 			BorderForeground(lipgloss.Color("63")).
	// 			Padding(1, 2).
	// 			Render(sb.String()),
	// 	)
	// }
	return nil
}
