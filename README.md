<a href="https://www.coditation.com"><img src="https://www.coditation.com/wp-content/uploads/2020/08/Small-Logo@4x-2.png" align="left" height="20" width="124" ></a></br>

# GraphQL Client/SDK Generator for Golang

A Golang client side code/SDK generator for consuming the GraphQL APIs.

## How to use?

### Build

```bash
go build gqlclientgen
```

### Run

```bash
gqlclientgen -config_path=<path of the directory containing the config.yaml> -plugin_path=<path of the directory where all the plugin of custom scalars> -query_path=<path of the all operations with fragments>
```

### Create Custom Scalar Plugin

```
The Custom Scalar Plugin should implement the methods:
    Type() string         // type of the custom scalar
	Code() *jen.Statement // code of dave/jennifer please refer [here](https://github.com/dave/jennifer)

The File name of the custom scalar plugin i.e TimeStamp.go

```

```
This file should be built in a form of plugin using the following command
```

    go build -buildmode=plugin TimeStamp.go TimeStamp.so

### Config (YAML)

```yaml
# Name of the package for generated client/SDK
packageName: "gql"
# Path of directory where generated code will be placed
outputDirectory: "/somedir"
# Default
sourceType: "file"
# Path of the graphql schema (IDL) file
sourceFilePath: "/somedir/schema.graphqls"
```

## Features

- [x] Golang model generation
- [x] Query
- [x] Mutations
- [x] Scalars - Map, Any, ID
- [x] Support for custom scalars
- [ ] GraphQL directives
- [x] GraphQL fragments
- [ ] Subscriptions
- [x] Load schema from URL

## License

This software is licensed under [AGPL v3.0](https://choosealicense.com/licenses/agpl-3.0/)
