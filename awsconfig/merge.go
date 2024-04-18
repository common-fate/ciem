package awsconfig

import (
	"strings"

	awsv1alpha1 "github.com/common-fate/sdk/gen/granted/registry/aws/v1alpha1"

	"gopkg.in/ini.v1"
)

type MergeOpts struct {
	Config            *ini.File
	ProfileName       string
	ProfileAttributes []*awsv1alpha1.ProfileAttributes
}

func Merge(opts MergeOpts) error {
	profileName := normalizeAccountName(opts.ProfileName)

	sectionName := "profile " + profileName

	opts.Config.DeleteSection(sectionName)

	newSection, err := opts.Config.NewSection(sectionName)
	if err != nil {
		return err
	}

	// add all the attributes returned from CF
	for _, item := range opts.ProfileAttributes {
		_, err := newSection.NewKey(item.Key, item.Value)
		if err != nil {
			return err
		}
	}

	return nil
}

func normalizeAccountName(accountName string) string {
	return strings.ReplaceAll(accountName, " ", "-")
}
