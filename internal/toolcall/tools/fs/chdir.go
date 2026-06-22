package fs

import (
	"encoding/json"
	"fmt"
	"os"
)

type ChangeDirTool struct{}

func (c ChangeDirTool) Name() string { return "change_dir" }

func (c ChangeDirTool) Description() string {
	return "Change the agent's working directory for this session."
}

func (c ChangeDirTool) Execute(line string) string {
	var args struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal([]byte(line), &args); err != nil {
		return "Error: invalid arguments: " + err.Error()
	}
	if err := os.Chdir(args.Path); err != nil {
		return "Error: " + err.Error()
	}
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Sprintf("Changed to %s", args.Path)
	}
	return fmt.Sprintf("Changed working directory to %s", cwd)
}

func (c ChangeDirTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "Directory path to change to.",
			},
		},
		"required": []string{"path"},
	}
}
