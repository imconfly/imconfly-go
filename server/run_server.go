package server

import (
	"github.com/imconfly/imconfly_go/configuration"
	"github.com/imconfly/imconfly_go/core/resolver"
	"github.com/imconfly/imconfly_go/core/transforms_conf"
	"github.com/imconfly/imconfly_go/server/handler"
	"github.com/imconfly/imconfly_go/server/middleware"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func RunServer(conf *configuration.Conf, trConf *transforms_conf.Conf) error {
	rs := resolver.NewResolver(
		conf.TransformConcurrency,
		conf.DataDir,
		conf.TmpDir,
		trConf)
	h := handler.NewHandler(rs)
	h = middleware.Logging(h)

	log.Infof("Server is listening on %s\n", conf.ServerAddr)
	return http.ListenAndServe(conf.ServerAddr, h)
}
