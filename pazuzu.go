package pazuzu

import (
	"archive/tar"
	"bytes"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/zalando-incubator/pazuzu/shared"
	"github.com/zalando-incubator/pazuzu/storageconnector"
)

// Pazuzu defines pazuzu config.
type Pazuzu struct {
	StorageReader  storageconnector.StorageReader
	Dockerfile     []byte
	TestSpec       []byte
	testSpec       string
	DockerEndpoint string
	docker         *docker.Client
	files          map[string]string
}

type PazuzuFile struct {
	Base     string
	Features []string
}

func Read(reader io.Reader) (PazuzuFile, error) {
	content, err := ioutil.ReadAll(reader)

	if err != nil {
		return PazuzuFile{}, err
	}

	pazuzuFile := &PazuzuFile{}
	err = yaml.Unmarshal(content, pazuzuFile)
	if err != nil {
		return PazuzuFile{}, err
	}

	return *pazuzuFile, nil

}

func Write(writer io.Writer, pazuzuFile PazuzuFile) error {
	data, err := yaml.Marshal(pazuzuFile)
	if err != nil {
		return err
	}

	_, err = writer.Write(data)

	return err
}

// Generate generates Dockfiler and test.spec file base on list of features
func (p *Pazuzu) Generate(baseimage string, features []string) error {
	var resolvedFeatures []string
	for _, feature := range features {
		// TODO: add proper error handling
		repoFeature, _ := p.StorageReader.GetFeature(feature)
		resolvedFeatures = append(resolvedFeatures, repoFeature.Meta.Name)
	}
	// TODO: add proper error handling
	var featureNamesWithDep []string
	featureNamesWithDep, featuresMap, _ := p.StorageReader.Resolve(resolvedFeatures...)
	featuresWithDep := make([]shared.Feature, 0, len(featuresMap))

	for _, featureName := range featureNamesWithDep {
		featuresWithDep = append(featuresWithDep, featuresMap[featureName])
	}

	err := p.generateDockerfile(baseimage, featuresWithDep)
	if err != nil {
		return err
	}

	return nil
}

// generate in-memory Dockerfile from list of features.
func (p *Pazuzu) generateDockerfile(baseimage string, features []shared.Feature) error {
	writer := NewDockerfileWriter()

	err := writer.AppendRaw(fmt.Sprintf("FROM %s\n", baseimage))
	if err != nil {
		return err
	}

	for _, feature := range features {
		err = writer.AppendRaw(fmt.Sprintf("# %s\n", feature.Meta.Name))
		if err != nil {
			return err
		}

		err = writer.AppendFeature(feature)
		if err != nil {
			return err
		}
	}

	err = writer.AppendRaw("CMD /bin/bash\n")
	if err != nil {
		return err
	}

	p.Dockerfile = writer.Bytes()

	return nil
}

// DockerBuild builds a docker image based on the generated Dockerfile.
func (p *Pazuzu) DockerBuild(name string) error {
	client, err := docker.NewClient(p.DockerEndpoint)
	if err != nil {
		return fmt.Errorf("Error: %s", err)
		return err
	}

	t := time.Now()
	inputBuf := bytes.NewBuffer(nil)
	tr := tar.NewWriter(inputBuf)
	err = tr.WriteHeader(&tar.Header{
		Name:       "Dockerfile",
		Size:       int64(len(p.Dockerfile)),
		ModTime:    t,
		AccessTime: t,
		ChangeTime: t,
	})
	if err != nil {
		return err
	}

	_, err = tr.Write(p.Dockerfile)
	if err != nil {
		return err
	}

	err = tr.Close()
	if err != nil {
		return err
	}

	opts := docker.BuildImageOptions{
		Name:         name,
		InputStream:  inputBuf,
		OutputStream: os.Stdout,
	}

	err2 := client.BuildImage(opts)
	if err2 != nil {
		fmt.Errorf("Error: %s", err2)
		return err
	}

	return nil
}
