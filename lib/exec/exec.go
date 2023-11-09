package exec

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

const (
	sourcePlaceholder = "{source}"
	targetPlaceholder = "{target}"
)

// @todo: ctx, env
func Exec(cmdStr, source, target string) error {
	cmdArgs := prepareCmd(cmdStr, source, target)
	fmt.Println(cmdArgs)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	fmt.Println(outb.String())
	fmt.Println(errb.String())
	return err
}

// prepareCmd returns args for go std os/exec.Command()
func prepareCmd(cmdStr, source, target string) []string {
	parts := strings.Split(cmdStr, " ")
	if len(parts) < 3 {
		panic("Minimum cmd template: cmd {source} {target} (3 parts)")
	}
	for k := range parts {
		parts[k] = strings.Replace(parts[k], sourcePlaceholder, source, -1)
		parts[k] = strings.Replace(parts[k], targetPlaceholder, target, -1)
	}
	return parts
}
