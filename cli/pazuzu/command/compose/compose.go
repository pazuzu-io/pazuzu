package compose

import (
	"github.com/urfave/cli"
	"fmt"
	"github.com/zalando-incubator/pazuzu"
	"github.com/zalando-incubator/pazuzu/config"
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
		{
			Name: "remove",
			Aliases: []string{"rm", "d"},
			Description: "Removes provided list of features from Pazuzufile",
			Usage: "remove feature,...",
			Action: removeFeaturesAction,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "d, directory",
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

	if baseImage == "" || c.String("init") != "" {
		baseImage = config.GetConfig().Base
	}

	err = generateFiles(destination, baseImage, featureNames)
	if err != nil {
		return err
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

var removeFeaturesAction = func(c *cli.Context) error {
	destination := c.String("directory")
	if !c.Args().Present(){
		return errors.New("ERROR: no features to remove")
	}
	features := strings.Split(c.Args().First(), ",")

	err := utils.CheckDestination(destination)
	if err != nil {
		return err
	}

	pazuzufilePath := utils.GetAbsoluteFilePath(destination, pazuzu.PazuzufileName)
	pazuzuFile, success := utils.ReadPazuzuFile(pazuzufilePath)
	if !success {
		return nil
	}

	newFeatures := pazuzuFile.Features
	baseImage := pazuzuFile.Base

	loop:
	for i := 0; i < len(newFeatures); i++ {
		f1 := newFeatures[i]
		for _, f2 := range features {
			if f1 == f2 {
				newFeatures = append(newFeatures[:i], newFeatures[i+1:]...)
				i--
				continue loop
			}
		}
	}

	err = generateFiles(destination, baseImage, newFeatures)
	if err != nil {
		return err
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

func generateFiles(destination string, baseImage string, featureNames []string) error {
	err := utils.CheckDestination(destination)
	if err != nil {
		return err
	}
	pazuzufilePath := utils.GetAbsoluteFilePath(destination, pazuzu.PazuzufileName)
	storageReader, err := config.GetStorageReader(*config.GetConfig())
	if err != nil {
		return err
	}
	fmt.Printf("Resolving the following features: %s\n", featureNames)
	features, err := utils.CheckFeaturesInRepository(featureNames, storageReader)
	if err != nil {
		return err
	}
	if baseImage == "" {
		baseImage = config.GetConfig().Base
	}

	fmt.Printf("Generating %s...\n", pazuzufilePath)
	pazuzuFile := &pazuzu.PazuzuFile{
		Base: baseImage,
		Features: features,
	}
	err = utils.WritePazuzuFile(pazuzufilePath, pazuzuFile)
	if err != nil {
		return err
	}
	fmt.Println("[DONE]")

	return nil
}
