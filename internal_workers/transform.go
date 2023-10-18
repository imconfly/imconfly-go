package internal_workers

import (
	"fmt"
	"github.com/imconfly/imconfly_go/lib/os_tools"
	"github.com/imconfly/imconfly_go/lib/tmp_file"
	"github.com/imconfly/imconfly_go/queue"
)

func TransformWorker(
	transformsQ,
	originQ *queue.Queue,
	dataDir,
	tmpDir os_tools.DirAbsPath,
) {
	task := transformsQ.Get()
	err := doTask(task, originQ, dataDir, tmpDir)
	transformsQ.Done(task.Request.Key, err)
}

func doTask(
	task *queue.Task,
	originQ *queue.Queue,
	dataDir,
	tmpDir os_tools.DirAbsPath,
) error {
	// check origin exist or get it
	var originPath os_tools.FileAbsPath
	if err := getOrigin(task, originQ, dataDir, &originPath); err != nil {
		return err
	}
	// if exactly origin requested - nothing to do
	if task.Request.IsOrigin() {
		return nil
	}
	return doTransform(task, originPath, dataDir, tmpDir)
}

// check origin exist, if not - get it (wait for ready)
func getOrigin(
	t *queue.Task,
	q *queue.Queue,
	dataDir os_tools.DirAbsPath,
	outOriginPath *os_tools.FileAbsPath,
) error {
	request := t.Request.GetOriginRequest()
	*outOriginPath = dataDir.FileAbsPath(t.Request.Key)

	if exist, err := os_tools.FileExist(*outOriginPath); err != nil {
		return err
	} else if exist {
		// origin exist - ok, here is nothing to do
		return nil
	}

	if t.Origin.GetType() == queue.OriginTypeFS {
		return fmt.Errorf("origin not found: `%s`", *outOriginPath)
	}

	// task for origin of current task
	var originTask *queue.Task
	if t.Request.IsOrigin() {
		// if current task is about origin - just get it
		originTask = t
	} else {
		// if not - create origin task from it
		originTask = &queue.Task{
			Request:   request,
			Origin:    t.Origin,
			Transform: nil,
		}
	}

	originReadyCh := q.Add(originTask)
	return <-originReadyCh
}

func doTransform(
	t *queue.Task,
	originPath os_tools.FileAbsPath,
	dataDir,
	tmpDir os_tools.DirAbsPath,
) error {
	targetPath := dataDir.FileAbsPath(t.Request.Key)
	// parallel running maybe
	if exist, err := os_tools.FileExist(targetPath); err != nil {
		return err
	} else if exist {
		return nil
	}

	tmpPath := t.Request.TmpPath(tmpDir)
	var tmpFile *tmp_file.TmpFile
	if err := tmp_file.NewTmpFile(tmpPath, tmpFile); err != nil {
		return err
	}
	defer tmpFile.Clean()

	if err := t.Transform.Execute(originPath, tmpPath); err != nil {
		return err
	}

	return tmpFile.Move(targetPath)
}
