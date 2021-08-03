package modelgen

import (
	"bytes"
	"gqlclientgen/config"
	"gqlclientgen/gen/context"
	"gqlclientgen/gen/utils"

	"os"
	"path"
	"strings"

	"github.com/spf13/viper"

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
	getAllTypes(parsedGql)
	build(parsedGql, context)
	jenFile := jen.NewFile(viper.GetViper().GetString(config.PackageNameKey))
	for _, v := range context.Model.Objects {
		jenFile.Add(v.CodeStatement)
	}
	for _, v := range context.Model.Enums {
		jenFile.Add(v.CodeStatement)
	}
	for _, v := range context.Model.Scalars {
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

func build(parsedGql *ast.Schema, c *context.Context) error {
	for _, def := range parsedGql.Types {
		if !def.BuiltIn {
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
			case ast.Interface:
				err := createInterface(def, c)
				if err != nil {
					return err
				}
			case ast.Union:
				err := createUnion(def, c)
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
	objStruct := jen.Type().Id(utils.ToPascalCase(def.Name))
	for _, field := range def.Fields {
		fieldType.Id(utils.ToPascalCase(field.Name))
		fieldType.Add(getType(field))
		fieldType.Tag(map[string]string{"json": utils.GetTags(field)}).Line()
	}
	objStruct.Struct(fieldType).Line()
	if def.Interfaces != nil && len(def.Interfaces) > 0 {
		for _, v := range def.Interfaces {
			objStruct.Add(jen.Func().Parens(jen.Id(def.Name)).Id("Is" + v).Call().Block()).Line()
		}
	}
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
	enumType := jen.Type().Id(utils.ToPascalCase(def.Name)).String().Line()
	enumConst := jen.Const()
	for _, enum := range def.EnumValues {
		enumName := utils.ToPascalCase(def.Name + utils.ToPascalCase(strings.ToLower(enum.Name)))
		enumValues = append(enumValues, enumName)
		enumConstValues.Add(jen.Id(enumName).Id(utils.ToPascalCase(def.Name)).Op("=").Lit(enum.Name)).Line()
		enumArrayValues.List(jen.Id(enumName).Op(","))
	}
	enumConst.Defs(enumConstValues).Line()
	enumType.Add(enumConst)
	enumArray := jen.Var().Id(utils.ToPascalCase("All" + def.Name)).Op("=").Index().Id(utils.ToPascalCase(def.Name)).Values(enumArrayValues).Line()
	enumType.Add(enumArray).Line()
	enumType.Add(getEnumMethods(utils.ToPascalCase(def.Name), enumValues)).Line()
	obj.GraphqlName = def.Name
	obj.MappedName = strings.ToLower(def.Name)
	obj.MappedType = strings.ToLower(string(def.Kind))
	obj.CodeStatement = enumType
	c.Model.Enums = append(c.Model.Enums, obj)
	return nil
}

func createFiles() (*os.File, error) {
	p := path.Join(utils.GetFilePath(), "model.go")
	if err := os.MkdirAll(path.Dir(p), os.ModePerm); err != nil {
		return nil, err
	}
	f, err := os.Create(p)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func getType(field *ast.FieldDefinition) *jen.Statement {
	fieldName := strings.ToLower(field.Type.NamedType)
	fieldType, ok := utils.TypeMappings[strings.ToLower(field.Type.Name())]
	if fieldName == "" && !ok {
		fieldName = utils.ToPascalCase(field.Type.Elem.Name())
		return jen.Index().Op("*").Id(fieldName)
	}
	if !ok {
		if !utils.StringInSlice(strings.ToLower(field.Type.Name()), utils.AllTypes) {
			return jen.Op("*").String()
		}
		if field.Type.NonNull {
			return jen.Op("*").Id(utils.ToPascalCase(field.Type.Name()))
		}
		return jen.Id(utils.ToPascalCase(field.Type.Name()))
	}
	if field.Type.Elem != nil {
		return jen.Add(checkForIndex(field.Type)).Add(fieldType.MappedType)
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

func createInterface(def *ast.Definition, c *context.Context) error {
	if utils.StringInSlice(def.Name, utils.TypeIgnoreList) {
		return nil
	}

	var (
		obj = &context.DataTypeInfo{}
	)
	jenType := jen.Type().Id(def.Name).Interface(
		jen.Id("Is" + def.Name).Call(),
	).Line()
	obj.GraphqlName = def.Name
	obj.MappedName = strings.ToLower(def.Name)
	obj.MappedType = strings.ToLower(string(def.Kind))
	obj.CodeStatement = jenType
	c.Model.Objects = append(c.Model.Objects, obj)
	return nil
}

func createUnion(def *ast.Definition, c *context.Context) error {
	if utils.StringInSlice(def.Name, utils.TypeIgnoreList) {
		return nil
	}

	var (
		obj = &context.DataTypeInfo{}
	)
	jenType := jen.Type().Id(def.Name).Interface(
		jen.Id("Is" + def.Name).Call(),
	).Line()
	if def.Types != nil {
		for _, t := range def.Types {
			jenType.Add(jen.Func().Parens(jen.Id(t)).Id("Is" + def.Name).Call().Block()).Line()
		}
	}
	obj.GraphqlName = def.Name
	obj.MappedName = strings.ToLower(def.Name)
	obj.MappedType = strings.ToLower(string(def.Kind))
	obj.CodeStatement = jenType
	c.Model.Objects = append(c.Model.Objects, obj)
	return nil
}

func checkForIndex(field *ast.Type) *jen.Statement {
	index := &jen.Statement{}
	if field.Elem != nil {
		index.Add(jen.Index())
		index.Add(checkForIndex(field.Elem))
	}
	return index
}

func getAllTypes(parsedGql *ast.Schema) {
	for _, def := range parsedGql.Types {
		if !def.BuiltIn {
			utils.AllTypes = append(utils.AllTypes, strings.ToLower(def.Name))
		}
	}
}
