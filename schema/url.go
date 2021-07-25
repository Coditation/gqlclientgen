package schema

import (
	"context"
	"fmt"
	"gqlclientgen/config"
	"net/http"

	"github.com/Yamashou/gqlgenc/introspection"
	"github.com/hasura/go-graphql-client"
	"github.com/spf13/viper"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/validator"
)

type UrlLoader struct{}

func (u UrlLoader) Load() (*ast.Schema, error) {
	var res introspection.Query
	ctx := context.Background()
	genClient := graphql.NewClient(viper.GetString(config.RemoteURL), http.DefaultClient)
	if err := genClient.Query(ctx, res, nil); err != nil {
		return nil, err
	}
	schema, err := validator.ValidateSchemaDocument(introspection.ParseIntrospectionQuery(viper.GetString(config.RemoteURL), res))
	if err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}
	return schema, nil
}
