package tools

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateFile(filename, content string) string {
	if filename == "" {
		filename = "note.txt" // default name
	}

	// Force all files inside 'output' folder
	fullPath := filepath.Join("output", filename)

	// Ensure the final path is not a directory
	info, err := os.Stat(fullPath)
	if err == nil && info.IsDir() {
		return fmt.Sprintf("Failed to create file: '%s' is a directory", fullPath)
	}

	err = os.WriteFile(fullPath, []byte(content), 0644)
	if err != nil {
		return fmt.Sprintf("Failed to create file: %v", err)
	}
	return fmt.Sprintf("File '%s' created successfully.", fullPath)
}
