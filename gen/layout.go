package gen

import (
	"gqlclientgen/config"
	"os"

	"github.com/spf13/viper"
)

const defaultPackageName = "gql"

func GenerateLayout() error {
	v := viper.GetViper()
	packageName := v.GetString(config.PackageNameKey)
	if packageName == "" {
		packageName = defaultPackageName
	}
	err := os.MkdirAll(v.GetString(config.OutputDirectoryKey)+"/"+packageName, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
