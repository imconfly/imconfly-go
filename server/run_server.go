package server

import (
	"fmt"
	"github.com/imconfly/imconfly_go/configuration"
	"github.com/imconfly/imconfly_go/core/queue"
	"github.com/imconfly/imconfly_go/internal_workers"
	"github.com/imconfly/imconfly_go/transforms_conf"
	"net/http"
)

func RunServer(conf *configuration.Conf, trConf *transforms_conf.Conf) error {
	transformsQ := queue.NewQueue()
	originQ := queue.NewQueue()
	defer transformsQ.Close()
	defer originQ.Close()

	for i := 0; i < conf.TransformConcurrency; i++ {
		go internal_workers.TransformWorker(transformsQ, originQ, conf.DataDir, conf.TmpDir)
		go internal_workers.OriginWorker(originQ, conf.DataDir, conf.TmpDir)
	}

	imconflyHandler := NewHandler(transformsQ, trConf, conf.DataDir, conf.TmpDir)

	fmt.Printf("Server is listening on %s\n", conf.ServerAddr)
	return http.ListenAndServe(conf.ServerAddr, imconflyHandler)
}
