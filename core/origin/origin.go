package origin

import (
	"errors"
	"fmt"
	"github.com/imconfly/imconfly_go/lib/exec"
	"github.com/imconfly/imconfly_go/lib/os_tools"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	oexec "os/exec"
	"path"
	"strings"
)

const (
	OriginName        = "origin"
	OriginTypeHTTP    = "HTTP"
	OriginTypeFS      = "LOCAL_FS"
	OriginTypeUnknown = "_error_"
)

var TypeError = errors.New("wrong origin type")

type Origin struct {
	Source    string `yaml:"Source"`
	Transport string `yaml:"Transport"`
	Access    bool   `yaml:"Access"`
}

func (o *Origin) GetType() string {
	l := strings.ToLower(o.Source)
	httpT := strings.HasPrefix(l, "http://")
	https := strings.HasPrefix(l, "https://")
	file := strings.HasPrefix(l, "/") // @todo: Windows

	if httpT || https {
		return OriginTypeHTTP
	} else if file {
		return OriginTypeFS
	} else {
		return OriginTypeUnknown
	}
}

func (o *Origin) GetTypeFsPath(suffix os_tools.FileRelativePath, out *os_tools.FileAbsPath) error {
	if o.GetType() != OriginTypeFS {
		return TypeError
	}
	*out = os_tools.FileAbsPath(path.Join(o.Source, string(suffix)))
	return nil
}

func (o *Origin) GetHTTPFile(suffix string, tmpPath os_tools.FileAbsPath) error {
	if o.GetType() != OriginTypeHTTP {
		return TypeError
	}
	url := o.Source + "/" + suffix
	fmt.Println("GET URL:", url)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.Create(string(tmpPath))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}

func (o *Origin) ExecTransport(suffix string, tmpPath os_tools.FileAbsPath) error {
	sourceURLorPath := o.Source + "/" + suffix

	if err := os_tools.MkdirFor(tmpPath); err != nil {
		return err
	}

	cmdName, cmdArgs, err := exec.WithSourceTarget(o.Transport, sourceURLorPath, string(tmpPath))
	if err != nil {
		return err
	}
	log.Debugf("exec command: %s %s", cmdName, strings.Join(cmdArgs, " "))

	cmd := oexec.Command(cmdName, cmdArgs...)
	//var outB, errB bytes.Buffer
	//cmd.Stdout = &outB
	//cmd.Stderr = &errB
	err = cmd.Run()
	//fmt.Println(outB.String())
	//_, _ = fmt.Fprintln(os.Stderr, errB.String())
	return err
}
