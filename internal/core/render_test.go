package core

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestRenderTemplateFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	defer os.RemoveAll(tempDir)

	t.Run("successful template rendering", func(t *testing.T) {
		// Create template file
		templateContent := `Hello {{.name}}!
Your age is {{.age}}.
Snake case: {{snake .camelCase}}
Upper snake: {{usnake .camelCase}}
Camel case: {{camel .snake_case}}
Lower camel: {{lcamel .snake_case}}`

		templatePath := filepath.Join(tempDir, "template.txt")
		err := os.WriteFile(templatePath, []byte(templateContent), 0755)
		if err != nil {
			t.Fatalf("Failed to create template file: %v", err)
		}

		// Prepare data
		data := map[string]any{
			"name":       "John",
			"age":        30,
			"camelCase":  "someVariableName",
			"snake_case": "some_variable_name",
		}

		// Render template
		destPath := filepath.Join(tempDir, "output.txt")
		err = RenderTemplateFile(templatePath, destPath, data)
		if err != nil {
			t.Fatalf("RenderTemplateFile failed: %v", err)
		}

		// Verify output
		output, err := os.ReadFile(destPath)
		if err != nil {
			t.Fatalf("Failed to read output file: %v", err)
		}

		expectedOutput := `Hello John!
Your age is 30.
Snake case: some_variable_name
Upper snake: SOME_VARIABLE_NAME
Camel case: SomeVariableName
Lower camel: someVariableName`

		if string(output) != expectedOutput {
			t.Errorf("Output mismatch:\nGot:\n%s\nWant:\n%s", string(output), expectedOutput)
		}

		// Verify permissions are preserved
		templateInfo, err := os.Stat(templatePath)
		if err != nil {
			t.Fatalf("Failed to stat template file: %v", err)
		}
		outputInfo, err := os.Stat(destPath)
		if err != nil {
			t.Fatalf("Failed to stat output file: %v", err)
		}
		if templateInfo.Mode() != outputInfo.Mode() {
			t.Errorf("Permission mismatch: got %v, want %v", outputInfo.Mode(), templateInfo.Mode())
		}
	})

	t.Run("template file does not exist", func(t *testing.T) {
		nonExistentPath := filepath.Join(tempDir, "nonexistent.txt")
		destPath := filepath.Join(tempDir, "output2.txt")
		data := map[string]any{"key": "value"}

		err := RenderTemplateFile(nonExistentPath, destPath, data)
		if err == nil {
			t.Error("Expected error when template file does not exist")
		}
		if !errors.Is(err, fs.ErrNotExist) {
			t.Errorf("Expected file not found error, got: %v", err)
		}
	})

	t.Run("invalid template syntax", func(t *testing.T) {
		invalidTemplate := `Hello {{.name}!` // Missing closing brace
		templatePath := filepath.Join(tempDir, "invalid.txt")
		err := os.WriteFile(templatePath, []byte(invalidTemplate), 0644)
		if err != nil {
			t.Fatalf("Failed to create invalid template file: %v", err)
		}

		destPath := filepath.Join(tempDir, "output3.txt")
		data := map[string]any{"name": "John"}

		err = RenderTemplateFile(templatePath, destPath, data)
		if err == nil {
			t.Error("Expected error for invalid template syntax")
		}

		expectedMsg := "could not parse template"
		if !contains(err.Error(), expectedMsg) {
			t.Errorf("Expected error message to contain %q, got: %v", expectedMsg, err.Error())
		}
	})

	t.Run("cannot create destination file", func(t *testing.T) {
		templateContent := `Hello {{.name}}!`
		templatePath := filepath.Join(tempDir, "template4.txt")
		err := os.WriteFile(templatePath, []byte(templateContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create template file: %v", err)
		}

		// Try to create destination in non-existent directory
		invalidDestPath := filepath.Join(tempDir, "nonexistent_dir", "output.txt")
		data := map[string]any{"name": "John"}

		err = RenderTemplateFile(templatePath, invalidDestPath, data)
		if err == nil {
			t.Error("Expected error when destination directory does not exist")
		}

		expectedMsg := "failed to create destination file"
		if !contains(err.Error(), expectedMsg) {
			t.Errorf("Expected error message to contain %q, got: %v", expectedMsg, err.Error())
		}
	})

	t.Run("template execution error", func(t *testing.T) {
		// Template that references non-existent field
		templateContent := `Hello {{.name.NonExistent}}!`
		templatePath := filepath.Join(tempDir, "template5.txt")
		err := os.WriteFile(templatePath, []byte(templateContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create template file: %v", err)
		}

		destPath := filepath.Join(tempDir, "output5.txt")
		data := map[string]any{"name": "John"}

		err = RenderTemplateFile(templatePath, destPath, data)
		if err == nil {
			t.Fatal("Expected error during template execution")
		}

		expectedMsg := "failed to render template"
		if !contains(err.Error(), expectedMsg) {
			t.Errorf("Expected error message to contain %q, got: %v", expectedMsg, err.Error())
		}
	})
}

func TestReplacePlaceholdersInPath(t *testing.T) {
	t.Run("successful path replacement", func(t *testing.T) {
		path := "/app/{{.service}}/{{snake .serviceName}}/config"
		data := map[string]any{
			"service":     "myapp",
			"serviceName": "MyAwesomeService",
		}

		result, err := ReplacePlaceholdersInPath(path, data)
		if err != nil {
			t.Fatalf("ReplacePlaceholdersInPath failed: %v", err)
		}

		expected := "/app/myapp/my_awesome_service/config"
		if result != expected {
			t.Errorf("Path replacement failed: got %q, want %q", result, expected)
		}
	})

	t.Run("path with all helper functions", func(t *testing.T) {
		path := "{{snake .name}}/{{usnake .name}}/{{camel .name}}/{{lcamel .name}}"
		data := map[string]any{
			"name": "someVariableName",
		}

		result, err := ReplacePlaceholdersInPath(path, data)
		if err != nil {
			t.Fatalf("ReplacePlaceholdersInPath failed: %v", err)
		}

		expected := "some_variable_name/SOME_VARIABLE_NAME/SomeVariableName/someVariableName"
		if result != expected {
			t.Errorf("Path replacement failed: got %q, want %q", result, expected)
		}
	})

	t.Run("invalid template syntax in path", func(t *testing.T) {
		path := "/app/{{.service}/config" // Missing closing brace
		data := map[string]any{
			"service": "myapp",
		}

		_, err := ReplacePlaceholdersInPath(path, data)
		if err == nil {
			t.Error("Expected error for invalid template syntax in path")
		}
	})

	t.Run("template execution error in path", func(t *testing.T) {
		path := "/app/{{.service.nonExistentField}}/config"
		data := map[string]any{
			"service": "myapp",
		}

		_, err := ReplacePlaceholdersInPath(path, data)
		if err == nil {
			t.Error("Expected error during path template execution")
		}
	})

	t.Run("empty path", func(t *testing.T) {
		path := ""
		data := map[string]any{}

		result, err := ReplacePlaceholdersInPath(path, data)
		if err != nil {
			t.Fatalf("ReplacePlaceholdersInPath failed for empty path: %v", err)
		}

		if result != "" {
			t.Errorf("Expected empty result for empty path, got %q", result)
		}
	})

	t.Run("path without placeholders", func(t *testing.T) {
		path := "/app/static/config"
		data := map[string]any{
			"service": "myapp",
		}

		result, err := ReplacePlaceholdersInPath(path, data)
		if err != nil {
			t.Fatalf("ReplacePlaceholdersInPath failed: %v", err)
		}

		if result != path {
			t.Errorf("Expected unchanged path %q, got %q", path, result)
		}
	})
}
