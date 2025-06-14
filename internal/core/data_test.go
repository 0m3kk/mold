package core

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDataFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	defer os.RemoveAll(tempDir)

	t.Run("load valid JSON file", func(t *testing.T) {
		// Create JSON test file
		jsonData := map[string]any{
			"name":    "test",
			"version": 1.0,
			"enabled": true,
		}
		jsonContent, err := json.Marshal(jsonData)
		if err != nil {
			t.Fatalf("Failed to marshal JSON: %v", err)
		}

		jsonPath := filepath.Join(tempDir, "test.json")
		err = os.WriteFile(jsonPath, jsonContent, 0644)
		if err != nil {
			t.Fatalf("Failed to write JSON file: %v", err)
		}

		// Load and verify
		result, err := LoadDataFile(jsonPath)
		if err != nil {
			t.Fatalf("LoadDataFile failed: %v", err)
		}

		if result["name"] != "test" {
			t.Errorf("Expected name 'test', got %v", result["name"])
		}
		if result["version"] != 1.0 {
			t.Errorf("Expected version 1.0, got %v", result["version"])
		}
		if result["enabled"] != true {
			t.Errorf("Expected enabled true, got %v", result["enabled"])
		}
	})

	t.Run("load valid YAML file with .yaml extension", func(t *testing.T) {
		yamlContent := `
name: test
version: 2
enabled: false
nested:
  key: value
`
		yamlPath := filepath.Join(tempDir, "test.yaml")
		err := os.WriteFile(yamlPath, []byte(yamlContent), 0644)
		if err != nil {
			t.Fatalf("Failed to write YAML file: %v", err)
		}

		result, err := LoadDataFile(yamlPath)
		if err != nil {
			t.Fatalf("LoadDataFile failed: %v", err)
		}

		if result["name"] != "test" {
			t.Errorf("Expected name 'test', got %v", result["name"])
		}
		if result["version"] != 2 {
			t.Errorf("Expected version 2, got %v", result["version"])
		}
		if result["enabled"] != false {
			t.Errorf("Expected enabled false, got %v", result["enabled"])
		}

		nested, ok := result["nested"].(map[string]any)
		if !ok {
			t.Errorf("Expected nested to be map[string]any, got %T", result["nested"])
		} else if nested["key"] != "value" {
			t.Errorf("Expected nested.key 'value', got %v", nested["key"])
		}
	})

	t.Run("load valid YAML file with .yml extension", func(t *testing.T) {
		yamlContent := `
name: yml_test
version: 3
`
		ymlPath := filepath.Join(tempDir, "test.yml")
		err := os.WriteFile(ymlPath, []byte(yamlContent), 0644)
		if err != nil {
			t.Fatalf("Failed to write YML file: %v", err)
		}

		result, err := LoadDataFile(ymlPath)
		if err != nil {
			t.Fatalf("LoadDataFile failed: %v", err)
		}

		if result["name"] != "yml_test" {
			t.Errorf("Expected name 'yml_test', got %v", result["name"])
		}
		if result["version"] != 3 {
			t.Errorf("Expected version 3, got %v", result["version"])
		}
	})

	t.Run("file does not exist", func(t *testing.T) {
		nonExistentPath := filepath.Join(tempDir, "nonexistent.json")

		_, err := LoadDataFile(nonExistentPath)
		if err == nil {
			t.Error("Expected error when file does not exist")
		}
		if !errors.Is(err, fs.ErrNotExist) {
			t.Errorf("Expected file not found error, got: %v", err)
		}
	})

	t.Run("unsupported file extension", func(t *testing.T) {
		txtPath := filepath.Join(tempDir, "test.txt")
		err := os.WriteFile(txtPath, []byte("some content"), 0644)
		if err != nil {
			t.Fatalf("Failed to write TXT file: %v", err)
		}

		_, err = LoadDataFile(txtPath)
		if err == nil {
			t.Error("Expected error for unsupported file extension")
		}

		expectedMsg := "unsupported data file format"
		if !contains(err.Error(), expectedMsg) {
			t.Errorf("Expected error message to contain %q, got: %v", expectedMsg, err.Error())
		}
	})

	t.Run("invalid JSON content", func(t *testing.T) {
		invalidJSONPath := filepath.Join(tempDir, "invalid.json")
		err := os.WriteFile(invalidJSONPath, []byte("{invalid json}"), 0644)
		if err != nil {
			t.Fatalf("Failed to write invalid JSON file: %v", err)
		}

		_, err = LoadDataFile(invalidJSONPath)
		if err == nil {
			t.Error("Expected error for invalid JSON content")
		}

		expectedMsg := "failed to parse JSON file"
		if !contains(err.Error(), expectedMsg) {
			t.Errorf("Expected error message to contain %q, got: %v", expectedMsg, err.Error())
		}
	})

	t.Run("invalid YAML content", func(t *testing.T) {
		invalidYamlPath := filepath.Join(tempDir, "invalid.yaml")
		err := os.WriteFile(invalidYamlPath, []byte("invalid:\n  yaml:\n    - content\n  - missing_indent"), 0644)
		if err != nil {
			t.Fatalf("Failed to write invalid YAML file: %v", err)
		}

		_, err = LoadDataFile(invalidYamlPath)
		if err == nil {
			t.Error("Expected error for invalid YAML content")
		}

		expectedMsg := "failed to parse YAML file"
		if !contains(err.Error(), expectedMsg) {
			t.Errorf("Expected error message to contain %q, got: %v", expectedMsg, err.Error())
		}
	})
}

// Helper function to check if a string contains a substring.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			containsAt(s, substr, 1)))
}

func containsAt(s, substr string, start int) bool {
	if start >= len(s) {
		return false
	}
	if start+len(substr) <= len(s) && s[start:start+len(substr)] == substr {
		return true
	}
	return containsAt(s, substr, start+1)
}
