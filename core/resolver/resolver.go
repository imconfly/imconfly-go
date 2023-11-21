package resolver

import (
	"github.com/imconfly/imconfly_go/core/queue"
	"github.com/imconfly/imconfly_go/core/request"
	"github.com/imconfly/imconfly_go/lib/os_tools"
	"github.com/imconfly/imconfly_go/transforms_conf"
)

type Resolver struct {
	transformsQ *queue.Queue
	trConf      *transforms_conf.Conf
	dataDir     os_tools.DirAbsPath
	tmpDir      os_tools.DirAbsPath
}

// NewResolver
// before create Resolver you must create both origin and transform queue,
// and also origin and transform workers (pools)
// Minimal example:
//
//	originQ := queue.NewQueue()
//	transformQ := queue.NewQueue()
//	defer transformQ.Close()
//	defer originQ.Close()
//	go internal_workers.OriginWorker(originQ, dataDir, tmpDir)
//	go internal_workers.TransformWorker(transformQ, originQ, dataDir, tmpDir)
func NewResolver(transformsQ *queue.Queue, trConf *transforms_conf.Conf, dataDir, tmpDir os_tools.DirAbsPath) *Resolver {
	return &Resolver{
		transformsQ: transformsQ,
		trConf:      trConf,
		dataDir:     dataDir,
		tmpDir:      tmpDir,
	}
}

func (r *Resolver) Request(requestStr string) (os_tools.FileAbsPath, error) {
	req, err := request.RequestFromString(requestStr)
	if err != nil {
		return "", err
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
