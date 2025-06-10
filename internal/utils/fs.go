package utils

import (
	"fmt"
	"io"
	"os"
)

// CopyFile copies a single file from a source path to a destination path.
// It creates the destination file and copies the content.
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file '%s': %w", src, err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file '%s': %w", dst, err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy content from '%s' to '%s': %w", src, dst, err)
	}

	// Preserve file permissions
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to stat source file '%s': %w", src, err)
	}
	return os.Chmod(dst, sourceInfo.Mode())
}
