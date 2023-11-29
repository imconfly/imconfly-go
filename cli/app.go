package cli

import "github.com/urfave/cli/v2"

var App = &cli.App{
	Name:                 "imconfly_go",
	Usage:                "Web server for on-the-fly data conversions",
	EnableBashCompletion: true,
	Commands: []*cli.Command{
		{
			Name:   "run",
			Usage:  "run HTTP server",
			Action: runAction,
		},
		{
			Name:      "exec",
			Usage:     "works like HTTP query but print target filename",
			ArgsUsage: "requestString - /container/transform/path string",
			Action:    execAction,
		},
		{
			Name:   "version",
			Usage:  "print version",
			Action: versionAction,
		},
		{
			Name:   "conf",
			Usage:  "print config in JSON format",
			Action: configAction,
		},
		{
			Name:   "tr-conf",
			Usage:  "print transforms config in JSON format",
			Action: trConfAction,
		},
	},
}
