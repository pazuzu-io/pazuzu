package main

import (
	"fmt"
	"github.com/urfave/cli"
	"github.com/zalando-incubator/pazuzu"
	"io/ioutil"
	"log"
	"os"
)

var Version = "0.1"

func main() {

	cli.VersionFlag = cli.BoolFlag{
		Name:  "version",
		Usage: "Print version",
	}

	app := cli.NewApp()
	app.Name = "pazuzu"
	app.Version = Version
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
		errCnf := pazuzu.NewConfig()

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
