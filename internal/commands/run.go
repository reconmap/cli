package commands

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/reconmap/cli/internal/containers"
	"github.com/reconmap/cli/internal/terminal"
	"github.com/reconmap/shared-lib/pkg/api"
)

// CreateNewContainer creates and starts a new container
func CreateNewContainer(command *api.Command, vars []string) (string, error) {
	cli, err := containers.CreateNewClient()
	if err != nil {
		return "", err
	}

	var updatedArgs = terminal.ReplaceArgs(command, vars)

	currentDir, err := os.Getwd()

	bgContext := context.Background()

	terminal.PrintYellowDot()
	fmt.Printf(" Downloading docker image '%s'\n", command.DockerImage)
	reader, err := cli.ImagePull(bgContext, command.DockerImage, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBufferString("")
	if _, err = io.Copy(buf, reader); err != nil {
		return "", err
	}
	fmt.Println(buf.String())

	commandLineArgs := strings.Split(updatedArgs, " ")
	terminal.PrintYellowDot()
	fmt.Printf(" Using command line args: %s\n", commandLineArgs)

	var containerName string = "reconmap-" + command.Name

	f := filters.NewArgs(filters.KeyValuePair{Key: "name", Value: containerName})
	containers, err := cli.ContainerList(bgContext, types.ContainerListOptions{
		Filters: f,
	})
	if err != nil {
		panic(err)
	}

	if len(containers) > 0 {
		for _, container := range containers {
			var timeout time.Duration = 5 * time.Second
			fmt.Printf("Container ID: %s, %s\n", container.Image, container.ID)
			if err = cli.ContainerStop(bgContext, container.ID, &timeout); err != nil {
				return "", err
			}
			err = cli.ContainerRemove(bgContext, container.ID, types.ContainerRemoveOptions{})
			if err != nil {
				return "", err
			}

		}
	} else {
		fmt.Println("There are no containers running")
	}

	cont, err := cli.ContainerCreate(
		bgContext,
		&container.Config{
			Image:        command.DockerImage,
			Cmd:          commandLineArgs,
			WorkingDir:   "/reconmap",
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
					Target: "/reconmap",
				},
			},
		}, nil, nil, containerName)
	if err != nil {
		panic(err)
	}

	resp, err := cli.ContainerAttach(bgContext, cont.ID, types.ContainerAttachOptions{
		Stream: true,
		Stdin:  false,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return "", err
	}

	reader, err = cli.ContainerLogs(bgContext, cont.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	if err != nil {
		return "", err
	}
	if _, err = stdcopy.StdCopy(os.Stdout, os.Stderr, reader); err != nil {
		return "", err
	}

	terminal.PrintYellowDot()
	fmt.Printf(" Starting container.\n")
	if err := cli.ContainerStart(bgContext, cont.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	go func() {
		defer resp.Close()
		if _, err = io.Copy(os.Stdout, resp.Reader); err != nil {
			fmt.Printf("Error logging response: %s\n", err)
		}
	}()

	terminal.PrintYellowDot()
	fmt.Printf(" Container started.\n")
	statusCh, errCh := cli.ContainerWait(bgContext, cont.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case status := <-statusCh:
		terminal.PrintYellowDot()
		fmt.Printf(" Container '%s' exited with code %d.\n", command.DockerImage, status.StatusCode)
	}

	if _, err = stdcopy.StdCopy(os.Stdout, os.Stderr, reader); err != nil {
		return "", err
	}

	terminal.PrintGreenTick()
	fmt.Printf(" Done\n")

	return cont.ID, nil
}
