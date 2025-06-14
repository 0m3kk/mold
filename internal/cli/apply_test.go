package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestApplyCmd(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		setupFunc      func(t *testing.T) (tempDir, templateDir, dataFile, outputDir string, cleanup func())
		expectedError  string
		validateOutput func(t *testing.T, outputDir string)
	}{
		{
			name: "successful_apply_with_json_data",
			args: []string{"template"},
			setupFunc: func(t *testing.T) (string, string, string, string, func()) {
				tempDir := t.TempDir()
				templateDir := filepath.Join(tempDir, "template")
				dataFile := filepath.Join(tempDir, "data.json")
				outputDir := filepath.Join(tempDir, "output")

				// Create template directory structure
				require.NoError(t, os.MkdirAll(templateDir, 0755))
				require.NoError(t, os.MkdirAll(filepath.Join(templateDir, "{{.project_name}}"), 0755))

				// Create template files
				templateContent := "package {{.package_name}}\n\nfunc Hello() string {\n\treturn \"{{.greeting}}\"\n}"
				require.NoError(
					t,
					os.WriteFile(filepath.Join(templateDir, "main.go.tmpl"), []byte(templateContent), 0644),
				)

				// Create regular file to copy
				require.NoError(
					t,
					os.WriteFile(filepath.Join(templateDir, "README.md"), []byte("# Project README"), 0644),
				)

				// Create nested template file
				nestedTemplate := "Name: {{.project_name}}"
				require.NoError(
					t,
					os.WriteFile(
						filepath.Join(templateDir, "{{.project_name}}", "config.yaml.tmpl"),
						[]byte(nestedTemplate),
						0644,
					),
				)

				// Create data file
				data := map[string]any{
					"project_name": "myproject",
					"package_name": "main",
					"greeting":     "Hello, World!",
				}
				dataBytes, _ := json.Marshal(data)
				require.NoError(t, os.WriteFile(dataFile, dataBytes, 0644))

				cleanup := func() {
					// Cleanup is handled by t.TempDir()
				}

				return tempDir, templateDir, dataFile, outputDir, cleanup
			},
			validateOutput: func(t *testing.T, outputDir string) {
				// Check rendered template file
				mainContent, err := os.ReadFile(filepath.Join(outputDir, "main.go"))
				require.NoError(t, err)
				assert.Contains(t, string(mainContent), "package main")
				assert.Contains(t, string(mainContent), "Hello, World!")

				// Check copied file
				readmeContent, err := os.ReadFile(filepath.Join(outputDir, "README.md"))
				require.NoError(t, err)
				assert.Equal(t, "# Project README", string(readmeContent))

				// Check nested rendered file
				configContent, err := os.ReadFile(filepath.Join(outputDir, "myproject", "config.yaml"))
				require.NoError(t, err)
				assert.Contains(t, string(configContent), "Name: myproject")
			},
		},
		{
			name: "successful_apply_with_yaml_data",
			args: []string{"template"},
			setupFunc: func(t *testing.T) (string, string, string, string, func()) {
				tempDir := t.TempDir()
				templateDir := filepath.Join(tempDir, "template")
				dataFile := filepath.Join(tempDir, "data.yaml")
				outputDir := filepath.Join(tempDir, "output")

				// Create template directory
				require.NoError(t, os.MkdirAll(templateDir, 0755))

				// Create template file
				templateContent := "version: {{.version}}"
				require.NoError(
					t,
					os.WriteFile(filepath.Join(templateDir, "config.yaml.tmpl"), []byte(templateContent), 0644),
				)

				// Create YAML data file
				data := map[string]any{
					"version": "1.0.0",
				}
				dataBytes, _ := yaml.Marshal(data)
				require.NoError(t, os.WriteFile(dataFile, dataBytes, 0644))

				return tempDir, templateDir, dataFile, outputDir, func() {}
			},
			validateOutput: func(t *testing.T, outputDir string) {
				configContent, err := os.ReadFile(filepath.Join(outputDir, "config.yaml"))
				require.NoError(t, err)
				assert.Contains(t, string(configContent), "version: 1.0.0")
			},
		},
		{
			name:          "missing_data_file_flag",
			args:          []string{"template"},
			expectedError: "the --data-file flag is required for rendering templates",
			setupFunc: func(t *testing.T) (string, string, string, string, func()) {
				tempDir := t.TempDir()
				templateDir := filepath.Join(tempDir, "template")
				require.NoError(t, os.MkdirAll(templateDir, 0755))
				return tempDir, templateDir, "", "", func() {}
			},
		},
		{
			name:          "missing_data_file_with_yaml_hint",
			args:          []string{"template"},
			expectedError: "the --data-file flag is required for rendering templates.\nHint: Found a",
			setupFunc: func(t *testing.T) (string, string, string, string, func()) {
				tempDir := t.TempDir()
				templateDir := filepath.Join(tempDir, "template")
				require.NoError(t, os.MkdirAll(templateDir, 0755))

				// Create example tmpl.yaml file
				require.NoError(t, os.WriteFile(filepath.Join(templateDir, "tmpl.yaml"), []byte("example: data"), 0644))

				return tempDir, templateDir, "", "", func() {}
			},
		},
		{
			name:          "missing_data_file_with_json_hint",
			args:          []string{"template"},
			expectedError: "the --data-file flag is required for rendering templates.\nHint: Found a",
			setupFunc: func(t *testing.T) (string, string, string, string, func()) {
				tempDir := t.TempDir()
				templateDir := filepath.Join(tempDir, "template")
				require.NoError(t, os.MkdirAll(templateDir, 0755))

				// Create example tmpl.json file (no tmpl.yaml so JSON will be found)
				require.NoError(
					t,
					os.WriteFile(filepath.Join(templateDir, "tmpl.json"), []byte(`{"example": "data"}`), 0644),
				)

				return tempDir, templateDir, "", "", func() {}
			},
		},
		{
			name:          "template_path_not_found",
			args:          []string{"nonexistent"},
			expectedError: "template path 'nonexistent' not found",
			setupFunc: func(t *testing.T) (string, string, string, string, func()) {
				tempDir := t.TempDir()
				dataFile := filepath.Join(tempDir, "data.json")
				data := map[string]any{"key": "value"}
				dataBytes, _ := json.Marshal(data)
				require.NoError(t, os.WriteFile(dataFile, dataBytes, 0644))
				return tempDir, "", dataFile, "", func() {}
			},
		},
		{
			name:          "invalid_data_file",
			args:          []string{"template"},
			expectedError: "failed to parse JSON file",
			setupFunc: func(t *testing.T) (string, string, string, string, func()) {
				tempDir := t.TempDir()
				templateDir := filepath.Join(tempDir, "template")
				dataFile := filepath.Join(tempDir, "invalid.json")

				require.NoError(t, os.MkdirAll(templateDir, 0755))
				require.NoError(t, os.WriteFile(dataFile, []byte("invalid json"), 0644))

				return tempDir, templateDir, dataFile, "", func() {}
			},
		},
		{
			name: "skip_tmpl_files",
			args: []string{"template"},
			setupFunc: func(t *testing.T) (string, string, string, string, func()) {
				tempDir := t.TempDir()
				templateDir := filepath.Join(tempDir, "template")
				dataFile := filepath.Join(tempDir, "data.json")
				outputDir := filepath.Join(tempDir, "output")

				require.NoError(t, os.MkdirAll(templateDir, 0755))

				// Create tmpl.json and tmpl.yaml files that should be skipped
				require.NoError(
					t,
					os.WriteFile(filepath.Join(templateDir, "tmpl.json"), []byte(`{"skip": "me"}`), 0644),
				)
				require.NoError(t, os.WriteFile(filepath.Join(templateDir, "tmpl.yaml"), []byte("skip: me"), 0644))

				// Create a regular file
				require.NoError(t, os.WriteFile(filepath.Join(templateDir, "regular.txt"), []byte("keep me"), 0644))

				// Create data file
				data := map[string]any{"key": "value"}
				dataBytes, _ := json.Marshal(data)
				require.NoError(t, os.WriteFile(dataFile, dataBytes, 0644))

				return tempDir, templateDir, dataFile, outputDir, func() {}
			},
			validateOutput: func(t *testing.T, outputDir string) {
				// tmpl.json and tmpl.yaml should not exist in output
				_, err := os.Stat(filepath.Join(outputDir, "tmpl.json"))
				assert.True(t, os.IsNotExist(err))

				_, err = os.Stat(filepath.Join(outputDir, "tmpl.yaml"))
				assert.True(t, os.IsNotExist(err))

				// Regular file should exist
				content, err := os.ReadFile(filepath.Join(outputDir, "regular.txt"))
				require.NoError(t, err)
				assert.Equal(t, "keep me", string(content))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global variables
			outputDir = "."
			dataFile = ""

			tempDir, templateDir, dataFileVar, outputDirVar, cleanup := tt.setupFunc(t)
			defer cleanup()

			// Set up command with flags
			cmd := &cobra.Command{}
			cmd.AddCommand(applyCmd)

			// Build command line args
			args := []string{"apply"}
			if templateDir != "" {
				// Use relative path from tempDir for the template
				if strings.HasPrefix(templateDir, tempDir) {
					relPath, _ := filepath.Rel(tempDir, templateDir)
					args = append(args, relPath)
				} else {
					args = append(args, templateDir)
				}
			} else {
				args = append(args, tt.args...)
			}

			if dataFileVar != "" {
				args = append(args, "--data-file", dataFileVar)
			}
			if outputDirVar != "" {
				args = append(args, "--output", outputDirVar)
			}

			// Change to temp directory for relative path resolution
			originalWd, _ := os.Getwd()
			if tempDir != "" {
				t.Chdir(tempDir)
				defer func() { t.Chdir(originalWd) }()
			}

			cmd.SetArgs(args)
			err := cmd.Execute()

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
				if tt.validateOutput != nil {
					tt.validateOutput(t, outputDirVar)
				}
			}
		})
	}
}

