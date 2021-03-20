package commands

import (
	"fmt"

	"github.com/reconmap/cli/internal/configuration"
	"github.com/reconmap/cli/internal/terminal"
)

func Configure(apiUrl string) error {
	config := configuration.Config{ApiUrl: apiUrl}

	filepath, err := configuration.SaveConfig(config)

	terminal.PrintGreenTick()
	fmt.Printf(" Configuration saved to '%s'\n", filepath)

	return err
}
