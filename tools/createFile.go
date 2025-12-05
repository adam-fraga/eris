package tools

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateFileTool() Tool {
	var t Tool
	t.Type = "function"
	t.Function.Name = "create_file"
	t.Function.Description = "Create a file on disk (parent folders auto-created)"
	t.Function.Parameters = map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "Full path of the file",
			},
			"content": map[string]any{
				"type":        "string",
				"description": "File contents",
			},
		},
		"required": []string{"path", "content"},
	}
	return t
}

func CreateFile(path, content string) string {
	if path == "" {
		return "error: path is required"
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Sprintf("error creating directories: %v", err)
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Sprintf("error writing file: %v", err)
	}

	return fmt.Sprintf("file created: %s", path)
}
