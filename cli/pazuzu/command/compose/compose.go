package compose

import (
	"github.com/urfave/cli"
	"fmt"
	"github.com/zalando-incubator/pazuzu"
	"github.com/zalando-incubator/pazuzu/config"
	"github.com/zalando-incubator/pazuzu/shared"
	"errors"
	"strings"
	"github.com/zalando-incubator/pazuzu/cli/pazuzu/utils"
)

var Command = cli.Command{
	Name:      "compose",
	Usage:     "Compose Pazuzufile and Dockerfile out of the selected features",
	ArgsUsage: " ", // Do not show arguments
	Description: "Compose step takes list of features as input, validates feature dependencies" +
		" and creates both Pazuzufile and Dockerfile.",
	Flags:  composeFlags,
	Action: composeAction,
	Subcommands: []cli.Command{
		{
			Name:  "list",
			Aliases: []string{"l"},
			Usage: "Lists features in Pazuzufile",
			Action: listFeaturesAction,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "d, directory",
					Usage: "Sets destination directory for Dockerfile and Pazuzufile to `DESTINATION`",
				},
			},
		},
	},
}




var composeFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "a, add",
		Usage: "Add features from comma-separated list of `FEATURES`",
	},
	cli.StringFlag{
		Name:  "i, init",
		Usage: "Init set of features from comma-separated list of `FEATURES`",
	},
	cli.StringFlag{
		Name:  "d, directory",
		Usage: "Sets destination directory for Docketfile and Pazuzufile to `DESTINATION`",
	},
	cli.StringFlag{
		Name:  "b, base",
		Usage: "Sets the base docker image to `BASE`, instead of the one from the configuration",
	},
}


var composeAction = func(c *cli.Context) error {

	initFeatures       := getFeaturesList(c.String("init"))
	addFeatures        := getFeaturesList(c.String("add"))
	destination        := c.String("directory")
	pazuzufileFeatures := []string {}
	baseImage          := c.String("base")


	if (c.String("add") == "") && (c.String("init") == "") {
		cli.ShowCommandHelp(c, "compose")
		return errors.New("ERROR: No feature specified. Please use at least one of -a or -i for the compose command.")
	}

	err := utils.CheckDestination(destination)
	if err != nil {
		return err
	}

	pazuzufilePath := utils.GetAbsoluteFilePath(destination, pazuzu.PazuzufileName)
	dockerfilePath := utils.GetAbsoluteFilePath(destination, pazuzu.DockerfileName)
	testSpecPath := utils.GetAbsoluteFilePath(destination, shared.TestSpecFilename)

	pazuzuFile, success := utils.ReadPazuzuFile(pazuzufilePath)
	if success {
		pazuzufileFeatures = pazuzuFile.Features
		if baseImage == "" {
			baseImage = pazuzuFile.Base
		}
	}

	featureNames, err := utils.GenerateFeaturesList(pazuzufileFeatures, initFeatures, addFeatures)
	if err != nil {
		return err
	}
	fmt.Printf("Resolving the following features: %s\n", featureNames)

	storageReader, err := config.GetStorageReader(*config.GetConfig())
	if err != nil {
		return err // TODO: process properly into human-readable message
	}

	features, err := utils.CheckFeaturesInRepository(featureNames, storageReader)
	if err != nil {
		return err
	}

	if baseImage == "" || c.String("init") != "" {
		baseImage = config.GetConfig().Base
	}

	fmt.Printf("Generating %s...", pazuzufilePath)

	pazuzuFile = &pazuzu.PazuzuFile{
		Base:     baseImage,
		Features: features,
	}

	err = utils.WritePazuzuFile(pazuzufilePath, pazuzuFile)
	if err != nil {
		return err
	} else {
		fmt.Println(" [DONE]")
	}

	fmt.Printf("Generating %s...", dockerfilePath)

	p := pazuzu.Pazuzu{StorageReader: storageReader}
	p.Generate(pazuzuFile.Base, pazuzuFile.Features)

	err = utils.WriteFile(dockerfilePath, p.Dockerfile)

	fmt.Printf("Generating %s...", testSpecPath)
	err = utils.WriteFile(testSpecPath, p.TestSpec)

	if err != nil {
		return err
	} else {
		fmt.Println(" [DONE]")
	}

	return nil
}

var listFeaturesAction = func(c *cli.Context) error {
	destination := c.String("directory")
	err := utils.CheckDestination(destination)
	if err != nil {
		return err
	}
	pazuzufilePath := utils.GetAbsoluteFilePath(destination, pazuzu.PazuzufileName)
	pazuzuFile, success := utils.ReadPazuzuFile(pazuzufilePath)
	if success {
		pazuzufileFeatures := pazuzuFile.Features
		for _, feature := range pazuzufileFeatures {
			fmt.Println(feature)
		}
	}
	return nil
}

func getFeaturesList(featureString string) []string {
	var features []string

	featureString = strings.Trim(featureString, ", ")
	if len(featureString) > 0 {
		for _, element := range strings.Split(featureString, ",") {
			features = append(features, strings.Trim(element, " "))
		}
	}

	return features
}
