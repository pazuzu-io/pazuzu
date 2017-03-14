package search

import (
	"github.com/urfave/cli"
	"github.com/zalando-incubator/pazuzu/config"
	"regexp"
	"fmt"
	"text/tabwriter"
	"os"
)

var Command = cli.Command{
	Name:      "search",
	Usage:     "search for features in registry",
	ArgsUsage: "[regexp] - Regexp to be used for feature lookup",
	Action: searchAction,
}

var searchAction = func(c *cli.Context) error {
	sc, err := config.GetStorageReader(*config.GetConfig())
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
}

