package main

import (
	"fmt"
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
				Name:    "login",
				Aliases: []string{"l"},
				Usage:   "initiates a session with the server",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "username", Aliases: []string{"u"}},
					&cli.StringFlag{Name: "password", Aliases: []string{"p"}},
				},
				Action: func(c *cli.Context) error {
					body, err := commands.Login(c.String("username"), c.String("password"))
					fmt.Printf("%s\n", body)
					return err
				},
			},
			{
				Name:    "run",
				Aliases: []string{"r"},
				Usage:   "runs a command and upload the results",
				Action: func(c *cli.Context) error {
					_, err := commands.CreateNewContainer("reconmap/pentest-container-tools-goohost")
					//_, err := commands.CreateNewContainer("nginx")
					return err
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
			{
				Name:    "upload-results",
				Aliases: []string{"u"},
				Usage:   "upload command results",
				Action: func(c *cli.Context) error {
					commands.UploadResults()
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
