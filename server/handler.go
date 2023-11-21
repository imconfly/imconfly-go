package server

import (
	"github.com/imconfly/imconfly_go/core/queue"
	"github.com/imconfly/imconfly_go/core/request"
	"github.com/imconfly/imconfly_go/lib/os_tools"
	"github.com/imconfly/imconfly_go/transforms_conf"
	"net/http"
)

type Handler struct {
	transformsQ *queue.Queue
	trConf      *transforms_conf.Conf
	dataDir     os_tools.DirAbsPath
	tmpDir      os_tools.DirAbsPath
}

// NewHandler
// before create Handler you must create both origin and transform queue,
// and also origin and transform workers (pools)
// Minimal example:
//
//	originQ := queue.NewQueue()
//	transformQ := queue.NewQueue()
//	defer transformQ.Close()
//	defer originQ.Close()
//	go internal_workers.OriginWorker(originQ, dataDir, tmpDir)
//	go internal_workers.TransformWorker(transformQ, originQ, dataDir, tmpDir)
func NewHandler(transformsQ *queue.Queue, trConf *transforms_conf.Conf, dataDir, tmpDir os_tools.DirAbsPath) *Handler {
	return &Handler{
		transformsQ: transformsQ,
		trConf:      trConf,
		dataDir:     dataDir,
		tmpDir:      tmpDir,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fileAbsPath, err := h.Request(r.RequestURI)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, string(fileAbsPath))
}

func (h *Handler) Request(requestStr string) (os_tools.FileAbsPath, error) {
	req, err := request.RequestFromString(requestStr)
	if err != nil {
		return "", err
	}
	task, err := h.trConf.ValidateRequest(req)
	if err != nil {
		return "", err
	}
	out := task.Request.LocalAbsPath(h.dataDir)
	exist, err := os_tools.FileExist(out)
	if err != nil {
		return "", err
	}
	if exist {
		return out, nil
	}

	err = <-h.transformsQ.Add(task)
	return out, err
}
