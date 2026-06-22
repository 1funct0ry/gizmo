package fs

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type EditTool struct{}

func (e EditTool) Name() string { return "edit_file" }

func (e EditTool) Description() string {
	return "Replace an exact string in a file. Safer than a full overwrite."
}

func (e EditTool) Execute(line string) string {
	var args struct {
		Path   string `json:"path"`
		OldStr string `json:"old_str"`
		NewStr string `json:"new_str"`
	}
	if err := json.Unmarshal([]byte(line), &args); err != nil {
		return "Error: invalid arguments: " + err.Error()
	}
	data, err := os.ReadFile(args.Path)
	if err != nil {
		return "Error: " + err.Error()
	}
	content := string(data)
	if !strings.Contains(content, args.OldStr) {
		return "Error: old_str not found in file"
	}
	updated := strings.Replace(content, args.OldStr, args.NewStr, 1)
	if err := os.WriteFile(args.Path, []byte(updated), 0644); err != nil {
		return "Error: " + err.Error()
	}
	return fmt.Sprintf("Edited %s", args.Path)
}

func (e EditTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "Path to the file to edit.",
			},
			"old_str": map[string]any{
				"type":        "string",
				"description": "Exact string to replace (must appear exactly once).",
			},
			"new_str": map[string]any{
				"type":        "string",
				"description": "String to substitute in place of old_str.",
			},
		},
		"required": []string{"path", "old_str", "new_str"},
	}
}
