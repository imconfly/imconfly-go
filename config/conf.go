package config

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"runtime"

	o "github.com/imconfly/imconfly_go/lib/os_tools"
)

type Conf struct {
	TransformConcurrency int          `yaml:"TransformConcurrency"`
	DataDir              o.DirAbsPath `yaml:"DataDir"`
	TmpDir               o.DirAbsPath `yaml:"TmpDir"`
	ServerHost           string       `yaml:"ServerHost"`
	ServerPort           int          `yaml:"ServerPort"`
	Containers           Containers   `yaml:"Containers"`
}

func ReadConf(reader io.Reader, yamlFormat bool) (*Conf, error) {
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

	if yamlFormat {
		if err = yaml.Unmarshal(b, conf); err != nil {
			return nil, fmt.Errorf("JSON parser error: %w", err)
		}
	} else {
		if err = json.Unmarshal(b, conf); err != nil {
			return nil, fmt.Errorf("JSON parser error: %w", err)
		}
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
