package transforms_conf

import (
	"errors"
	"fmt"
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

func (c *Conf) ValidateRequest(r *queue.Request, out *queue.Task) error {
	var container Container
	var origin *queue.Origin
	var transform *queue.Transform

	{
		container, found := c.Containers[r.Container]
		if !found {
			return fmt.Errorf("bad request: container `%s` not exist", r.Container)
		}
		origin = &container.Origin
	}

	if r.IsOrigin() {
		transform = nil
		if !origin.Access {
			return errors.New("forbidden")
		}
	} else {
		t, found := container.Transforms[r.Transform]
		if !found {
			return fmt.Errorf("bad request: transform name `%s` not exist", r.Transform)
		}
		transform = &t
	}

	out = &queue.Task{
		Request:   r,
		Origin:    origin,
		Transform: transform,
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
