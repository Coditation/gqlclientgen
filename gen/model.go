package gen

import (
	"os"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/vektah/gqlparser/v2/ast"
)

type typeMapping struct {
	mappedType *jen.Statement
	mappedTag  *jen.Statement
}

const (
	gqlClientPackageName = "github.com/hasura/go-graphql-client"
)

var (
	typeIgnoreList []string = []string{"Query", "Subscription", "Mutation"}
	typeMappings            = map[string]typeMapping{
		"String":  {mappedType: jen.Qual(gqlClientPackageName, "String"), mappedTag: nil},
		"ID":      {mappedType: jen.Qual(gqlClientPackageName, "ID"), mappedTag: nil},
		"Int":     {mappedType: jen.Qual(gqlClientPackageName, "Int"), mappedTag: nil},
		"Float":   {mappedType: jen.Qual(gqlClientPackageName, "Float"), mappedTag: nil},
		"Boolean": {mappedType: jen.Qual(gqlClientPackageName, "Boolean"), mappedTag: nil},
	}
)

func GenerateModel(parsedGql *ast.Schema) error {
	context := Create()
	f, err := createFiles()
	if err != nil {
		return err
	}
	defer f.Close()
	build(parsedGql, context)
	return nil
}

func build(parsedGql *ast.Schema, context *Context) error {
	for _, def := range parsedGql.Types {
		if def.BuiltIn != true {
			kind := def.Kind
			switch kind {
			case ast.Scalar:
				err := createScalar(def, context)
				if err != nil {
					panic("error in creating the scalar")
				}
			case ast.Object, ast.InputObject:
				err := createObject(def, context)
				if err != nil {
					panic("error in creating the object")
				}
			}
		}
	}
	return nil
}

func createScalar(def *ast.Definition, context *Context) error {
	context.Model.Scalars = append(context.Model.Scalars, &DataTypeInfo{
		GraphqlName: def.Name,
		MappedName:  strings.ToLower(def.Name),
		MappedType:  strings.ToLower(string(def.Kind)),
	})
	return nil
}

func createObject(def *ast.Definition, context *Context) error {
	if StringInSlice(def.Name, typeIgnoreList) == true {
		return nil
	}
	obj := &DataTypeInfo{}
	objStruct := jen.Type().Id(def.Name)
	for _, field := range def.Fields {
		objStruct.Struct().Append(jen.Id(strings.Title(field.Name)).Tag(map[string]string{"json": getTags(field)}))
	}
	obj.GraphqlName = def.Name
	obj.MappedName = strings.ToLower(def.Name)
	obj.MappedType = strings.ToLower(string(def.Kind))
	obj.CodeStatement = objStruct
	context.Model.Objects = append(context.Model.Objects, obj)
	return nil
}

func createQuery(def *ast.Definition, context *Context) error {

	return nil
}

func createMutation(def *ast.Definition, context *Context) error {
	return nil
}

func createEnum(def *ast.Definition, context *Context) error {
	// obj := &DataTypeInfo{}

	return nil
}

func createFiles() (*os.File, error) {
	path := GetPackagePath() + "/models.go"
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	} else {
		return f, nil
	}
}
