package cli

import (
	"fmt"
	"os"

	"github.com/imconfly/imconfly_go/constants"
)

func Version() {
	fmt.Println(constants.Version)
	os.Exit(0)
}
