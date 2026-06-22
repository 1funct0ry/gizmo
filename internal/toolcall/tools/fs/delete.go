package fs

import (
	"encoding/json"
	"fmt"
	"os"
)

type DeleteFileTool struct{}

func (d DeleteFileTool) Name() string { return "delete_file" }

func (d DeleteFileTool) Description() string {
	return "Delete a file from disk. Requires user confirmation before executing."
}

func (d DeleteFileTool) Execute(line string) string {
	var args struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal([]byte(line), &args); err != nil {
		return "Error: invalid arguments: " + err.Error()
	}
	if err := os.Remove(args.Path); err != nil {
		return "Error: " + err.Error()
	}
	return fmt.Sprintf("Deleted %s", args.Path)
}

func (d DeleteFileTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "Path of the file to delete.",
			},
		},
		"required": []string{"path"},
	}
}
