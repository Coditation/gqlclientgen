package main

import (
	config2 "coditation.com/gqlclientgen/config"
	"coditation.com/gqlclientgen/schema"
	"github.com/vektah/gqlparser/v2"
	"os"
)

func main() {
	arg := os.Args[1]
	if arg == "" {
		panic("please specify a configuration YAML file")
	}
	err := config2.LoadConfig(arg)
	if err != nil {
		panic("cannot read configuration")
	}
	loader := schema.GetLoader()
	sources, err := loader.Load()
	parsedSchema := gqlparser.MustLoadSchema(sources...)
	_ = parsedSchema
}
