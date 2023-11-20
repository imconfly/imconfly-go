package server

import (
	"fmt"
	"github.com/imconfly/imconfly_go/configuration"
	"github.com/imconfly/imconfly_go/internal_workers"
	"github.com/imconfly/imconfly_go/lib/os_tools"
	"github.com/imconfly/imconfly_go/queue"
	"github.com/imconfly/imconfly_go/transforms_conf"
	"net/http"
)

type core struct {
	transformsQ *queue.Queue
	trConf      *transforms_conf.Conf
	dataDir     os_tools.DirAbsPath
	tmpDir      os_tools.DirAbsPath
}

func (c *core) Request(requestStr string) (os_tools.FileAbsPath, error) {
	request, err := queue.RequestFromString(requestStr)
	if err != nil {
		return "", err
	}
	task, err := c.trConf.ValidateRequest(request)
	if err != nil {
		return "", err
	}
	out := task.Request.LocalAbsPath(c.dataDir)
	exist, err := os_tools.FileExist(out)
	if err != nil {
		return "", err
	}
	if exist {
		return out, nil
	}

	err = <-c.transformsQ.Add(task)
	return out, err
}

func RunServer(conf *configuration.Conf, trConf *transforms_conf.Conf) error {
	transformsQ := queue.NewQueue()
	originQ := queue.NewQueue()
	defer transformsQ.Close()
	defer originQ.Close()

	for i := 0; i < conf.TransformConcurrency; i++ {
		go internal_workers.TransformWorker(transformsQ, originQ, conf.DataDir, conf.TmpDir)
		go internal_workers.OriginWorker(originQ, conf.DataDir, conf.TmpDir)
	}

	c := &core{
		transformsQ: transformsQ,
		trConf:      trConf,
		dataDir:     conf.DataDir,
		tmpDir:      conf.TmpDir,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fileAbsPath, err := c.Request(r.RequestURI)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.ServeFile(w, r, string(fileAbsPath))
	})

	fmt.Printf("Server is listening on %s\n", conf.ServerAddr)
	return http.ListenAndServe(conf.ServerAddr, nil)
}
