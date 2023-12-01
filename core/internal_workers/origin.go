package internal_workers

import (
	"fmt"
	"github.com/imconfly/imconfly_go/core/queue"
	o "github.com/imconfly/imconfly_go/lib/os_tools"
	"github.com/imconfly/imconfly_go/lib/tmp_file"
	log "github.com/sirupsen/logrus"
)

func OriginWorker(q *queue.Queue, dataDir, tmpDir o.DirAbsPath) {
	for {
		task := q.Get()
		// queue channel closed
		if task == nil {
			break
		}
		log.Debugf("get origin for task: %s", task.Request.Key)
		err := doOneOriginTask(task, dataDir, tmpDir)
		if err != nil {
			err = fmt.Errorf("get origin error: %w", err)
		}
		q.TaskDone(task.Request.Key, err)
	}
}

func doOneOriginTask(task *queue.Task, dataDir, tmpDir o.DirAbsPath) error {
	if !task.Request.IsOrigin() {
		panic("impossible")
	}

	tmpPath := task.Request.TmpPath(tmpDir)
	targetPath := task.Request.LocalAbsPath(dataDir)

	tmpFile, err := tmp_file.NewTmpFile(tmpPath)
	if err != nil {
		return err
	}
	defer tmpFile.Clean()

	if task.Origin.Transport == "" {
		// Internal HTTP transport
		if err := task.Origin.GetHTTPFile(string(task.Request.PathLastPart), tmpPath); err != nil {
			return err
		}
	} else {
		// User custom transport
		if err := task.Origin.ExecTransport(string(task.Request.PathLastPart), tmpPath); err != nil {
			return err
		}
	}

	return tmpFile.Move(targetPath)
}
