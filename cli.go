package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

var buildCmd = cli.Command{
	Name:   "build",
	Usage:  "build docker image",
	Action: buildFeatures,
}

func buildFeatures(c *cli.Context) {
	pazuzu := Pazuzu{
		registry:   "http://localhost:8080/api",
		dockerfile: "Dockerfile",
		testScript: "test.spec",
	}

	if len(c.Args()) == 0 {
		fmt.Println("no features specified")
		os.Exit(1)
	}

	err := pazuzu.Generate(c.Args())
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "pazuzu"
	app.Version = "0.1"
	app.Usage = "Build Docker features from pazuzu-registry"
	app.Commands = []cli.Command{
		buildCmd,
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
