package main

import (
	"gqlclientgen/config"
	"gqlclientgen/gen/merge"

	"github.com/Coditation/skael-connectors-shared/logger"
)

func main() {
	if loadErr := config.LoadConfig("config"); loadErr != nil {
		logger.LogError(loadErr)
	}
	if genErr := merge.Generate(); genErr != nil {
		logger.LogError(genErr)
	}
}
