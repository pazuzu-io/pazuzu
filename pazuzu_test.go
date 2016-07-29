package main

import "testing"

type TestRegistry struct{}

func (r TestRegistry) GetFeatures(features []string) ([]Feature, error) {
	return []Feature{
		Feature{
			Name:            "python",
			DockerData:      "RUN apt-get update && apt-get install python --yes",
			TestInstruction: "python -V",
		},
	}, nil
}

// Test generating a Dockerfile from a list of features.
func TestGenerate(t *testing.T) {
	pazuzu := Pazuzu{
		registry: TestRegistry{},
		testSpec: "test_spec.json",
	}
	err := pazuzu.Generate("ubuntu", []string{"python"})
	if err != nil {
		t.Errorf("should not fail: %s", err)
	}
}

// Test building a generated Dockerfile.
func TestDockerBuild(t *testing.T) {
	pazuzu := Pazuzu{
		dockerEndpoint: "unix:///var/run/docker.sock",
		dockerfile: []byte(`FROM ubuntu:latest
RUN apt-get update && apt-get install python --yes`),
		testSpec: "test_spec.json",
	}

	err := pazuzu.DockerBuild("test")
	if err != nil {
		t.Errorf("should not fail: %s", err)
	}
}

// Test verifying a docker image.
func TestRunTestSpec(t *testing.T) {
	pazuzu := Pazuzu{
		dockerEndpoint: "unix:///var/run/docker.sock",
		testSpec:       "test_spec.json",
	}

	err := pazuzu.RunTestSpec("test")
	if err != nil {
		t.Errorf("should not fail: %s", err)
	}
}
