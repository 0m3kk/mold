package cli

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/om3kk/mold/internal/core"
	"github.com/om3kk/mold/internal/utils"
)

//nolint:gochecknoglobals // this is cmd flag
var (
	outputDir string
	dataFile  string
)

// applyCmd represents the apply command, renamed from createCmd.
//
//nolint:gochecknoglobals // this is command definition
var applyCmd = &cobra.Command{
	Use:   "apply <template_name>",
	Short: "Applies a template to generate a project",
	Long: `Generates a new project structure based on a specified template.
It copies files from the template directory, processes files ending in '.tmpl'
by filling in placeholders, and saves the result to the output directory.`,
	Args: cobra.ExactArgs(1), // Ensures exactly one argument (template_name) is passed.
	RunE: func(_ *cobra.Command, args []string) error {
		templateName := args[0]
		// The path to the specific template is now built using the configurable templatesDir
		templatePath := filepath.Join(templatesDir, templateName)
		placeholders := make(map[string]bool) // Using a map to store unique placeholders

		// 1. Validate Template
		if _, err := os.Stat(templatePath); os.IsNotExist(err) {
			return fmt.Errorf("template '%s' not found in '%s' directory", templateName, templatesDir)
		}

		fmt.Printf("🚀 Applying template: %s\n", templateName)

		// 2. Create output directory if it doesn't exist.
		if err := os.MkdirAll(outputDir, 0750); err != nil {
			return fmt.Errorf("failed to create output directory '%s': %w", outputDir, err)
		}

		// 3. Walk the template directory to copy files and identify placeholders
		err := filepath.WalkDir(templatePath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			// Get the relative path to construct the destination path
			relPath, err := filepath.Rel(templatePath, path)
			if err != nil {
				return fmt.Errorf("failed to get relative path for '%s': %w", path, err)
			}
			destPath := filepath.Join(outputDir, relPath)

			if d.IsDir() {
				// Create corresponding directory in the destination
				return os.MkdirAll(destPath, d.Type().Perm())
			}

			// Check if it's a template file
			if strings.HasSuffix(d.Name(), ".tmpl") {
				// This is a template file, identify placeholders
				fmt.Printf("🔍 Identifying placeholders in: %s\n", relPath)
				vars, innerErr := core.IdentifyPlaceholders(path)
				if innerErr != nil {
					return fmt.Errorf("failed to parse template '%s': %w", path, innerErr)
				}
				for _, v := range vars {
					placeholders[v] = true
				}
				// For now, just copy the file. Rendering will happen later.
				// We'll rename it by removing the .tmpl extension in the real render step.
				// For this skeleton, we will just copy it.
				// destPath = strings.TrimSuffix(destPath, ".tmpl")
			}

			// Copy the file content
			return utils.CopyFile(path, destPath)
		})

		if err != nil {
			return fmt.Errorf("error during template processing: %w", err)
		}

		// 4. Gather Data (Placeholder for now)
		fmt.Println("\n-------------------------------------------")
		if len(placeholders) > 0 {
			fmt.Println("✨ Found the following placeholders:")
			for p := range placeholders {
				fmt.Printf("  - {{.%s}}\n", p)
			}
			fmt.Println("\n(Next step: Gather data for these placeholders and render templates)")
		} else {
			fmt.Println("✨ No '.tmpl' files with placeholders found. All files copied as-is.")
		}
		fmt.Println("-------------------------------------------")

		// 5. Render Files (To be implemented)

		// 6. Success Message
		fmt.Printf("\n✅ Successfully applied template to: %s\n", outputDir)
		return nil
	},
}

//nolint:gochecknoinits // The command 'init' is acceptable.
func init() {
	// Add flags to the 'apply' command.
	applyCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "Output directory for the new project")
	applyCmd.Flags().StringVarP(&dataFile, "data-file", "d", "", "Path to a JSON or YAML file with placeholder data")
}
