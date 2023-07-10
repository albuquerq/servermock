package generator

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go/format"
	"io/fs"
	"os"
	"path"
	"regexp"
	"strings"
	"text/template"

	"github.com/albuquerq/servermock/schema"
	"github.com/albuquerq/servermock/template1"
)

type HandlerType int

const (
	TypeDefault HandlerType = iota
	TypeQueue
	TypeStack
)

func ParseHandlerType(s string) (HandlerType, error) {
	switch strings.ToLower(s) {
	case "default":
		return TypeDefault, nil
	case "queue":
		return TypeQueue, nil
	case "stack":
		return TypeStack, nil
	default:
		return 0, errors.New("unknown handler type")
	}
}

func (ht HandlerType) String() string {
	switch ht {
	case TypeDefault:
		return "default"
	case TypeQueue:
		return "queue"
	case TypeStack:
		return "stack"
	default:
		return "unknown"
	}
}

type environment struct {
	Module      string
	Package     string
	DataPackage string
	TypeName    string
	Data        any
	HandlerType HandlerType
}

type Config struct {
	ModulePath  string
	Package     string
	DataPackage string
	TypeName    string
	HandlerType HandlerType
}

func createEnvironment(cfg Config) environment {
	return environment{
		Module:      cfg.ModulePath,
		Package:     cfg.Package,
		TypeName:    cfg.TypeName,
		DataPackage: cfg.DataPackage,
		HandlerType: cfg.HandlerType,
	}
}

func Generate(cfg Config, s schema.Schema, folder string) error {
	env := createEnvironment(cfg)

	packagePath := path.Join(folder, cfg.Package)

	if err := os.MkdirAll(packagePath, fs.ModePerm); err != nil {
		return fmt.Errorf("creating the destination folder: %w", err)
	}

	tmpl := template.New("root").Funcs(template.FuncMap{
		"file":       readFile,
		"minifyjson": minifyjson,
	})

	tmpl, err := tmpl.ParseFS(template1.FS, "root/*.gotmpl")
	if err != nil {
		return fmt.Errorf("parsing template files: %w", err)
	}

	tmpl, err = tmpl.ParseFS(template1.FS, fmt.Sprintf("parts/%s.gotmpl", cfg.HandlerType))
	if err != nil {
		return fmt.Errorf("parting template definitions files: %w", err)
	}

	var b bytes.Buffer

	if err := tmpl.Lookup("server.gotmpl").Execute(&b, env.With(s.Server)); err != nil {
		return fmt.Errorf("generating server file: %w", err)
	}

	if err := writeGoFile(path.Join(packagePath, "server.gen.go"), b.Bytes()); err != nil {
		return fmt.Errorf("writing server file: %w", err)
	}

	if err := os.MkdirAll(path.Join(packagePath, cfg.DataPackage), fs.ModePerm); err != nil {
		return fmt.Errorf("creating %s folder: %w", cfg.DataPackage, err)
	}

	for _, h := range s.Server.Handlers {
		b.Reset()

		filePrefix := toSnakeCase(h.Name)

		handlerEnv := env.With(h)

		if err := tmpl.Lookup("handle_functions.gotmpl").Execute(&b, handlerEnv); err != nil {
			return fmt.Errorf("executing handler functions template: %w", err)
		}

		handlerFileName := filePrefix + "_handlers.gen.go"

		if err := writeGoFile(path.Join(packagePath, handlerFileName), b.Bytes()); err != nil {
			return fmt.Errorf("writing handler file \"%s\": %w", handlerFileName, err)
		}

		b.Reset()

		if err := tmpl.Lookup("responses.gotmpl").Execute(&b, handlerEnv); err != nil {
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
		HandlerType: e.HandlerType,
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
