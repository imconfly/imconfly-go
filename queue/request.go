package queue

import (
	"crypto/md5"
	"fmt"
	o "github.com/imconfly/imconfly_go/lib/os_tools"
	"path"
	"strings"
)

type Request struct {
	Container    string
	Transform    string
	PathLastPart o.FileRelativePath
	// Key is also local relative path to files (path.Join(Container, Transform, PathLastPart))
	Key o.FileRelativePath
}

func (r *Request) TmpPath(tmpDir o.DirAbsPath) o.FileAbsPath {
	// container/transform/MD5SUM
	relative := o.FileRelativePath(path.Join(r.Container, r.Transform, md5string([]byte(r.PathLastPart))))
	return tmpDir.FileAbsPath(relative)
}

func (r *Request) LocalAbsPath(dataDir o.DirAbsPath) o.FileAbsPath {
	return dataDir.FileAbsPath(r.Key)
}

func (r *Request) IsOrigin() bool {
	return r.Transform == OriginName
}

func (r *Request) GetOriginRequest() *Request {
	if r.IsOrigin() {
		return r
	}
	return newRequest(r.Container, OriginName, r.PathLastPart)
}

func newRequest(c, t string, p o.FileRelativePath) *Request {
	return &Request{
		Container:    c,
		Transform:    t,
		PathLastPart: p,
		Key:          o.FileRelativePath(path.Join(c, t, string(p))),
	}
}

func RequestFromString(httpGet string, out *Request) error {
	parts := strings.Split(httpGet, "/")
	if len(parts) < 4 {
		return fmt.Errorf("bad request: no `/container/transform/path` pattern in `%s`", httpGet)
	}

	out = newRequest(parts[1], parts[2], "todo") // @todo: o.FileRelativePath(path.Join(...parts[3:]))
	return nil
}

func md5string(b []byte) string {
	sum := md5.Sum(b)
	return string(sum[:])
}
