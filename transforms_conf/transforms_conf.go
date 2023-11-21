package transforms_conf

import (
	"errors"
	"fmt"
	"github.com/imconfly/imconfly_go/core/origin"
	"github.com/imconfly/imconfly_go/core/queue"
	"github.com/imconfly/imconfly_go/core/request"
	"github.com/imconfly/imconfly_go/core/transform"
	"gopkg.in/yaml.v2"
	"io"
)

type Container struct {
	Origin     origin.Origin
	Transforms map[string]transform.Transform
}

// Conf - containers, origins and transforms configuration
// All requests have `container/transform/path` type
// if `transform` == "origin" - exactly origin requested
type Conf struct {
	Containers map[string]*Container
}

// HaveNonLocalOrigins - returns true if originQueue and origin workers needed
func (c *Conf) HaveNonLocalOrigins() bool {
	return true // @todo
}

func (c *Conf) ValidateRequest(r *request.Request) (*queue.Task, error) {
	var container *Container
	var orig *origin.Origin
	var trans *transform.Transform

	{
		var found bool
		container, found = c.Containers[r.Container]
		if !found {
			return nil, fmt.Errorf("bad request: container `%s` not exist", r.Container)
		}
		orig = &container.Origin
	}

	if r.IsOrigin() {
		trans = nil
		if !orig.Access {
			return nil, errors.New("forbidden")
		}
	} else {
		t, found := container.Transforms[r.Transform]
		if !found {
			return nil, fmt.Errorf("bad request: transform name %q not exist", r.Transform)
		}
		trans = &t
	}

	task := queue.Task{
		Request:   r,
		Origin:    orig,
		Transform: trans,
	}
	return &task, nil
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
			if transformName == origin.OriginName {
				return nil, fmt.Errorf(
					"transform name cat`t be %q (in %q container)",
					origin.OriginName,
					containerName,
				)
			}
		}

	}

	return conf, nil
}
