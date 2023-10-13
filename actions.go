package main

import (
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/imconfly/imconfly_go/conf"
	"github.com/imconfly/imconfly_go/transforms_conf"
	"github.com/imconfly/imconfly_go/version"
)

func runAction(ctx *cli.Context) error {
	fmt.Println("starts HTTP server here")
	return nil
}

func transformAction(ctx *cli.Context) error {
	fmt.Println("ext_worker action here")
	return nil
}

func versionAction(ctx *cli.Context) error {
	fmt.Println(version.VERSION)
	return nil
}

func configAction(ctx *cli.Context) error {
	var c conf.Conf
	if err := conf.GetConf(&c); err != nil {
		return err
	}

	j, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(j))
	return nil
}

func trConfAction(ctx *cli.Context) error {
	var coreConf conf.Conf
	if err := conf.GetConf(&coreConf); err != nil {
		return err
	}

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
