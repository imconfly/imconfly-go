package main

import (
	"github.com/imconfly/imconfly_go/cli"
	"log"
	"os"
)

func main() {
	if err := cli.App.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
