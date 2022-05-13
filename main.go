package main

import (
	"github.com/enustah/protoc-gen-gin-gateway/generate"
	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	protogen.Options{}.Run(func(gen *protogen.Plugin) error {
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			if err := generate.GenerateFile(gen, f); err != nil {
				return err
			}
		}
		return nil
	})
}
