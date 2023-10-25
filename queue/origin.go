package queue

import (
	"errors"
	"fmt"
	"github.com/imconfly/imconfly_go/lib/os_tools"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

const (
	OriginName        = "origin"
	OriginTypeHTTP    = "HTTP"
	OriginTypeFS      = "Local filesystem"
	OriginTypeUnknown = "_error_"
)

var TypeError = errors.New("wrong origin type")

type Origin struct {
	Source string
	Access bool
}

func (o *Origin) GetType() string {
	l := strings.ToLower(o.Source)
	http := strings.HasPrefix(l, "http://")
	https := strings.HasPrefix(l, "https://")
	file := strings.HasPrefix(l, "/") // @todo: Windows

	if http || https {
		return OriginTypeHTTP
	} else if file {
		return OriginTypeFS
	} else {
		return OriginTypeUnknown
	}
}

func (o *Origin) GetPath(suffix os_tools.FileRelativePath, out *os_tools.FileAbsPath) error {
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
