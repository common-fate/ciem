package access

import (
	"strings"

	awsv1alpha1 "github.com/common-fate/sdk/gen/granted/registry/aws/v1alpha1"

	"gopkg.in/ini.v1"
)

type MergeOpts struct {
	Config              *ini.File
	Prefix              string
	ProfileName         string
	ProfileAttributes   []*awsv1alpha1.ProfileAttributes
	SectionNameTemplate string
	NoCredentialProcess bool
}

func AddProfileToConfig(opts MergeOpts) error {
	if opts.SectionNameTemplate == "" {
		opts.SectionNameTemplate = "{{ .AccountName }}/{{ .RoleName }}"
	}

	profileName := normalizeAccountName(opts.ProfileName)

	sectionName := "profile " + profileName

	opts.Config.DeleteSection(sectionName)

	newSection, err := opts.Config.NewSection(sectionName)
	if err != nil {
		return err
	}
	//add all the attributes returned from CF
	for _, item := range opts.ProfileAttributes {
		newSection.NewKey(item.Key, item.Value)
	}

	return nil
}

func normalizeAccountName(accountName string) string {
	return strings.ReplaceAll(accountName, " ", "-")
}
