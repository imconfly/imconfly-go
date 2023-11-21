package exec

import (
	"github.com/imconfly/imconfly_go/internal_workers"
	o "github.com/imconfly/imconfly_go/lib/os_tools"
	"github.com/imconfly/imconfly_go/queue"
	"github.com/imconfly/imconfly_go/server"
	"github.com/imconfly/imconfly_go/transforms_conf"
)

func DoOneTask(
	requestString string,
	dataDir,
	tmpDir o.DirAbsPath,
	trConf *transforms_conf.Conf,
	out *string,
) error {
	// Prepare to start Handler
	originQ := queue.NewQueue()
	transformQ := queue.NewQueue()
	defer transformQ.Close()
	defer originQ.Close()
	go internal_workers.OriginWorker(originQ, dataDir, tmpDir)
	go internal_workers.TransformWorker(transformQ, originQ, dataDir, tmpDir)

	imconflyHandler := server.NewHandler(transformQ, trConf, dataDir, tmpDir)
	targetFileAbsPath, err := imconflyHandler.Request(requestString)
	*out = string(targetFileAbsPath)

	return err
}
