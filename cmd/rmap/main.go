package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/reconmap/cli/internal/build"
	"github.com/reconmap/cli/internal/commands"
	"github.com/reconmap/cli/internal/logging"
	"github.com/urfave/cli/v2"
)

func main() {
	logger := logging.GetLoggerInstance()
	defer logger.Sync()

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("Version=%s\nBuildDate=%s\nGitCommit=%s\n", c.App.Version, build.BuildTime, build.BuildCommit)
	}

	app := cli.NewApp()
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:     "banner",
			Usage:    "show nice ASCII art banner",
			Aliases:  []string{"b"},
			Required: false,
		},
	}
	app.Before = func(c *cli.Context) error {
		if c.Bool("banner") {
			banner := "ICBfX19fICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICANCiB8ICBfIFwgX19fICBfX18gX19fICBfIF9fICBfIF9fIF9fXyAgIF9fIF8gXyBfXyAgDQogfCB8XykgLyBfIFwvIF9fLyBfIFx8ICdfIFx8ICdfIGAgXyBcIC8gX2AgfCAnXyBcIA0KIHwgIF8gPCAgX18vIChffCAoXykgfCB8IHwgfCB8IHwgfCB8IHwgKF98IHwgfF8pIHwNCiB8X3wgXF9cX19ffFxfX19cX19fL3xffCB8X3xffCB8X3wgfF98XF9fLF98IC5fXy8gDQogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgfF98ICAgIA0KDQo="
			sDec, _ := base64.StdEncoding.DecodeString(banner)
			color.Set(color.FgHiRed)
			fmt.Print(string(sDec))
			color.Unset()
		}
		return nil
	}
	app.Version = build.BuildVersion
	app.Copyright = "Apache License v2.0"
	app.Usage = "Reconmap's CLI"
	app.Description = "Reconmap's command line interface"
	app.Authors = []*cli.Author{
		{
			Name: "Reconmap (https://github.com/reconmap)",
		},
	}
	app.Commands = commands.CommandList

	err := app.Run(os.Args)
	if err != nil {
		logger.Error(err)
	}
}
