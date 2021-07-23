package codegen

import (
	"gqlclientgen/gen/context"
	"gqlclientgen/gen/utils"
	"os"
	"path"

	"github.com/dave/jennifer/jen"
)

func Render(m context.ModelInfo) error {

	return nil
}

func createFiles() (*os.File, error) {
	p := utils.GetPackagePath() + "/client.go"
	if err := os.MkdirAll(path.Dir(p), os.ModePerm); err != nil {
		return nil, err
	}
	f, err := os.Create(p)
	if err != nil {
		return nil, err
	} else {
		return f, nil
	}
}

func CreateClientStruct() {}

func BuildClientCode(c *context.Context) {
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
