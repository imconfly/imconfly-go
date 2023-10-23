package cli

import (
	"encoding/json"
	"fmt"
	"github.com/imconfly/imconfly_go/cli/exec"
	"github.com/imconfly/imconfly_go/conf"
	"github.com/imconfly/imconfly_go/transforms_conf"
	"github.com/imconfly/imconfly_go/version"
	"github.com/urfave/cli/v2"
)

func runAction(_ *cli.Context) error {
	fmt.Println("starts HTTP server here")
	return nil
}

func execAction(ctx *cli.Context) error {
	if ctx.Args().Len() != 1 {
		return cli.Exit("Usage: need one arg", 64)
	}
	arg := ctx.Args().First()
	//fmt.Printf("%s\n", arg)

	c := conf.GetConf()
	var trC *transforms_conf.Conf
	if err := transforms_conf.GetConf(trC, c.ConfigFile); err != nil {
		return cli.Exit(err.Error(), 78)
	}

	var target string
	if err := exec.Exec(arg, c.DataDir, c.TmpDir, trC, &target); err != nil {
		return err
	}

	fmt.Println(target)
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
