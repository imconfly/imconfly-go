package cli

import (
	"encoding/json"
	"fmt"
	"github.com/imconfly/imconfly_go/conf"
	"github.com/imconfly/imconfly_go/transforms_conf"
	"github.com/imconfly/imconfly_go/version"
	"github.com/urfave/cli/v2"
)

func runAction(_ *cli.Context) error {
	fmt.Println("starts HTTP server here")
	return nil
}

func transformAction(_ *cli.Context) error {
	fmt.Println("ext_worker action here")
	return nil
}

func versionAction(_ *cli.Context) error {
	fmt.Println(version.Version)
	return nil
}

func configAction(_ *cli.Context) error {
	c := conf.GetConf()

	j, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(j))
	return nil
}

func trConfAction(_ *cli.Context) error {
	coreConf := conf.GetConf()
	var c transforms_conf.Conf
	if err := transforms_conf.GetConf(&c, coreConf.ConfigFile); err != nil {
		return err
	}

	j, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(j))
	return nil
}
