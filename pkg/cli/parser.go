package cli

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/charmbracelet/huh"
	"gopkg.in/yaml.v3"
)

// Template is a git commit template
type Template struct {
	Name        string     `yaml:"name"`
	Description string     `yaml:"description"`
	Text        string     `yaml:"text"`
	Variables   []Variable `yaml:"variables"`
}

func (t Template) validate(id int) error {
	var errBuilder strings.Builder

	if t.Name == "" {
		fmt.Fprintf(&errBuilder, "Template %d: name is empty\n", id)
	}

	if t.Text == "" {
		fmt.Fprintf(&errBuilder, "Template %d: text is empty\n", id)
	}

	for varId, variable := range t.Variables {
		if !strings.Contains(t.Text, fmt.Sprintf("%%{%s}", variable.Name)) {
			fmt.Fprintf(&errBuilder, "Template %d - Variable %d: variable %s is not used in text\n", id, varId, variable.Name)
		}

		if validationErr := variable.validate(id, varId); validationErr != nil {
			fmt.Fprintf(&errBuilder, "%v\n", validationErr)
		}
	}

	if errBuilder.Len() > 0 {
		return fmt.Errorf(errBuilder.String())
	}

	return nil

}

// Varialbe is a variable in a git commit template
type Variable struct {
	Name    string   `yaml:"name"`
	Type    string   `yaml:"type"`
	Options []string `yaml:"options"`
}

func (v Variable) validate(TemplateId, VarId int) error {
	inputTypes := []string{"", "input", "text", "select"}
	var errBuilder strings.Builder

	if v.Name == "" {
		fmt.Fprintf(&errBuilder, "Template %d - Variable %d: variable name is empty\n", TemplateId, VarId)
	}

	if !slices.Contains(inputTypes, v.Type) {
		fmt.Fprintf(&errBuilder, "Template %d - Variable %d: variable %s has invalid type %s\n", TemplateId, VarId, v.Name, v.Type)
	}

	if errBuilder.Len() > 0 {
		return fmt.Errorf(errBuilder.String())
	}

	return nil
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
		err = template.validate(i)
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
- name: 1
  description: Simple commit message
  text: |
    %{description}

    %{body}
  variables:
    - name: description 
    - name: body
      type: text
- name: 2
  description: Commit message with type
  text: |
    [%{type}] %{description}

    %{body}
  variables:
    - name: type
      type: select
      options:
        - ✨ feat
        - 🐛 fix
        - ♻️  refactor
        - 📝 docs
        - 🎨 style
        - ✅ test
        - ⚡️ perf
    - name: description
    - name: body
      type: text
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
		var input huh.Field
		switch variable.Type {
		case "input", "":
			input = huh.NewInput().
				Title(variable.Name).
				Value(&inputValues[i])
		case "text":
			input = huh.NewText().
				Title(variable.Name).
				Value(&inputValues[i])
		case "select":
			var options = make([]huh.Option[string], len(variable.Options))
			for i, option := range variable.Options {
				options[i] = huh.NewOption(option, option)
			}
			input = huh.NewSelect[string]().
				Title(variable.Name).
				Options(options...).
				Value(&inputValues[i])
		default:
			return "", fmt.Errorf("unknown variable type: %s", variable.Type)
		}
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
