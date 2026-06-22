package sysinfo

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

type SysProcessesTool struct{}

func (t SysProcessesTool) Name() string { return "sys_processes" }

func (t SysProcessesTool) Description() string {
	return "Return the top N processes by CPU usage (default: 10). Args: n (optional int)."
}

func (t SysProcessesTool) Execute(args string) string {
	n := 10
	if strings.TrimSpace(args) != "" {
		var params struct {
			N int `json:"n"`
		}
		if err := json.Unmarshal([]byte(args), &params); err == nil && params.N > 0 {
			n = params.N
		}
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "darwin" {
		// ps on macOS doesn't support --sort; -r sorts by CPU descending
		cmd = exec.Command("ps", "aux", "-r")
	} else {
		cmd = exec.Command("ps", "aux", "--sort=-%cpu")
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "Error: " + err.Error() + "\n" + string(out)
	}

	lines := strings.Split(strings.TrimRight(string(out), "\n"), "\n")
	// lines[0] is the header; keep header + top N process lines
	limit := 1 + n
	if limit > len(lines) {
		limit = len(lines)
	}
	return fmt.Sprintf("Top %d processes by CPU:\n", n) + strings.Join(lines[:limit], "\n") + "\n"
}

func (t SysProcessesTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"n": map[string]any{
				"type":        "integer",
				"description": "Number of top processes to return (default: 10).",
			},
		},
	}
}
