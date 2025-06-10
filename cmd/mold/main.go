package main

import (
	"os"

	"github.com/om3kk/mold/internal/cli"
)

// main is the entry point for the Mold CLI application.
// It executes the root command and handles any errors, exiting
// with a non-zero status code if an error occurs.
func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
