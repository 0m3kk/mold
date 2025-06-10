package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// LoadDataFile reads a JSON or YAML file from the given path and unmarshals it
// into a map that can be used for template rendering.
func LoadDataFile(path string) (map[string]any, error) {
	// Read the file content.
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read data file '%s': %w", path, err)
	}

	data := make(map[string]any)

	// Determine the file type by extension and unmarshal accordingly.
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json":
		if err = json.Unmarshal(content, &data); err != nil {
			return nil, fmt.Errorf("failed to parse JSON file '%s': %w", path, err)
		}
	case ".yaml", ".yml":
		if err = yaml.Unmarshal(content, &data); err != nil {
			return nil, fmt.Errorf("failed to parse YAML file '%s': %w", path, err)
		}
	default:
		return nil, fmt.Errorf("unsupported data file format: '%s'. Please use .json, .yaml, or .yml", ext)
	}

	return data, nil
}
