package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/moby/moby/client"
)

type ContainerListTool struct {
}

func (c ContainerListTool) Name() string {
	return "docker_containers"
}

func (c ContainerListTool) Description() string {
	return "List docker containers with their status and ports. By default only running containers are shown; pass all=true to include stopped ones."
}

func (c ContainerListTool) Execute(line string) string {
	var args struct {
		All bool `json:"all"`
	}
	// Arguments are optional; ignore unmarshal errors on empty input.
	if strings.TrimSpace(line) != "" {
		if err := json.Unmarshal([]byte(line), &args); err != nil {
			return "error: invalid arguments: " + err.Error()
		}
	}

	cli, err := newClient()
	if err != nil {
		return fmt.Sprintf("Error creating Docker client: %v", err)
	}
	defer func(cli *client.Client) {
		_ = cli.Close()
	}(cli)

	containers, err := cli.ContainerList(context.Background(), client.ContainerListOptions{All: args.All})
	if err != nil {
		return fmt.Sprintf("Error listing containers: %v", err)
	}

	if len(containers.Items) == 0 {
		if args.All {
			return "No Docker containers found on the local machine."
		}
		return "No running Docker containers. Pass all=true to include stopped containers."
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Found %d container(s):\n\n", len(containers.Items)))

	for _, ctr := range containers.Items {
		id := ctr.ID
		if len(id) > 12 {
			id = id[:12]
		}
		result.WriteString(fmt.Sprintf("ID: %s\n", id))

		names := make([]string, 0, len(ctr.Names))
		for _, n := range ctr.Names {
			names = append(names, strings.TrimPrefix(n, "/"))
		}
		if len(names) > 0 {
			result.WriteString(fmt.Sprintf("Name: %s\n", strings.Join(names, ", ")))
		}

		result.WriteString(fmt.Sprintf("Image: %s\n", ctr.Image))
		result.WriteString(fmt.Sprintf("State: %s\n", ctr.State))
		result.WriteString(fmt.Sprintf("Status: %s\n", ctr.Status))

		if len(ctr.Ports) > 0 {
			// Docker reports IPv4 and IPv6 bindings as separate entries with the
			// same public/private/type, so dedupe to avoid "X->Y, X->Y".
			seen := make(map[string]struct{})
			ports := make([]string, 0, len(ctr.Ports))
			for _, p := range ctr.Ports {
				var s string
				if p.PublicPort != 0 {
					s = fmt.Sprintf("%d->%d/%s", p.PublicPort, p.PrivatePort, p.Type)
				} else {
					s = fmt.Sprintf("%d/%s", p.PrivatePort, p.Type)
				}
				if _, ok := seen[s]; ok {
					continue
				}
				seen[s] = struct{}{}
				ports = append(ports, s)
			}
			result.WriteString(fmt.Sprintf("Ports: %s\n", strings.Join(ports, ", ")))
		}

		result.WriteString("\n")
	}

	return result.String()
}

func (c ContainerListTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"all": map[string]any{
				"type":        "boolean",
				"description": "Include stopped containers, not just running ones. Defaults to false.",
			},
		},
		"required": []string{},
	}
}
