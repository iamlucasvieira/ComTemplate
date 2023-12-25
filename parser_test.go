package main

import (
	"os"
	"testing"
)

var mockData = `
- name: First template
  text: |
    %{title}

    %{body}

  variables:
    - name: title
    - name: body

- name: Second template
  text: |
    [%{identifier}]%{title}

    %{body}

  variables:
    - name: identifier
    - name: title
    - name: body
`

func TestParse(t *testing.T) {
	templates, err := Parse(mockData)

	if err != nil {
		t.Errorf("error parsing yaml: %v", err)
	}

	t.Run("should have two templates", func(t *testing.T) {
		if len(templates) != 2 {
			t.Errorf("expected two templates, got %d", len(templates))
		}
	})

	t.Run("should have a name", func(t *testing.T) {
		if templates[0].Name != "First template" {
			t.Errorf("expected name to be 'First template', got '%s'", templates[0].Name)
		}
	})

	t.Run("should have a text", func(t *testing.T) {
		if templates[0].Text != "%{title}\n\n%{body}\n" {
			t.Errorf("expected text to be '%%{title}\n\n%%{body}\n', got '%s'", templates[0].Text)
		}
	})

	t.Run("should have two variables", func(t *testing.T) {
		if len(templates[0].Variables) != 2 {
			t.Errorf("expected two variables, got %d", len(templates[0].Variables))
		}
	})

}

func TestOpen(t *testing.T) {
	// Create a temporary file
	f, err := os.CreateTemp("", "test")
	defer os.Remove(f.Name())

	if err != nil {
		t.Errorf("error creating temp file: %v", err)
	}

	// Write mock data to file
	_, err = f.Write([]byte(mockData))
	if err != nil {
		t.Fatalf("error writing to temp file: %v", err)
	}

	// Close the file before reopening in open function
	err = f.Close()
	if err != nil {
		t.Fatalf("error closing file: %v", err)
	}

	// Open the file
	data, err := open(f.Name())

	if err != nil {
		t.Errorf("error opening file: %v", err)
	}

	if data != mockData {
		t.Errorf("expected data to be '%s', got '%s'", mockData, data)
	}

	if err != nil {
		t.Errorf("error closing file: %v", err)
	}
}
