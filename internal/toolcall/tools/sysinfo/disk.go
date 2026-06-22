package sysinfo

import (
	"encoding/json"
	"os/exec"
	"strings"
)

type SysDiskTool struct{}

func (t SysDiskTool) Name() string { return "sys_disk" }

func (t SysDiskTool) Description() string {
	return "Return disk usage for a path (default: /). Args: path (optional string)."
}

func (t SysDiskTool) Execute(args string) string {
	path := "/"
	if strings.TrimSpace(args) != "" {
		var params struct {
			Path string `json:"path"`
		}
		if err := json.Unmarshal([]byte(args), &params); err == nil && params.Path != "" {
			path = params.Path
		}
	}

	out, err := exec.Command("df", "-h", path).CombinedOutput()
	if err != nil {
		return "Error: " + err.Error() + "\n" + string(out)
	}
	return string(out)
}

func (t SysDiskTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "Filesystem path to check disk usage for (default: /).",
			},
		},
	}
}
