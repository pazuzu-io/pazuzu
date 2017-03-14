package config

import (
	"github.com/urfave/cli"
	"github.com/zalando-incubator/pazuzu"
	"github.com/zalando-incubator/pazuzu/config"
	"fmt"
)

var Command = cli.Command{
	Name:  "config",
	Usage: "Configure pazuzu",
	Subcommands: []cli.Command{
		{
			Name:   "get",
			Usage:  "Get pazuzu configuration",
			Action: getConfig,
		},
		{
			Name:   "set",
			Usage:  "Set pazuzu configuration",
			Action: setConfig,
		},
		{
			Name:   "help",
			Usage:  "Print help on configuration",
			Action: helpConfigs,
		},
		{
			Name:   "list",
			Usage:  "List current effective configuration",
			Action: listConfigs,
		},
	},
}

func setConfig(c *cli.Context) error {
	a := c.Args()
	if len(a) != 2 {
		return pazuzu.ErrTooFewOrManyParameters
	}
	//
	givenPath := a.Get(0)
	givenValRepr := a.Get(1)
	cfg := config.GetConfig()
	cfgMirror := config.GetConfigMirror()
	errSet := cfgMirror.SetConfig(givenPath, givenValRepr)
	if errSet == nil {
		// Oh, it's nice.
		_ = cfg.Save()
		return nil
	}
	return errSet
}

func getConfig(c *cli.Context) error {
	a := c.Args()
	if len(a) != 1 {
		return pazuzu.ErrTooFewOrManyParameters
	}
	//
	givenPath := a.Get(0)
	cfgMirror := config.GetConfigMirror()
	repr, err := cfgMirror.GetRepr(givenPath)
	if err == nil {
		fmt.Println(repr)
		return nil
	}
	return pazuzu.ErrNotFound
}

func helpConfigs(c *cli.Context) error {
	cfgMirror := config.GetConfigMirror()
	fmt.Println("Pazuzu CLI Config related commands:")
	fmt.Println("\tpazuzu config list\t -- Listing of configuration.")
	fmt.Println("\tpazuzu config help\t-- This help documentation.")
	fmt.Println("\tpazuzu config get KEY\t-- Get specific configuration value.")
	fmt.Println("\tpazuzu config set KEY VALUE\t-- Set configuration.")
	fmt.Printf("\nConfiguration keys and its descriptions:\n")
	for _, k := range cfgMirror.GetKeys() {
		help, errHelp := cfgMirror.GetHelp(k)
		if errHelp == nil {
			fmt.Printf("\t%s\t\t%s\n", k, help)
		}
	}
	return nil
}

func listConfigs(c *cli.Context) error {
	cfgMirror := config.GetConfigMirror()
	for _, k := range cfgMirror.GetKeys() {
		repr, errRepr := cfgMirror.GetRepr(k)
		if errRepr == nil {
			fmt.Printf("%s=%s\n", k, repr)
		}
	}
	return nil
}
