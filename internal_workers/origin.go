package internal_workers

import (
	o "github.com/imconfly/imconfly_go/lib/os_tools"
	"github.com/imconfly/imconfly_go/lib/tmp_file"
	"github.com/imconfly/imconfly_go/queue"
)

func OriginWorker(q *queue.Queue, dataDir, tmpDir o.DirAbsPath) {
	task := q.Get()
	err := Action(task, dataDir, tmpDir)
	q.Done(task.Request.Key, err)
}

func Action(t *queue.Task, dataDir, tmpDir o.DirAbsPath) error {
	if !t.Request.IsOrigin() {
		panic("wat?")
	}

	tmpPath := t.Request.TmpPath(tmpDir)
	targetPath := t.Request.LocalAbsPath(dataDir)

	var tmpFile *tmp_file.TmpFile
	if err := tmp_file.NewTmpFile(tmpPath, tmpFile); err != nil {
		return err
	}
	defer tmpFile.Clean()

	if err := t.Origin.GetHTTPFile(string(t.Request.Key), tmpPath); err != nil {
		return err
	}

	return tmpFile.Move(targetPath)
}