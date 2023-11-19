package exec

import (
	"errors"
	"strings"
)

const (
	sourcePlaceholder = "{source}"
	targetPlaceholder = "{target}"
)

// WithSourceTarget returns args for go std os/exec.Command()
func WithSourceTarget(cmdStr, source, target string) (string, []string, error) {
	parts := strings.Split(cmdStr, " ")
	if len(parts) < 3 {
		return "", nil, errors.New("minimum cmd template: cmd {source} {target} (3 parts)")
	}
	for k := range parts {
		parts[k] = strings.Replace(parts[k], sourcePlaceholder, source, -1)
		parts[k] = strings.Replace(parts[k], targetPlaceholder, target, -1)
	}
	return parts[0], parts[1:], nil
}
