package tmp_file

import (
	"github.com/imconfly/imconfly_go/lib/os_tools"
)

type TmpFile struct {
	path  os_tools.FileAbsPath
	moved bool
}

func NewTmpFile(p os_tools.FileAbsPath) (*TmpFile, error) {
	out := TmpFile{
		path:  p,
		moved: false,
	}
	return &out, os_tools.MkdirFor(out.path)
}

// Clean
// use defer t.Clean()
func (t *TmpFile) Clean() error {
	if t.moved {
		return nil
	}
	return os_tools.Remove(t.path)
}

func (t *TmpFile) Move(path os_tools.FileAbsPath) error {
	if err := os_tools.Rename(t.path, path); err != nil {
		return err
	}
	t.moved = true
	return nil
}
