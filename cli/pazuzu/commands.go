package main

import (
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"os"
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
	Action: func(с *cli.Context) error {
		return ErrNotImplemented
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

var cnfGetCmd = cli.Command{
	Name:  "get",
	Usage: "Get pazuzu configuration",
	Action: func(с *cli.Context) error {
		// log.Print("Getting pazuzu configuration")
		// return nil
		return errors.New(ERROR_NOT_IMPLEMENTED)
	},
}
var cnfSetCmd = cli.Command{
	Name:  "set",
	Usage: "Set pazuzu configuration",
	Action: func(с *cli.Context) error {
		// log.Print("Setting pazuzu configuration")
		// return nil
		return errors.New(ERROR_NOT_IMPLEMENTED)
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

func main() {

	cli.VersionFlag = cli.BoolFlag{
		Name:  "version",
		Usage: "Print version",
	}

	app := cli.NewApp()
	app.Name = "pazuzu"
	app.Version = VERSION
	app.Usage = "Build Docker features from pazuzu-registry"
	app.Commands = []cli.Command{
		searchCmd,
		composeCmd,
		buildCmd,
		configCmd,
	}

	// global flags
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose, v",
			Usage: "Verbose output",
		},
	}
	app.Before = func(c *cli.Context) error {
		// remove formatting for log module
		// and suppress logging output if not set explicitly
		log.SetFlags(0)
		if c.Bool("verbose") {
			log.SetOutput(os.Stdout)
		} else {
			log.SetOutput(ioutil.Discard)
		}

		//TODO: Init config struct
		errCnf := NewConfig()

		if errCnf != nil {
			fmt.Println(errCnf)
			os.Exit(1)
		}
		// Sample reating conf values
		// log.Printf("Using URL: %v", config.Git.Url)

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
