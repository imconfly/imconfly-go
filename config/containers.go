package config

import (
	"errors"
	"fmt"
	"github.com/imconfly/imconfly_go/core/origin"
	"github.com/imconfly/imconfly_go/core/queue"
	"github.com/imconfly/imconfly_go/core/request"
	"github.com/imconfly/imconfly_go/core/transform"
)

type Container struct {
	Origin     origin.Origin                  `yaml:"Origin"`
	Transforms map[string]transform.Transform `yaml:"Transforms"`
}

// Containers - origins and transforms configuration
// All requests have `container/transform/path` type
// if `transform` == "origin" - exactly origin requested
type Containers map[string]*Container

// Check for correct containers configuration
func (c Containers) Check() error {
	// check what no "origin" transforms names
	for containerName, container := range c {
		for transformName := range container.Transforms {
			if transformName == origin.OriginName {
				return fmt.Errorf(
					"transform name cat`t be %q (in %q container)",
					origin.OriginName,
					containerName,
				)
			}
		}

	}
	return nil
}

// HaveNonLocalOrigins - returns true if originQueue and origin workers needed
func (c Containers) HaveNonLocalOrigins() bool {
	for _, container := range c {
		if container.Origin.GetType() != origin.OriginTypeFS {
			return true
		}
	}
	return false
}

func (c Containers) ValidateRequest(r *request.Request) (*queue.Task, error) {
	var container *Container
	var orig *origin.Origin
	var trans *transform.Transform

	{
		var found bool
		container, found = c[r.Container]
		if !found {
			return nil, fmt.Errorf("container `%s` not exist", r.Container)
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
			return nil, fmt.Errorf("transform name %q not exist", r.Transform)
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
