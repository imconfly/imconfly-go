package transforms_conf

import (
	"errors"
	"fmt"
	"github.com/imconfly/imconfly_go/lib/os_tools"
	"github.com/imconfly/imconfly_go/queue"
	"gopkg.in/yaml.v2"
	"os"
)

type Container struct {
	Origin     queue.Origin
	Transforms map[string]queue.Transform
}

// Conf - containers, origins and transforms configuration
// All requests have `container/ext_worker/path` type
// if `ext_worker` == "origin" - exactly origin requested
type Conf struct {
	Containers map[string]Container
}

func (c *Conf) ValidateRequest(r *queue.Request) (*queue.Task, error) {
	var container Container
	var origin *queue.Origin
	var transform *queue.Transform

	{
		container, found := c.Containers[r.Container]
		if !found {
			return nil, fmt.Errorf("bad request: container `%s` not exist", r.Container)
		}
		origin = &container.Origin
	}

	if r.IsOrigin() {
		transform = nil
		if !origin.Access {
			return nil, errors.New("forbidden")
		}
	} else {
		t, found := container.Transforms[r.Transform]
		if !found {
			return nil, fmt.Errorf("bad request: transform name `%s` not exist", r.Transform)
		}
		transform = &t
	}

	task := queue.Task{
		Request:   r,
		Origin:    origin,
		Transform: transform,
	}
	return &task, nil
}

func GetConf(filePath os_tools.FileAbsPath) (*Conf, error) {
	b, err := os.ReadFile(string(filePath))
	if err != nil {
		return nil, err
	}

	conf := new(Conf)
	if err := yaml.Unmarshal(b, conf); err != nil {
		return nil, err
	}
	// @todo: check what no "origin" transforms names
	return conf, nil
}
