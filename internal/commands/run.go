package commands

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

// CreateNewContainer creates and starts a new container
func CreateNewContainer(command Command) (string, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println("Unable to create docker client")
		panic(err)
	}

	currentDir, err := os.Getwd()

	bgContext := context.Background()

	reader, err := cli.ImagePull(bgContext, command.DockerImage, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)

	cont, err := cli.ContainerCreate(
		bgContext,
		&container.Config{
			Image:        command.DockerImage,
			Cmd:          strings.Split(command.ContainerArgs, " "),
			WorkingDir:   "/tools/qqq",
			AttachStdout: true,
			AttachStderr: true,
			Tty:          true,
			AttachStdin:  true,
			OpenStdin:    true,
		},
		&container.HostConfig{
			AutoRemove: true,
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: currentDir,
					Target: "/tools/qqq",
				},
			},
		}, nil, nil, "")
	if err != nil {
		panic(err)
	}

	reader, err = cli.ContainerLogs(bgContext, cont.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		Follow:     true,
	})
	if err != nil {
		return "", err
	}
	io.Copy(os.Stdout, reader)
	defer reader.Close()

	if err := cli.ContainerStart(bgContext, cont.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	fmt.Printf("Container %s is started\n", cont.ID)
	statusCh, errCh := cli.ContainerWait(bgContext, cont.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	fmt.Println("Exit wait")
	time.Sleep(3000)

	return cont.ID, nil
}
