package main


import (
	"fmt"

	"github.com/urfave/cli"
	"github.com/zalando-incubator/pazuzu"
)


var composeAction = func(c *cli.Context) error {
	var initFeatures = getFeaturesList(c.String("init"))
	var addFeatures = getFeaturesList(c.String("add"))
	var pazuzufileFeatures []string
	var baseImage string

	pazuzuFile, success := readPazuzuFile()
	if success {
		pazuzufileFeatures = pazuzuFile.Features
		baseImage = pazuzuFile.Base
	}

	featureNames, err := generateFeaturesList(pazuzufileFeatures, initFeatures, addFeatures)
	if err != nil {
		return err
	}
	fmt.Printf("Resolving the following features: %s\n", featureNames)

	config := pazuzu.GetConfig()
	storageReader, err := pazuzu.GetStorageReader(*config)
	if err != nil {
		return err // TODO: process properly into human-readable message
	}

	features, err := checkFeaturesInRepository(featureNames, storageReader)
	if err != nil {
		return err
	}

	if baseImage == "" {
		baseImage = config.Base
	}

	fmt.Print("Generating Pazuzufile...")

	pazuzuFile = &pazuzu.PazuzuFile{
		Base:     baseImage,
		Features: features,
	}

	err = writePazuzuFile(pazuzuFile)
	if err != nil {
		return err
	} else {
		fmt.Println(" [DONE]")
	}

	fmt.Print("Generating Dockerfile...")

	p := pazuzu.Pazuzu{StorageReader: storageReader}
	p.Generate(pazuzuFile.Base, pazuzuFile.Features)

	err = writeDockerFile(p.Dockerfile)
	if err != nil {
		return err
	} else {
		fmt.Println(" [DONE]")
	}

	return nil
}
