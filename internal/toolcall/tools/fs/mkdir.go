package fs

import (
	"encoding/json"
	"fmt"
	"os"
)

type MakeDirTool struct{}

func (m MakeDirTool) Name() string { return "make_dir" }

func (m MakeDirTool) Description() string {
	return "Create a directory and any necessary parent directories."
}

func (m MakeDirTool) Execute(line string) string {
	var args struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal([]byte(line), &args); err != nil {
		return "Error: invalid arguments: " + err.Error()
	}
	if err := os.MkdirAll(args.Path, 0755); err != nil {
		return "Error: " + err.Error()
	}
	return fmt.Sprintf("Created directory %s", args.Path)
}

func (m MakeDirTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "Path of the directory to create.",
			},
		},
		"required": []string{"path"},
	}
}
