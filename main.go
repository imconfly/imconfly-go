package main

import (
	"fmt"
	"log"
	"os"

	"github.com/imconfly/imconfly_go/version"

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
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func runAction(ctx *cli.Context) error {
	fmt.Println("starts HTTP server here")
	return nil
}

func transformAction(ctx *cli.Context) error {
	fmt.Println("transform action here")
	return nil
}

func versionAction(ctx *cli.Context) error {
	fmt.Println(version.VERSION)
	return nil
}

func configAction(ctx *cli.Context) error {
	fmt.Println("yama conf")
	return nil
}
