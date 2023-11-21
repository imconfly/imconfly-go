package transforms_conf

import (
	"github.com/imconfly/imconfly_go/core/request"
	"github.com/imconfly/imconfly_go/testdata"
	"strings"
	"testing"
)

const testYAML = `
containers:
    wikimedia:
        # https://upload.wikimedia.org/wikipedia/commons/4/41/Inter-Con_Kabul.jpg
        origin:
            source: https://upload.wikimedia.org/wikipedia/commons
            transport: 'curl {source} --output {target}'
            access: true
        transforms:
            dummy: 'cp {source} {target}'
`
const transformRequestString = "/wikimedia/dummy/4/41/Inter-Con_Kabul.jpg"

func TestConf_ValidateRequest(t *testing.T) {
	trConf, err := GetConf(strings.NewReader(testYAML))
	if err != nil {
		t.Fatal(err)
	}
	testdata.LogJSON(t, trConf)

	req, err := request.RequestFromString(transformRequestString)
	if err != nil {
		t.Fatal(err)
	}
	testdata.LogJSON(t, req)

	task, err := trConf.ValidateRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	testdata.LogJSON(t, task)
}
