package commands

import (
	"github.com/reconmap/cli/internal/api"
)

func RunCommand(command *api.Command, vars []string) error {

	_, err := CreateNewContainer(command, vars)

	return err
}
