package transforms_conf

import (
	"errors"
	"fmt"
	"github.com/imconfly/imconfly_go/lib/os_tools"
	"github.com/imconfly/imconfly_go/queue"
	"gopkg.in/yaml.v2"
	"io"
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
	Containers map[string]*Container
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

// MustGetConf is shortcut for GetConf, uses for tests only
func MustGetConf(confFile os_tools.FileAbsPath) *Conf {
	f, err := os.Open(string(confFile))
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()

	conf, err := GetConf(f)
	if err != nil {
		panic(err.Error())
	}
	return conf
}

func GetConf(reader io.Reader) (*Conf, error) {
	b, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	conf := new(Conf)
	if err := yaml.Unmarshal(b, conf); err != nil {
		return nil, err
	}

	// check what no "origin" transforms names
	for containerName, container := range conf.Containers {
		for transformName, _ := range container.Transforms {
			if transformName == queue.OriginName {
				return nil, fmt.Errorf(
					"transform name cat`t be %q (in %q container)",
					queue.OriginName,
					containerName,
				)
			}
		}

	}

	return conf, nil
}
