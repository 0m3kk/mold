package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// initCmd represents the init command.
//
//nolint:gochecknoglobals // this is command definition
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a directory to store templates",
	Long: `Creates a directory to store your template sets.
By default, this is the 'templates' directory, but this can be changed
globally using the --dir flag.`,
	Run: func(_ *cobra.Command, _ []string) {
		// Check if the directory already exists. It now uses the value from the --dir flag.
		if _, err := os.Stat(templatesDir); !os.IsNotExist(err) {
			fmt.Printf("Directory '%s' already exists. Nothing to do.\n", templatesDir)
			return
		}

		// Create the directory using the path from the --dir flag.
		err := os.Mkdir(templatesDir, 0750)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating directory '%s': %v\n", templatesDir, err)
			return
		}

		// Create a placeholder file to ensure the directory is added to git.
		placeholderPath := filepath.Join(templatesDir, ".gitkeep")
		file, err := os.Create(placeholderPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not create .gitkeep file: %v\n", err)
		} else {
			err = file.Close()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: close file error: %v\n", err)
			}
		}

		fmt.Printf("✅ Successfully created directory: %s\n", templatesDir)
		fmt.Println("You can now add your project templates inside this directory.")
	},
}
