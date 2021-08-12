package codegen

import (
	"gqlclientgen/gen/context"
	"gqlclientgen/gen/utils"
	"strings"
	"sync"

	"github.com/dave/jennifer/jen"
	"github.com/vektah/gqlparser/v2/ast"
)

func buildQuery(def *ast.Schema, queryDoc []*ast.QueryDocument, c *context.Context) error {
	queries := def.Query
	var so sync.Once
	for _, query := range queries.Fields {
		q := &jen.Statement{}
		so.Do(func() {
			q.Add(createOperationWithFragments(def, queryDoc))
		})
		if query.Position != nil {
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
		queryArgs = &jen.Statement{}
		varDict   = jen.Dict{}
		tags      []string
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
	returnType := jen.Var().Id("query")
	if strings.ToLower(d.Type.Name()) != "any" {
		returnType.Struct(
			jen.Id(utils.ToPascalCase(d.Name)).Struct(
				jen.Id(d.Type.Name()).Add(utils.GetRequestType(d)).Tag(map[string]string{"json": getStructTags(d), "graphql": getStructTags(d)}),
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
		jen.List(jen.Id("resp"), jen.Err()).Op(":=").Id("c").Dot("client").Dot("QueryRaw").Params(jen.List(jen.Id("ctx"), jen.Id("&query"), jen.Id("variables"))),
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

func getStructTags(d *ast.FieldDefinition) string {
	name := utils.ToCamelCase(d.Type.Name())
	if !d.Type.NonNull {
		return name + ",omitempty"
	}
	if d.Type.Elem != nil {
		return name + ",omitempty"
	}
	return name
}

func createOperationWithFragments(schema *ast.Schema, queryDoc []*ast.QueryDocument) *jen.Statement {
	var (
		op          = &jen.Statement{}
		uniqueFrags = make(map[string]*ast.FragmentDefinition)
		fragsList   ast.FragmentDefinitionList
	)
	if queryDoc != nil {
		for _, d := range queryDoc {
			if d.Fragments != nil && fragsList == nil {
				for _, f := range d.Fragments {
					_, ok := uniqueFrags[f.Name]
					if !ok {
						uniqueFrags[f.Name] = f
					}
				}
				for _, f := range uniqueFrags {
					fragsList = append(fragsList, f)
				}
				op.Add(createFrags(fragsList)).Line()
			}
			if d.Operations != nil && len(d.Operations) > 0 {
				for _, operation := range d.Operations {
					op.Add(createOperations(operation))
				}
			}
		}
	}
	return op
}

func createFrags(frags ast.FragmentDefinitionList) *jen.Statement {
	var (
		frag = &jen.Statement{}
	)
	for _, f := range frags {
		frag.Add(createFragsStruct(f)).Line()
		utils.AllTypes = append(utils.AllTypes, f.Name)
	}
	return frag.Line()
}

func createOperations(op *ast.OperationDefinition) *jen.Statement {
	var (
		queryArgs = &jen.Statement{}
		varDict   = jen.Dict{}
	)
	qFunc := getOpFragStruct(op)
	qFunc.Add(jen.Func().Params(utils.GetClientParams()).Id(utils.ToPascalCase(op.Name)))
	queryArgs.Add(jen.Add(jen.Id("ctx"), jen.Qual("context", "Context")).Op(","))
	for _, arg := range op.VariableDefinitions {
		queryArgs.Add(jen.Id(utils.ToCamelCase(arg.Variable)).Add(getOperationsArgsType(arg)).Op(","))
		varDict[jen.Lit(arg.Variable)] = jen.Id(utils.ToCamelCase(arg.Variable))
	}
	qFunc.Parens(queryArgs)
	qFunc.Parens(jen.List(jen.Op("*").Id(op.Name), jen.Error()))
	variables := jen.Id("variables").Op(":=").Map(jen.String()).Interface()
	returnType := jen.Var().Id("res").Struct(jen.Id(op.Name).Id(op.Name).Tag(getOpRespTags(op)))
	qFunc.Block(
		variables.Values(varDict).Line(),
		returnType,
		jen.List(jen.Id("resp"), jen.Err()).Op(":=").Id("c").Dot("client").Dot(utils.ToPascalCase(string(op.Operation))+"Raw").Params(jen.List(jen.Id("ctx"), jen.Op("&").Id("res"), jen.Id("variables"))),
		jen.If(
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Return(jen.Nil(), jen.Err()),
		),
		jen.If(
			jen.Id("resp").Op("!=").Nil(),
		).Block(
			jen.Return(jen.Op("&").Id("res").Dot(op.Name), jen.Nil()),
		),
		jen.Return(jen.Nil(), jen.Nil()),
	)
	return qFunc.Line()
}

func createFragsStruct(f *ast.FragmentDefinition) *jen.Statement {
	var frag = &jen.Statement{}
	frag.Type().Id(utils.ToPascalCase(f.Name))
	getFragFields(f.SelectionSet, frag)
	return frag
}

func getFragFields(f ast.SelectionSet, c *jen.Statement) *jen.Statement {
	var ftype = &jen.Statement{}
	if f != nil && len(f) > 0 {
		for _, set := range f {
			ftype.Add(createStructFields(set))
		}
	}
	c.Struct(ftype)
	return c
}

func createStructFields(f ast.Selection) *jen.Statement {
	var field = &jen.Statement{}
	switch fieldType := f.(type) {
	case *ast.Field:
		field.Id(utils.ToPascalCase(fieldType.Definition.Name))
		if fieldType.SelectionSet != nil && len(fieldType.SelectionSet) > 0 {
			fields := &jen.Statement{}
			for _, f := range fieldType.SelectionSet {
				fields.Add(createStructFields(f))
			}
			if fieldType.Definition.Type.Elem != nil {
				field.Add(jen.Index().Op("*"))
			}
			field.Struct(fields).Tag(map[string]string{"graphql": fieldType.Definition.Name, "json": fieldType.Definition.Name})
			return field
		}
		field.Add(fragmentFieldType(fieldType.Definition.Type.NamedType, fieldType.Definition.Type)).Tag(map[string]string{"graphql": fieldType.Definition.Name, "json": fieldType.Definition.Name}).Line()
		return field
	case *ast.FragmentSpread:
		field.Id(utils.ToPascalCase(fieldType.Definition.Name))
		if fieldType.Definition.SelectionSet != nil && len(fieldType.Definition.SelectionSet) > 0 {
			fields := &jen.Statement{}
			for _, f := range fieldType.Definition.SelectionSet {
				fields.Add(createStructFields(f))
			}
			field.Struct(fields)
			return field
		}
		return field.Line()
	case *ast.InlineFragment:
		return field.Line()
	}
	return field.Line()
}

func fragmentFieldType(name string, f *ast.Type) *jen.Statement {
	fieldName := utils.ToPascalCase(f.NamedType)
	if fieldName == "" && f.Elem != nil {
		fieldName = utils.ToPascalCase(f.Elem.Name())
		if f.NonNull {
			jen.Index().Op("*").Id(utils.ToPascalCase(fieldName))
		}
		return jen.Index().Id(utils.ToPascalCase(fieldName))
	}
	fieldType, ok := utils.TypeMappings[strings.ToLower(fieldName)]
	if !ok {
		if f.NonNull {
			return jen.Op("*").Id(utils.ToPascalCase(fieldName))
		}
		return jen.Id(utils.ToPascalCase(fieldName))
	}
	if !f.NonNull {
		return jen.Op("*").Add(fieldType.MappedType)
	}
	return fieldType.MappedType
}

func getOperationsArgsType(v *ast.VariableDefinition) *jen.Statement {
	fieldName := utils.ToPascalCase(v.Type.NamedType)
	if fieldName == "" {
		fieldName = utils.ToPascalCase(v.Type.Elem.Name())
		return jen.Index().Op("*").Id(utils.ToPascalCase(fieldName))
	}
	fieldType, ok := utils.TypeMappings[strings.ToLower(fieldName)]
	if !ok {
		if v.Type.NonNull {
			return jen.Op("*").Id(utils.ToPascalCase(v.Type.Name()))
		}
		return jen.Id(utils.ToPascalCase(v.Type.Name()))
	}
	if !v.Type.NonNull {
		return jen.Op("*").Add(fieldType.MappedType)
	}
	return fieldType.MappedType
}

func getOpFragStruct(op *ast.OperationDefinition) *jen.Statement {
	opStruct := jen.Type().Id(op.Name)
	fields := &jen.Statement{}
	for _, field := range op.SelectionSet {
		fieldType := field.(*ast.Field)
		f := fieldType.SelectionSet[0].(*ast.FragmentSpread).Name
		fields.Id(utils.ToPascalCase(fieldType.Alias)).Op("*").Id(utils.ToPascalCase(f)).Tag(map[string]string{"json": utils.ToCamelCase(fieldType.Alias), "graphql": utils.ToCamelCase(fieldType.Alias)}).Line()
	}
	opStruct.Struct(fields)
	return opStruct.Line()
}

func getOpRespTags(op *ast.OperationDefinition) map[string]string {
	m := make(map[string]string)
	tags := []string{}
	operation := op.Name
	if op.VariableDefinitions != nil {
		for _, v := range op.VariableDefinitions {
			if v.DefaultValue != nil {
				tags = append(tags, v.Variable+`:\"`+v.DefaultValue.Raw+`\"`)
			} else {
				tags = append(tags, v.Variable+`:&`+v.Variable)
			}
		}
		if len(tags) > 0 {
			operation = operation + "(" + strings.Join(tags, ",") + ")"
		}
	}
	m["graphql"] = operation
	return m
}
