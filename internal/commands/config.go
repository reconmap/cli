package commands

import (
	"fmt"

	"github.com/reconmap/cli/internal/configuration"
	"github.com/reconmap/cli/internal/terminal"
)

func Configure(authUrl string, apiUrl string) error {
	config := configuration.Config{AuthUrl: authUrl, ApiUrl: apiUrl}

	filepath, err := configuration.SaveConfig(config)

	terminal.PrintGreenTick()
	fmt.Printf(" Configuration saved to '%s'\n", filepath)

	return err
}
