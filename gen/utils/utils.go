package utils

import (
	"bytes"
	"gqlclientgen/config"
	"os"
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
	TypeIgnoreList []string = []string{"Query", "Subscription", "Mutation"}
	TypeMappings            = map[string]TypeMapping{
		"String":   {MappedType: jen.String()},
		"Id":       {MappedType: jen.String()},
		"Int":      {MappedType: jen.Int()},
		"Float":    {MappedType: jen.Float64()},
		"Boolean":  {MappedType: jen.Bool()},
		"Any":      {MappedType: jen.Interface()},
		"Map":      {MappedType: jen.Map(jen.String()).Interface()},
		"Date":     {MappedType: jen.Qual("time", "Time")},
		"Time":     {MappedType: jen.Qual("time", "Time")},
		"DateTime": {MappedType: jen.Qual("time", "Time")},
	}
)

const (
	DefaultPackageName   = "gql"
	GqlClientPackageName = "github.com/hasura/go-graphql-client"
)

func GetPackagePath() string {
	gopath := os.Getenv("GOPATH")
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

	outDir, err = filepath.Rel(path.Join(gopath, "src"), outDir)
	if err != nil {
		logger.LogError(err)
		return ""
	}
	// currDir, _ := os.Getwd()
	// outDir, err = filepath.Rel(currDir, outDir)
	// if err != nil {
	// 	logger.LogError(err)
	// 	return ""
	// }
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

func ToCamelCase(s string) string {
	var g []string
	p := strings.Fields(s)
	for _, value := range p {
		g = append(g, strings.Title(value))
	}
	return strings.Join(g, "")
}

func GetMethodParams(s, str string) *jen.Statement {
	return jen.Id(s).Op("*").Id(ToCamelCase(str))
}

func GetEnumParam(s, str string, t bool) *jen.Statement {
	e := jen.Id(s)
	if t {
		e.Op("*")
	}
	return e.Id(ToCamelCase(str))
}

func GetClientParams() *jen.Statement {
	return GetMethodParams("c", "client")
}

func getModelPath() string {
	return path.Join(GetPackagePath(), "model")
}

func GetArgsType(arg *ast.ArgumentDefinition) *jen.Statement {
	fieldName := ToCamelCase(strings.ToLower(arg.Type.NamedType))
	if fieldName == "" {
		fieldName = ToCamelCase(arg.Type.Elem.Name())
		return jen.Index().Op("*").Qual(getModelPath(), ToCamelCase(fieldName))
	}
	fieldType, ok := TypeMappings[fieldName]
	if !ok {
		if arg.Type.NonNull {
			return jen.Qual(getModelPath(), ToCamelCase(arg.Type.Name()))
		}
		return jen.Qual(getModelPath(), ToCamelCase(arg.Type.Name()))
	}
	return fieldType.MappedType
}

func GetReturnType(field *ast.FieldDefinition) *jen.Statement {
	fieldName := ToCamelCase(strings.ToLower(field.Type.NamedType))
	if fieldName == "" {
		fieldName = ToCamelCase(field.Type.Elem.Name())
		return jen.Index().Op("*").Qual(getModelPath(), ToCamelCase(fieldName))
	}
	fieldType, ok := TypeMappings[fieldName]
	if !ok {
		if field.Type.NonNull {
			return jen.Op("*").Qual(getModelPath(), ToCamelCase(fieldName))
		}
		return jen.Op("*").Qual(getModelPath(), ToCamelCase(fieldName))
	}
	return fieldType.MappedType
}

func GetRequestType(field *ast.FieldDefinition) *jen.Statement {
	fieldName := ToCamelCase(strings.ToLower(field.Type.NamedType))
	if fieldName == "" {
		fieldName = ToCamelCase(field.Type.Elem.Name())
		return jen.Index().Qual(getModelPath(), ToCamelCase(fieldName))
	}
	fieldType, ok := TypeMappings[fieldName]
	if !ok {
		if field.Type.NonNull {
			return jen.Qual(getModelPath(), ToCamelCase(fieldName))
		}
		return jen.Qual(getModelPath(), ToCamelCase(fieldName))
	}
	return fieldType.MappedType
}

func GetRequestTags(operation string, arr []string) map[string]string {
	m := make(map[string]string)
	tag := operation
	v := []string{}
	for _, k := range arr {
		v = append(v, k+": &"+k)
		tag = "(" + strings.Join(v, ",") + ")"
	}
	m["graphql"] = tag
	return m
}

func ToSmallCamelCase(s string) string {
	if len(s) < 2 {
		return strings.ToLower(s)
	}
	bts := []byte(s)
	lc := bytes.ToLower([]byte{bts[0]})
	rest := bts[1:]
	return string(bytes.Join([][]byte{lc, rest}, nil))
}
