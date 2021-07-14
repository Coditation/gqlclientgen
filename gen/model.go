package gen

import (
	"github.com/dave/jennifer/jen"
	"github.com/vektah/gqlparser/v2/ast"
	"os"
)

type typeMapping struct {
	mappedType *jen.Statement
	mappedTag  *jen.Statement
}

const (
	scalarKind           = "SCALAR"
	objectKind           = "OBJECT"
	inputObjectKind      = "INPUT_OBJECT"
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
			case scalarKind:
				err := createScalar(def, context)
				if err != nil {
					panic("error in creating the scalar")
				}
			case objectKind:
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

	return nil
}

func createObject(def *ast.Definition, context *Context) error {
	if StringInSlice(def.Name, typeIgnoreList) == true {
		return nil
	}
	//var fields []*jen.Statement
	//for _, fd := range def.Fields {
	//	fd.Type.
	//
	//}
	//c := jen.Type().Id(strings.Title(strings.ToLower(def.Name))).Struct(
	//
	//	)
	return nil
}

func createQuery(parsedGql []*ast.Source, file *os.File) error {
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
