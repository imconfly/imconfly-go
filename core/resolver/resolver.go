package resolver

import (
	"errors"
	"fmt"
	"github.com/imconfly/imconfly_go/core/internal_workers"
	"github.com/imconfly/imconfly_go/core/queue"
	"github.com/imconfly/imconfly_go/core/request"
	"github.com/imconfly/imconfly_go/core/resolver/resolver_errors"
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
		originQ = queue.NewQueue()
		for i := 0; i < concurrency; i++ {
			go internal_workers.OriginWorker(originQ, dataDir, tmpDir)
		}
		log.Debugf("NewResolver(): Found non local origins. Origin workers started: %d", concurrency)
	} else {
		log.Debug("NewResolver(): Only local origins - no origin queue/workers started.")
	}

	transformsQ := queue.NewQueue()
	for i := 0; i < concurrency; i++ {
		go internal_workers.TransformWorker(transformsQ, originQ, dataDir, tmpDir)
	}
	log.Debugf("NewResolver(): Transforms workers started: %d", concurrency)

	return &Resolver{
		transformsQ: transformsQ,
		trConf:      trConf,
		dataDir:     dataDir,
		tmpDir:      tmpDir,
	}
}

func (r *Resolver) Request(requestStr string) (result os_tools.FileAbsPath, err error) {
	logName := fmt.Sprintf("Resolver.Request(%s)", requestStr)

	req, err := request.RequestFromString(requestStr)
	if err != nil {
		err = resolver_errors.New(http.StatusBadRequest, fmt.Errorf("request format error: %w", err))
		log.Errorf("%s: %s", logName, err.Error())
		return
	}
	task, err := r.trConf.ValidateRequest(req)
	if err != nil {
		err = resolver_errors.New(
			http.StatusBadRequest,
			fmt.Errorf("request don`t match transforms conf: %w", err))
		log.Errorf("%s: %s", logName, err.Error())
		return
	}
	result = task.Request.LocalAbsPath(r.dataDir)
	exist, err := os_tools.FileExist(result)
	// fs error probably or something like this...
	if err != nil {
		log.Errorf("%s: %s", logName, err.Error())
		return
	}
	// nothing to do
	if exist {
		log.Debugf("%s: file exist: %s. Return.", logName, result)
		return
	}
	log.Debugf("%s: file not exist: %s. Add task to transforms queue and waiting.", logName, result)

	err = <-r.transformsQ.Add(task)
	if err != nil {
		var rsError *resolver_errors.Error
		if errors.As(err, &rsError) {
			return
		}
		err = resolver_errors.New(http.StatusInternalServerError, fmt.Errorf("transform error: %w", err))
		log.Errorf("%s: %s", logName, err.Error())
	}
	log.Debugf("%s: Task done.", logName)
	return
}
