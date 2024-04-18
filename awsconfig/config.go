package awsconfig

import (
	"fmt"
	"os"
	"path"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/common-fate/clio"
	"gopkg.in/ini.v1"
)

const (
	// permission for user to read/write/execute.
	USER_READ_WRITE_PERM = 0700
)

func loadAWSConfigFileFromPath(filepath string) (*ini.File, error) {
	awsConfig, err := ini.LoadSources(ini.LoadOptions{
		SkipUnrecognizableLines: true,
		AllowNonUniqueSections:  true,
		AllowNestedValues:       true,
	}, filepath)
	if err != nil {
		return nil, err
	}

	return awsConfig, nil
}

// GetAWSConfigPath will return default AWS config file path unless $AWS_CONFIG_FILE
// environment variable is set
func GetAWSConfigPath() string {
	file := os.Getenv("AWS_CONFIG_FILE")
	if file != "" {
		clio.Debugf("using aws config filepath: %s", file)
		return file
	}

	return config.DefaultSharedConfigFilename()
}

// loadAWSConfigFile loads the `~/.aws/config` file, and creates it if it doesn't exist.
func Load() (*ini.File, string, error) {
	filepath := GetAWSConfigPath()

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		clio.Infof("created AWS config file: %s", filepath)

		// create all parent directory if necessary.
		err := os.MkdirAll(path.Dir(filepath), USER_READ_WRITE_PERM)
		if err != nil {
			return nil, "", err
		}

		_, err = os.Create(filepath)
		if err != nil {
			return nil, "", fmt.Errorf("unable to create AWS config file: %w", err)
		}
	}

	awsConfig, err := loadAWSConfigFileFromPath(filepath)
	if err != nil {
		return nil, "", err
	}
	return awsConfig, filepath, nil
}
