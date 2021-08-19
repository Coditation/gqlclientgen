package pluginmapper

import (
	"errors"
	"gqlclientgen/config"
	"gqlclientgen/gen/utils"
	"os"

	"path/filepath"
	"plugin"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/spf13/viper"
)

type CustomScalarMapper interface {
	Type() string         // type of the custom scalar
	Code() *jen.Statement // jen code ot custom scalar
}

type Plugin interface {
	GetCustomScalarMapper() CustomScalarMapper
}

func LoadPlugins() error {
	pluginPath := viper.GetViper().GetString(config.PluginPath)
	if pluginPath != "" && strings.TrimSpace(pluginPath) != "" {
		pluginsPath, err := filepath.Abs(pluginPath)
		if err != nil {
			return err
		}
		plugins, err := os.ReadDir(pluginsPath)
		if err != nil {
			return err
		}
		for _, pluginType := range plugins {
			if filepath.Ext(pluginType.Name()) == ".so" {
				p, err := plugin.Open(pluginPath + "/" + pluginType.Name())
				if err != nil {
					return err
				}
				pluginName := utils.ToPascalCase(fileNameWithoutExtension(pluginType.Name()))
				v, err := p.Lookup(pluginName)
				if err != nil {
					return err
				}
				customScalarMapper, ok := v.(CustomScalarMapper)
				if !ok {
					return errors.New("Custom Scalar type doesn't implement CustomScalarMapper")
				}
				utils.TypeMappings[customScalarMapper.Type()] = utils.TypeMapping{
					MappedType: customScalarMapper.Code(),
				}
			}
		}
	}
	return nil
}

func fileNameWithoutExtension(fileName string) string {
	fileName = filepath.Base(fileName)
	if pos := strings.LastIndexByte(fileName, '.'); pos != -1 {
		return fileName[:pos]
	}
	return fileName
}
