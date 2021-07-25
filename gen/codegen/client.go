package codegen

import (
	"bytes"
	"gqlclientgen/gen/context"
	"gqlclientgen/gen/utils"
	"os"
	"path"

	"github.com/dave/jennifer/jen"
	"github.com/vektah/gqlparser/v2/ast"
)

func GenerateClientCode(parsedGql *ast.Schema) error {
	context := context.Create()
	f, err := createFiles()
	if err != nil {
		return err
	}
	defer f.Close()
	buildClientCode(context)
	buildMutation(parsedGql, context)
	buildQuery(parsedGql, context)
	jenFile := jen.NewFile("client")
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
	p := path.Join(utils.GetPackagePath(), "client.go")
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
		jen.Id("Client").Add(jen.Op("*").Qual(utils.GqlClientPackageName, "Client")),
	)
	client.Line()

	newClient := jen.Func().Id("NewClient").Params(
		jen.Id("url").String(),
		jen.Id("httpClient").Op("*").Qual("net/http", "Client"),
	).Op("*").Qual("github.com/hasura/go-graphql-client", "Client").Block(
		jen.Return(jen.Qual("github.com/hasura/go-graphql-client", "NewClient").Parens(jen.List(jen.Id("url"), jen.Id("httpClient")))),
	)
	client.Add(newClient)
	c.Client.Client = append(c.Client.Client, &context.DataTypeInfo{
		CodeStatement: client,
	})
}
