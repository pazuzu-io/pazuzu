package main

import (
	"bufio"
	"fmt"
	"log"
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
	Action: func(c *cli.Context) error {
		config := pazuzu.GetConfig()
		sc, err := pazuzu.GetStorageReader(config)
		if err != nil {
			return err // TODO: process properly into human-readable message
		}

		var features []string

		// Check if feature actually exists in repository
		for _, v := range c.Args() {
			log.Printf("Checking: %v\n", v)

			_, err := sc.GetMeta(v)
			if err != nil {
				log.Printf("could not find feature \"%v\" in repository.", v)
				return err
			}
			features = append(features, fmt.Sprintf("%v", v))

		}

		log.Printf("features: %v", features)

		f, err := os.Create("Pazuzufile")
		if err != nil {
			log.Print("could not create Pazuzufile")
			return err
		}

		defer f.Close()
		w := bufio.NewWriter(f)

		pazuzu.Write(w, pazuzu.PazuzuFile{
			Base:     config.Base,
			Features: features})

		w.Flush()

		return nil
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
	// TODO: Make file configurable via CLI args (GH Issue #102)
	fileName := "Pazuzufile"

	file, err := os.Open(fileName)
	if err != nil {
		log.Printf("Cannot find %s", fileName)
		return err
	}
	defer file.Close()

	config := pazuzu.GetConfig()
	storageReader, err := pazuzu.GetStorageReader(config)

	reader := bufio.NewReader(file)
	pazuzuFile, err := pazuzu.Read(reader)

	p := pazuzu.Pazuzu{StorageReader: storageReader}
	p.Generate(pazuzuFile.Base, pazuzuFile.Features)

	f, err := os.Create("Dockerfile")
	if err != nil {
		log.Print("could not create Dockerfile")
		return err
	}

	defer f.Close()

	w := bufio.NewWriter(f)
	w.Write(p.Dockerfile)
	w.Flush()

	return nil

	// TODO: In case of -f/--feature-set option slice of features
	// should be used instead of Pazuzufile

	// log.Print("Building Dockerfile out of the features")
	// // TODO: check number of c.NArgs() and throw error if nothing was passed
	// if c.NArg() == 0 {

	// 	return errors.New(ERROR_NO_VALID_PAZUZU_FILE)

	// }

	// return nil
}
