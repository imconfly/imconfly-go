package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "imconfly_go",
		Usage: "Web server for on-the-fly data conversions",
		Commands: []*cli.Command{
			{
				Name:   "run",
				Usage:  "run HTTP server",
				Action: runAction,
			},
			{
				Name:   "transform",
				Usage:  "works like HTTP query but print target filename",
				Action: transformAction,
			},
			{
				Name:   "version",
				Usage:  "print version",
				Action: versionAction,
			},
			{
				Name:   "conf",
				Usage:  "print config in JSON format",
				Action: configAction,
			},
			{
				Name:   "tr-conf",
				Usage:  "print transforms config in JSON format",
				Action: trConfAction,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
