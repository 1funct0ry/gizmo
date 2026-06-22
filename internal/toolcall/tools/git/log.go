package git

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type LogTool struct{}

func (t LogTool) Name() string { return "git_log" }

func (t LogTool) Description() string {
	return "Show the last N commits in short oneline format (default 10)."
}

func (t LogTool) Execute(line string) string {
	var args struct {
		N int `json:"n"`
	}
	args.N = 10
	if strings.TrimSpace(line) != "" {
		_ = json.Unmarshal([]byte(line), &args)
	}
	if args.N <= 0 {
		args.N = 10
	}

	out, err := exec.Command("git", "log", "--oneline", fmt.Sprintf("-n%d", args.N)).CombinedOutput()
	if err != nil {
		return "Error: " + err.Error() + "\n" + string(out)
	}
	return string(out)
}

func (t LogTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"n": map[string]any{
				"type":        "integer",
				"description": "Number of commits to show (default 10).",
			},
		},
	}
}
