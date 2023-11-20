package cli

import (
	"encoding/json"
	"fmt"
	"github.com/imconfly/imconfly_go/cli/exec"
	"github.com/imconfly/imconfly_go/configuration"
	"github.com/imconfly/imconfly_go/lib/os_tools"
	"github.com/imconfly/imconfly_go/server"
	"github.com/imconfly/imconfly_go/transforms_conf"
	"github.com/imconfly/imconfly_go/version"
	"github.com/urfave/cli/v2"
	"os"
)

func runAction(_ *cli.Context) error {
	conf := configuration.GetConf()
	trConf, err := _getTrConf(conf.ConfigFile)
	if err != nil {
		return err
	}
	return server.RunServer(conf, trConf)
}

func execAction(ctx *cli.Context) error {
	if ctx.Args().Len() != 1 {
		return cli.Exit("Usage: need one arg", 64)
	}
	arg := ctx.Args().First()

	coreConf := configuration.GetConf()
	trConf, err := _getTrConf(coreConf.ConfigFile)
	if err != nil {
		return cli.Exit(err.Error(), 78)
	}

	var target string
	if err := exec.Exec(arg, coreConf.DataDir, coreConf.TmpDir, trConf, &target); err != nil {
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
	c := configuration.GetConf()

	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(b))
	return nil
}

func trConfAction(_ *cli.Context) error {
	coreConf := configuration.GetConf()
	transformsConf, err := _getTrConf(coreConf.ConfigFile)
	if err != nil {
		return err
	}

	j, err := json.MarshalIndent(transformsConf, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(j))
	return nil
}

func _getTrConf(path os_tools.FileAbsPath) (*transforms_conf.Conf, error) {
	f, err := os.Open(string(path))
	if err != nil {
		return nil, cli.Exit(err.Error(), 78)
	}
	defer f.Close()

	trConf, err := transforms_conf.GetConf(f)
	if err != nil {
		return nil, cli.Exit(err.Error(), 78)
	}
	return trConf, nil
}
