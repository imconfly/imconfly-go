package resolver

import (
	"encoding/json"
	"github.com/imconfly/imconfly_go/core/transforms_conf"
	"github.com/imconfly/imconfly_go/testdata"
	"os"
	"path/filepath"
	"testing"
)

const trConfFile = "../../testdata/imconfly.yaml"

// https://upload.wikimedia.org/wikipedia/commons/4/41/Inter-Con_Kabul.jpg
const originRequestString = "/wikimedia/origin/4/41/Inter-Con_Kabul.jpg"
const transformRequestString = "/wikimedia/dummy/4/41/Inter-Con_Kabul.jpg"

func TestExec_originDefaultHTTPTransport(t *testing.T) {
	trConf := getTrConf(t)
	trConf.Containers["wikimedia"].Origin.Transport = ""
	testExec(t, trConf, originRequestString)
}

func TestExec_originCustomTransport(t *testing.T) {
	trConf := getTrConf(t)
	testExec(t, trConf, originRequestString)
}

func TestExec_transform(t *testing.T) {
	trConf := getTrConf(t)
	testExec(t, trConf, transformRequestString)
}

func clean(t *testing.T) {
	t.Logf("Clean: rm dir %s", testdata.TestDir)
	if err := os.RemoveAll(testdata.TestDir); err != nil {
		t.Error(err)
	}
}

func jsonMarshal(v any) string {
	b, _ := json.MarshalIndent(v, "  ", "  ")
	return string(b)
}

func testExec(t *testing.T, trConf *transforms_conf.Conf, request string) {
	defer clean(t)
	t.Logf("TrConf: %+v, %s", trConf, jsonMarshal(trConf))
	t.Log("Exec(), test params are:")
	t.Logf("Request (rStr): %s", request)
	t.Logf("Data dir (dDir): %s", testdata.TestDataDir)
	t.Logf("Tmp dir (tDir): %s", testdata.TestTmpDir)
	var result string
	err := OneRequest(
		request,
		testdata.TestDataDir,
		testdata.TestTmpDir,
		trConf,
		&result,
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("RESULT: %s", result)
}

func getTrConf(t *testing.T) *transforms_conf.Conf {
	trConfAbsPath, err := filepath.Abs(trConfFile)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Get transforms configuration from %s...", trConfAbsPath)

	f, err := os.Open(trConfAbsPath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	c, err := transforms_conf.GetConf(f)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Ok.")
	return c
}
