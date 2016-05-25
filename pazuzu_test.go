package main

import "testing"

func TestGenerate(t *testing.T) {
	pazuzu := Pazuzu{
		registry:   "http://localhost:8080/api",
		dockerfile: "Dockerfile",
		testScript: "test.spec",
	}
	err := pazuzu.Generate([]string{"python"})
	if err != nil {
		t.Errorf("should not fail: %s", err)
	}
}

func TestDockerBuild(t *testing.T) {
	pazuzu := Pazuzu{
		registry:       "http://localhost:8080/api",
		dockerEndpoint: "unix:///var/run/docker.sock",
		dockerfile:     "Dockerfile",
		testScript:     "test.spec",
	}

	err := pazuzu.DockerBuild("test")
	if err != nil {
		t.Errorf("should not fail: %s", err)
	}
}
