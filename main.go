package main

import (
	"flag"
	"gqlclientgen/config"
	"gqlclientgen/gen/codegen"
	"gqlclientgen/gen/merge"
	"gqlclientgen/gen/pluginmapper"

	"github.com/Coditation/skael-connectors-shared/logger"
)

//execute file as "go run main.go -config_path=CONFIG_DIR_PATH"

func main() {
	var configUrl = flag.String("config_path", "", "gqlclienthgen config path")
	var pluginPath = flag.String("plugin_path", "", "gqlclienthgen plugins path")
	var queryPath = flag.String("query_path", "", "gqlclienthgen plugins path")
	flag.Parse()
	if *configUrl == "" {
		logger.LogFatal("Please enter config path to read config from specified path")
	}
	if *queryPath == "" {
		codegen.QueryPath = ""
	}
	if *queryPath != "" {
		codegen.QueryPath = *queryPath
	}
	if loadErr := config.LoadConfig(*configUrl); loadErr != nil {
		logger.LogError(loadErr)
	}
	if loadErr := pluginmapper.LoadPlugins(*pluginPath); loadErr != nil {
		logger.LogFatal(loadErr)
	}
	if genErr := merge.Generate(); genErr != nil {
		logger.LogError(genErr)
	}
}
