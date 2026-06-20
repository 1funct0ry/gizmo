package docker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/1funct0ry/gizmo/internal/utils"
	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/client"
)

type ContainerLogsTool struct {
}

func (l ContainerLogsTool) Name() string {
	return "docker_logs"
}

func (l ContainerLogsTool) Description() string {
	return "Fetch the last N lines of a container's logs (stdout and stderr). Defaults to the last 50 lines."
}

func (l ContainerLogsTool) Execute(line string) string {
	var args struct {
		Container string `json:"container"`
		Tail      string `json:"tail"`
	}
	if err := json.Unmarshal([]byte(line), &args); err != nil {
		return "error: invalid arguments: " + err.Error()
	}
	if args.Container == "" {
		return "error: 'container' is required"
	}
	if args.Tail == "" {
		args.Tail = "50"
	}

	cli, err := newClient()
	if err != nil {
		return fmt.Sprintf("Error creating Docker client: %v", err)
	}
	defer func(cli *client.Client) {
		_ = cli.Close()
	}(cli)

	rc, err := cli.ContainerLogs(context.Background(), args.Container, client.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       args.Tail,
	})
	if err != nil {
		return fmt.Sprintf("Error fetching logs: %v", err)
	}
	defer func() { _ = rc.Close() }()

	// Non-TTY container streams are multiplexed; demultiplex stdout and stderr
	// into a single buffer for the model to read.
	var out bytes.Buffer
	if _, err := stdcopy.StdCopy(&out, &out, rc); err != nil {
		return fmt.Sprintf("Error reading logs: %v", err)
	}

	if out.Len() == 0 {
		return fmt.Sprintf("No logs for container %q.", args.Container)
	}

	return utils.Truncate(out.String(), 20000)
}

func (l ContainerLogsTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"container": map[string]any{
				"type":        "string",
				"description": "Name or ID of the container.",
			},
			"tail": map[string]any{
				"type":        "string",
				"description": "Number of lines from the end of the logs to show, or 'all'. Defaults to '50'.",
			},
		},
		"required": []string{"container"},
	}
}
