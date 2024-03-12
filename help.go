package main

import (
	"fmt"
	"os"

	"github.com/imconfly/imconfly_go/constants"
)

const usageTxt = `imconfly is a web server for on-the-fly data conversion.

Usage:
  imconfly <command>

Commands:
  serve    run HTTP server
  exec     works like HTTP query but print target filename in stdout
  version  print version
  conf     print current configuration in JSON format

Use "imconfly help <command>" for more information about a command.
`

// wrongUsage
// show usage text end exit with EX_USAGE code
func wrongUsage() {
	fmt.Print(usageTxt)
	os.Exit(constants.ExUsage)
}

func helpCommand(subCommand string) {
	if subCommand == "" { // just "imconfly help"
		fmt.Print(usageTxt)
		os.Exit(0)
	}
	switch subCommand {
	case "serve":
		fallthrough
	case "exec":
		fallthrough
	case "version":
		fallthrough
	case "conf":
		fmt.Fprintf(os.Stderr, "not implemented yet\n")
		os.Exit(constants.ExSoftware)
	}
	wrongUsage()
}
