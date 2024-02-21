package server

import (
	"fmt"
	"github.com/imconfly/imconfly_go/config"
	"github.com/imconfly/imconfly_go/core/resolver"
	"github.com/imconfly/imconfly_go/server/handler"
	log "github.com/sirupsen/logrus"

	"net/http"
)

func RunServer(conf *config.Conf) error {
	rs := resolver.NewResolver(
		conf.TransformConcurrency,
		conf.DataDir,
		conf.TmpDir,
		conf.Containers)
	h := handler.NewHandler(rs)

	addr := fmt.Sprintf("%s:%d", conf.ServerHost, conf.ServerPort)
	log.Infof("Server is listening on %s\n", addr)
	return http.ListenAndServe(addr, h)
}
