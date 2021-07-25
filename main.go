package main

import (
	"gqlclientgen/gen/merge"

	"github.com/Coditation/skael-connectors-shared/logger"
)

func main() {
	if genErr := merge.Generate(); genErr != nil {
		logger.LogError(genErr)
	}
}
