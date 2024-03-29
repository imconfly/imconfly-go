package internal_workers

import (
	"fmt"
	"github.com/imconfly/imconfly_go/core/origin"
	"github.com/imconfly/imconfly_go/core/queue"
	"github.com/imconfly/imconfly_go/lib/os_tools"
	"github.com/imconfly/imconfly_go/lib/tmp_file"
	log "github.com/sirupsen/logrus"
)

func TransformWorker(
	transformsQ,
	originQ *queue.Queue,
	dataDir,
	tmpDir os_tools.DirAbsPath,
) {
	logName := "TransformWorker()"
	for {
		task := transformsQ.Get()
		// queue channel closed
		if task == nil {
			log.Debugf("%s: queue channel closed. Return.", logName)
			break
		}
		log.Debugf("%s: task recieved from queue: %s", logName, task.Request.Key)
		err := doTask(task, originQ, dataDir, tmpDir)
		if err != nil {
			err = fmt.Errorf("transform error: %w", err)
		}
		transformsQ.TaskDone(task.Request.Key, err)
	}
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
	originQ *queue.Queue,
	dataDir os_tools.DirAbsPath,
	outOriginPath *os_tools.FileAbsPath,
) error {
	if t.Request.IsOrigin() {
		// check file exist
		*outOriginPath = t.Request.LocalAbsPath(dataDir)
		if exist, err := os_tools.FileExist(*outOriginPath); err != nil {
			return err
		} else if exist {
			// origin exist - ok, here is nothing to do
			return nil
		}
		// not exist - add current task in origin queue
		if t.Origin.GetType() == origin.OriginTypeFS {
			return fmt.Errorf("origin not found: `%s`", *outOriginPath)
		}
		if originQ == nil {
			panic("never")
		}
		return <-originQ.Add(t)
	} else {
		originRequest := t.Request.GetOriginRequest()
		*outOriginPath = originRequest.LocalAbsPath(dataDir)
		if exist, err := os_tools.FileExist(*outOriginPath); err != nil {
			return err
		} else if exist {
			// origin exist - ok, here is nothing to do
			return nil
		}
		// not exist - create origin task and add to origin queue
		if t.Origin.GetType() == origin.OriginTypeFS {
			return fmt.Errorf("origin not found: `%s`", *outOriginPath)
		}
		if originQ == nil {
			panic("never")
		}
		originTask := &queue.Task{
			Request:   originRequest,
			Origin:    t.Origin,
			Transform: nil,
		}
		return <-originQ.Add(originTask)
	}
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
	tmpFile, err := tmp_file.NewTmpFile(tmpPath)
	if err != nil {
		return err
	}
	defer tmpFile.Clean()

	if err := t.Transform.Execute(originPath, tmpPath); err != nil {
		return err
	}

	return tmpFile.Move(targetPath)
}
