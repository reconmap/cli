package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/reconmap/cli/internal/api"
	"github.com/reconmap/cli/internal/configuration"
	"github.com/rodaine/table"
	"github.com/urfave/cli/v2"
	"golang.org/x/term"
)

func preActionChecks(c *cli.Context) error {
	if !configuration.HasConfig() {
		return errors.New("Rmap has not been configured. Please call the 'rmap configure' command first.")
	}
	return nil
}

var CommandArguments []*cli.Command = []*cli.Command{
	{
		Name:  "login",
		Usage: "Initiate session with the server",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "username", Aliases: []string{"u"}, Required: true},
			&cli.StringFlag{Name: "password", Aliases: []string{"p"}, Required: false},
		},
		Before: preActionChecks,
		Action: func(c *cli.Context) error {
			var password string
			if c.IsSet("password") {
				password = c.String("password")
			} else {
				fmt.Print("Password: ")
				passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
				if err != nil {
					return err
				}
				password = string(passwordBytes)
				println()
			}
			err := Login(c.String("username"), password)
			return err
		},
	},
	{
		Name:   "logout",
		Usage:  "Terminate session with the server",
		Flags:  []cli.Flag{},
		Before: preActionChecks,
		Action: func(c *cli.Context) error {
			err := Logout()
			return err
		},
	},
	{
		Name:    "config",
		Aliases: []string{"configure"},
		Usage:   "Configure server settings",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "api-url", Aliases: []string{"url"}, Required: true},
		},
		Action: func(c *cli.Context) error {
			err := Configure(c.String("api-url"))
			return err
		},
	},
	{
		Name:    "command",
		Aliases: []string{"c"},
		Usage:   "Command related options",
		Before:  preActionChecks,
		Subcommands: []*cli.Command{
			{
				Name:  "search",
				Usage: "Search commands by keywords",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "keywords", Aliases: []string{"k"}, Required: true},
				},
				Action: func(c *cli.Context) error {
					keywords := c.String("keywords")
					commands, err := api.GetCommandsByKeywords(keywords)
					if err != nil {
						return err
					}

					fmt.Printf("%d commands matching '%s'\n", len(*commands), keywords)
					fmt.Println()

					headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
					columnFmt := color.New(color.FgYellow).SprintfFunc()

					tbl := table.New("ID", "Name", "Description", "Output parser", "Executable type", "Executable path", "Arguments")
					tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

					for _, command := range *commands {
						tbl.AddRow(command.ID, command.Name, command.Description, command.OutputParser, command.ExecutableType, command.ExecutablePath, command.ContainerArgs)

					}
					tbl.Print()

					return err
				},
			},
			{
				Name:  "run",
				Usage: "Run a command and upload the results",
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "commandId", Aliases: []string{"cid"}, Required: true},
					&cli.StringSliceFlag{Name: "var", Required: false},
					&cli.IntFlag{Name: "taskId", Aliases: []string{"tid"}, Required: false},
				},
				Action: func(c *cli.Context) error {
					taskId := c.Int("taskId")
					command, err := api.GetCommandById(c.Int("commandId"))
					if err != nil {
						return err
					}
					err = RunCommand(command, c.StringSlice("var"))
					if err != nil {
						return err
					}

					err = UploadResults(command, taskId)
					return err
				},
			},
		},
	},
	{
		Name:    "debug",
		Aliases: []string{"d"},
		Usage:   "Debug related options",
		Subcommands: []*cli.Command{
			{
				Name:    "list-containers",
				Aliases: []string{"l"},
				Usage:   "List all Reconmap containers",
				Action: func(c *cli.Context) error {
					return ListContainer()
				},
			},
		},
	},
}
