package fs

import (
	"encoding/json"
	"fmt"
	"os"
)

type WriteTool struct{}

func (w WriteTool) Name() string {
	return "write_file"
}

func (w WriteTool) Description() string {
	return "Write a string to a text file on disk."
}

func (w WriteTool) Execute(line string) string {
	var args struct {
		Path    string `json:"path"`
		Content string `json:"content"`
	}
	if err := json.Unmarshal([]byte(line), &args); err != nil {
		return "error: invalid arguments: " + err.Error()
	}
	err := os.WriteFile(args.Path, []byte(args.Content), 0644)

	if err != nil {
		return "error reading file: " + err.Error()
	}
	return fmt.Sprintf("Wrote %s to %s", args.Content, args.Path)
}

func (w WriteTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "Path of the file to write to.",
			},
			"content": map[string]any{
				"type":        "string",
				"description": "Content to write to the file.",
			},
		},
		"required": []string{"path", "content"},
	}
}
