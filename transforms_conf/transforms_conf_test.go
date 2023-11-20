package transforms_conf

import (
	"github.com/imconfly/imconfly_go/queue"
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

	request, err := queue.RequestFromString(transformRequestString)
	if err != nil {
		t.Fatal(err)
	}
	testdata.LogJSON(t, request)

	task, err := trConf.ValidateRequest(request)
	if err != nil {
		t.Fatal(err)
	}
	testdata.LogJSON(t, task)
}
