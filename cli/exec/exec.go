package exec

import (
	"errors"
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

	return do(task, dDir, tDir)
}

func do(t *queue.Task, dDir, tDir o.DirAbsPath) error {
	// currently origins only
	if !t.Request.IsOrigin() {
		return errors.New("origin pls :))")
	}

	oQ := queue.NewQueue()
	go internal_workers.OriginWorker(oQ, dDir, tDir)
	ch := oQ.Add(t)

	return <-ch
}
