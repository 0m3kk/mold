package cli

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/om3kk/mold/internal/core"
	"github.com/om3kk/mold/internal/utils"

	"github.com/spf13/cobra"
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
	Use:   "apply <template_path>",
	Short: "Applies a template directory to generate a project using a data file",
	Long: `Generates a project structure from a template directory.
This command requires a data file (JSON or YAML) to render templates.
It processes files ending in '.tmpl' by filling in placeholders from the data file
and saves the result to the output directory. All other files are copied as-is.`,
	Args: cobra.ExactArgs(1), // Requires exactly one argument: the path to the template.
	RunE: func(_ *cobra.Command, args []string) error {
		var err error
		templatePath := args[0]

		// 1. Validate the --data-file flag. It is now mandatory.
		if dataFile == "" {
			// Check if an example data file exists to provide a helpful hint.
			exampleHint := ""
			exampleYAML := filepath.Join(templatePath, "template.yaml")
			exampleJSON := filepath.Join(templatePath, "template.json")

			if _, err = os.Stat(exampleYAML); err == nil {
				exampleHint = fmt.Sprintf(
					"\nHint: Found a '%s' file. You can copy and edit it for your data.",
					exampleYAML,
				)
			} else if _, err = os.Stat(exampleJSON); err == nil {
				exampleHint = fmt.Sprintf("\nHint: Found a '%s' file. You can copy and edit it for your data.", exampleJSON)
			}
			return fmt.Errorf("the --data-file flag is required for rendering templates.%s", exampleHint)
		}

		// 2. Validate Template Path
		if _, err = os.Stat(templatePath); os.IsNotExist(err) {
			return fmt.Errorf("template path '%s' not found", templatePath)
		}
		fmt.Printf("🚀 Applying template from: %s\n", templatePath)

		// 3. Load data from the specified file.
		fmt.Printf("📖 Loading data from: %s\n", dataFile)
		var data map[string]any
		data, err = core.LoadDataFile(dataFile)
		if err != nil {
			return err // Error is already descriptive.
		}

		// 4. Create output directory if it doesn't exist.
		if err = os.MkdirAll(outputDir, 0750); err != nil {
			return fmt.Errorf("failed to create output directory '%s': %w", outputDir, err)
		}

		// 5. Walk the template directory to render/copy files.
		err = filepath.WalkDir(templatePath, func(path string, d fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}

			// Determine the destination path for the file or directory.
			relPath, innerErr := filepath.Rel(templatePath, path)
			if innerErr != nil {
				return fmt.Errorf("failed to get relative path for '%s': %w", path, innerErr)
			}
			destPath := filepath.Join(outputDir, relPath)

			if d.IsDir() {
				// Create the corresponding directory in the destination.
				return os.MkdirAll(destPath, d.Type().Perm())
			}

			// Decide whether to render or copy the file.
			if strings.HasSuffix(d.Name(), ".tmpl") {
				// This is a template file that needs to be rendered.
				finalDestPath := strings.TrimSuffix(destPath, ".tmpl")
				fmt.Printf("✨ Rendering: %s -> %s\n", relPath, strings.TrimSuffix(relPath, ".tmpl"))
				return core.RenderTemplateFile(path, finalDestPath, data)
			}

			// This is a regular file, so just copy it.
			fmt.Printf("📄 Copying: %s\n", relPath)
			return utils.CopyFile(path, destPath)
		})

		if err != nil {
			return fmt.Errorf("error during template processing: %w", err)
		}

		// 6. Success Message
		fmt.Printf("\n✅ Successfully applied template to: %s\n", outputDir)
		return nil
	},
}

//nolint:gochecknoinits // The command 'init' is acceptable.
func init() {
	// Add flags to the 'apply' command.
	applyCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "Output directory for the new project")
	applyCmd.Flags().
		StringVarP(&dataFile, "data-file", "d", "", "Path to a JSON or YAML file with placeholder data (required)")
}
