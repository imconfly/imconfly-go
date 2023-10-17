package queue

import (
	"github.com/imconfly/imconfly_go/lib/os_tools"
)

type Transform string

func (t Transform) Execute(source, target os_tools.FileAbsPath) error {
	return nil
}
