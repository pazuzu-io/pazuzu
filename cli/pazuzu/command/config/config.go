package config

import (
	"github.com/urfave/cli"
	"github.com/pazuzu-io/pazuzu/cli/pazuzu/actions"
)

var Command = cli.Command{
	Name:  "config",
	Usage: "Configure global settings of Pazuzu",
	Subcommands: []cli.Command{
		{
			Name:   "get",
			Usage:  "Get Pazuzu configuration key's value",
			Action: actions.ConfigGet,
		},
		{
			Name:   "set",
			Usage:  "Set Pazuzu configuration key's value",
			Action: actions.ConfigSet,
		},
		{
			Name:   "show",
			Usage:  "Show Pazuzu configuration keys and values",
			Action: actions.ConfigShow,
		},
	},
}
