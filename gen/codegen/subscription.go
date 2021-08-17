package codegen

import (
	"gqlclientgen/gen/context"
	"gqlclientgen/gen/utils"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/vektah/gqlparser/v2/ast"
)

func buildSubscription(def *ast.Schema, c *context.Context) error {
	subscription := def.Subscription
	if subscription != nil {
		for _, subscribe := range subscription.Fields {
			if subscription.Position != nil {
				q := &jen.Statement{}
				q.Add(createSubscriptionFunc(subscribe))
				c.Model.Subscriptions = append(c.Model.Subscriptions, &context.DataTypeInfo{
					GraphqlName:   subscription.Name,
					MappedName:    strings.ToLower(subscription.Name),
					MappedType:    strings.ToLower(string(subscription.Kind)),
					CodeStatement: q,
				})
			}

		}
	}
	return nil
}

func createSubscriptionFunc(d *ast.FieldDefinition) *jen.Statement {
	var (
		muArgs  = &jen.Statement{}
		varDict = jen.Dict{}
		tags    []string
	)
	qFunc := jen.Func().Params(utils.GetClientParams()).Id(utils.ToPascalCase(d.Name))
	muArgs.Add(jen.Add(jen.Id("serverUrl").String().Op(",")))
	for _, arg := range d.Arguments {
		muArgs.Add(jen.Id(utils.ToCamelCase(arg.Name)).Add(utils.GetArgsType(arg)).Op(","))
		varDict[jen.Lit(arg.Name)] = jen.Id(utils.ToCamelCase(arg.Name))
		tags = append(tags, utils.ToCamelCase(arg.Name))
	}
	qFunc.Parens(muArgs)
	qFunc.Parens(jen.List(jen.Add(utils.GetReturnType(d)), jen.Error()))
	returnType := jen.Var().Id("subscribe")
	if strings.ToLower(d.Type.Name()) != "any" {
		returnType.Struct(
			jen.Id(utils.ToPascalCase(d.Name)).Struct(
				jen.Id(d.Type.Name()).Add(utils.GetRequestType(d)).Tag(map[string]string{"json": utils.GetTags(d), "graphql": utils.GetTags(d)}),
			).Tag(utils.GetRequestTags(utils.ToCamelCase(d.Name), tags)),
		)
	} else {
		returnType.Struct(
			jen.Id(utils.ToPascalCase(d.Name)).Interface().Tag(utils.GetRequestTags(utils.ToCamelCase(d.Name), tags)),
		)
	}
	variables := jen.Id("variables").Op(":=").Map(jen.String()).Interface()
	subscribeClient := jen.Id("subscribeClient").Op(":=").Qual("github.com/hasura/go-graphql-client", "NewSubscriptionClient").Params(jen.Id("serverUrl")).Line()
	subscribeClient.Id("subscribeClient").Dot("Run").Call()
	qFunc.Block(
		returnType,
		subscribeClient,
		variables.Values(varDict).Line(),
		jen.Var().Id("res").Add(utils.GetVarType(d)),
		jen.List(jen.Id("_"), jen.Err()).Op(":=").Id("subscribeClient").Dot("Subscribe").Params(jen.List(jen.Id("&subscribe"), jen.Id("variables"), jen.Func().Params(jen.Id("dataValue").Op("*").Qual("encoding/json", "RawMessage"), jen.Id("errorValue").Error()).Error().Block(
			jen.If(
				jen.Id("errorValue").Op("!=").Nil(),
			).Block(
				jen.Return(jen.Id("errorValue")),
			),
			jen.List(jen.Id("byteData"), jen.Err()).Op(":=").Id("dataValue").Dot("MarshalJSON").Call(),
			jen.If(
				jen.Err().Op("!=").Nil(),
			).Block(
				jen.Return(jen.Err()),
			),
			jen.Err().Op("=").Qual("encoding/json", "Unmarshal").Params(jen.Id("byteData"), jen.Op("&").Id("res")),
			jen.If(
				jen.Err().Op("!=").Nil(),
			).Block(
				jen.Return(jen.Err()),
			),
			jen.Return(jen.Nil()),
		))),
		jen.If(
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Return(jen.Nil(), jen.Err()),
		),
		jen.Return(jen.Add(utils.GetReturnFuncType(d)), jen.Nil()),
	).Line()

	return qFunc.Line()
}
