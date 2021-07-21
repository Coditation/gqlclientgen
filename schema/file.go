package schema

import (
	"gqlclientgen/config"
	"io/ioutil"

	"github.com/spf13/viper"
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

type SdlFileLoader struct{}

func (s SdlFileLoader) Load() (*ast.Schema, error) {
	v := viper.GetViper()
	sfp := v.GetString(config.SourceFilePathKey)
	b, err := ioutil.ReadFile(sfp)
	if err != nil {
		return nil, err
	}
	var sources = []*ast.Source{
		{Name: "default", Input: string(b)},
	}
	schema, loadErr := gqlparser.LoadSchema(sources...)
	if loadErr != nil {
		return nil, err
	}
	return schema, nil
}
