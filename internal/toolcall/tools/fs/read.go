package fs

import (
	"encoding/json"
	"os"
)

type ReadTool struct{}

func (r ReadTool) Name() string {
	return "read_file"
}

func (r ReadTool) Description() string {
	return "Read and return the contents of a text file from disk."
}

func (r ReadTool) Execute(line string) string {
	var args struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal([]byte(line), &args); err != nil {
		return "error: invalid arguments: " + err.Error()
	}
	data, err := os.ReadFile(args.Path)
	if err != nil {
		return "error reading file: " + err.Error()
	}
	return string(data)
}

func (r ReadTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "Path of the file to read.",
			},
		},
		"required": []string{"path"},
	}
}
