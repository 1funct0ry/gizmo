package docker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/moby/moby/client"
)

type ContainerInspectTool struct {
}

func (i ContainerInspectTool) Name() string {
	return "docker_inspect"
}

func (i ContainerInspectTool) Description() string {
	return "Return the full JSON inspect output for a container (config, state, network, mounts, etc.)."
}

func (i ContainerInspectTool) Execute(line string) string {
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

	res, err := cli.ContainerInspect(context.Background(), args.Container, client.ContainerInspectOptions{})
	if err != nil {
		return fmt.Sprintf("Error inspecting container: %v", err)
	}

	var pretty bytes.Buffer
	if err := json.Indent(&pretty, res.Raw, "", "  "); err != nil {
		// Fall back to the raw payload if it can't be re-indented.
		return string(res.Raw)
	}

	return pretty.String()
}

func (i ContainerInspectTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"container": map[string]any{
				"type":        "string",
				"description": "Name or ID of the container to inspect.",
			},
		},
		"required": []string{"container"},
	}
}
