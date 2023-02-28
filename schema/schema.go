package schema

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

// TODO: YAML schema validation.

type Schema struct {
	Version int    `yaml:"version"`
	Server  Server `yaml:"server"`
}

type Server struct {
	BaseURL  string    `yaml:"base_url"`
	Handlers []Handler `yaml:"handlers"`
}

type Handler struct {
	Name      string            `yaml:"name"`
	Method    string            `yaml:"method"`
	Path      string            `yaml:"path"`
	Headers   map[string]string `yaml:"headers"`
	Requests  []Request         `yaml:"requests"`
	Responses []Response        `yaml:"responses"`
}

type Request struct {
	Name     string `yaml:"name"`
	BodyFile string `yaml:"body"`
}

type Response struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	StatusCode  int               `yaml:"status_code"`
	Headers     map[string]string `yaml:"headers"`
	Body        string            `yaml:"body"`
}

func Parse(r io.Reader) (Schema, error) {
	var s Schema
	if err := yaml.NewDecoder(r).Decode(&s); err != nil {
		return Schema{}, fmt.Errorf("decoding server schema: %w", err)
	}
	// TODO: validate the schema.
	return s, nil
}
