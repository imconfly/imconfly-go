package os_tools

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
)

func Rename(source, target FileAbsPath) error {
	if exist, err := FileExist(source); err != nil {
		return err
	} else if !exist {
		return fmt.Errorf("file `%s` not exist", source)
	}

	if exist, err := FileExist(target); err != nil {
		return err
	} else if exist {
		return fmt.Errorf("file `%s` already exist", target)
	}

	if err := MkdirFor(target); err != nil {
		return err
	}
	// @todo: invalid cross-device link
	return os.Rename(string(source), string(target))
}

func MkdirFor(pa FileAbsPath) error {
	return os.MkdirAll(path.Dir(string(pa)), os.ModePerm)
}

func FileExist(pa FileAbsPath) (bool, error) {
	stringPath := string(pa)

	if fileInfo, err := os.Stat(stringPath); err == nil {
		if fileInfo.IsDir() {
			return false, fmt.Errorf("this is a directory: %s", stringPath)
		} else {
			return true, nil
		}
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		// floating error: stat xxx: not a directory
		var pathError *fs.PathError
		// @todo
		if errors.As(err, &pathError) && pathError.Err.Error() == "not a directory" {
			return false, nil
		}

		return false, err

	}
}

func Remove(pa FileAbsPath) error {
	return os.RemoveAll(string(pa))
}
