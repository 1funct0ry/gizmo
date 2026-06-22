package fs

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type ListDirTool struct{}

func (l ListDirTool) Name() string { return "list_dir" }

func (l ListDirTool) Description() string {
	return "List files and directories at the given path."
}

func (l ListDirTool) Execute(line string) string {
	var args struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal([]byte(line), &args); err != nil {
		return "Error: invalid arguments: " + err.Error()
	}
	entries, err := os.ReadDir(args.Path)
	if err != nil {
		return "Error: " + err.Error()
	}
	var lines []string
	for _, e := range entries {
		if e.IsDir() {
			lines = append(lines, fmt.Sprintf("%s/ [dir]", e.Name()))
		} else {
			info, err := e.Info()
			if err != nil {
				lines = append(lines, fmt.Sprintf("%s [file]", e.Name()))
			} else {
				lines = append(lines, fmt.Sprintf("%s [file, %d bytes]", e.Name(), info.Size()))
			}
		}
	}
	if len(lines) == 0 {
		return "(empty directory)"
	}
	return strings.Join(lines, "\n")
}

func (l ListDirTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "Path of the directory to list.",
			},
		},
		"required": []string{"path"},
	}
}
