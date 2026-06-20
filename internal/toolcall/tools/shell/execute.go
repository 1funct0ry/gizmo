package shell

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type ExecuteTool struct {
}

func (s *ExecuteTool) Name() string {
	return "run_shell_command"
}

func (s *ExecuteTool) Description() string {
	return "Execute a shell command on the local machine and return its combined stdout/stderr."
}

func (s *ExecuteTool) Execute(line string) string {
	var args struct {
		Command string `json:"command"`
	}
	if err := json.Unmarshal([]byte(line), &args); err != nil {
		return "error: invalid arguments: " + err.Error()
	}
	out, err := exec.Command("sh", "-c", args.Command).CombinedOutput()
	if err != nil {
		return fmt.Sprintf("error running command: %v\noutput:\n%s", err, out)
	}
	return string(out)
}

func (s *ExecuteTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"command": map[string]any{
				"type":        "string",
				"description": "The shell command to run, e.g. 'ls -la'.",
			},
		},
		"required": []string{"command"},
	}
}
