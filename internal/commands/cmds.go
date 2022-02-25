package commands

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/reconmap/cli/internal/api"
)

func copyAndCapture(w io.Writer, r io.Reader) ([]byte, error) {
	var out []byte
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			out = append(out, d...)
			_, err := w.Write(d)
			if err != nil {
				return out, err
			}
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			return out, err
		}
	}
}

func RunCommand(command *api.Command, vars []string) error {
	var err error
	if command.ExecutableType == "custom" {
		println("Command to run: " + command.ExecutablePath + " " + command.ContainerArgs)
		cmd := exec.Command(command.ExecutablePath, strings.Fields(command.ContainerArgs)...)
		var stdout, stderr []byte
		var errStdout, errStderr error
		stdoutIn, _ := cmd.StdoutPipe()
		stderrIn, _ := cmd.StderrPipe()
		err := cmd.Start()
		if err != nil {
			log.Fatalf("cmd.Start() failed with '%s'\n", err)
		}
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			stdout, errStdout = copyAndCapture(os.Stdout, stdoutIn)
			wg.Done()
		}()

		stderr, errStderr = copyAndCapture(os.Stderr, stderrIn)

		wg.Wait()

		err = cmd.Wait()
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}
		if errStdout != nil || errStderr != nil {
			log.Fatal("failed to capture stdout or stderr\n")
		}
		outStr, errStr := string(stdout), string(stderr)
		fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)
	} else {
		_, err = CreateNewContainer(command, vars)
	}

	return err
}
