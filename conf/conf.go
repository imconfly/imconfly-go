package conf

import (
	"os"
	"path"
)

const envPrefix = "IF_"

const defaultConfigFileName = "imconfly.yaml"
const defaultDataDir = "DATA"
const defaultHost = "localhost"
const defaultPort = 80

type Conf struct {
	TransformConcurrency int
	RelativePathsFrom    string
	ConfigFile           string
	DataDir              string
	Host                 string
	Port                 int
}

func GetConf(conf *Conf) error {
	{
		// @todo
		conf.TransformConcurrency = 24
	}
	{
		conf.RelativePathsFrom = os.Getenv(envPrefix + "RELATIVE_PATHS_FROM")
		if conf.RelativePathsFrom == "" {
			var err error
			conf.RelativePathsFrom, err = os.Getwd()
			if err != nil {
				return err
			}
		}
	}
	conf.ConfigFile = path.Join(conf.RelativePathsFrom, defaultConfigFileName)
	conf.DataDir = path.Join(conf.RelativePathsFrom, defaultDataDir)
	conf.Host = defaultHost
	conf.Port = defaultPort

	return nil
}
