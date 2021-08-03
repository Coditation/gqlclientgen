package codegen

import (
	"gqlclientgen/gen/context"
	"gqlclientgen/gen/utils"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/vektah/gqlparser/v2/ast"
)

func buildQuery(def *ast.Schema, c *context.Context) error {
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
	qFunc := jen.Func().Params(utils.GetClientParams()).Id(utils.ToPascalCase(d.Name))
	queryArgs.Add(jen.Add(jen.Id("ctx"), jen.Qual("context", "Context")).Op(","))
	for _, arg := range d.Arguments {
		queryArgs.Add(jen.Id(utils.ToCamelCase(arg.Name)).Add(utils.GetArgsType(arg)).Op(","))
		varDict[jen.Lit(arg.Name)] = jen.Id(utils.ToCamelCase(arg.Name))
		tags = append(tags, utils.ToCamelCase(arg.Name))
	}
	qFunc.Parens(queryArgs)
	qFunc.Parens(jen.List(jen.Add(utils.GetReturnType(d)), jen.Error()))
	returnType.Var().Id("query")
	if strings.ToLower(d.Type.Name()) != "any" {
		returnType.Struct(
			jen.Id(utils.ToPascalCase(d.Name)).Struct(
				jen.Id(d.Type.Name()).Add(utils.GetRequestType(d)),
			).Tag(utils.GetRequestTags(utils.ToCamelCase(d.Name), tags)),
		)
	} else {
		returnType.Struct(
			jen.Id(utils.ToPascalCase(d.Name)).Interface().Tag(utils.GetRequestTags(utils.ToCamelCase(d.Name), tags)),
		)
	}
	variables := jen.Id("variables").Op(":=").Map(jen.String()).Interface()
	qFunc.Block(
		returnType,
		variables.Values(varDict).Line(),
		jen.List(jen.Id("resp"), jen.Err()).Op(":=").Id("c").Dot("Client").Dot("QueryRaw").Params(jen.List(jen.Id("ctx"), jen.Id("&query"), jen.Id("variables"))),
		jen.If(
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Return(jen.Nil(), jen.Err()),
		),
		jen.Var().Id("res").Add(utils.GetVarType(d)),
		jen.If(
			jen.Id("resp").Op("!=").Nil(),
		).Block(
			jen.List(jen.Id("byteData"), jen.Err()).Op(":=").Id("resp").Dot("MarshalJSON").Call(),
			jen.If(
				jen.Err().Op("!=").Nil(),
			).Block(
				jen.Return(jen.Nil(), jen.Err()),
			),
			jen.If(
				jen.Id("unMarshalErr").Op(":=").Qual("encoding/json", "Unmarshal").Params(jen.Id("byteData"), jen.Op("&").Id("res")),
				jen.Id("unMarshalErr").Op("!=").Nil(),
			).Block(
				jen.Return(jen.Nil(), jen.Id("unMarshalErr")),
			),
		),
		jen.Return(jen.Add(utils.GetReturnFuncType(d)), jen.Nil()),
	).Line()

	return qFunc.Line()
}
