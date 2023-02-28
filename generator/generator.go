package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"io/fs"
	"os"
	"path"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/albuquerq/servermock/schema"
	"github.com/albuquerq/servermock/template1"
)

type environment struct {
	Module      string
	Package     string
	DataPackage string
	TypeName    string
	Data        any
}

type Config struct {
	ModulePath  string
	Package     string
	DataPackage string
	TypeName    string
}

func createEnvironment(cfg Config) environment {
	return environment{
		Module:      cfg.ModulePath,
		Package:     cfg.Package,
		TypeName:    cfg.TypeName,
		DataPackage: cfg.DataPackage,
	}
}

func Generate(cfg Config, s schema.Schema, folder string) error {
	env := createEnvironment(cfg)

	packagePath := path.Join(folder, cfg.Package)

	if err := os.MkdirAll(packagePath, fs.ModePerm); err != nil {
		return fmt.Errorf("creating the destination folder: %w", err)
	}

	t := template.New("root").Funcs(template.FuncMap{
		"file":       readFile,
		"minifyjson": minifyjson,
		"now":        now,
	})

	t, err := t.ParseFS(template1.FS, "*.gotmpl")
	if err != nil {
		return fmt.Errorf("parsing template files: %w", err)
	}

	var b bytes.Buffer

	if err := t.Lookup("server.gotmpl").Execute(&b, env.With(s.Server)); err != nil {
		return fmt.Errorf("generating fake server file: %w", err)
	}

	if err := writeGoFile(path.Join(packagePath, "server.gen.go"), b.Bytes()); err != nil {
		return fmt.Errorf("writing fake server file: %w", err)
	}

	if err := os.MkdirAll(path.Join(packagePath, cfg.DataPackage), fs.ModePerm); err != nil {
		return fmt.Errorf("creating %s folder: %w", cfg.DataPackage, err)
	}

	for _, h := range s.Server.Handlers {
		b.Reset()

		filePrefix := toSnakeCase(h.Name)

		handlerEnv := env.With(h)

		if err := t.Lookup("handlers.gotmpl").Execute(&b, handlerEnv); err != nil {
			return fmt.Errorf("executing handlers template: %w", err)
		}

		handlerFileName := filePrefix + "_handlers.gen.go"

		if err := writeGoFile(path.Join(packagePath, handlerFileName), b.Bytes()); err != nil {
			return fmt.Errorf("writing handler file \"%s\": %w", handlerFileName, err)
		}

		b.Reset()

		if err := t.Lookup("responses.gotmpl").Execute(&b, handlerEnv); err != nil {
			return fmt.Errorf("executing responses template: %w", err)
		}

		responseFileName := filePrefix + "_responses.gen.go"

		if err := writeGoFile(path.Join(packagePath, cfg.DataPackage, responseFileName), b.Bytes()); err != nil {
			return fmt.Errorf("writing handler file \"%s\": %w", responseFileName, err)
		}
	}

	return nil
}

func (e environment) With(data any) environment {
	return environment{
		Package:     e.Package,
		DataPackage: e.DataPackage,
		Module:      e.Module,
		TypeName:    e.TypeName,
		Data:        data,
	}
}

func readFile(fileName string) string {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return ""
	}
	return string(data)
}

func now() string {
	return time.Now().Format(time.RFC3339)
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func minifyjson(s string) string {
	var b bytes.Buffer
	if err := json.Compact(&b, []byte(s)); err != nil {
		return fmt.Sprintf("invalid json: %v", err)
	}
	return b.String()
}

func writeGoFile(name string, data []byte) error {
	formated, err := format.Source(data)
	if err != nil {
		return err
	}
	return os.WriteFile(name, formated, fs.ModePerm)
}
