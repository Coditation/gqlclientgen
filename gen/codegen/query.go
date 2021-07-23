package codegen

import (
	"gqlclientgen/gen/context"
	"gqlclientgen/gen/utils"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/vektah/gqlparser/v2/ast"
)

func BuildQuery(def *ast.Schema, c *context.Context) error {
	queries := def.Query
	for _, query := range queries.Fields {
		if query.Position != nil {
			q := &jen.Statement{}
			q.Add(createQueryFunc(query))
			c.Model.Queries = append(c.Model.Queries, &context.DataTypeInfo{
				GraphqlName:   queries.Name,
				MappedName:    strings.ToLower(queries.Name),
				MappedType:    strings.ToLower(string(queries.Kind)),
				CodeStatement: q,
			})
		}

	}
	return nil
}

func createQueryFunc(d *ast.FieldDefinition) *jen.Statement {
	var (
		queryArgs  = &jen.Statement{}
		varDict    = jen.Dict{}
		tags       []string
		returnType = &jen.Statement{}
	)
	qFunc := jen.Func().Params(utils.GetClientParams()).Id(utils.ToCamelCase(d.Name))
	queryArgs.Add(jen.Add(jen.Id("ctx"), jen.Qual("context", "Context")).Op(","))
	for _, arg := range d.Arguments {
		queryArgs.Add(jen.Id(utils.ToSmallCamelCase(arg.Name)).Add(utils.GetArgsType(arg)).Op(","))
		varDict[jen.Lit(arg.Name)] = jen.Id(utils.ToSmallCamelCase(arg.Name))
		tags = append(tags, utils.ToSmallCamelCase(arg.Name))
	}
	qFunc.Parens(queryArgs)
	qFunc.Parens(jen.List(jen.Add(utils.GetReturnType(d)), jen.Error()))
	if strings.ToLower(d.Type.Name()) != "any" {
		returnType.Id(d.Type.Name()).Add(utils.GetRequestType(d))
	} else {
		returnType.Id(d.Type.Name()).Interface()
	}
	variables := jen.Id("variables").Op(":=").Map(jen.String()).Interface()
	qFunc.Block(
		jen.Var().Id("query").Struct(
			jen.Id(utils.ToCamelCase(d.Name)).Struct(
				returnType,
			).Tag(utils.GetRequestTags(utils.ToSmallCamelCase(d.Name), tags)),
		),
		variables.Values(varDict).Line(),
		jen.If(
			jen.Err().Op(":=").Id("c").Dot("Client").Dot("Query").Params(jen.List(jen.Id("ctx"), jen.Id("&query"), jen.Id("variables"))),
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Return(jen.Nil(), jen.Err()),
		),
		jen.Return(jen.Op("&").Id("query").Dot(utils.ToCamelCase(d.Name)), jen.Nil()),
	).Line()

	return qFunc.Line()
}
