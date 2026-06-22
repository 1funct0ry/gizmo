package git

import (
	"encoding/json"
	"os/exec"
	"strings"
)

type DiffTool struct{}

func (t DiffTool) Name() string { return "git_diff" }

func (t DiffTool) Description() string {
	return "Show uncommitted changes as a unified diff. Optionally scope to a single file."
}

func (t DiffTool) Execute(line string) string {
	var args struct {
		File string `json:"file"`
	}
	if strings.TrimSpace(line) != "" {
		_ = json.Unmarshal([]byte(line), &args)
	}

	cmd := []string{"diff"}
	if args.File != "" {
		cmd = append(cmd, "--", args.File)
	}

	out, err := exec.Command("git", cmd...).CombinedOutput()
	if err != nil {
		return "Error: " + err.Error() + "\n" + string(out)
	}
	result := string(out)
	if len(result) > 20000 {
		result = result[:20000] + "\n... (truncated)"
	}
	return result
}

func (t DiffTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"file": map[string]any{
				"type":        "string",
				"description": "Optional file path to scope the diff to a single file.",
			},
		},
	}
}
