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
	var r *queue.Request
	if err := queue.RequestFromString(rStr, r); err != nil {
		return err
	}

	var t *queue.Task
	if err := trConf.ValidateRequest(r, t); err != nil {
		return err
	}

	target := dDir.FileAbsPath(r.Key)
	*out = string(target)
	if found, err := o.FileExist(target); err != nil {
		return err
	} else if found {
		return nil
	}

	return do(t, dDir, tDir)
}

func do(t *queue.Task, dDir, tDir o.DirAbsPath) error {
	// currently origins only
	if !t.Request.IsOrigin() {
		return errors.New("origin pls :))")
	}

	oQ := queue.NewQueue()
	ch := oQ.Add(t)
	internal_workers.OriginWorker(oQ, dDir, tDir)

	return <-ch
}
