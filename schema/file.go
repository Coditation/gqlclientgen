package schema

import (
	"coditation.com/gqlclientgen/config"
	"fmt"
	"github.com/spf13/viper"
	"github.com/vektah/gqlparser/v2/ast"
	"io/ioutil"
)

type SdlFileLoader struct {
}

func (s SdlFileLoader) Load() ([]*ast.Source, error) {
	v := viper.GetViper()
	sfp := v.GetString(config.SourceFilePathKey)
	b, err := ioutil.ReadFile(sfp)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	str := string(b) // convert content to a 'string'
	var sources = []*ast.Source{
		{Name: "default", Input: str},
	}
	return sources, nil
}
