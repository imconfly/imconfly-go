package testdata

import (
	"encoding/json"
	"testing"
)

const (
	TestDir     = "/tmp/imconfly_tests"
	TestDataDir = TestDir + "/data"
	TestTmpDir  = TestDir + "/tmp"
)

func MarshalJSON(t *testing.T, value any) string {
	b, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func LogJSON(t *testing.T, value any) {
	t.Log(MarshalJSON(t, value))
}
