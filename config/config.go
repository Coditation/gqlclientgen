package config

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/Coditation/skael-connectors-shared/logger"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var GqlConfig GqlClientGenConfig

const (
	ConfigName              = "config"
	ConfigType              = "yaml"
	PackageNameKey          = "packageName"
	OutputDirectoryKey      = "outputDirectory"
	SourceTypeKey           = "sourceType"
	GraphqlServerBaseUrlKey = "graphqlServerBaseUrl"
	SourceFilePathKey       = "sourceFilePath"
	FileSourceType          = "file"
	UrlSourceType           = "remote"
	RemoteURL               = "url"
	QueryPath               = "queryPath"
	PluginPath              = "pluginPath"
)

type GqlClientGenConfig struct {
	OutputDirectory      string
	PackageName          string
	SourceType           string
	GraphQLServerBaseUrl string
	SourceFilePath       string
	QueryPath            string
	PluginPath           string
}

func LoadConfig(configFile string) error {
	v := viper.GetViper()
	v.SetConfigName(ConfigName)
	v.SetConfigType(ConfigType)
	v.AddConfigPath(configFile)
	err := v.ReadInConfig()
	if err != nil {
		logger.LogPanic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	byteData, err := ioutil.ReadFile(path.Join(configFile, ConfigName+"."+ConfigType))
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(byteData, &GqlConfig); err != nil {
		return err
	}
	return nil
}
