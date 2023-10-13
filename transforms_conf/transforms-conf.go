package transforms_conf

import (
	"fmt"
	"github.com/imconfly/imconfly_go/task"
	"gopkg.in/yaml.v2"
	"os"
)

type Container struct {
	Origin     task.Origin
	Transforms map[string]string
}

// Conf - containers, origins and transforms configuration
// All requests have `container/ext_worker/path` type
// if `ext_worker` == "origin" - exactly origin requested
type Conf struct {
	Containers map[string]Container
}

func (c *Conf) NewTask(tReq *task.Request, out *task.Task) error {
	var container Container
	{
		container, found := c.Containers[tReq.Container]
		if !found {
			return fmt.Errorf("bad request: container `%s` not exist", tReq.Container)
		}
		out.Origin = &container.Origin
	}

	{
		if tReq.Transform == task.OriginName {
			out.Transform = nil
		} else {
			t, found := container.Transforms[tReq.Transform]
			if !found {
				return fmt.Errorf("bad request: ext_worker name `%s` not exist", tReq.Transform)
			}
			*out.Transform = task.Transform(t)
		}
	}

	return nil
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
