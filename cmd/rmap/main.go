package main

import (
	"log"
	"os"

	"github.com/reconmap/cli/internal/commands"
	"github.com/urfave/cli/v2"
)

func main() {

	app := cli.App{
		Name: "Reconmap CLI",
		Authors: []*cli.Author{
			{
				Name:  "Reconmap developers",
				Email: "devs@reconmap.org",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "run",
				Aliases: []string{"r"},
				Usage:   "runs a command and upload the results",
				Action: func(c *cli.Context) error {
					commands.CreateNewContainer("hello-world")
					return nil
				},
			},
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "list Reconmap containers",
				Action: func(c *cli.Context) error {
					commands.ListContainer()
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
