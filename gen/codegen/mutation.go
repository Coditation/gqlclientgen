package codegen

import (
	"gqlclientgen/gen/context"
	"gqlclientgen/gen/utils"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/vektah/gqlparser/v2/ast"
)

func buildMutation(def *ast.Schema, c *context.Context) error {
	mutation := def.Mutation
	for _, mutate := range mutation.Fields {
		if mutate.Position != nil {
			q := &jen.Statement{}
			q.Add(createMutationFunc(mutate))
			c.Model.Queries = append(c.Model.Queries, &context.DataTypeInfo{
				GraphqlName:   mutation.Name,
				MappedName:    strings.ToLower(mutation.Name),
				MappedType:    strings.ToLower(string(mutation.Kind)),
				CodeStatement: q,
			})
		}

	}
	return nil
}

func createMutationFunc(d *ast.FieldDefinition) *jen.Statement {
	var (
		muArgs     = &jen.Statement{}
		varDict    = jen.Dict{}
		tags       []string
		returnType = &jen.Statement{}
	)
	qFunc := jen.Func().Params(utils.GetClientParams()).Id(utils.ToPascalCase(d.Name))
	muArgs.Add(jen.Add(jen.Id("ctx"), jen.Qual("context", "Context")).Op(","))
	for _, arg := range d.Arguments {
		muArgs.Add(jen.Id(utils.ToSmallPascalCase(arg.Name)).Add(utils.GetArgsType(arg)).Op(","))
		varDict[jen.Lit(arg.Name)] = jen.Id(utils.ToSmallPascalCase(arg.Name))
		tags = append(tags, utils.ToSmallPascalCase(arg.Name))
	}
	qFunc.Parens(muArgs)
	qFunc.Parens(jen.List(jen.Add(utils.GetReturnType(d)), jen.Error()))
	returnType.Var().Id("mutate")
	if strings.ToLower(d.Type.Name()) != "any" {
		returnType.Struct(
			jen.Id(utils.ToPascalCase(d.Name)).Struct(
				jen.Id(d.Type.Name()).Add(utils.GetRequestType(d)),
			).Tag(utils.GetRequestTags(utils.ToSmallPascalCase(d.Name), tags)),
		)
	} else {
		returnType.Struct(
			jen.Id(utils.ToPascalCase(d.Name)).Interface().Tag(utils.GetRequestTags(utils.ToSmallPascalCase(d.Name), tags)),
		)
	}
	variables := jen.Id("variables").Op(":=").Map(jen.String()).Interface()
	qFunc.Block(
		returnType,
		variables.Values(varDict).Line(),
		jen.List(jen.Id("resp"), jen.Err()).Op(":=").Id("c").Dot("Client").Dot("QueryRaw").Params(jen.List(jen.Id("ctx"), jen.Id("&mutate"), jen.Id("variables"))),
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
