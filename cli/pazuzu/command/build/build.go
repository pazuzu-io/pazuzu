package build

import (
	"github.com/urfave/cli"
	"github.com/zalando-incubator/pazuzu"
	"github.com/zalando-incubator/pazuzu/config"
	"github.com/zalando-incubator/pazuzu/shared"
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

	pazuzufilePath := utils.GetAbsoluteFilePath(directory, pazuzu.PazuzufileName)
	dockerfilePath := utils.GetAbsoluteFilePath(directory, pazuzu.DockerfileName)
	testSpecPath := utils.GetAbsoluteFilePath(directory, shared.TestSpecFilename)
	pazuzuFile, success := utils.ReadPazuzuFile(pazuzufilePath)
	if !success {
		return fmt.Errorf("Can not read configuration: %s\n", pazuzufilePath)
	}

	p := pazuzu.Pazuzu{StorageReader: storageReader}
	p.Generate(pazuzuFile.Base, pazuzuFile.Features)
	fmt.Printf("Generating %s...\n", dockerfilePath)
	err = utils.WriteFile(dockerfilePath, p.Dockerfile)
	if err != nil {
		return fmt.Errorf("Can not write Dockerfile: %s\n%s", dockerfilePath, err)
	}
	fmt.Printf("Generating %s...\n", testSpecPath)
	err = utils.WriteFile(testSpecPath, p.TestSpec)
	if err != nil {
		return fmt.Errorf("Can not write TestSpec: %s\n%s", testSpecPath, err)
	}

	dat, err := ioutil.ReadFile(dockerfilePath)
	if err != nil {
		return fmt.Errorf("Error during attempt to read docker file:%s", err)
	}

	p.DockerEndpoint = "unix:///var/run/docker.sock"
	p.Dockerfile = dat

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
