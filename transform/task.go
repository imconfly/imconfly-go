package transform

import (
	"github.com/imconfly/imconfly_go/files"
)

type Origin struct {
	Source string
	Access bool
}

type Transform string

type Task struct {
	File      *files.File
	Origin    *Origin
	Transform *Transform
}
