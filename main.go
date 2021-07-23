package main

import (
	config "gqlclientgen/config"
	"gqlclientgen/gen/modelgen"
	"gqlclientgen/schema"

	"github.com/Coditation/skael-connectors-shared/logger"
	"github.com/spf13/viper"
)

func main() {
	// arg := os.Args[1]
	// if arg == "" {
	// 	panic("please specify a configuration YAML file")
	// }
	err := config.LoadConfig(".")
	if err != nil {
		panic("cannot read configuration")
	}
	viper.SetDefault("sourceType", "file")
	viper.SetDefault("sourceFilePath", "/home/sahilp/WorkSpace/Go/src/gqlclientgen/schema/schema_test/schema/schema.graphqls")
	loader := schema.GetLoader()
	source, err := loader.Load()
	if err != nil {
		logger.LogError(err)
	}
	modelgen.GenerateModel(source)
}