func TestApplyCmdFlags(t *testing.T) {
	// Test that flags are properly registered
	assert.True(t, applyCmd.Flags().HasFlags())

	outputFlag := applyCmd.Flags().Lookup("output")
	require.NotNil(t, outputFlag)
	assert.Equal(t, "o", outputFlag.Shorthand)
	assert.Equal(t, ".", outputFlag.DefValue)

	dataFileFlag := applyCmd.Flags().Lookup("data-file")
	require.NotNil(t, dataFileFlag)
	assert.Equal(t, "d", dataFileFlag.Shorthand)
	assert.Empty(t, dataFileFlag.DefValue)
}

func TestApplyCmdBasicProperties(t *testing.T) {
	assert.Equal(t, "apply <template_path>", applyCmd.Use)
	assert.Equal(t, "Applies a template directory to generate a project using a data file", applyCmd.Short)
	assert.Contains(t, applyCmd.Long, "Generates a project structure from a template directory")
}

func TestApplyCmdArgumentValidation(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name:        "no_arguments",
			args:        []string{},
			expectError: true,
		},
		{
			name:        "too_many_arguments",
			args:        []string{"arg1", "arg2"},
			expectError: true,
		},
		{
			name:        "exactly_one_argument",
			args:        []string{"template"},
			expectError: false, // Args validation passes, but command will fail later
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global variables
			outputDir = "."
			dataFile = ""

			cmd := &cobra.Command{}
			cmd.AddCommand(applyCmd)

			args := append([]string{"apply"}, tt.args...)
			cmd.SetArgs(args)

			err := cmd.Execute()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				// Command will still error due to missing data file, but not due to args validation
				require.Error(t, err)
				assert.Contains(t, err.Error(), "data-file flag is required")
			}
		})
	}
}

