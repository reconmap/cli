package commands

import (
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/reconmap/cli/internal/api"
	"github.com/reconmap/cli/internal/terminal"
)

// replace with https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html

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
		argsRendered := terminal.ReplaceArgs(command, vars)
		println("Command to run: " + command.ExecutablePath + " " + argsRendered)
		cmd := exec.Command(command.ExecutablePath, strings.Fields(argsRendered)...)
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

		outputFilename := strconv.Itoa(command.ID) + ".out"
		f, err := os.Create(outputFilename)
		defer f.Close()
		f.WriteString(outStr)
		command.OutputFileName = outputFilename

		if len(errStr) > 0 {
			log.Println(errStr)
		}
	} else {
		_, err = CreateNewContainer(command, vars)
	}

	return err
}
