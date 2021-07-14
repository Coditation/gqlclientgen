package schema

import "github.com/vektah/gqlparser/ast"

type UrlLoader struct{}

func (u UrlLoader) Load(source string, params map[string]interface{}) ([]*ast.Source, error) {
	if params == nil {

	}
	panic("implement me")
}