func TestApplyCmdErrorHandling(t *testing.T) {
	t.Run("output_directory_creation_failure", func(t *testing.T) {
		tempDir := t.TempDir()
		templateDir := filepath.Join(tempDir, "template")
		dataFileVar := filepath.Join(tempDir, "data.json")

		require.NoError(t, os.MkdirAll(templateDir, 0755))
		require.NoError(t, os.WriteFile(filepath.Join(templateDir, "file.txt"), []byte("content"), 0644))

		data := map[string]any{"key": "value"}
		dataBytes, _ := json.Marshal(data)
		require.NoError(t, os.WriteFile(dataFileVar, dataBytes, 0644))

		// Create a file where we want to create the output directory
		invalidOutputDir := filepath.Join(tempDir, "existing_file")
		require.NoError(t, os.WriteFile(invalidOutputDir, []byte("block"), 0644))

		// Reset global variables
		outputDir = "."
		dataFile = ""

		cmd := &cobra.Command{}
		cmd.AddCommand(applyCmd)

		originalWd, _ := os.Getwd()
		// require.NoError(t, os.Chdir(tempDir))
		t.Chdir(tempDir)
		defer func() { t.Chdir(originalWd) }()

		relTemplatePath, _ := filepath.Rel(tempDir, templateDir)
		args := []string{"apply", relTemplatePath, "--data-file", dataFileVar, "--output", invalidOutputDir}
		cmd.SetArgs(args)

		err := cmd.Execute()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create output directory")
	})
}

// TestInit verifies the init function runs without panicking.
func TestInit(t *testing.T) {
	// The init function should have already run when the package was loaded
	// We just verify the command has the expected flags
	assert.NotNil(t, applyCmd.Flags().Lookup("output"))
	assert.NotNil(t, applyCmd.Flags().Lookup("data-file"))
}
