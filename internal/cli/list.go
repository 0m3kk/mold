package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
//
//nolint:gochecknoglobals // this is command definition
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all available templates",
	Long:  `Scans the templates directory and lists all available template sets (subdirectories).`,
	Run: func(_ *cobra.Command, _ []string) {
		// Check if the templates directory (specified by --dir flag) exists.
		if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
			fmt.Printf("Directory '%s' not found.\n", templatesDir)
			fmt.Printf("Run 'mold init --dir %s' to create it.\n", templatesDir)
			return
		}

		// Read the contents of the templates directory.
		entries, err := os.ReadDir(templatesDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading directory '%s': %v\n", templatesDir, err)
			return
		}

		var templates []string
		for _, entry := range entries {
			// We are only interested in directories.
			if entry.IsDir() {
				templates = append(templates, entry.Name())
			}
		}

		if len(templates) == 0 {
			fmt.Printf("No templates found in the '%s' directory.\n", templatesDir)
			fmt.Printf("Add a new directory inside '%s' to create a template set.\n", templatesDir)
			return
		}

		fmt.Println("Available templates:")
		for _, t := range templates {
			fmt.Printf("  - %s\n", t)
		}
	},
}
