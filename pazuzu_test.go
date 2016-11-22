package pazuzu

import (
	"bytes"
	"regexp"
	"strings"
	"testing"

	"github.com/zalando-incubator/pazuzu/storageconnector"
)

type TestStorage struct{}

func (s *TestStorage) GetFeature(feature string) (storageconnector.Feature, error) {
	return storageconnector.Feature{
		Meta: storageconnector.FeatureMeta{
			Name:        "python",
			Description: "Use `python -V`",
		},
		Snippet: "RUN apt-get update && apt-get install python --yes",
	}, nil
}

func (s *TestStorage) GetMeta(name string) (storageconnector.FeatureMeta, error) {
	return storageconnector.FeatureMeta{
		Name:        "python",
		Description: "Use `python -V`",
	}, nil
}

func (s *TestStorage) SearchMeta(name *regexp.Regexp) ([]storageconnector.FeatureMeta, error) {
	return make([]storageconnector.FeatureMeta, 0), nil
}

func (s *TestStorage) Resolve(names ...string) (map[string]storageconnector.Feature, error) {
	return make(map[string]storageconnector.Feature), nil
}

// Test generating a Dockerfile from a list of features.
func TestGenerate(t *testing.T) {
	pazuzu := Pazuzu{
		StorageReader: &TestStorage{},
		testSpec:      "test_spec.json",
	}

	err := pazuzu.Generate("ubuntu", []string{"python"})
	defer pazuzu.Cleanup()
	if err != nil {
		t.Errorf("should not fail: %s", err)
	}
}

func TestRead(t *testing.T) {
	bufferedReader := strings.NewReader(`---
base: ubuntuCommon
features:
  - Java8
  - anotherFeature
  - oneMoreFeature`)

	pazuzuFile, err := Read(bufferedReader)

	if err != nil {
		t.Errorf("should not fail: %s", err)
	}

	if strings.Compare(pazuzuFile.Base, "ubuntuCommon") != 0 {
		t.Errorf("wrong base: %s", pazuzuFile.Base)
	}
}

func TestWrite(t *testing.T) {
	pazuzuFile := PazuzuFile{
		Base:     "ubuntuCommon",
		Features: []string{"java8", "anotherFeature", "oneMoreFeature"},
	}

	b := []byte{}
	ioWriter := bytes.NewBuffer(b)
	err := Write(ioWriter, pazuzuFile)

	if err != nil {
		t.Errorf("should not fail: %s", err)
	}
}

// Test building a generated Dockerfile.
func TestDockerBuild(t *testing.T) {
	pazuzu := Pazuzu{
		dockerEndpoint: "unix:///var/run/docker.sock",
		Dockerfile: []byte(`FROM ubuntu:latest
RUN apt-get update && apt-get install python --yes`),
		testSpec: "test_spec.json",
	}

	err := pazuzu.DockerBuild("test")
	if err != nil {
		t.Errorf("should not fail: %s", err)
	}
}
