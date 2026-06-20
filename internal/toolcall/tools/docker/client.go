package docker

import (
	"context"
	"os"
	"path/filepath"

	"github.com/moby/moby/client"
)

// newClient returns a Docker API client connected to the local daemon.
//
// It honours DOCKER_HOST when set; otherwise it probes the common engine
// socket locations and returns the first one that responds to a ping. The API
// version is negotiated with the daemon rather than pinned, so it works across
// daemon versions. (On Docker Desktop for macOS, /var/run/docker.sock is often
// a symlink to the CLI socket and 404s on engine endpoints — the per-user
// ~/.docker/run/docker.sock is the real engine socket, hence the fallback.)
func newClient() (*client.Client, error) {
	// If the environment points at a daemon explicitly, trust it.
	if os.Getenv("DOCKER_HOST") != "" {
		return client.New(client.FromEnv, client.WithAPIVersionNegotiation())
	}

	candidates := []string{"unix:///var/run/docker.sock"}
	if home, err := os.UserHomeDir(); err == nil {
		candidates = append(candidates, "unix://"+filepath.Join(home, ".docker", "run", "docker.sock"))
	}

	for _, host := range candidates {
		c, err := client.New(client.WithHost(host), client.WithAPIVersionNegotiation())
		if err != nil {
			continue
		}
		if _, err := c.Ping(context.Background(), client.PingOptions{}); err != nil {
			_ = c.Close()
			continue
		}
		return c, nil
	}

	// Nothing pinged successfully. Return a client for the default host anyway
	// so callers surface a meaningful "Error response from daemon" to the model.
	return client.New(client.WithHost(candidates[0]), client.WithAPIVersionNegotiation())
}
