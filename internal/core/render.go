package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/stoewer/go-strcase"
)

//nolint:gochecknoglobals // helper function use when render templates
var helperFunc = template.FuncMap{
	"snake":  strcase.SnakeCase,
	"usnake": strcase.UpperSnakeCase,
	"camel":  strcase.UpperCamelCase,
	"lcamel": strcase.LowerCamelCase,
}

// RenderTemplateFile reads a template file, executes it with the provided data,
// and writes the output to the destination path.
func RenderTemplateFile(templatePath, destPath string, data map[string]any) error {
	// Read the template content.
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("could not read template file '%s': %w", templatePath, err)
	}

	// Create a new template, parse the content, and execute it.
	tmpl, err := template.New(filepath.Base(templatePath)).Funcs(helperFunc).Parse(string(content))
	if err != nil {
		return fmt.Errorf("could not parse template '%s': %w", templatePath, err)
	}

	// Create the destination file.
	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file '%s': %w", destPath, err)
	}
	defer destFile.Close()

	// Execute the template and write the output directly to the file.
	if err = tmpl.Execute(destFile, data); err != nil {
		return fmt.Errorf("failed to render template '%s': %w", templatePath, err)
	}

	// Preserve file permissions from the original template
	sourceInfo, err := os.Stat(templatePath)
	if err != nil {
		return fmt.Errorf("failed to stat source file '%s': %w", templatePath, err)
	}
	return os.Chmod(destPath, sourceInfo.Mode())
}

// ReplacePlaceholdersInPath replace placeholders in directory names.
func ReplacePlaceholdersInPath(path string, data map[string]any) (string, error) {
	tmpl, err := template.New("path").Funcs(helperFunc).Parse(path)
	if err != nil {
		return "", err
	}
	var result strings.Builder
	if err = tmpl.Execute(&result, data); err != nil {
		return "", err
	}
	return result.String(), nil
}
