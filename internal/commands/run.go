package commands

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/reconmap/cli/internal/api"
	"github.com/reconmap/cli/internal/terminal"
)

// CreateNewContainer creates and starts a new container
func CreateNewContainer(command *api.Command, vars []string) (string, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return "", err
	}

	var updatedArgs = command.ContainerArgs
	for _, v := range vars {
		var tokens = strings.Split(v, "=")
		var validID = regexp.MustCompile("{{{" + tokens[0] + ".*}}}")
		updatedArgs = validID.ReplaceAllString(updatedArgs, tokens[1])
	}

	currentDir, err := os.Getwd()

	bgContext := context.Background()

	terminal.PrintYellowDot()
	fmt.Printf(" Downloading docker image '%s'\n", command.DockerImage)
	reader, err := cli.ImagePull(bgContext, command.DockerImage, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBufferString("")
	io.Copy(buf, reader)
	fmt.Println(buf.String())

	commandLineArgs := strings.Split(updatedArgs, " ")
	terminal.PrintYellowDot()
	fmt.Printf(" Using command line args: %s\n", commandLineArgs)

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
		}, nil, nil, "reconmap-"+command.ShortName)
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
	stdcopy.StdCopy(os.Stdout, os.Stderr, reader)

	terminal.PrintYellowDot()
	fmt.Printf(" Starting container.\n")
	if err := cli.ContainerStart(bgContext, cont.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	go func() {
		defer resp.Close()
		io.Copy(os.Stdout, resp.Reader)
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

	stdcopy.StdCopy(os.Stdout, os.Stderr, reader)

	terminal.PrintGreenTick()
	fmt.Printf(" Done\n")

	return cont.ID, nil
}
