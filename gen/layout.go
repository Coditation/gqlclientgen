package gen

import (
	"coditation.com/gqlclientgen/config"
	"github.com/spf13/viper"
	"os"
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
