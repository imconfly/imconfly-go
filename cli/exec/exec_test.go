package exec

import (
	"github.com/imconfly/imconfly_go/lib/os_tools"
	"github.com/imconfly/imconfly_go/transforms_conf"
	"os"
	"path/filepath"
	"testing"
)

const trConfFile = "../../testdata/imconfly.yaml"
const requestString = "/wikimedia/origin/4/41/Inter-Con_Kabul.jpg"
const testDir = "/tmp/imconfly_tests"
const dataDir = testDir + "/data"
const tmpDir = testDir + "/tmp"

func TestExec(t *testing.T) {
	trConf := getTrConf(t)
	t.Logf("TrConf: %+v", trConf)
	t.Log("Exec(), test params are:")
	t.Logf("Request (rStr): %s", requestString)
	t.Logf("Data dir (dDir): %s", dataDir)
	t.Logf("Tmp dir (tDir): %s", tmpDir)
	var result string
	defer os.RemoveAll(testDir)
	err := Exec(
		requestString,
		dataDir,
		tmpDir,
		trConf,
		&result,
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("RESULT: %s", result)
}

func getTrConf(t *testing.T) *transforms_conf.Conf {
	trConfPath, err := filepath.Abs(trConfFile)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Get transforms conf from %s...", trConfPath)
	c, err := transforms_conf.GetConf(os_tools.FileAbsPath(trConfPath))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Ok.")
	return c
}
