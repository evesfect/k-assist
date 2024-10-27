package dirutil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GetCurrentDirectoryContents returns a string containing information about the current directory
func GetCurrentDirectoryContents(dir string) (string, error) {
	dirContents, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	var info strings.Builder
	info.WriteString(fmt.Sprintf("Current directory: %s\nDirectory contents:\n", dir))
	for _, entry := range dirContents {
		info.WriteString(fmt.Sprintf("- %s\n", entry.Name()))
	}
	return info.String(), nil
}

// GetAllDirectoryContents returns a string containing information about all subdirectories and files
func GetAllDirectoryContents(dir string) (string, error) {
	var info strings.Builder
	info.WriteString(fmt.Sprintf("Current directory: %s\nAll directory contents:\n", dir))

	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		info.WriteString(fmt.Sprintf("- %s\n", relPath))
		return nil
	})

	if err != nil {
		return "", err
	}

	return info.String(), nil
}
