package commands

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func ListContainer() error {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	if len(containers) > 0 {
		for _, container := range containers {
			fmt.Printf("Container ID: %s", container.ID)
		}
	} else {
		fmt.Println("There are no containers running")
	}
	return nil
}
