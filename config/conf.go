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

	// @todo: check

	return conf, nil
}
