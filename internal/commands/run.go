package commands

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

// CreateNewContainer creates and starts a new container
func CreateNewContainer(image string) (string, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println("Unable to create docker client")
		panic(err)
	}
	/*

		hostBinding := nat.PortBinding{
			HostIP:   "0.0.0.0",
			HostPort: "8000",
		}
			containerPort, err := nat.NewPort("tcp", "80")
			if err != nil {
				panic("Unable to get the port")
			}*/

	currentDir, err := os.Getwd()
	fmt.Println(currentDir)

	bgContext := context.Background()
	/*
		vol, err := cli.VolumeCreate(bgContext, volume.VolumeCreateBody{
			Name: "foo",
		})
	*/

	reader, err := cli.ImagePull(bgContext, image, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)

	//portBinding := nat.PortMap{containerPort: []nat.PortBinding{hostBinding}}
	cont, err := cli.ContainerCreate(
		bgContext,
		&container.Config{
			Image:        image,
			Cmd:          []string{"-t", "reconmap.org"},
			WorkingDir:   "/tools/qqq",
			AttachStdout: true,
			AttachStderr: true,
			Tty:          true,
			AttachStdin:  true,
			OpenStdin:    true,
		},
		&container.HostConfig{
			//PortBindings: portBinding,
			AutoRemove: true,
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: currentDir,
					Target: "/tools/qqq",
				},
			},
			/*

				Binds: []string{
					currentDir + ":/tools:rw",
				},
			*/}, nil, nil, "")
	if err != nil {
		panic(err)
	}

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

	reader, err = cli.ContainerLogs(bgContext, cont.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		Follow:     true,
	})
	if err != nil {
		return "", err
	}
	io.Copy(os.Stdout, reader)
	defer reader.Close()

	fmt.Println("Exit wait")
	time.Sleep(1000)

	return cont.ID, nil
}
