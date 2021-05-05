package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/reconmap/cli/internal/api"
	"github.com/reconmap/cli/internal/build"
	"github.com/reconmap/cli/internal/commands"
	"github.com/reconmap/cli/internal/configuration"
	"github.com/reconmap/cli/internal/terminal"
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

func main() {

	banner := "ICBfX19fICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICANCiB8ICBfIFwgX19fICBfX18gX19fICBfIF9fICBfIF9fIF9fXyAgIF9fIF8gXyBfXyAgDQogfCB8XykgLyBfIFwvIF9fLyBfIFx8ICdfIFx8ICdfIGAgXyBcIC8gX2AgfCAnXyBcIA0KIHwgIF8gPCAgX18vIChffCAoXykgfCB8IHwgfCB8IHwgfCB8IHwgKF98IHwgfF8pIHwNCiB8X3wgXF9cX19ffFxfX19cX19fL3xffCB8X3xffCB8X3wgfF98XF9fLF98IC5fXy8gDQogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgfF98ICAgIA0KDQo="
	sDec, _ := base64.StdEncoding.DecodeString(banner)
	color.Set(color.FgHiRed)
	fmt.Print(string(sDec))
	color.Unset()

	app := cli.NewApp()
	app.Version = build.BuildVersion
	app.Copyright = "Reconmap license"
	app.Usage = "Reconmap's command line interface"
	app.Description = "Reconmap's command line interface"
	app.Authors = []*cli.Author{
		{
			Name:  "Reconmap contributors",
			Email: "devs@reconmap.org",
		},
	}
	app.Commands = []*cli.Command{
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
				err := commands.Login(c.String("username"), password)
				return err
			},
		},
		{
			Name:   "logout",
			Usage:  "Terminate session with the server",
			Flags:  []cli.Flag{},
			Before: preActionChecks,
			Action: func(c *cli.Context) error {
				err := commands.Logout()
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
				err := commands.Configure(c.String("api-url"))
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

						tbl := table.New("ID", "Short name", "Description", "Executable type", "Executable path", "Arguments")
						tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

						for _, command := range *commands {
							tbl.AddRow(command.ID, command.ShortName, command.Description, command.ExecutableType, command.ExecutablePath, command.ContainerArgs)

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
						err = commands.RunCommand(command, c.StringSlice("var"))
						if err != nil {
							return err
						}

						err = commands.UploadResults(command, taskId)
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
						return commands.ListContainer()
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		terminal.PrintRedCross()
		fmt.Printf(" %s\n", err)
	}
}
