package actions

import (
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"github.com/zalando-incubator/pazuzu/config"
	"os"
	"regexp"
	"text/tabwriter"
)

func Search(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.New("Search regexp is not provided")
	}
	feature := c.Args().First()
	featureRegexp, err := regexp.Compile(feature)
	if err != nil {
		return errors.New("Can't compile search regexp")
	}

	cfg := config.GetConfig()
	storage, err := config.GetStorageReader(*cfg)
	if err != nil {
		return errors.New("Can't create storage reader")
	}
	features, err := storage.SearchMeta(featureRegexp)
	if err != nil {
		return errors.New("Can't execute search request")
	}
	if len(features) == 0 {
		fmt.Println("No features found")
		return nil
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintf(writer, "Name\tAuthor\tDescription\n")
	for _, f := range features {
		fmt.Fprintf(writer, "%s\t%s\t%s\n", f.Name, f.Author, f.Description)
	}
	writer.Flush()
	return nil
}
