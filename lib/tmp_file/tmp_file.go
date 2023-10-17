package tmp_file

import (
	"fmt"
	"github.com/imconfly/imconfly_go/lib/os_tools"
)

type TmpFile struct {
	path  os_tools.FileAbsPath
	moved bool
}

func NewTmpFile(p os_tools.FileAbsPath, out *TmpFile) error {
	out.path = p
	out.moved = false
	return os_tools.MkdirFor(out.path)
}

// Clean
// use defer t.Clean()
func (t *TmpFile) Clean() {
	if t.moved {
		return
	}
	if err := os_tools.Remove(t.path); err != nil {
		panic(fmt.Sprintf("Can`t remove tmp file `%s`: %s", t.path, err))
	}
}

func (t *TmpFile) Move(path os_tools.FileAbsPath) error {
	if err := os_tools.Rename(t.path, path); err != nil {
		return err
	}
	t.moved = true
	return nil
}
