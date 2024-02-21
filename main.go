package main

// @todo: Windows support

import (
	"encoding/json"
	"fmt"
	"github.com/imconfly/imconfly_go/config"
	"github.com/imconfly/imconfly_go/constants"
	"github.com/imconfly/imconfly_go/server"
	log "github.com/sirupsen/logrus"
	"os"
)

const (
	defaultConfFile = "/usr/local/etc/imconfly.json"
	confFileEnvVar  = "IMCONFLY_CONF_FILE"
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

// main
// functions like *Command must manage output to stdout and stderr by itself
// and call os.Exit() finally
func main() {
	if len(os.Args) < 2 {
		wrongUsage()
	}

	cmd := os.Args[1]
	if cmd == "help" {
		if len(os.Args) == 2 { // imconfly help
			helpCommand("")
		}
		if len(os.Args) == 3 { // imconfly help <command>
			helpCommand(os.Args[2])
		}
		// imconfly help ... .. ..(incorrect args length)
		wrongUsage()
	}

	log.SetLevel(log.TraceLevel)
	switch cmd {
	case "conf":
		confCommand()
	case "version":
		versionCommand()
	case "serve":
		serveCommand()
	default:
		wrongUsage()
	}
}

// wrongUsage
// show usage text end exit with EX_USAGE code
func wrongUsage() {
	fmt.Print(usageTxt)
	os.Exit(constants.ExUsage)
}

func serveCommand() {
	c := getConf()
	err := server.RunServer(c)
	fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	os.Exit(constants.ExSoftware)
}

func versionCommand() {
	fmt.Println(constants.Version)
	os.Exit(0)
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
		os.Exit(70)
	}
	wrongUsage()
}

func confCommand() {
	c := getConf()

	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
	os.Exit(0)
}

func getConf() *config.Conf {
	confFile := os.Getenv(confFileEnvVar)
	if confFile == "" {
		confFile = defaultConfFile
	}

	f, err := os.Open(confFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(constants.ExConfig)
	}
	defer f.Close()

	c, err := config.ReadConf(f)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(constants.ExConfig)
	}
	return c
}
