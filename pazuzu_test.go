package pazuzu

import (
	"bytes"
	"strings"
	"testing"

	"github.com/zalando-incubator/pazuzu/mock"
	"io/ioutil"
)

// Test generating a Dockerfile from a list of features.
func TestGenerate(t *testing.T) {
	pazuzu := Pazuzu{
		StorageReader: &mock.TestStorage{},
		testSpec:      "test_spec.json",
	}

	err := pazuzu.Generate("ubuntu", []string{"python"})
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
		DockerEndpoint: "unix:///var/run/docker.sock",
		Dockerfile: []byte(`FROM ubuntu:latest
RUN apt-get update && apt-get install python --yes`),
		TestSpec: []byte(`#!/usr/bin/env bats
@test "Check echo" {
	command echo
}`),
	}

	// usually, this is composed in `pazuzu compose` phase
	err := ioutil.WriteFile(TestSpecFilename, pazuzu.TestSpec, 0644)
	if err != nil {
		t.Errorf("could not create tests.bat file: %s", err)
	}

	err = pazuzu.DockerBuild("test")
	if err != nil {
		t.Errorf("should not fail: %s", err)
	}
}
