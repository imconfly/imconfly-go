package transforms_conf

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Conf struct {
	Containers map[string]struct {
		Origin     Origin
		Transforms map[string]Transform
	}
}

type Origin struct {
	Remote string
	Local  string
}

type Transform struct {
	Transform string
	Local     string
}

func GetConf(conf *Conf, confFilePath string) error {
	b, err := os.ReadFile(confFilePath)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(b, conf); err != nil {
		return err
	}
	return nil
}
