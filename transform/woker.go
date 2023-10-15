package transform

import (
	"crypto/md5"
	"github.com/imconfly/imconfly_go/files"
	"github.com/imconfly/imconfly_go/origin"
	"os"
	"path"
)

func Worker(q *Queue, originQueue *origin.Queue, fs *files.FS) {
	task := q.Get()

	// get origin if not exist
	{
		originLocal := path.Join(task.File.Container, files.OriginName, task.File.Path)
		if _, err := os.Stat(path.Join(dataDir, originLocal)); err != nil {
			originTask := &origin.Task{
				Remote: task.Origin.Source,
				Key:    origin.Key(originLocal),
			}
			ch := make(chan error, 1)
			originQueue.Add(originTask, ch)
			if err := <-ch; err != nil {
				q.Done(task.File.Key, err)
				return
			}
		}
	}

	// exactly origin requested
	if task.Transform == nil {
		q.Done(task.File.Key, nil)
		return
	}

	var cmd string
	{
		sum := md5.Sum([]byte(ta.Request.Key))
		target := path.Join(tmpDir, string(sum[:]))
		cmd = makeExecutorString(ta.Transform, originLocalPath, target)
	}

	var err error = execute(cmd)

	q.Done(ta.Request.Key, err)
}

func execute(str string) error {
	return nil
}

func makeExecutorString(transform *Transform, source, target string) string {
	return ""
}
