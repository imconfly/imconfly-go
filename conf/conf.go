package conf

import (
	"os"
	"path"
)

const configFileName = "imconfly.yaml"

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
		conf.RelativePathsFrom = os.Getenv("IF_RELATIVE_PATHS_FROM")
		if conf.RelativePathsFrom == "" {
			var err error
			conf.RelativePathsFrom, err = os.Getwd()
			if err != nil {
				return err
			}
		}
	}
	conf.ConfigFile = path.Join(conf.RelativePathsFrom, configFileName)
	conf.DataDir = path.Join(conf.RelativePathsFrom, "DATA")
	conf.Host = "localhost"
	conf.Port = 80

	return nil
}
