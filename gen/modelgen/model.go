package modelgen

import (
	"gqlclientgen/gen/context"
	"gqlclientgen/gen/utils"

	"os"
	"path"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/vektah/gqlparser/v2/ast"
)

func GenerateModel(parsedGql *ast.Schema) error {
	context := context.Create()
	f, err := createFiles()
	if err != nil {
		return err
	}
	defer f.Close()
	Build(parsedGql, context)
	return nil
}

func Build(parsedGql *ast.Schema, c *context.Context) error {
	for _, def := range parsedGql.Types {
		if def.BuiltIn != true {
			kind := def.Kind
			switch kind {
			case ast.Scalar:
				err := createScalar(def, c)
				if err != nil {
					return err
				}
			case ast.Object, ast.InputObject:
				err := createObject(def, c)
				if err != nil {
					return err
				}
			case ast.Enum:
				err := createEnum(def, c)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func createScalar(def *ast.Definition, c *context.Context) error {
	if utils.StringInSlice(def.Name, utils.TypeIgnoreList) {
		return nil
	}
	c.Model.Scalars = append(c.Model.Scalars, &context.DataTypeInfo{
		GraphqlName: def.Name,
		MappedName:  strings.ToLower(def.Name),
		MappedType:  strings.ToLower(string(def.Kind)),
	})
	return nil
}

func createObject(def *ast.Definition, c *context.Context) error {
	if utils.StringInSlice(def.Name, utils.TypeIgnoreList) {
		return nil
	}
	var (
		obj       = &context.DataTypeInfo{}
		fieldType = &jen.Statement{}
	)
	objStruct := jen.Type().Id(utils.ToCamelCase(def.Name))
	for _, field := range def.Fields {
		fieldType.Id(utils.ToCamelCase(field.Name))
		fieldType.Add(getType(field))
		fieldType.Tag(map[string]string{"json": utils.GetTags(field)}).Line()
	}
	objStruct.Struct(fieldType).Line()
	obj.GraphqlName = def.Name
	obj.MappedName = strings.ToLower(def.Name)
	obj.MappedType = strings.ToLower(string(def.Kind))
	obj.CodeStatement = objStruct
	c.Model.Objects = append(c.Model.Objects, obj)
	return nil
}

func createEnum(def *ast.Definition, c *context.Context) error {
	if utils.StringInSlice(def.Name, utils.TypeIgnoreList) {
		return nil
	}
	var (
		obj             = &context.DataTypeInfo{}
		enumArrayValues = &jen.Statement{}
		enumConstValues = &jen.Statement{}
		enumValues      []string
	)
	enumType := jen.Type().Id(utils.ToCamelCase(def.Name)).String().Line()
	enumConst := jen.Const()
	for _, enum := range def.EnumValues {
		enumName := utils.ToCamelCase(def.Name + utils.ToCamelCase(strings.ToLower(enum.Name)))
		enumValues = append(enumValues, enumName)
		enumConstValues.Add(jen.Id(enumName).Id(utils.ToCamelCase(def.Name)).Op("=").Lit(enum.Name)).Line()
		enumArrayValues.List(jen.Id(enumName).Op(","))
	}
	enumConst.Defs(enumConstValues).Line()
	enumType.Add(enumConst)
	enumArray := jen.Var().Id(utils.ToCamelCase("All" + def.Name)).Op("=").Index().Id(utils.ToCamelCase(def.Name)).Values(enumArrayValues).Line()
	enumType.Add(enumArray).Line()
	enumType.Add(getEnumMethods(utils.ToCamelCase(def.Name), enumValues)).Line()
	obj.GraphqlName = def.Name
	obj.MappedName = strings.ToLower(def.Name)
	obj.MappedType = strings.ToLower(string(def.Kind))
	obj.CodeStatement = enumType
	c.Model.Enums = append(c.Model.Enums, obj)
	return nil
}

func createFiles() (*os.File, error) {
	p := utils.GetPackagePath() + "/model/models.go"
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

func getType(field *ast.FieldDefinition) *jen.Statement {
	fieldName := utils.ToCamelCase(strings.ToLower(field.Type.NamedType))
	if fieldName == "" {
		fieldName = utils.ToCamelCase(field.Type.Elem.Name())
		return jen.Index().Op("*").Id(fieldName)
	}
	fieldType, ok := utils.TypeMappings[fieldName]
	if !ok {
		if field.Type.NonNull {
			return jen.Op("*").Id(utils.ToCamelCase(field.Type.Name()))
		}
		return jen.Id(utils.ToCamelCase(field.Type.Name()))
	}
	return fieldType.MappedType
}

func getEnumMethods(s string, arr []string) *jen.Statement {
	c := &jen.Statement{}
	startChar := strings.ToLower(string(s[0]))

	c.Add(jen.Func().Params(utils.GetEnumParam(startChar, s, false)).Id("IsValid").Call().Bool().Block(
		jen.Switch(jen.Id(startChar)).Block(
			jen.Case(jen.Id(strings.Join(arr, ","))).Block(
				jen.Return(jen.Lit(true)),
			),
		),
		jen.Return(jen.Lit(false)),
	)).Line()

	c.Add(jen.Func().Params(utils.GetEnumParam(startChar, s, false)).Id("String").Call().String().Block(
		jen.Return(jen.String().Params(jen.Id(startChar))),
	)).Line()

	c.Add(jen.Func().Params(utils.GetEnumParam(startChar, s, true)).Id("UnmarshalGQL").Params(jen.Id("v").Interface()).Error().Block(
		jen.List(jen.Id("str"), jen.Id("ok")).Op(":=").Id("v").Assert(jen.String()),
		jen.If(jen.Op("!").Id("ok")).Block(
			jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("enums must be strings"))),
		),
		jen.Op("*").Id(startChar).Op("=").Id(s).Params(jen.Id("str")),
		jen.If(jen.Op("!").Id(startChar).Dot("IsValid").Call().Block(
			jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("enums must be strings"))),
		)),
		jen.Return(jen.Nil()),
	)).Line()

	c.Add(jen.Func().Params(utils.GetEnumParam(startChar, s, false)).Id("MarshalGQL").Params(jen.Id("w").Qual("io", "Writer")).Block(
		jen.Qual("fmt", "Fprint").Call(jen.List(jen.Id("w"), jen.Qual("strconv", "Quote").Params(jen.Id(startChar).Dot("String").Call()))),
	)).Line()
	return c
}

func getModelPath() string {
	return path.Join(utils.GetPackagePath(), "model")
}
