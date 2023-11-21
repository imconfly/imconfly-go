package server

import (
	"github.com/imconfly/imconfly_go/configuration"
	"github.com/imconfly/imconfly_go/core/internal_workers"
	"github.com/imconfly/imconfly_go/core/queue"
	"github.com/imconfly/imconfly_go/core/resolver"
	"github.com/imconfly/imconfly_go/core/transforms_conf"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func RunServer(conf *configuration.Conf, trConf *transforms_conf.Conf) error {
	var originQ *queue.Queue
	if trConf.HaveNonLocalOrigins() {
		log.Debug("Non-local origins found.")
		originQ = queue.NewQueue()
		defer originQ.Close()
		for i := 0; i < conf.TransformConcurrency; i++ {
			go internal_workers.OriginWorker(originQ, conf.DataDir, conf.TmpDir)
		}
	}

	transformsQ := queue.NewQueue()
	defer transformsQ.Close()
	for i := 0; i < conf.TransformConcurrency; i++ {
		go internal_workers.TransformWorker(transformsQ, originQ, conf.DataDir, conf.TmpDir)
	}

	handler := &Handler{
		Resolver: resolver.NewResolver(transformsQ, trConf, conf.DataDir, conf.TmpDir),
	}

	log.Infof("Server is listening on %s\n", conf.ServerAddr)
	return http.ListenAndServe(conf.ServerAddr, handler)
}
