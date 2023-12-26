package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
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

	var validTemplates []Template
	for i, template := range templates {
		err = validateTemplate(i, template)
		if err != nil {
			fmt.Println(err)
		} else {
			validTemplates = append(validTemplates, template)
		}

	}
	return validTemplates, nil

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

// CreateDefault creates a default template file
func CreateDefault() error {
	defaultName := "comtemplate.yml"

	defaultTemplate := []byte(`
- name: simple
  text: |
    %{title}

    %{body}
  variables:
    - name: title
    - name: body
- name: type
  text: |
    [%{type}] %{title}

    %{body}
  variables:
    - name: type
    - name: title
    - name: body
`)
	// Check if file exists
	if _, err := os.Stat(defaultName); err == nil {
		return fmt.Errorf("File %s already exists", defaultName)
	}

	err := os.WriteFile(defaultName, defaultTemplate, 0644)

	if err != nil {
		return fmt.Errorf("error creating default template: %v", err)
	}

	return nil
}

// validateTemplate checks if a template is valid
func validateTemplate(id int, t Template) error {

	var err string
	if t.Name == "" {
		err += fmt.Sprintf("Template %d: name is empty\n", id)
	}

	if t.Text == "" {
		err += fmt.Sprintf("Template %d: text is empty\n", id)
	}

	for varId, variable := range t.Variables {
		if variable.Name == "" {
			err += fmt.Sprintf("Template %d - Variable %d: variable name is empty\n", id, varId)
		}

		if !strings.Contains(t.Text, fmt.Sprintf("%%{%s}", variable.Name)) {
			err += fmt.Sprintf("Template %d - Variable %d: variable %s not found in template text\n", id, varId, variable.Name)
		}
	}

	if err != "" {
		return fmt.Errorf(err)
	}

	return nil
}

// PopulateTemplate replaces variables in a template with values
func PopulateTemplate(template Template, variables map[string]string) (string, error) {
	text := template.Text

	for _, variable := range template.Variables {
		value, ok := variables[variable.Name]

		if !ok {
			return "", fmt.Errorf("variable %s not found", variable.Name)
		}

		varName := fmt.Sprintf("%%{%s}", variable.Name)

		text = strings.Replace(text, varName, value, -1)
	}

	return text, nil
}

func PopulateFromForm(template Template) (string, error) {
	var variables = make(map[string]string)
	var inputList []huh.Field

	// Create a slice for intermediate storage
	inputValues := make([]string, len(template.Variables))

	for i, variable := range template.Variables {
		input := huh.NewInput().
			Title(variable.Name).
			Value(&inputValues[i])

		inputList = append(inputList, input)
	}

	form := huh.NewForm(huh.NewGroup(inputList...))

	err := form.Run()

	if err != nil {
		return "", fmt.Errorf("error running form: %v", err)
	}

	// Update the map with the values from the form
	for i, variable := range template.Variables {
		variables[variable.Name] = inputValues[i]
	}

	return PopulateTemplate(template, variables)
}
