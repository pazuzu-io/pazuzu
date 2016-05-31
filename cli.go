package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

var version = "0.1"

var buildCmd = cli.Command{
	Name:   "build",
	Usage:  "build docker image",
	Action: buildFeatures,
}

// Fetches and builds features into a docker image.
func buildFeatures(c *cli.Context) error {
	pazuzu := Pazuzu{
		registry:       "http://localhost:8080/api",
		testScript:     "test.spec",
		dockerEndpoint: c.GlobalString("docker-endpoint"),
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
	app.Version = version
	app.Usage = "Build Docker features from pazuzu-registry"
	app.Commands = []cli.Command{
		buildCmd,
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "docker-endpoint, e",
			Value: "unix:///var/run/docker.sock",
			Usage: "Set the docker endpoint",
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
