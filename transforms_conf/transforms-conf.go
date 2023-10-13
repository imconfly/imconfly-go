package transforms_conf

import (
	"fmt"
	"github.com/imconfly/imconfly_go/transform"
	"gopkg.in/yaml.v2"
	"os"
)

type Container struct {
	Origin     transform.Origin
	Transforms map[string]string
}

// Conf - containers, origins and transforms configuration
// All requests have `container/transform/path` type
// if `transform` == "origin" - exactly origin requested
type Conf struct {
	Containers map[string]Container
}

func (c *Conf) GetTransformTask(tReq *transform.TaskRequest) (*transform.Task, error) {
	var origin transform.Origin
	var container Container
	{
		container, found := c.Containers[tReq.Container]
		if !found {
			return nil, fmt.Errorf("bad request: container `%s` not exist", tReq.Container)
		}
		origin = container.Origin
	}

	var tr *string
	{
		if tReq.Transform == transform.OriginName {
			tr = nil
		} else {
			t, found := container.Transforms[tReq.Transform]
			if !found {
				return nil, fmt.Errorf("bad request: transform name `%s` not exist", tReq.Transform)
			}
			tr = &t
		}
	}

	return transform.NewTask(tReq, &origin, tr), nil
}

func GetConf(conf *Conf, confFilePath string) error {
	b, err := os.ReadFile(confFilePath)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(b, conf); err != nil {
		return err
	}
	// @todo: check what no "origin" transforms names
	return nil
}
