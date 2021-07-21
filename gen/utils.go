package gen

import (
	"gqlclientgen/config"

	"github.com/spf13/viper"
	"github.com/vektah/gqlparser/v2/ast"
)

func GetPackagePath() string {
	v := viper.GetViper()
	packageName := v.GetString(config.PackageNameKey)
	if packageName == "" {
		packageName = defaultPackageName
	}
	return v.GetString(config.OutputDirectoryKey) + "/" + packageName
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func getTags(f *ast.FieldDefinition) string {
	if f.Type.NonNull {
		return f.Name + ",omitempty"
	}
	return f.Name
}
