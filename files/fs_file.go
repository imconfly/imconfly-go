package files

import (
	"crypto/md5"
	"fmt"
	"github.com/imconfly/imconfly_go/lib/os_tools"
	"path"
)

type FS struct {
	DataDir os_tools.DirAbsPath
	TmpDir  os_tools.DirAbsPath
}

// FsFile - file in fs adventures ))
// 1. fsf := NewFsFile(..)
// 2. defer fsf.Clean() // removes tmp file in some conditions
// 3. exist, err := fsf.TargetExist() // if exist - nothing to do (just serve)
// 4. fsf.PrepareTmp(&tmpName)
// 5. (external `tmpName` creation)
// 6. fsf.MoveTmpToTarget()
type FsFile struct {
	f       *File
	fs      *FS
	tmpPath os_tools.FileAbsPath
}

func (fsf *FsFile) targetAbsPath() os_tools.FileAbsPath {
	return os_tools.FileAbsPath(
		path.Join(
			string(fsf.fs.DataDir),
			string(fsf.f.Key),
		),
	)
}

func (fsf *FsFile) TargetExist() (bool, error) {
	return os_tools.FileExist(fsf.targetAbsPath())
}

// PrepareTmp builds temp filename and creates dir (-p) for this file
func (fsf *FsFile) PrepareTmp(out *os_tools.FileAbsPath) error {
	sum := md5.Sum([]byte(fsf.f.Key))
	result := os_tools.FileAbsPath(
		path.Join(
			string(fsf.fs.TmpDir),
			fsf.f.Container,
			fsf.f.Transform,
			string(sum[:]),
		),
	)

	if err := os_tools.MkdirFor(result); err != nil {
		return err
	}

	fsf.tmpPath = result
	out = &result
	return nil
}

// Clean - use `defer fsf.Clean()`
// removes tmp file, which lives when some error occurs between
// PrepareTmp() and MoveTmpToTarget()
func (fsf *FsFile) Clean() {
	if fsf.tmpPath != "" {
		if err := os_tools.Remove(fsf.tmpPath); err != nil {
			panic(fmt.Sprintf("Can`t remove tmp file `%s`: %s", fsf.tmpPath, err))
		}
	}
}

func (fsf *FsFile) MoveTmpToTarget() error {
	if fsf.tmpPath == "" {
		panic("wat?")
	}
	if err := os_tools.Rename(fsf.tmpPath, fsf.targetAbsPath()); err != nil {
		return err
	}

	fsf.tmpPath = ""
	return nil
}

func NewFsFile(f *File, fs *FS) *FsFile {
	return &FsFile{
		f:  f,
		fs: fs,
	}
}
