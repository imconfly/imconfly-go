package main

import (
	"github.com/imconfly/imconfly_go/cli"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	log.SetLevel(log.TraceLevel)
	if err := cli.App.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
