package merge

import (
	"gqlclientgen/gen/codegen"
	"gqlclientgen/gen/modelgen"
	"gqlclientgen/schema"

	"github.com/Coditation/skael-connectors-shared/logger"
)

func Generate() error {
	loader := schema.GetLoader()
	s, err := loader.Load()
	if err != nil {
		logger.LogError("failed to load data ", err)
		return err
	}
	if genErr := codegen.GenerateClientCode(s); genErr != nil {
		logger.LogError("failed to generate client code ", genErr)
		return genErr
	}
	if genErr := modelgen.GenerateModel(s); genErr != nil {
		logger.LogError("failed to generate model code ", genErr)
		return genErr
	}
	return nil
}
