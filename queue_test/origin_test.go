package queue_test

import (
	"github.com/imconfly/imconfly_go/lib/os_tools"
	"github.com/imconfly/imconfly_go/queue"
	"github.com/imconfly/imconfly_go/testdata"
	"github.com/imconfly/imconfly_go/transforms_conf"
	"os"
	"path"
	"strings"
	"testing"
)

// https://upload.wikimedia.org/wikipedia/commons/4/41/Inter-Con_Kabul.jpg
const trConf1 = `
containers:
    wikimedia:
        origin:
            source: https://upload.wikimedia.org/wikipedia/commons
            transport: 'curl {source} --output {target}'
`

const suffix = "4/41/Inter-Con_Kabul.jpg"

func getOrigin(t *testing.T, yamlString string) *queue.Origin {
	reader := strings.NewReader(yamlString)
	trConf, err := transforms_conf.GetConf(reader)
	if err != nil {
		t.Fatal(err)
	}

	container, ok := trConf.Containers["wikimedia"]
	if !ok {
		t.Fatal("wat?")
	}
	return &container.Origin
}

func TestOrigin_GetType(t *testing.T) {
	origin := getOrigin(t, trConf1)
	originType := origin.GetType()
	if originType != queue.OriginTypeHTTP {
		t.Fatalf("Origin unexpected type: %q", originType)
	}
}

func TestOrigin_GetHTTPFile(t *testing.T) {
	origin := getOrigin(t, trConf1)

	tmpFName := os_tools.FileAbsPath(path.Join(testdata.TestTmpDir, "testFile"))
	if err := os_tools.MkdirFor(tmpFName); err != nil {
		t.Fatal()
	}
	if err := origin.GetHTTPFile(suffix, tmpFName); err != nil {
		t.Fatal(err)
	}

	if err := os.RemoveAll(testdata.TestDir); err != nil {
		t.Fatal()
	}
}

func TestOrigin_ExecTransport(t *testing.T) {
	origin := getOrigin(t, trConf1)

	tmpFName := os_tools.FileAbsPath(path.Join(testdata.TestTmpDir, "testFile"))
	if err := os_tools.MkdirFor(tmpFName); err != nil {
		t.Fatal()
	}
	if err := origin.ExecTransport(suffix, tmpFName); err != nil {
		t.Fatal(err)
	}

	if err := os.RemoveAll(testdata.TestDir); err != nil {
		t.Fatal()
	}
}
