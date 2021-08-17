package context

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
	PackageName   string
	Queries       []*DataTypeInfo
	Mutations     []*DataTypeInfo
	Subscriptions []*DataTypeInfo
	Objects       []*DataTypeInfo
	Enums         []*DataTypeInfo
	Scalars       []*DataTypeInfo
}

type ClientInfo struct {
	Client []*DataTypeInfo
}

type Context struct {
	Model  *ModelInfo
	Client *ClientInfo
}

var context *Context
var once sync.Once

func Create() *Context {
	once.Do(func() {
		context = NewContext()
		context.Model.Scalars = make([]*DataTypeInfo, 0)
		context.Model.Queries = make([]*DataTypeInfo, 0)
		context.Model.Mutations = make([]*DataTypeInfo, 0)
		context.Model.Enums = make([]*DataTypeInfo, 0)
		context.Model.Objects = make([]*DataTypeInfo, 0)
		context.Model.Subscriptions = make([]*DataTypeInfo, 0)
		context.Client.Client = make([]*DataTypeInfo, 0)
	})
	return context
}

func NewContext() *Context {
	return &Context{
		Model:  &ModelInfo{},
		Client: &ClientInfo{},
	}
}
