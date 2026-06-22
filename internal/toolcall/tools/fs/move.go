package fs

import (
	"encoding/json"
	"fmt"
	"os"
)

type MoveFileTool struct{}

func (m MoveFileTool) Name() string { return "move_file" }

func (m MoveFileTool) Description() string {
	return "Move or rename a file or directory."
}

func (m MoveFileTool) Execute(line string) string {
	var args struct {
		Src string `json:"src"`
		Dst string `json:"dst"`
	}
	if err := json.Unmarshal([]byte(line), &args); err != nil {
		return "Error: invalid arguments: " + err.Error()
	}
	if err := os.Rename(args.Src, args.Dst); err != nil {
		return "Error: " + err.Error()
	}
	return fmt.Sprintf("Moved %s to %s", args.Src, args.Dst)
}

func (m MoveFileTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"src": map[string]any{
				"type":        "string",
				"description": "Source path of the file or directory.",
			},
			"dst": map[string]any{
				"type":        "string",
				"description": "Destination path.",
			},
		},
		"required": []string{"src", "dst"},
	}
}
