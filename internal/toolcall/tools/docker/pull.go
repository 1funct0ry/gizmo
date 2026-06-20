package docker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/moby/moby/client"
)

type PullTool struct {
}

func (p PullTool) Name() string {
	return "docker_pull"
}

func (p PullTool) Description() string {
	return "Pull an image from a registry onto the local machine."
}

func (p PullTool) Execute(line string) string {
	var args struct {
		Image string `json:"image"`
	}
	if err := json.Unmarshal([]byte(line), &args); err != nil {
		return "error: invalid arguments: " + err.Error()
	}
	if args.Image == "" {
		return "error: 'image' is required"
	}

	cli, err := newClient()
	if err != nil {
		return fmt.Sprintf("Error creating Docker client: %v", err)
	}
	defer func(cli *client.Client) {
		_ = cli.Close()
	}(cli)

	ctx := context.Background()
	resp, err := cli.ImagePull(ctx, args.Image, client.ImagePullOptions{})
	if err != nil {
		return fmt.Sprintf("Error pulling image: %v", err)
	}
	defer func() { _ = resp.Close() }()

	// Block until the pull finishes (or fails).
	if err := resp.Wait(ctx); err != nil {
		return fmt.Sprintf("Error pulling image: %v", err)
	}

	return fmt.Sprintf("Pulled %s", args.Image)
}

func (p PullTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"image": map[string]any{
				"type":        "string",
				"description": "Image reference to pull, e.g. 'alpine:latest'.",
			},
		},
		"required": []string{"image"},
	}
}
