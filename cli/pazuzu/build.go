package main

import (
	"github.com/urfave/cli"
	"github.com/zalando-incubator/pazuzu"
	"fmt"
	"io/ioutil"
)

// Fetches and builds features into a docker image.
func buildFeatures(c *cli.Context) error {
	config := pazuzu.GetConfig()
	storageReader, err := pazuzu.GetStorageReader(*config)
	if err != nil {
		return fmt.Errorf("Error during storage setup:%s", err)
	}

	directory := c.String(directoryOption)
	err = checkDestination(directory)
	if err != nil {
		return fmt.Errorf("Error to access directory:%s\n%s", directory, err)
	}
	dockerFileName := getAbsoluteFilePath(directory, DockerfileName)
	dat, err := ioutil.ReadFile(dockerFileName)
	if err != nil {
		return fmt.Errorf("Error during attempt to read docker file:%s", err)
	}

	p := pazuzu.Pazuzu{StorageReader: storageReader,
		DockerEndpoint: "unix:///var/run/docker.sock",
		Dockerfile:dat,
	}

	name := "pazuzu image"
	if c.String("name") == "" {
		name = c.String("name")
	}
	err2 := p.DockerBuild(name)
	if err2 != nil {
		return fmt.Errorf("should not fail: %s", err)
	}
	return nil
}

