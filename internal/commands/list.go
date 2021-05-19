package commands

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/reconmap/cli/internal/containers"
)

// ListContainer list all the containers created by Reconmap
func ListContainer() error {
	containerContext := context.Background()
	cli, err := containers.CreateNewClient()
	if err != nil {
		panic(err)
	}

	f := filters.NewArgs(filters.KeyValuePair{Key: "label", Value: "reconmap"})
	containers, err := cli.ContainerList(containerContext, types.ContainerListOptions{
		Filters: f,
	})
	if err != nil {
		panic(err)
	}

	if len(containers) > 0 {
		for _, container := range containers {
			fmt.Printf("Container ID: %s, %s\n", container.Image, container.ID)
		}
	} else {
		fmt.Println("There are no containers running")
	}
	return nil
}
