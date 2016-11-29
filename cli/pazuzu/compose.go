package main

import (
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"github.com/zalando-incubator/pazuzu"
)

var composeAction = func(c *cli.Context) error {
	var (
		initFeatures       = getFeaturesList(c.String("init"))
		addFeatures        = getFeaturesList(c.String("add"))
		destination        = c.String(directoryOption)
		pazuzufileFeatures []string
		baseImage          string
	)

	if (c.String("add") == "") && (c.String("init") == "") {
		cli.ShowCommandHelp(c, "compose")
		return errors.New("ERROR: No feature specified. Please use at least one of -a or -i for the compose command.")
	}

	err := checkDestination(destination)
	if err != nil {
		return err
	}

	pazuzufilePath := getAbsoluteFilePath(destination, PazuzufileName)
	dockerfilePath := getAbsoluteFilePath(destination, DockerfileName)
	testSpecPath := getAbsoluteFilePath(destination, TestSpecFileName)

	pazuzuFile, success := readPazuzuFile(pazuzufilePath)
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

	fmt.Printf("Generating %s...", pazuzufilePath)

	pazuzuFile = &pazuzu.PazuzuFile{
		Base:     baseImage,
		Features: features,
	}

	err = writePazuzuFile(pazuzufilePath, pazuzuFile)
	if err != nil {
		return err
	} else {
		fmt.Println(" [DONE]")
	}

	fmt.Printf("Generating %s...", dockerfilePath)

	p := pazuzu.Pazuzu{StorageReader: storageReader}
	p.Generate(pazuzuFile.Base, pazuzuFile.Features)

	err = writeFile(dockerfilePath, p.Dockerfile)
	if err != nil {
		return err
	} else {
		fmt.Println(" [DONE]")
	}

	fmt.Printf("Generating %s...", testSpecPath)
	err = writeFile(testSpecPath, p.TestSpec)
	if err != nil {
		return err
	} else {
		fmt.Println(" [DONE]")
	}

	return nil
}
