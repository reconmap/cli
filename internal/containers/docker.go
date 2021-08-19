package containers

import (
	"github.com/docker/docker/client"
)

func CreateNewClient() (*client.Client, error) {
	return client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
}
