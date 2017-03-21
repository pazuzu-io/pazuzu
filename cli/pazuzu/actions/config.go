package actions

import (
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"github.com/zalando-incubator/pazuzu"
	"github.com/zalando-incubator/pazuzu/config"
	"os"
	"text/tabwriter"
)

func ConfigShow(c *cli.Context) error {
	if c.NArg() != 0 {
		return errors.New("Wrong number of arguments")
	}
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintf(writer, "Key\tValue\n")
	cfgMirror := config.GetConfigMirror()
	for _, k := range cfgMirror.GetKeys() {
		repr, errRepr := cfgMirror.GetRepr(k)
		if errRepr == nil {
			fmt.Fprintf(writer, "%s\t%s\n", k, repr)
		}
	}
	writer.Flush()
	return nil
}

func ConfigGet(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.New("Wrong number of arguments")
	}
	key := c.Args().First()
	cfgMirror := config.GetConfigMirror()
	repr, err := cfgMirror.GetRepr(key)
	if err != nil {
		return pazuzu.ErrNotFound
	}
	fmt.Printf("%s => %s\n", key, repr)
	return nil
}

func ConfigSet(c *cli.Context) error {
	if c.NArg() != 2 {
		return errors.New("Wrong number of arguments")
	}
	key := c.Args()[0]
	value := c.Args()[1]
	cfg := config.GetConfig()
	cfgMirror := config.GetConfigMirror()
	err := cfgMirror.SetConfig(key, value)
	if err != nil {
		return errors.New("Can't set config key")
	}
	err = cfg.Save()
	if err != nil {
		return errors.New("Can't save configuration")
	}
	return nil
}
