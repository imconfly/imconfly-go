package handler

import (
	"github.com/imconfly/imconfly_go/core/internal_workers"
	"github.com/imconfly/imconfly_go/core/queue"
	"github.com/imconfly/imconfly_go/core/resolver"
	"github.com/imconfly/imconfly_go/core/transforms_conf"
	"github.com/imconfly/imconfly_go/lib/os_tools"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func BuildHandler(
	concurrency int, // for both origins and transforms - it`s ok yet
	dataDir os_tools.DirAbsPath,
	tmpDir os_tools.DirAbsPath,
	trConf *transforms_conf.Conf,
) http.Handler {
	var originQ *queue.Queue
	if trConf.HaveNonLocalOrigins() {
		log.Debug("Non-local origins found.")
		originQ = queue.NewQueue()
		for i := 0; i < concurrency; i++ {
			go internal_workers.OriginWorker(originQ, dataDir, tmpDir)
		}
		log.Debugf("Origin workers started: %d", concurrency)
	} else {
		log.Debug("Only local origins - no origin queue/workers started.")
	}

	transformsQ := queue.NewQueue()
	for i := 0; i < concurrency; i++ {
		go internal_workers.TransformWorker(transformsQ, originQ, dataDir, tmpDir)
	}
	log.Debugf("Transforms workers started: %d", concurrency)

	return &Handler{
		Resolver: resolver.NewResolver(transformsQ, trConf, dataDir, tmpDir),
	}
}
