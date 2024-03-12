package cli

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/imconfly/imconfly_go/config"
	"github.com/imconfly/imconfly_go/constants"
	"github.com/imconfly/imconfly_go/server"
)

func Serve() {
	c := getConf()
	if err := config.CheckDirs(c.DataDir, c.TmpDir); err != nil {
		log.Errorf("%s", err)
	}
	log.Debugf("DataDir: %s, TmpDir: %s. Ok.", c.DataDir, c.TmpDir)

	err := server.RunServer(c)
	fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	os.Exit(constants.ExSoftware)
}
