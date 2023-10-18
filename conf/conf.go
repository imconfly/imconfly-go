package conf

import (
	"github.com/imconfly/imconfly_go/lib/env_conf"
	o "github.com/imconfly/imconfly_go/lib/os_tools"
	"os"
	"path"
)

const (
	envPrefix             = "IF_"
	defaultConfigFileName = "imconfly.yaml"
	defaultDataDir        = "DATA"
	defaultHost           = "localhost"
	defaultPort           = 80
	defaultTmpDir         = "/tmp/imconfly" // @todo: Windows
)

type Conf struct {
	TransformConcurrency int
	RelativePathsFrom    o.DirAbsPath
	ConfigFile           o.FileAbsPath
	DataDir              o.DirAbsPath
	TmpDir               o.DirAbsPath
	Host                 string
	Port                 int
}

func GetConf() *Conf {
	e := env_conf.New(envPrefix)

	rp := e.Str("RELATIVE_PATHS_FROM", env_conf.Must(os.Getwd()))
	return &Conf{
		TransformConcurrency: e.Int("TRANSFORM_CONCURRENCY", 24), // @todo: default
		RelativePathsFrom:    o.DirAbsPath(rp),
		ConfigFile:           o.FileAbsPath(e.Str("CONFIG_FILE", path.Join(rp, defaultConfigFileName))),
		DataDir:              o.DirAbsPath(e.Str("DATA_DIR", path.Join(rp, defaultDataDir))),
		TmpDir:               o.DirAbsPath(e.Str("TMP_DIR", defaultTmpDir)),
		Host:                 e.Str("HOST", defaultHost),
		Port:                 e.Int("PORT", defaultPort),
	}
}
