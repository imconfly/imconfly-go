package resolver

import (
	"github.com/imconfly/imconfly_go/core/queue"
	"github.com/imconfly/imconfly_go/internal_workers"
	o "github.com/imconfly/imconfly_go/lib/os_tools"
	"github.com/imconfly/imconfly_go/transforms_conf"
)

func OneRequest(
	requestString string,
	dataDir,
	tmpDir o.DirAbsPath,
	trConf *transforms_conf.Conf,
	out *string,
) error {
	// Prepare to start Resolver
	transformQ := queue.NewQueue()
	{
		var originQ *queue.Queue
		if trConf.HaveNonLocalOrigins() {
			originQ = queue.NewQueue()
			defer originQ.Close()
			go internal_workers.OriginWorker(originQ, dataDir, tmpDir)
		}
		defer transformQ.Close()
		go internal_workers.TransformWorker(transformQ, originQ, dataDir, tmpDir)
	}

	resolver := NewResolver(transformQ, trConf, dataDir, tmpDir)
	targetFileAbsPath, err := resolver.Request(requestString)
	*out = string(targetFileAbsPath)

	return err
}
