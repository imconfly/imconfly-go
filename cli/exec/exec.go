package exec

import (
	"github.com/imconfly/imconfly_go/internal_workers"
	o "github.com/imconfly/imconfly_go/lib/os_tools"
	"github.com/imconfly/imconfly_go/queue"
	"github.com/imconfly/imconfly_go/transforms_conf"
)

func Exec(
	rStr string,
	dDir,
	tDir o.DirAbsPath,
	trConf *transforms_conf.Conf,
	out *string,
) error {
	r, err := queue.RequestFromString(rStr)
	if err != nil {
		return err
	}

	task, err := trConf.ValidateRequest(r)
	if err != nil {
		return err
	}

	target := dDir.FileAbsPath(r.Key)
	*out = string(target)
	if found, err := o.FileExist(target); err != nil {
		return err
	} else if found {
		return nil
	}

	return doOneTask(task, dDir, tDir)
}

func doOneTask(task *queue.Task, dataDir, tmpDir o.DirAbsPath) error {
	originQ := queue.NewQueue()
	transformQ := queue.NewQueue()
	defer transformQ.Close()
	defer originQ.Close()

	go internal_workers.OriginWorker(originQ, dataDir, tmpDir)
	go internal_workers.TransformWorker(transformQ, originQ, dataDir, tmpDir)

	return <-transformQ.Add(task)
}
