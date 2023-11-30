package resolver

import (
	"github.com/imconfly/imconfly_go/core/internal_workers"
	"github.com/imconfly/imconfly_go/core/queue"
	"github.com/imconfly/imconfly_go/core/request"
	"github.com/imconfly/imconfly_go/core/resolver/errors"
	"github.com/imconfly/imconfly_go/core/transforms_conf"
	"github.com/imconfly/imconfly_go/lib/os_tools"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Resolver struct {
	transformsQ *queue.Queue
	trConf      *transforms_conf.Conf
	dataDir     os_tools.DirAbsPath
	tmpDir      os_tools.DirAbsPath
}

func NewResolver(
	concurrency int, // for both origins and transforms - it`s ok yet
	dataDir os_tools.DirAbsPath,
	tmpDir os_tools.DirAbsPath,
	trConf *transforms_conf.Conf,
) *Resolver {
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

	return &Resolver{
		transformsQ: transformsQ,
		trConf:      trConf,
		dataDir:     dataDir,
		tmpDir:      tmpDir,
	}
}

// @todo: ctx
func (r *Resolver) Request(requestStr string) (os_tools.FileAbsPath, error) {
	req, err := request.RequestFromString(requestStr)
	if err != nil {
		return "", &errors.ResolverError{
			HTTPCode: http.StatusBadRequest,
			Err:      err,
		}
	}
	task, err := r.trConf.ValidateRequest(req)
	if err != nil {
		return "", err
	}
	out := task.Request.LocalAbsPath(r.dataDir)
	exist, err := os_tools.FileExist(out)
	if err != nil {
		return "", err
	}
	if exist {
		return out, nil
	}

	err = <-r.transformsQ.Add(task)
	return out, err
}
