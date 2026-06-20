package docker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/client"
)

// ExecTool runs a command inside a running container. SPEC marks this as
// "requires confirmation", but the confirmation layer (destructive-tool gating,
// --no-confirm, shell prompt) is not built yet — that wiring is deferred to the
// shell layer described in SPEC §Agentic Loop.
type ExecTool struct {
}

func (e ExecTool) Name() string {
	return "docker_exec"
}

func (e ExecTool) Description() string {
	return "Run a command inside a running container and return its output."
}

func (e ExecTool) Execute(line string) string {
	var args struct {
		Container string `json:"container"`
		Command   string `json:"command"`
	}
	if err := json.Unmarshal([]byte(line), &args); err != nil {
		return "error: invalid arguments: " + err.Error()
	}
	if args.Container == "" {
		return "error: 'container' is required"
	}
	if args.Command == "" {
		return "error: 'command' is required"
	}

	cli, err := newClient()
	if err != nil {
		return fmt.Sprintf("Error creating Docker client: %v", err)
	}
	defer func(cli *client.Client) {
		_ = cli.Close()
	}(cli)

	ctx := context.Background()

	created, err := cli.ExecCreate(ctx, args.Container, client.ExecCreateOptions{
		Cmd:          []string{"sh", "-c", args.Command},
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		return fmt.Sprintf("Error creating exec: %v", err)
	}

	att, err := cli.ExecAttach(ctx, created.ID, client.ExecAttachOptions{})
	if err != nil {
		return fmt.Sprintf("Error attaching to exec: %v", err)
	}
	defer att.Close()

	// The exec stream is multiplexed (no TTY); demultiplex into one buffer.
	var out bytes.Buffer
	if _, err := stdcopy.StdCopy(&out, &out, att.Reader); err != nil {
		return fmt.Sprintf("Error reading exec output: %v", err)
	}

	result := out.String()

	if inspect, err := cli.ExecInspect(ctx, created.ID, client.ExecInspectOptions{}); err == nil && inspect.ExitCode != 0 {
		result += fmt.Sprintf("\n[exit code: %d]", inspect.ExitCode)
	}

	if result == "" {
		return "(command produced no output)"
	}
	return result
}

func (e ExecTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"container": map[string]any{
				"type":        "string",
				"description": "Name or ID of the running container.",
			},
			"command": map[string]any{
				"type":        "string",
				"description": "The command to run inside the container, e.g. 'ls -la /app'.",
			},
		},
		"required": []string{"container", "command"},
	}
}
