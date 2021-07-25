package context

import (
	"gqlclientgen/config"
	"gqlclientgen/gen/utils"
	"os"
	"path"

	"github.com/spf13/viper"
)

func GenerateLayout() error {
	v := viper.GetViper()
	packageName := v.GetString(config.PackageNameKey)
	if packageName == "" {
		packageName = utils.DefaultPackageName
	}
	err := os.MkdirAll(path.Join(v.GetString(config.OutputDirectoryKey), packageName), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
