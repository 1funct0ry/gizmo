package docker

import (
	"context"
	"fmt"
	"strings"

	"github.com/moby/moby/client"
)

type ImageListTool struct {
}

func (i ImageListTool) Name() string {
	return "docker_images"
}

func (i ImageListTool) Description() string {
	return "List all docker images on the local machine."
}

func (i ImageListTool) Execute(line string) string {
	cli, err := newClient()
	if err != nil {
		return fmt.Sprintf("Error creating Docker client: %v", err)
	}
	defer func(cli *client.Client) {
		_ = cli.Close()
	}(cli)

	images, err := cli.ImageList(context.Background(), client.ImageListOptions{})
	if err != nil {
		return fmt.Sprintf("Error listing images: %v", err)
	}

	if len(images.Items) == 0 {
		return "No Docker images found on the local machine."
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Found %d Docker image(s):\n\n", len(images.Items)))

	for _, img := range images.Items {
		result.WriteString(fmt.Sprintf("ID: %s\n", img.ID[:12]))
		if len(img.RepoTags) > 0 {
			result.WriteString(fmt.Sprintf("Tags: %s\n", strings.Join(img.RepoTags, ", ")))
		} else {
			result.WriteString("Tags: <none>\n")
		}
		result.WriteString(fmt.Sprintf("Size: %.2f MB\n", float64(img.Size)/1024/1024))
		result.WriteString(fmt.Sprintf("Created: %d\n\n", img.Created))
	}

	return result.String()
}

func (i ImageListTool) Parameters() map[string]any {
	return map[string]any{
		"type":       "object",
		"properties": map[string]any{},
		"required":   []string{},
	}
}
