package codegen

import (
	"bytes"
	"gqlclientgen/config"
	"gqlclientgen/gen/context"
	"gqlclientgen/gen/queryparser"
	"gqlclientgen/gen/utils"
	"os"
	"path"
	"strings"

	"github.com/spf13/viper"

	"github.com/dave/jennifer/jen"
	"github.com/vektah/gqlparser/v2/ast"
)

var (
	queries []*ast.QueryDocument
)

func GenerateClientCode(parsedGql *ast.Schema) error {
	context := context.Create()
	f, err := createFiles()
	if err != nil {
		return err
	}
	defer f.Close()
	queryPath := viper.GetViper().GetString(config.QueryPath)
	if queryPath != "" && strings.TrimSpace(queryPath) != "" {
		queryDocument, err := queryparser.ParseQueryDocuments(queryPath, parsedGql)
		if err != nil {
			return err
		}
		queries, err = queryparser.QueryDocumentsByOperations(parsedGql, queryDocument.Operations)
		if err != nil {
			return err
		}
	}
	buildClientCode(context)
	buildMutation(parsedGql, context)
	buildQuery(parsedGql, queries, context)
	buildSubscription(parsedGql, context)
	jenFile := jen.NewFile(viper.GetViper().GetString(config.PackageNameKey))
	jenFile.ImportAlias(utils.GqlClientPackageName, "graphql")
	for _, v := range context.Client.Client {
		jenFile.Add(v.CodeStatement)
	}
	for _, v := range context.Model.Queries {
		jenFile.Add(v.CodeStatement)
	}
	for _, v := range context.Model.Mutations {
		jenFile.Add(v.CodeStatement)
	}
	for _, v := range context.Model.Subscriptions {
		jenFile.Add(v.CodeStatement)
	}
	buf := &bytes.Buffer{}
	err = jenFile.Render(buf)
	if err != nil {
		return err
	}
	if writeErr := os.WriteFile(f.Name(), buf.Bytes(), os.ModePerm); writeErr != nil {
		return writeErr
	}
	return nil
}

func createFiles() (*os.File, error) {
	p := path.Join(utils.GetFilePath(), "client.go")
	if err := os.MkdirAll(path.Dir(p), os.ModePerm); err != nil {
		return nil, err
	}
	f, err := os.Create(p)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func buildClientCode(c *context.Context) {
	client := jen.Type().Id("Client").Struct(
		jen.Id("client").Add(jen.Op("*").Qual(utils.GqlClientPackageName, "Client")),
	).Line()

	newClient := jen.Func().Id("NewClient").Params(
		jen.Id("url").String(),
		jen.Id("httpClient").Op("*").Qual("net/http", "Client"),
	).Op("*").Id("Client").Block(
		jen.Return(jen.Op("&").Id("Client").Block(
			jen.Id("client").Op(":").Qual("github.com/hasura/go-graphql-client", "NewClient").Params(jen.List(jen.Id("url"), jen.Id("httpClient"))).Op(","),
		)),
	)
	client.Add(newClient).Line()
	c.Client.Client = append(c.Client.Client, &context.DataTypeInfo{
		CodeStatement: client,
	})
}
