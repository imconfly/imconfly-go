package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/imconfly/imconfly_go/config"
	"github.com/imconfly/imconfly_go/constants"
)

const (
	defaultConfFile = "/usr/local/etc/imconfly.json"
	confFileEnvVar  = "IMCONFLY_CONF_FILE"
)

func Conf() {
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

	var yamlFormat bool
	if strings.HasSuffix(confFile, ".json") {
		yamlFormat = false
	} else if strings.HasSuffix(confFile, ".yaml") {
		yamlFormat = true
	} else {
		fmt.Fprintln(os.Stderr, `Config file name must ends with ".json" or ".yaml"`)
		os.Exit(constants.ExConfig)
	}

	f, err := os.Open(confFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(constants.ExConfig)
	}
	defer f.Close()

	c, err := config.ReadConf(f, yamlFormat)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(constants.ExConfig)
	}
	return c
}
