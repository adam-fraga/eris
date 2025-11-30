package tools

import (
	"fmt"
	"os"
)

func CreateFileTool() Tool {
	return Tool{
		Type: "function",
		Function: struct {
			Name        string         `json:"name"`
			Description string         `json:"description"`
			Parameters  map[string]any `json:"parameters"`
		}{
			Name:        "create_file",
			Description: "Create a file on disk",
			Parameters: map[string]any{
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
			},
		},
	}
}

func CreateFile(path, content string) string {
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return "file created"
}
