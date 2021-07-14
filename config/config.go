package config

import (
	"fmt"
	"github.com/spf13/viper"
)

const (
	ConfigName              = "config"
	ConfigType              = "yaml"
	PackageNameKey          = "packageName"
	OutputDirectoryKey      = "outputDirectory"
	SourceTypeKey           = "sourceType"
	GraphqlServerBaseUrlKey = "graphqlServerBaseUrl"
	SourceFilePathKey       = "sourceFilePath"
	FileSourceType          = "file"
	UrlSourceType           = "url"
)

type GqlClientGenConfig struct {
	OutputDirectory      string
	PackageName          string
	SourceType           string
	GraphQLServerBaseUrl string
	SourceFilePath       string
}

func LoadConfig(configFile string) error {
	v := viper.GetViper()
	v.SetConfigName(ConfigName)
	v.SetConfigType(ConfigType)
	v.AddConfigPath(configFile)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	return nil
	//v.AutomaticEnv()
	//c.Server.Url = v.GetString("CONNECTOR_URL")
}
