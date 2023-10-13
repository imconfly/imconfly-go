package transform

import "path"

const OriginName = "origin"

// TaskRequest - info from request only
type TaskRequest struct {
	Container string
	Transform string
	// Path from HTTP request with only "/" slashes
	Path string
}

// GetLocalPath - returns RELATIVE path to local file
func (t *TaskRequest) GetLocalPath() string {
	// @todo: Windows not supported
	return path.Join(t.Container, t.Transform, t.Path)
}
