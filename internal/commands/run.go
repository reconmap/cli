package commands

import (
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

// CreateNewContainer creates and starts a new container
func CreateNewContainer(command Command, vars []string) (string, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println("Unable to create docker client")
		panic(err)
	}

	var updatedArgs = command.ContainerArgs
	for _, v := range vars {
		var tokens = strings.Split(v, "=")
		var validID = regexp.MustCompile("{{{" + tokens[0] + ".*}}}")
		updatedArgs = validID.ReplaceAllString(updatedArgs, tokens[1])
	}

	currentDir, err := os.Getwd()

	bgContext := context.Background()

	fmt.Printf("> Downloading docker image '%s'\n", command.DockerImage)
	reader, err := cli.ImagePull(bgContext, command.DockerImage, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	cont, err := cli.ContainerCreate(
		bgContext,
		&container.Config{
			Image:        command.DockerImage,
			Cmd:          strings.Split(updatedArgs, " "),
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

	fmt.Printf("> Starting container.\n")
	if err := cli.ContainerStart(bgContext, cont.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	fmt.Printf("> Container started.\n")
	statusCh, errCh := cli.ContainerWait(bgContext, cont.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	time.Sleep(3000)
	fmt.Printf("> Container '%s' exited.\n", command.DockerImage)

	return cont.ID, nil
}
