package os_tools

import (
	"errors"
	"fmt"
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

	return os.Rename(string(source), string(target))
}

func MkdirFor(pa FileAbsPath) error {
	return os.MkdirAll(path.Dir(string(pa)), os.ModePerm)
}

func FileExist(pa FileAbsPath) (bool, error) {
	_, err := os.Stat(string(pa))
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func Remove(pa FileAbsPath) error {
	return os.RemoveAll(string(pa))
}
