package schema

import (
	"gqlclientgen/config"

	"github.com/spf13/viper"
	"github.com/vektah/gqlparser/v2/ast"
)

type Loader interface {
	Load() (*ast.Schema, error)
}

func GetLoader() Loader {
	switch sourceType := viper.GetViper().Get(config.SourceTypeKey); sourceType {
	case config.FileSourceType:
		return new(SdlFileLoader)
	case config.UrlSourceType:
		return new(UrlLoader)
	default:
		panic("incorrect source type configured")
	}
}
