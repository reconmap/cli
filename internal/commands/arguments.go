package commands

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/reconmap/cli/internal/api"
	"github.com/reconmap/cli/internal/configuration"
	"github.com/rodaine/table"
	"github.com/urfave/cli/v2"
	"golang.org/x/term"
)

func preActionChecks(c *cli.Context) error {
	if !configuration.HasConfig() {
		return errors.New("Rmap has not been configured. Please call the 'rmap config' command first.")
	}
	return nil
}

var CommandArguments []*cli.Command = []*cli.Command{
	{
		Name:  "login",
		Usage: "Initiates session with the server",
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
		Usage:  "Terminates session with the server",
		Flags:  []cli.Flag{},
		Before: preActionChecks,
		Action: func(c *cli.Context) error {
			err := Logout()
			return err
		},
	},
	{
		Name:    "configure",
		Aliases: []string{"config"},
		Usage:   "Configures server settings",
		Subcommands: []*cli.Command{
			{
				Name:  "set",
				Usage: "Configures server settings",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "api-url", Required: true},
				},
				Action: func(c *cli.Context) error {
					err := Configure(c.String("api-url"))
					return err
				},
			},
			{
				Name:  "view",
				Usage: "View server settings",
				Action: func(c *cli.Context) error {
					if configuration.HasConfig() {
						config, err := configuration.ReadConfig()
						if err != nil {
							return err
						}
						fmt.Printf("api-url: %s\n", config.ApiUrl)
					}
					return nil
				},
			},
		},
	},
	{
		Name:    "command",
		Aliases: []string{"cmd"},
		Usage:   "Works with commands",
		Before:  preActionChecks,
		Subcommands: []*cli.Command{
			{
				Name:  "search",
				Usage: "Search commands by keywords",
				Action: func(c *cli.Context) error {
					if c.Args().Len() == 0 {
						return errors.New("no keywords were entered after the search command")
					}
					var keywords string = strings.Join(c.Args().Slice(), " ")
					commands, err := api.GetCommandsByKeywords(keywords)
					if err != nil {
						return err
					}

					var numCommands int = len(*commands)
					fmt.Printf("%d commands matching '%s'\n", numCommands, keywords)

					if numCommands > 0 {
						fmt.Println()

						headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
						columnFmt := color.New(color.FgYellow).SprintfFunc()

						tbl := table.New("ID", "Name", "Description", "Output parser", "Executable type", "Executable path", "Arguments")
						tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

						for _, command := range *commands {
							tbl.AddRow(command.ID, command.Name, command.Description, command.OutputParser, command.ExecutableType, command.ExecutablePath, command.ContainerArgs)

						}
						tbl.Print()
					}

					return err
				},
			},
			{
				Name:  "run",
				Usage: "Run a command and upload its output to the server",
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "commandId", Aliases: []string{"cid"}, Required: true},
					&cli.IntFlag{Name: "taskId", Aliases: []string{"tid"}, Required: false},
					&cli.StringSliceFlag{Name: "var", Required: false},
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
			{
				Name:  "upload-output",
				Usage: "Upload command output to server",
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "taskId", Aliases: []string{"tid"}, Required: true},
					&cli.StringFlag{Name: "outputFile", Aliases: []string{"of"}, Required: true},
				},
				Action: func(c *cli.Context) error {
					taskId := c.Int("taskId")
					outputFileName := c.String("outputFile")

					return UploadCommandOutputUsingFileName(outputFileName, taskId)
				},
			},
		},
	},
	{
		Name:    "task",
		Aliases: []string{"tsk"},
		Usage:   "Lists tasks",
		Before:  preActionChecks,
		Subcommands: []*cli.Command{
			{
				Name:  "search",
				Usage: "Search tasks by keywords",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "keywords", Aliases: []string{"k"}, Required: true},
				},
				Action: func(c *cli.Context) error {
					keywords := c.String("keywords")
					tasks, err := api.GetTasksByKeywords(keywords)
					if err != nil {
						return err
					}

					var numTasks int = len(*tasks)
					fmt.Printf("%d tasks matching '%s'\n", numTasks, keywords)

					if numTasks > 0 {
						fmt.Println()

						headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
						columnFmt := color.New(color.FgYellow).SprintfFunc()

						tbl := table.New("ID", "Summary", "Description", "Status")
						tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

						for _, task := range *tasks {
							tbl.AddRow(task.ID, task.Summary, task.Description, task.Status)

						}
						tbl.Print()
					}

					return err
				},
			},
		},
	},
	{
		Name:    "vulnerability",
		Aliases: []string{"vuln"},
		Usage:   "Works with vulnerabilities",
		Before:  preActionChecks,
		Subcommands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List open vulnerabilities",
				Action: func(c *cli.Context) error {
					vulnerabilities, err := api.GetVulnerabilities()
					if err != nil {
						return err
					}

					var numVulnerabilities int = len(*vulnerabilities)
					fmt.Printf("%d vulnerabilities\n", numVulnerabilities)

					if numVulnerabilities > 0 {
						fmt.Println()

						headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
						columnFmt := color.New(color.FgYellow).SprintfFunc()

						tbl := table.New("ID", "Summary", "Risk", "Status", "Category")
						tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

						for _, vuln := range *vulnerabilities {
							tbl.AddRow(vuln.ID, vuln.Summary, vuln.Risk, vuln.Status, vuln.CategoryName)

						}
						tbl.Print()
					}

					return err
				},
			},
		},
	},
	{
		Name:    "debug",
		Aliases: []string{"dbg"},
		Usage:   "Shows debugging info",
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
