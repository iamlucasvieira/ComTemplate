package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Template is a git commit template
type Template struct {
	Name      string     `yaml:"name"`
	Text      string     `yaml:"text"`
	Variables []Variable `yaml:"variables"`
}

// Varialbe is a variable in a git commit template
type Variable struct {
	Name string `yaml:"name"`
}

// Parse turns a list of bites into a list of templates
func Parse(s string) ([]Template, error) {

	templates := []Template{}

	err := yaml.Unmarshal([]byte(s), &templates)

	if err != nil {
		return nil, fmt.Errorf("error parsing yaml: %v", err)
	}

	return templates, nil

}

// Open opens a file and returns the contents as a string
func open(path string) (string, error) {
	data, err := os.ReadFile(path)

	if err != nil {
		return "", fmt.Errorf("error reading file: %v", err)
	}

	return string(data), nil
}
