package configuration

import (
	"github.com/imconfly/imconfly_go/lib/env_conf"
	o "github.com/imconfly/imconfly_go/lib/os_tools"
	"os"
	"path"
	"runtime"
)

const (
	envPrefix             = "IF_"
	defaultConfigFileName = "imconfly.yaml"
	defaultDataDir        = "DATA"
	defaultServerAddr     = "localhost:8081"
	defaultTmpDir         = "/tmp/imconfly" // @todo: Windows
)

type Conf struct {
	TransformConcurrency int
	RelativePathsFrom    o.DirAbsPath
	ConfigFile           o.FileAbsPath
	DataDir              o.DirAbsPath
	TmpDir               o.DirAbsPath
	ServerAddr           string
}

func GetConf() *Conf {
	e := env_conf.New(envPrefix)

	rp := e.Str("RELATIVE_PATHS_FROM", env_conf.Must(os.Getwd()))

	return &Conf{
		TransformConcurrency: e.Int("TRANSFORM_CONCURRENCY", runtime.NumCPU()),
		RelativePathsFrom:    o.DirAbsPath(rp),
		ConfigFile:           o.FileAbsPath(e.Str("CONFIG_FILE", path.Join(rp, defaultConfigFileName))),
		DataDir:              o.DirAbsPath(e.Str("DATA_DIR", path.Join(rp, defaultDataDir))),
		TmpDir:               o.DirAbsPath(e.Str("TMP_DIR", defaultTmpDir)),
		ServerAddr:           e.Str("SERVER_ADDR", defaultServerAddr),
	}
}
