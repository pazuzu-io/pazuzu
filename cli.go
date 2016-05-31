package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

var buildCmd = cli.Command{
	Name:   "build",
	Usage:  "build docker image",
	Action: buildFeatures,
}

func buildFeatures(c *cli.Context) error {
	pazuzu := Pazuzu{
		registry:       "http://localhost:8080/api",
		testScript:     "test.spec",
		dockerEndpoint: "unix:///var/run/docker.sock",
	}

	if len(c.Args()) == 0 {
		return fmt.Errorf("no features specified")
	}

	err := pazuzu.Generate(c.Args())
	if err != nil {
		return err
	}

	err = pazuzu.DockerBuild("test")
	if err != nil {
		return err
	}

	return nil
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
