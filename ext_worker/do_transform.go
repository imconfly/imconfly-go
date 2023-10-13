package ext_worker

import (
	"github.com/imconfly/imconfly_go/queue"
	"github.com/imconfly/imconfly_go/task"
	"github.com/imconfly/imconfly_go/transforms_conf"
	"os"
	"path"
)

func WaitTransform(
	req *task.Request,
	trConf *transforms_conf.Conf,
	dataDir string,
	q *queue.Queue,
) error {
	// no wait if file exist
	{
		localFilePath := path.Join(dataDir, string(req.Key))
		if _, err := os.Stat(localFilePath); err == nil {
			return nil
		}
	}

	var ta *task.Task
	if err := trConf.NewTask(req, ta); err != nil {
		return err
	}

	resultChan := make(chan error, 1)

	q.Add(ta, resultChan)

	return <-resultChan
}
