package main

import (
	"encoding/json"
	"fmt"
	"gqlclientgen/config"
	"gqlclientgen/schema"

	"github.com/Coditation/skael-connectors-shared/logger"
	"github.com/spf13/viper"
)

func main() {
	err := config.LoadConfig(".")
	if err != nil {
		panic("cannot read configuration")
	}
	v := viper.GetViper()
	v.SetDefault("sourceType", "file")
	v.SetDefault("sourceFilePath", "schema/schema.graphqls")
	loader := schema.GetLoader()
	s, err := loader.Load()
	if err != nil {
		logger.LogError(err)
	}
	data, _ := json.MarshalIndent(s, "", "\t")
	fmt.Println(string(data))
}
