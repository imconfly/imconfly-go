package transform

import (
	"bytes"
	"fmt"
	"github.com/imconfly/imconfly_go/lib/exec"
	"github.com/imconfly/imconfly_go/lib/os_tools"
	"os"
	oexec "os/exec"
)

type Transform string

func (t Transform) Execute(source, target os_tools.FileAbsPath) error {
	if err := os_tools.MkdirFor(target); err != nil {
		return err
	}

	cmdName, cmdArgs, err := exec.WithSourceTarget(string(t), string(source), string(target))
	if err != nil {
		return err
	}

	cmd := oexec.Command(cmdName, cmdArgs...)
	var outB, errB bytes.Buffer
	cmd.Stdout = &outB
	cmd.Stderr = &errB
	err = cmd.Run()
	fmt.Println(outB.String())
	_, _ = fmt.Fprintln(os.Stderr, errB.String())
	return err
}
