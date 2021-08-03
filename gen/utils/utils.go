package utils

import (
	"bytes"
	"gqlclientgen/config"

	"path"
	"path/filepath"
	"strings"

	"github.com/Coditation/skael-connectors-shared/logger"
	"github.com/dave/jennifer/jen"
	"github.com/spf13/viper"
	"github.com/vektah/gqlparser/v2/ast"
)

type TypeMapping struct {
	MappedType *jen.Statement
}

var (
	AllTypes       []string
	TypeIgnoreList []string = []string{"Query", "Subscription", "Mutation"}
	TypeMappings            = map[string]TypeMapping{
		"string":   {MappedType: jen.String()},
		"id":       {MappedType: jen.String()},
		"int":      {MappedType: jen.Int()},
		"float":    {MappedType: jen.Float64()},
		"boolean":  {MappedType: jen.Bool()},
		"any":      {MappedType: jen.Interface()},
		"map":      {MappedType: jen.Map(jen.String()).Interface()},
		"date":     {MappedType: jen.Op("*").Qual("time", "Time")},
		"time":     {MappedType: jen.Op("*").Qual("time", "Time")},
		"datetime": {MappedType: jen.Op("*").Qual("time", "Time")},
	}
)

const (
	DefaultPackageName   = "gql"
	GqlClientPackageName = "github.com/hasura/go-graphql-client"
)

func GetPackagePath() string {
	v := viper.GetViper()
	packageName := v.GetString(config.PackageNameKey)
	if packageName == "" {
		packageName = DefaultPackageName
	}
	outDir, err := filepath.Abs(v.GetString(config.OutputDirectoryKey))
	if err != nil {
		logger.LogError(err)
		return ""
	}
	//
	//outDir, err = filepath.Rel(path.Join(gopath, "src"), outDir)
	//if err != nil {
	//	logger.LogError(err)
	//	return ""
	//}
	return path.Join(outDir, packageName)
}

func GetFilePath() string {
	v := viper.GetViper()
	packageName := v.GetString(config.PackageNameKey)
	if packageName == "" {
		packageName = DefaultPackageName
	}
	outDir, err := filepath.Abs(v.GetString(config.OutputDirectoryKey))
	if err != nil {
		logger.LogError(err)
		return ""
	}
	return path.Join(outDir, packageName)
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func GetTags(f *ast.FieldDefinition) string {
	if f.Type.Elem != nil {
		return f.Name + ",omitempty"
	}
	if f.Type.NonNull {
		return f.Name + ",omitempty"
	}
	return f.Name
}

func ToPascalCase(s string) string {
	var g []string
	p := strings.Fields(s)
	for _, value := range p {
		g = append(g, strings.Title(value))
	}
	return strings.Join(g, "")
}

func GetMethodParams(s, str string) *jen.Statement {
	return jen.Id(s).Op("*").Id(ToPascalCase(str))
}

func GetEnumParam(s, str string, t bool) *jen.Statement {
	e := jen.Id(s)
	if t {
		e.Op("*")
	}
	return e.Id(ToPascalCase(str))
}

func GetClientParams() *jen.Statement {
	return GetMethodParams("c", "client")
}

func getModelPath() string {
	return ""
}

func GetArgsType(arg *ast.ArgumentDefinition) *jen.Statement {
	fieldName := ToPascalCase(arg.Type.NamedType)
	if fieldName == "" {
		fieldName = ToPascalCase(arg.Type.Elem.Name())
		return jen.Index().Op("*").Qual(getModelPath(), ToPascalCase(fieldName))
	}
	fieldType, ok := TypeMappings[strings.ToLower(fieldName)]
	if !ok {
		if arg.Type.NonNull {
			return jen.Qual(getModelPath(), ToPascalCase(arg.Type.Name()))
		}
		return jen.Qual(getModelPath(), ToPascalCase(arg.Type.Name()))
	}
	return fieldType.MappedType
}

func GetReturnType(field *ast.FieldDefinition) *jen.Statement {
	fieldName := ToPascalCase(field.Type.NamedType)
	if fieldName == "" {
		fieldName = ToPascalCase(field.Type.Elem.Name())
		return jen.Index().Op("*").Qual(getModelPath(), ToPascalCase(fieldName))
	}
	fieldType, ok := TypeMappings[strings.ToLower(fieldName)]
	if !ok {
		if field.Type.NonNull {
			return jen.Op("*").Qual(getModelPath(), ToPascalCase(fieldName))
		}
		return jen.Op("*").Qual(getModelPath(), ToPascalCase(fieldName))
	}
	return fieldType.MappedType
}

func GetVarType(field *ast.FieldDefinition) *jen.Statement {
	fieldName := ToPascalCase(field.Type.NamedType)
	if fieldName == "" {
		fieldName = ToPascalCase(field.Type.Elem.Name())
		return jen.Index().Op("*").Qual(getModelPath(), ToPascalCase(fieldName))
	}
	fieldType, ok := TypeMappings[strings.ToLower(fieldName)]
	if !ok {
		if field.Type.NonNull {
			return jen.Qual(getModelPath(), ToPascalCase(fieldName))
		}
		return jen.Qual(getModelPath(), ToPascalCase(fieldName))
	}
	return fieldType.MappedType
}

func GetReturnFuncType(field *ast.FieldDefinition) *jen.Statement {
	fieldName := ToPascalCase(field.Type.NamedType)
	if fieldName == "" {

		return jen.Id("res")
	}
	return jen.Op("&").Id("res")
}

func GetRequestType(field *ast.FieldDefinition) *jen.Statement {
	fieldName := ToPascalCase(field.Type.NamedType)
	if fieldName == "" {
		fieldName = ToPascalCase(field.Type.Elem.Name())
		return jen.Index().Qual(getModelPath(), ToPascalCase(fieldName))
	}
	fieldType, ok := TypeMappings[strings.ToLower(fieldName)]
	if !ok {
		if field.Type.NonNull {
			return jen.Qual(getModelPath(), ToPascalCase(fieldName))
		}
		return jen.Qual(getModelPath(), ToPascalCase(fieldName))
	}
	return fieldType.MappedType
}

func GetRequestTags(operation string, arr []string) map[string]string {
	m := make(map[string]string)

	v := []string{}
	for _, k := range arr {
		v = append(v, k+": &"+k)
	}
	operation = operation + "(" + strings.Join(v, ",") + ")"
	m["graphql"] = operation
	return m
}

func ToCamelCase(s string) string {
	if len(s) < 2 {
		return strings.ToLower(s)
	}
	bts := []byte(s)
	lc := bytes.ToLower([]byte{bts[0]})
	rest := bts[1:]
	return string(bytes.Join([][]byte{lc, rest}, nil))
}
