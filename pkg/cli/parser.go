package cli

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

// parse turns a list of bites into a list of templates
func parse(s string) ([]Template, error) {

	templates := []Template{}

	err := yaml.Unmarshal([]byte(s), &templates)

	if err != nil {
		return nil, fmt.Errorf("error parsing yaml: %v", err)
	}

	return templates, nil

}

// open opens a file and returns the contents as a string
func open(path string) (string, error) {
	data, err := os.ReadFile(path)

	if err != nil {
		return "", fmt.Errorf("error reading file: %v", err)
	}

	return string(data), nil
}

// read reads a file and returns a list of templates
func read(path string) ([]Template, error) {
	data, err := open(path)

	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	var templates []Template
	templates, err = parse(data)

	if err != nil {
		return nil, fmt.Errorf("error parsing file: %v", err)
	}

	return templates, nil
}

// ReadDefault reads the default template file
func ReadDefault() ([]Template, error) {
	possible := []string{
		"comtemplate.yml",
		"comtemplate.yaml",
	}

	for _, path := range possible {
		templates, err := read(path)
		if err == nil {
			return templates, nil
		}
	}

	return nil, fmt.Errorf("Read Default: no default template found")

}

func CreateDefault() error {
	defaultName := "comtemplate.yml"

	defaultTemplate := []byte(`
- name: Commit Template 1
  text: |
    %{title}

    %{body}
  variables:
    - name: title
    - name: body

`)

	err := os.WriteFile(defaultName, defaultTemplate, 0644)

	if err != nil {
		return fmt.Errorf("error creating default template: %v", err)
	}

	return nil
}
