package config

import (
	"encoding/json"
	"fmt"
	"io"
	"runtime"

	o "github.com/imconfly/imconfly_go/lib/os_tools"
)

type Conf struct {
	TransformConcurrency int
	DataDir              o.DirAbsPath
	TmpDir               o.DirAbsPath
	ServerHost           string
	ServerPort           int
	Containers           Containers
}

func ReadConf(reader io.Reader) (*Conf, error) {
	// defaults
	conf := &Conf{
		TransformConcurrency: runtime.NumCPU(),
		DataDir:              "/var/local/imconfly/data",
		TmpDir:               "/var/local/imconfly/tmp",
		ServerHost:           "localhost",
		ServerPort:           8081,
	}

	b, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(b, conf); err != nil {
		return nil, fmt.Errorf("JSON parser error: %w", err)
	}

	// check containers for correct configuration
	if err := conf.Containers.Check(); err != nil {
		return nil, err
	}

	return conf, nil
}

func CheckDirs(dataDir, tmpDir o.DirAbsPath) error {
	if err := checkDir(dataDir); err != nil {
		return fmt.Errorf("data dir error: %s", err)
	}
	if err := checkDir(tmpDir); err != nil {
		return fmt.Errorf("tmp dir error: %s", err)
	}
	// @todo: check mv ability
	return nil
}

func checkDir(dir o.DirAbsPath) error {
	exist, err := dir.CheckExist()
	if err != nil {
		return err
	}
	if exist {
		return nil
	}
	return dir.Mkdir()
}
