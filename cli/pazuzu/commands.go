package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
	"github.com/zalando-incubator/pazuzu"
	"regexp"
	"text/tabwriter"
)

var cnfGetCmd = cli.Command{
	Name:  "get",
	Usage: "Get pazuzu configuration",
	Action: func(с *cli.Context) error {
		// log.Print("Getting pazuzu configuration")
		// return nil
		return ErrNotImplemented
	},
}
var cnfSetCmd = cli.Command{
	Name:  "set",
	Usage: "Set pazuzu configuration",
	Action: func(с *cli.Context) error {
		// log.Print("Setting pazuzu configuration")
		// return nil
		return ErrNotImplemented
	},
}

var configCmd = cli.Command{
	Name:  "config",
	Usage: "Configure pazuzu",
	// Action: configure,
	Subcommands: []cli.Command{
		cnfGetCmd,
		cnfSetCmd,
	},
}

var searchCmd = cli.Command{
	Name:      "search",
	Usage:     "search for features in registry",
	ArgsUsage: "[regexp] - Regexp to be used for feature lookup",
	Action: func(c *cli.Context) error {
		sc, err := pazuzu.GetStorageReader(pazuzu.GetConfig())
		if err != nil {
			return err // TODO: process properly into human-readable message
		}

		arg := c.Args().Get(0)
		searchRegexp, err := regexp.Compile(arg)

		if err != nil {
			return fmt.Errorf("could not process search regexp '%s': %s", arg, err.Error())
		}
		features, err := sc.SearchMeta(searchRegexp)
		if err != nil {
			return fmt.Errorf("could not search for features: %s", err.Error())
		}

		if len(features) == 0 {
			fmt.Println("no features found")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight)
		fmt.Fprintf(w, "Name \tAuthor \tDescription\n")
		for _, f := range features {
			fmt.Fprintf(w, "%s \t%s \t%s\n", f.Name, f.Author, f.Description)
		}

		w.Flush()

		return nil
	},
}

var composeCmd = cli.Command{
	Name:        "compose",
	Usage:       "Compose Pazuzufile out of the selected features",
	ArgsUsage:   "[features] - Space separated feature names",
	Description: "Compose step takes list of features as input, validates feature dependencies and creates Pazuzufile.",
	Action: func(с *cli.Context) error {
		return ErrNotImplemented
	},
	// TODO: add -o/--out option according to README file
}

var buildCmd = cli.Command{
	Name:      "build",
	Usage:     "build Dockerfile out of Pazuzufile",
	ArgsUsage: "[features] - This can be either path to Pazuzufile or a space separated feature names",
	Action:    buildFeatures,
}

// Fetches and builds features into a docker image.
func buildFeatures(c *cli.Context) error {
	return ErrNotImplemented
	// TODO: In case of -f/--feature-set option slice of features
	// should be used instead of Pazuzufile

	// log.Print("Building Dockerfile out of the features")
	// // TODO: check number of c.NArgs() and throw error if nothing was passed
	// if c.NArg() == 0 {

	// 	return errors.New(ERROR_NO_VALID_PAZUZU_FILE)

	// }

	// return nil
}
