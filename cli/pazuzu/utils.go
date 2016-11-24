package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/zalando-incubator/pazuzu"
	"github.com/zalando-incubator/pazuzu/storageconnector"
)

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

func generateFeaturesList(pazuzufileFeatures []string, featuresToInit []string, featuresToAdd []string) ([]string, error) {
	var features []string

	if len(featuresToInit) > 0 && len(featuresToAdd) > 0 {
		return features, pazuzu.ErrInitAndAddAreSpecified
	}

	if len(featuresToInit) > 0 {
		return featuresToInit, nil
	}

	if len(featuresToAdd) > 0 {
		features = pazuzufileFeatures
		for _, feature := range featuresToAdd {
			features = appendIfMissing(features, feature)
		}
		return features, nil
	}

	return features, nil
}

func appendIfMissing(slice []string, element string) []string {
	for _, next := range slice {
		if next == element {
			return slice
		}
	}
	return append(slice, element)
}

// Reads Pazuzufile
// returns PazuzuFile struct and a success flag
func readPazuzuFile(path string) (*pazuzu.PazuzuFile, bool) {
	file, err := os.Open(path)
	if err != nil {
		return nil, false
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	pazuzuFile, err := pazuzu.Read(reader)
	if err != nil {
		return nil, false
	}

	return &pazuzuFile, true
}

func writePazuzuFile(path string, pazuzuFile *pazuzu.PazuzuFile) error {
	// TODO: do it safer way (#108)
	file, err := os.Create(path)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not create %v", PazuzufileName))
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	pazuzu.Write(writer, *pazuzuFile)

	writer.Flush()
	return nil
}

func writeDockerFile(path string, contents []byte) error {
	file, err := os.Create(path)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not create %v", DockerfileName))
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.Write(contents)
	writer.Flush()

	return nil
}

func checkFeaturesInRepository(names []string, storage storageconnector.StorageReader) ([]string, error) {
	var features []string

	for _, name := range names {
		log.Printf("Checking: %v\n", name)

		_, err := storage.GetMeta(name)
		if err != nil {
			return features, errors.New(fmt.Sprintf("Feature %v not found", name))
		}
		features = append(features, fmt.Sprintf("%v", name))
	}

	return features, nil
}

// Gets absolute file paths for Pazuzufile and Dockerfile
// returns Pazuzufile path, Dockerfile path and an error
func getAbsoluteFilePaths(destination string) (string, string, error) {
	var pazuzufilePath = PazuzufileName
	var dockerfilePath = DockerfileName

	if destination != "" {
		destination, err := filepath.Abs(destination)
		if err != nil {
			return "", "", err
		}

		_, err = os.Stat(destination)
		if err != nil {
			err = errors.New(fmt.Sprintf("Destination path %s is not found", destination))
			return "", "", err
		}

		pazuzufilePath = filepath.Join(destination, PazuzufileName)
		dockerfilePath = filepath.Join(destination, DockerfileName)
	}
	return pazuzufilePath, dockerfilePath, nil
}
