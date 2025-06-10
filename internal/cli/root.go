package cli

import (
	"github.com/spf13/cobra"
)

// templatesDir will be populated by the persistent flag.
// It holds the path to the directory containing the templates.
//
//nolint:gochecknoglobals // this is command flag
var templatesDir string

// rootCmd represents the base command when called without any subcommands.
//
//nolint:gochecknoglobals // this is command definition
var rootCmd = &cobra.Command{
	Use:   "mold",
	Short: "A CLI tool for scaffolding projects from templates.",
	Long: `Mold is a powerful and simple command-line tool that helps you
generate project structures, files, and configurations from predefined templates.

Use 'mold init' to create a templates directory, 'mold list' to see
available templates, and 'mold create' to generate a new project.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

// init function is called by Go when the package is initialized.
//
//nolint:gochecknoinits // The command 'init' is acceptable.
func init() {
	// Add a persistent flag to the root command for the templates directory.
	// This flag will be available to all subcommands that descend from rootCmd.
	rootCmd.PersistentFlags().
		StringVarP(&templatesDir, "dir", "t", "templates", "Directory to store and read templates from")

	// Add subcommands to the root command.
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(applyCmd)
}
