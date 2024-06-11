package awsconfig

import (
	awsv1alpha1 "github.com/common-fate/sdk/gen/granted/registry/aws/v1alpha1"

	"gopkg.in/ini.v1"
)

type MergeOpts struct {
	Config            *ini.File
	ProfileName       string
	ProfileAttributes []*awsv1alpha1.ProfileAttributes
}

func Merge(opts MergeOpts) error {

	opts.Config.DeleteSection(opts.ProfileName)

	newSection, err := opts.Config.NewSection(opts.ProfileName)
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
