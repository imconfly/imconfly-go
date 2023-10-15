package files

import (
	"errors"
	o "github.com/imconfly/imconfly_go/lib/os_tools"
	"path"
	"strings"
)

const OriginName = "origin"

type File struct {
	Container    string
	Transform    string
	PathLastPart o.FileRelativePath
	// Key is also local relative path to files (path.Join(Container, Transform, PathLastPart))
	Key o.FileRelativePath
}

func (f *File) IsOrigin() bool {
	return f.Transform == OriginName
}

func NewFile(httpGet string, out *File) error {
	parts := strings.Split(httpGet, "/")
	if len(parts) < 4 {
		return errors.New("bad request: no `/container/transform/path` pattern")
	}
	out.Container = parts[1]
	out.Transform = parts[2]
	out.PathLastPart = o.FileRelativePath(parts[3])

	out.Key = o.FileRelativePath(path.Join(
		out.Container,
		out.Transform,
		string(out.PathLastPart),
	))

	return nil
}
