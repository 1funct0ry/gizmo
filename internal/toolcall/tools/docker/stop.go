package docker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/moby/moby/client"
)

// StopTool stops a running container. SPEC marks this as "requires
// confirmation", but the confirmation layer is not built yet — that wiring is
// deferred to the shell layer described in SPEC §Agentic Loop.
type StopTool struct {
}

func (s StopTool) Name() string {
	return "docker_stop"
}

func (s StopTool) Description() string {
	return "Stop a running container (graceful SIGTERM, then SIGKILL after the default 10s timeout)."
}

func (s StopTool) Execute(line string) string {
	var args struct {
		Container string `json:"container"`
	}
	if err := json.Unmarshal([]byte(line), &args); err != nil {
		return "error: invalid arguments: " + err.Error()
	}
	if args.Container == "" {
		return "error: 'container' is required"
	}

	cli, err := newClient()
	if err != nil {
		return fmt.Sprintf("Error creating Docker client: %v", err)
	}
	defer func(cli *client.Client) {
		_ = cli.Close()
	}(cli)

	if _, err := cli.ContainerStop(context.Background(), args.Container, client.ContainerStopOptions{}); err != nil {
		return fmt.Sprintf("Error stopping container: %v", err)
	}

	return fmt.Sprintf("Stopped container %s", args.Container)
}

func (s StopTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"container": map[string]any{
				"type":        "string",
				"description": "Name or ID of the container to stop.",
			},
		},
		"required": []string{"container"},
	}
}
