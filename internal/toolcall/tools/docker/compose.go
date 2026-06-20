package docker

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// ComposePsTool reports the status of services in a docker compose project.
// The Moby client has no compose API, so this shells out to `docker compose`,
// mirroring internal/toolcall/tools/shell/execute.go.
type ComposePsTool struct {
}

func (c ComposePsTool) Name() string {
	return "docker_compose_ps"
}

func (c ComposePsTool) Description() string {
	return "Show the status of services in a docker compose project. Optionally point at a specific compose file."
}

func (c ComposePsTool) Execute(line string) string {
	var args struct {
		File string `json:"file"`
	}
	if strings.TrimSpace(line) != "" {
		if err := json.Unmarshal([]byte(line), &args); err != nil {
			return "error: invalid arguments: " + err.Error()
		}
	}

	cmdArgs := []string{"compose"}
	if args.File != "" {
		cmdArgs = append(cmdArgs, "-f", args.File)
	}
	cmdArgs = append(cmdArgs, "ps")

	out, err := exec.Command("docker", cmdArgs...).CombinedOutput()
	if err != nil {
		return fmt.Sprintf("error running docker compose ps: %v\noutput:\n%s", err, out)
	}
	return string(out)
}

func (c ComposePsTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"file": map[string]any{
				"type":        "string",
				"description": "Path to a compose file (passed to 'docker compose -f'). Optional; defaults to compose files in the current directory.",
			},
		},
		"required": []string{},
	}
}
