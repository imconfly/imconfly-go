package transform

import (
	"github.com/imconfly/imconfly_go/lib/exec"
	"github.com/imconfly/imconfly_go/lib/os_tools"
	log "github.com/sirupsen/logrus"
	oexec "os/exec"
	"strings"
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

	log.Debugf("exec command: %s %s", cmdName, strings.Join(cmdArgs, " "))
	cmd := oexec.Command(cmdName, cmdArgs...)
	// var outB, errB bytes.Buffer
	// cmd.Stdout = &outB
	// cmd.Stderr = &errB
	err = cmd.Run()
	// fmt.Println(outB.String())
	// _, _ = fmt.Fprintln(os.Stderr, errB.String())
	return err
}
