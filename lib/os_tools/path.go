package os_tools

import (
	"fmt"
	"os"
	"path"
)

type _path string
type absPath _path
type relativePath _path

type FileAbsPath absPath
type DirAbsPath absPath

type FileRelativePath relativePath

func (d DirAbsPath) FileAbsPath(suffix FileRelativePath) FileAbsPath {
	return FileAbsPath(path.Join(string(d), string(suffix)))
}

func (d DirAbsPath) CheckExist() (bool, error) {
	if fileInfo, err := os.Stat(string(d)); err == nil {
		if fileInfo.IsDir() {
			return true, nil
		} else {
			return false, fmt.Errorf("this is not a directory: %s", d)
		}
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

func (d DirAbsPath) Mkdir() error {
	return os.MkdirAll(path.Dir(string(d)), os.ModePerm)
}
