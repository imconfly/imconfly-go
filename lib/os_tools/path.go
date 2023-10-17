package os_tools

import "path"

type _path string
type absPath _path
type relativePath _path

type FileAbsPath absPath
type DirAbsPath absPath

type FileRelativePath relativePath

func (d DirAbsPath) FileAbsPath(suffix FileRelativePath) FileAbsPath {
	return FileAbsPath(path.Join(string(d), string(suffix)))
}
