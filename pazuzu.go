package main

import (
	"archive/tar"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fsouza/go-dockerclient"
)

// Pazuzu defines pazuzu config.
type Pazuzu struct {
	registry       string
	dockerfile     []byte
	testScript     string
	dockerEndpoint string
}

// Feature defines a feature fetched from pazuzu-registry.
type Feature struct {
	Name            string
	DockerData      string `json:"docker_data"`
	TestInstruction string `json:"test_instruction"`
}

// Generate generates Dockfiler and test.spec file base on list of features
func (p *Pazuzu) Generate(features []string) error {
	fs, err := p.getFeatures(features)
	if err != nil {
		return err
	}

	err = p.generateDockerfile(fs)
	if err != nil {
		return err
	}

	err = p.generateTestScript(fs)
	if err != nil {
		return err
	}

	return nil
}

// generate in-memory Dockerfile from list of features.
func (p *Pazuzu) generateDockerfile(features []Feature) error {
	var buf bytes.Buffer

	_, err := buf.WriteString("FROM ubuntu:latest\n")
	if err != nil {
		return err
	}

	for _, feature := range features {
		_, err = buf.WriteString(fmt.Sprintf("# %s\n", feature.Name))
		if err != nil {
			return err
		}

		_, err = buf.WriteString(fmt.Sprintf("%s\n", feature.DockerData))
		if err != nil {
			return err
		}
	}

	_, err = buf.WriteString("CMD /bin/bash\n")
	if err != nil {
		return err
	}

	p.dockerfile = buf.Bytes()

	return nil
}

// generate test script from list of features.
func (p *Pazuzu) generateTestScript(features []Feature) error {
	f, err := os.Create(p.testScript)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	for _, feature := range features {
		_, err = w.WriteString(fmt.Sprintf("# %s\n", feature.Name))
		if err != nil {
			return err
		}

		_, err = w.WriteString(fmt.Sprintf("%s;\n", feature.DockerData))
		if err != nil {
			return err
		}
	}

	return w.Flush()
}

// APIError defines error response from pazuzu-registry.
type APIError struct {
	Code            string
	Message         string
	DetailedMessage string
}

// get a list of features given the feature names.
func (p *Pazuzu) getFeatures(features []string) ([]Feature, error) {
	resp, err := http.Get(fmt.Sprintf("%s/features?name=%s", p.registry, strings.Join(features, ",")))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp APIError

		dec := json.NewDecoder(resp.Body)
		err = dec.Decode(&errResp)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf(errResp.Message)
	}

	var res []Feature

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// DockerBuild builds a docker image based on the generated Dockerfile.
func (p *Pazuzu) DockerBuild(name string) error {
	client, err := docker.NewClient(p.dockerEndpoint)
	if err != nil {
		return err
	}

	t := time.Now()
	inputBuf := bytes.NewBuffer(nil)
	tr := tar.NewWriter(inputBuf)
	tr.WriteHeader(&tar.Header{
		Name:       "Dockerfile",
		Size:       int64(len(p.dockerfile)),
		ModTime:    t,
		AccessTime: t,
		ChangeTime: t,
	})
	tr.Write(p.dockerfile)
	tr.Close()

	opts := docker.BuildImageOptions{
		Name:         name,
		InputStream:  inputBuf,
		OutputStream: os.Stdout,
	}

	err = client.BuildImage(opts)
	if err != nil {
		return err
	}

	return nil
}
