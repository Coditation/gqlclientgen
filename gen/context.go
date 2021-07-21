package gen

import (
	"sync"

	"github.com/dave/jennifer/jen"
)

type DataTypeInfo struct {
	GraphqlName   string
	MappedName    string
	MappedType    string
	CodeStatement *jen.Statement
}

type ModelInfo struct {
	PackageName string
	Types       []*DataTypeInfo
	Queries     []*DataTypeInfo
	Mutations   []*DataTypeInfo
	Interfaces  []*DataTypeInfo
	Objects     []*DataTypeInfo
	Enums       []*DataTypeInfo
	Scalars     []*DataTypeInfo
}

type ClientInfo struct{}

type Context struct {
	Model  *ModelInfo
	Client *ClientInfo
}

var context *Context
var once sync.Once

func Create() *Context {
	once.Do(func() {
		context = new(Context)
		context.Model.Scalars = make([]*DataTypeInfo, 0)
		context.Model.Queries = make([]*DataTypeInfo, 0)
		context.Model.Mutations = make([]*DataTypeInfo, 0)
		context.Model.Types = make([]*DataTypeInfo, 0)
		context.Model.Interfaces = make([]*DataTypeInfo, 0)
		context.Model.Enums = make([]*DataTypeInfo, 0)
		context.Model.Objects = make([]*DataTypeInfo, 0)
	})
	return context
}
