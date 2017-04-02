package actions

import (
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"github.com/zalando-incubator/pazuzu/config"
	"github.com/zalando-incubator/pazuzu/shared"
	storage "github.com/zalando-incubator/pazuzu/storageconnector"
	"os"
	"regexp"
	"text/tabwriter"
)

func Search(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.New("Search regexp is not provided")
	}
	featureName := c.Args().First()
	cfg := config.GetConfig()
	storage, err := config.GetStorageReader(*cfg)
	if err != nil {
		return errors.New("Can't create storage reader")
	}

	features, err := SearchHandler(featureName, storage)
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

func SearchHandler(feature string, storage storage.StorageReader) ([]shared.FeatureMeta, error) {
	featureRegexp, err := regexp.Compile(feature)
	if err != nil {
		return nil, errors.New("Can't compile search regexp")
	}

	features, err := storage.SearchMeta(featureRegexp)
	if err != nil {
		return nil, errors.New("Can't execute search request")
	}

	return features, nil
}
