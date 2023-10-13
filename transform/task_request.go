package transform

import (
	"errors"
	"path"
	"strings"
)

const OriginName = "origin"

// TaskRequest - info from request only
type TaskRequest struct {
	Container string
	Transform string
	// Path from HTTP request with only "/" slashes
	Path string
	// RELATIVE path to local file
	LocalPath string
}

func (tr *TaskRequest) NewTaskRequest(httpGet string) (*TaskRequest, error) {
	var res *TaskRequest
	parts := strings.Split(httpGet, "/")
	if len(parts) < 4 {
		return nil, errors.New("bad request: no `/container/transform/path` pattern")
	}
	res.Container = parts[1]
	res.Transform = parts[2]
	res.Path = parts[3]
	res.LocalPath = path.Join(
		res.Container,
		res.Transform,
		res.Path,
	)
	return res, nil
}
