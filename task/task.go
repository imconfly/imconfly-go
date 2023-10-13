package task

import (
	"errors"
	"path"
	"strings"
)

const OriginName = "origin"

type Origin struct {
	Source string
	Access bool
}

type Transform string

type Request struct {
	Container string
	Transform string
	Path      string
	// Key is also local relative path to file
	Key Key
}

type Key string

type Task struct {
	Request   *Request
	Origin    *Origin
	Transform *Transform
}

func NewRequest(httpGet string, out *Request) error {
	parts := strings.Split(httpGet, "/")
	if len(parts) < 4 {
		return errors.New("bad request: no `/container/ext_worker/path` pattern")
	}
	out.Container = parts[1]
	out.Transform = parts[2]
	out.Path = parts[3]

	out.Key = Key(path.Join(
		out.Container,
		out.Transform,
		out.Path,
	))

	return nil
}
