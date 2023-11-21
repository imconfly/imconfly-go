package queue

import (
	"github.com/imconfly/imconfly_go/lib/os_tools"
	"github.com/imconfly/imconfly_go/testdata"
	"os"
	"path"
	"testing"
)

const (
	// https://upload.wikimedia.org/wikipedia/commons/4/41/Inter-Con_Kabul.jpg
	originTestSource        = "https://upload.wikimedia.org/wikipedia/commons"
	originTestTransport     = "curl {source} --output {target}"
	originTestRequestSuffix = "4/41/Inter-Con_Kabul.jpg"
)

var origin = &Origin{
	Source:    originTestSource,
	Transport: originTestTransport,
}

func TestOrigin_GetType(t *testing.T) {
	originType := origin.GetType()
	if originType != OriginTypeHTTP {
		t.Fatalf("Origin unexpected type: %q", originType)
	}
}

func TestOrigin_GetHTTPFile(t *testing.T) {
	tmpFName := os_tools.FileAbsPath(path.Join(testdata.TestTmpDir, "testFile"))
	if err := os_tools.MkdirFor(tmpFName); err != nil {
		t.Fatal()
	}
	if err := origin.GetHTTPFile(originTestRequestSuffix, tmpFName); err != nil {
		t.Fatal(err)
	}

	if err := os.RemoveAll(testdata.TestDir); err != nil {
		t.Fatal()
	}
}

func TestOrigin_ExecTransport(t *testing.T) {
	tmpFName := os_tools.FileAbsPath(path.Join(testdata.TestTmpDir, "testFile"))
	if err := os_tools.MkdirFor(tmpFName); err != nil {
		t.Fatal()
	}
	if err := origin.ExecTransport(originTestRequestSuffix, tmpFName); err != nil {
		t.Fatal(err)
	}

	if err := os.RemoveAll(testdata.TestDir); err != nil {
		t.Fatal()
	}
}
