package main

// @todo: Windows support

import (
	"github.com/imconfly/imconfly_go/cli"
	log "github.com/sirupsen/logrus"
	"os"
)

// main
// functions like *Command must manage output to stdout and stderr by itself
// and call os.Exit() finally
func main() {
	if len(os.Args) < 2 {
		cli.WrongUsage()
	}

	cmd := os.Args[1]
	if cmd == "help" {
		if len(os.Args) == 2 { // imconfly help
			cli.Help("")
		}
		if len(os.Args) == 3 { // imconfly help <command>
			cli.Help(os.Args[2])
		}
		// imconfly help ... .. ..(incorrect args length)
		cli.WrongUsage()
	}

	log.SetLevel(log.TraceLevel)
	switch cmd {
	case "conf":
		cli.Conf()
	case "version":
		cli.Version()
	case "serve":
		cli.Serve()
	default:
		cli.WrongUsage()
	}
}
