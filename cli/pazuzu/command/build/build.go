package build

import (
	"github.com/urfave/cli"
	"github.com/zalando-incubator/pazuzu"
	"github.com/zalando-incubator/pazuzu/config"
	"fmt"
	"io/ioutil"
	"strings"
	"github.com/satori/go.uuid"
	"github.com/zalando-incubator/pazuzu/cli/pazuzu/utils"
)

var Command = cli.Command{
	Name:      "build",
	Usage:     "Builds and tests Docker image from Dockerfile",
	ArgsUsage: " ",
	Flags:     buildFlags,
	Action:    buildFeatures,
}


var buildFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "d, directory",
		Usage: "Sets source path where Docketfile are located.",
	},
	cli.StringFlag{
		Name:  "n, name",
		Usage: "Sets a name for docker image",
	},
}


// Fetches and builds features into a docker image.
func buildFeatures(c *cli.Context) error {
	storageReader, err := config.GetStorageReader(*config.GetConfig())
	if err != nil {
		return fmt.Errorf("Error during storage setup:%s", err)
	}

	directory := c.String("directory")
	err = utils.CheckDestination(directory)
	if err != nil {
		return fmt.Errorf("Error to access directory:%s\n%s", directory, err)
	}
	dockerFileName := utils.GetAbsoluteFilePath(directory, pazuzu.DockerfileName)
	dat, err := ioutil.ReadFile(dockerFileName)
	if err != nil {
		return fmt.Errorf("Error during attempt to read docker file:%s", err)
	}

	p := pazuzu.Pazuzu{StorageReader: storageReader,
		DockerEndpoint: "unix:///var/run/docker.sock",
		Dockerfile:     dat,
	}

	name := ""
	if c.String("name") != "" {
		name = c.String("name")
	} else {
		name = strings.Replace(uuid.NewV1().String(), "-", "", -1)
	}
	err2 := p.DockerBuild(name)
	if err2 != nil {
		return fmt.Errorf("should not fail: %s", err2)
	}
	return nil
}
