package actions

import (
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/urfave/cli"
	"github.com/zalando-incubator/pazuzu"
	"github.com/zalando-incubator/pazuzu/cli/pazuzu/utils"
	"github.com/zalando-incubator/pazuzu/config"
	"io/ioutil"
	"os"
	"strings"
)

func ProjectClean(c *cli.Context) error {
	err := os.Remove(pazuzu.PazuzufileName)
	if err != nil {
		fmt.Println(err)
	}
	err = os.Remove(pazuzu.DockerfileName)
	if err != nil {
		fmt.Println(err)
	}
	err = os.Remove(pazuzu.TestSpecFilename)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func ProjectBuild(c *cli.Context) error {
	directory := c.String("directory")
	err := utils.CheckDestination(directory)
	if err != nil {
		return fmt.Errorf("Error to access directory:%s\n%s", directory, err)
	}

	pazuzufilePath := utils.GetAbsoluteFilePath(directory, pazuzu.PazuzufileName)
	dockerfilePath := utils.GetAbsoluteFilePath(directory, pazuzu.DockerfileName)
	testSpecPath := utils.GetAbsoluteFilePath(directory, pazuzu.TestSpecFilename)
	pazuzuFile, success := utils.ReadPazuzuFile(pazuzufilePath)
	if !success {
		return fmt.Errorf("Can not read configuration: %s\n", pazuzufilePath)
	}

	storageReader, err := config.GetStorageReader(*config.GetConfig())
	if err != nil {
		return fmt.Errorf("Error during storage setup:%s", err)
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

	p.DockerEndpoint = pazuzu.DefaultDockerEndpoint
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

func ProjectListFeatures(c *cli.Context) error {
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

func ProjectRemoveFeatures(c *cli.Context) error {
	destination := c.String("directory")
	if !c.Args().Present() {
		return errors.New("ERROR: no features to remove")
	}
	features := getFeaturesList(c.Args().First())

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

func ProjectAddFeatures(c *cli.Context) error {
	destination := c.String("directory")
	if !c.Args().Present() {
		return errors.New("ERROR: no features to add")
	}
	features := getFeaturesList(c.Args().First())

	err := utils.CheckDestination(destination)
	if err != nil {
		return err
	}

	pazuzufilePath := utils.GetAbsoluteFilePath(destination, pazuzu.PazuzufileName)
	pazuzuFile, success := utils.ReadPazuzuFile(pazuzufilePath)
	baseImage := config.GetConfig().Base
	var currentFeatures []string
	if success {
		baseImage = pazuzuFile.Base
		currentFeatures = pazuzuFile.Features
	}
	for _, f := range features {
		if !isFeatureInList(currentFeatures, f) {
			currentFeatures = append(currentFeatures, f)
		}
	}

	err = generateFiles(destination, baseImage, currentFeatures)
	if err != nil {
		return err
	}

	return nil
}

func ProjectShow(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.New("Wrong number of arguments")
	}
	key := c.Args().First()
	destination := c.String("directory")

	err := utils.CheckDestination(destination)
	if err != nil {
		return err
	}

	pazuzufilePath := utils.GetAbsoluteFilePath(destination, pazuzu.PazuzufileName)
	pazuzuFile, success := utils.ReadPazuzuFile(pazuzufilePath)
	if !success {
		return errors.New("Project doesn't have configuration yet")
	}
	if key == "base" {
		fmt.Printf("%s => %s\n", key, pazuzuFile.Base)
	}
	return nil
}

func ProjectSet(c *cli.Context) error {
	if c.NArg() != 2 {
		return errors.New("Wrong number of arguments")
	}
	key := c.Args()[0]
	value := c.Args()[1]
	destination := c.String("directory")

	err := utils.CheckDestination(destination)
	if err != nil {
		return err
	}

	var features []string
	base := config.GetConfig().Base

	pazuzufilePath := utils.GetAbsoluteFilePath(destination, pazuzu.PazuzufileName)
	pazuzuFile, success := utils.ReadPazuzuFile(pazuzufilePath)

	if success {
		base = pazuzuFile.Base
		features = pazuzuFile.Features
	}

	if key != "base" {
		return errors.New("Key is not supported")
	}

	if base == value {
		// nothing changes so return
		return nil
	}

	err = generateFiles(destination, base, features)
	if err != nil {
		return err
	}

	return nil
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
		Base:     baseImage,
		Features: features,
	}
	err = utils.WritePazuzuFile(pazuzufilePath, pazuzuFile)
	if err != nil {
		return err
	}
	fmt.Println("[DONE]")

	return nil
}

func getFeaturesList(featuresString string) []string {
	var features []string
	featuresString = strings.Trim(featuresString, ", ")
	if len(featuresString) > 0 {
		for _, element := range strings.Split(featuresString, ",") {
			features = append(features, strings.Trim(element, " "))
		}
	}

	return features
}

func isFeatureInList(features []string, feature string) bool {
	for _, f := range features {
		if f == feature {
			return true
		}
	}
	return false
}
