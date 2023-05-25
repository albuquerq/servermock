package main

import (
	"errors"
	"flag"
	"log"
	"os"

	"github.com/albuquerq/servermock/generator"
	"github.com/albuquerq/servermock/schema"
)

func init() {
	flag.StringVar(&arg.SchemaFilePath, "schema", "api-schema.yaml", "schema file path.")
	flag.StringVar(&arg.OutputFolder, "output-folder", ".", "destination folder of the generated code.")
	flag.StringVar(&arg.ProjectModule, "project-module", "", "project module uri ex.: \"github.com/albuquerq/fake\".")
	flag.StringVar(&arg.Package, "pkg", "apimock", "name for the geneated package.")
	flag.StringVar(&arg.DataPackage, "pkgdata", "testdata", "name for the generated data package.")
	flag.StringVar(&arg.TypeName, "gotype", "ServerMock", "go type name for the mock.")
	flag.BoolVar(&arg.StackHandler, "stack", false, "enable stack handler.")
}

type args struct {
	SchemaFilePath string
	OutputFolder   string
	ProjectModule  string
	Package        string
	DataPackage    string
	TypeName       string
	StackHandler   bool
}

func (a args) validate() error {
	if a.SchemaFilePath == "" {
		return errors.New("missing value of argument \"schema\"")
	}
	if a.OutputFolder == "" {
		return errors.New("missing value of argument \"output-folder\"")
	}
	if a.ProjectModule == "" {
		return errors.New("missing value of argument \"project-module\"")
	}
	if a.Package == "" {
		return errors.New("missing value of argument \"pkg\"")
	}
	if a.DataPackage == "" {
		return errors.New("missing value of argument \"pkgdata\"")
	}
	if a.TypeName == "" {
		return errors.New("missing value of argument \"gotype\"")
	}
	return nil
}

var arg args

func main() {
	flag.Parse()

	if err := arg.validate(); err != nil {
		log.Printf("validating arguments: %v", err)
		return
	}

	f, err := os.Open(arg.SchemaFilePath)
	if err != nil {
		log.Printf("failed to open file \"%s\": %v", arg.SchemaFilePath, err)
		return
	}
	defer f.Close()

	s, err := schema.Parse(f)
	if err != nil {
		log.Printf("failed parsing schema: %v", err)
		return
	}

	err = generator.Generate(
		generator.Config{
			ModulePath:  arg.ProjectModule,
			Package:     arg.Package,
			DataPackage: arg.DataPackage,
			TypeName:    arg.TypeName,
		},
		s,
		arg.OutputFolder,
	)
	if err != nil {
		log.Printf("failed go generate code: %v", err)
		return
	}
}
