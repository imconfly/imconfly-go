package origin

import (
	"crypto/md5"
	"os"
	"path"
)

func WorkerAction(q *Queue, dataDir, tmpDir string) {
	task := q.Get()
	err := Action(task, dataDir, tmpDir)
	q.Done(task.Key, err)
}

func Action(task *Task, dataDir, tmpDir string) error {
	tmpTarget := tmpPath(tmpDir, task.Key)
	if err := get(task.Remote, tmpTarget); err != nil {
		return err
	}

	target := path.Join(dataDir, string(task.Key))
	return os.Rename(tmpTarget, target)
}

func get(url, path string) error {
	return nil
}

func tmpPath(tmpDir string, key Key) string {
	sum := md5.Sum([]byte(key))
	return path.Join(tmpDir, string(sum[:]))
}
