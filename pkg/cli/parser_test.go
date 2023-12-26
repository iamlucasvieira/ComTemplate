package cli

import (
	"os"
	"path/filepath"
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
	templates, err := parse(mockData)

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

func TestRead(t *testing.T) {
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
	templates, err := read(f.Name())

	if err != nil {
		t.Errorf("error reading file: %v", err)
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

func TestReadDefault(t *testing.T) {
	// Create a temporaty directory
	tempDir, err := os.MkdirTemp("", "testDir")
	defer os.RemoveAll(tempDir)

	if err != nil {
		t.Errorf("error creating temp directory: %v", err)
	}

	tempFilePath := filepath.Join(tempDir, "comtemplate.yml")
	f, err := os.Create(tempFilePath)

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

	// Change the working directory to the temporary directory
	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("error changing working directory: %v", err)
	}

	templates, err := ReadDefault()

	if err != nil {
		t.Errorf("error reading file: %v", err)
	}

	t.Run("should have two templates", func(t *testing.T) {
		if len(templates) != 2 {
			t.Errorf("expected two templates, got %d", len(templates))
		}
	})

}

func TestCreateDefault(t *testing.T) {
	// Create a temporaty directory
	tempDir, err := os.MkdirTemp("", "testDir")

	if err != nil {
		t.Errorf("error creating temp directory: %v", err)
	}

	defer os.RemoveAll(tempDir)

	// Change the working directory to the temporary directory
	err = os.Chdir(tempDir)

	if err != nil {
		t.Fatalf("error changing working directory: %v", err)
	}

	err = CreateDefault()

	if err != nil {
		t.Errorf("error creating default file: %v", err)
	}

	t.Run("should have created a file", func(t *testing.T) {
		_, err := os.Stat("comtemplate.yml")

		if err != nil {
			t.Errorf("expected file to exist, got %v", err)
		}
	})

	// Check content of file
	data, err := read("comtemplate.yml")

	if err != nil {
		t.Errorf("error opening file: %v", err)
	}

	t.Run("should have default content with 1 template", func(t *testing.T) {
		if len(data) != 1 {
			t.Errorf("expected one template, got %d", len(data))
		}
	})

}
