package pazuzu

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// Pazuzu defines pazuzu config.
type Pazuzu struct {
	registry       PazuzuRegistry
	dockerfile     []byte
	testSpec       string
	dockerEndpoint string
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
	fs, err := p.registry.GetFeatures(features)
	if err != nil {
		return err
	}

	err = p.fetchFiles(fs)
	if err != nil {
		return err
	}

	err = p.generateDockerfile(baseimage, fs)
	if err != nil {
		return err
	}

	err = p.generateTestSpec(fs)
	if err != nil {
		return err
	}

	return nil
}

// Delete downloaded files
func (p *Pazuzu) Cleanup() {
	for _, filePath := range p.files {
		err := os.Remove(filePath)
		if err != nil {
			log.Printf("Can't remove file '%s': %s", filePath, err)
		}
	}
}

func (p *Pazuzu) fetchFiles(features []Feature) error {
	files := make(map[string]string)
	for _, feature := range features {
		featureFiles, err := p.registry.FetchFiles(feature)
		if err != nil {
			return err
		}

		for origFName, tmpFName := range featureFiles {
			files[feature.Name+"/"+origFName] = tmpFName
		}
	}

	p.files = files

	return nil
}

// generate in-memory Dockerfile from list of features.
func (p *Pazuzu) generateDockerfile(baseimage string, features []Feature) error {
	writer := NewDockerfileWriter()

	err := writer.AppendRaw(fmt.Sprintf("FROM %s\n", baseimage))
	if err != nil {
		return err
	}

	for _, feature := range features {
		err = writer.AppendRaw(fmt.Sprintf("# %s\n", feature.Name))
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

	p.dockerfile = writer.Bytes()

	return nil
}

type TestSpec struct {
	Feature string `json:"feature"`
	Cmd     string `json:"cmd"`
}

// generate test spec from list of features.
func (p *Pazuzu) generateTestSpec(features []Feature) error {
	f, err := os.Create(p.testSpec)
	if err != nil {
		return err
	}
	defer f.Close()

	var specs []TestSpec

	for _, feature := range features {
		spec := TestSpec{
			Feature: feature.Name,
			Cmd:     feature.TestInstruction,
		}

		specs = append(specs, spec)
	}

	enc := json.NewEncoder(f)
	return enc.Encode(specs)
}

// read test specs from file.
func (p *Pazuzu) readTestSpec() ([]TestSpec, error) {
	f, err := os.Open(p.testSpec)
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(f)

	var specs []TestSpec

	err = dec.Decode(&specs)
	if err != nil {
		return nil, err
	}

	return specs, nil
}

func (p *Pazuzu) appendFeaturesFiles(tr *tar.Writer) error {
	for fileName, filePath := range p.files {
		tmpFile, err := os.Open(filePath)
		if err != nil {
			return err
		}

		stat, err := tmpFile.Stat()
		if err != nil {
			return err
		}

		now := time.Now()

		err = tr.WriteHeader(&tar.Header{
			Name:       fileName,
			Size:       stat.Size(),
			ModTime:    now,
			AccessTime: now,
			ChangeTime: now,
		})

		if err != nil {
			return err
		}

		_, err = io.Copy(tr, tmpFile)
		if err != nil {
			return err
		}

		err = tmpFile.Close()
		if err != nil {
			return err
		}
	}
	return nil
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
	err = tr.WriteHeader(&tar.Header{
		Name:       "Dockerfile",
		Size:       int64(len(p.dockerfile)),
		ModTime:    t,
		AccessTime: t,
		ChangeTime: t,
	})
	if err != nil {
		return err
	}

	_, err = tr.Write(p.dockerfile)
	if err != nil {
		return err
	}

	err = p.appendFeaturesFiles(tr)
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

	err = client.BuildImage(opts)
	if err != nil {
		return err
	}

	return nil
}

// Start a docker container running /bin/bash.
func (p *Pazuzu) dockerStart(image string) (*docker.Container, error) {
	var err error
	p.docker, err = docker.NewClient(p.dockerEndpoint)
	if err != nil {
		return nil, err
	}

	opts := docker.CreateContainerOptions{
		Config: &docker.Config{
			Image: image,
			Tty:   true,
			Cmd: []string{
				"/bin/bash",
			},
		},
	}

	container, err := p.docker.CreateContainer(opts)
	if err != nil {
		return nil, err
	}

	err = p.docker.StartContainer(container.ID, nil)
	if err != nil {
		return nil, err
	}

	return container, nil
}

// Execute command in docker container.
// The command will run in /bin/bash -c ''.
func (p *Pazuzu) dockerExec(ID string, cmd string) error {
	execOpts := docker.CreateExecOptions{
		Container:    ID,
		AttachStdin:  false,
		AttachStdout: true,
		AttachStderr: true,
		Cmd: []string{
			"/bin/bash",
			"-c",
			cmd,
		},
		Tty: false,
	}

	exec, err := p.docker.CreateExec(execOpts)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	var errBuf bytes.Buffer

	startExecOpts := docker.StartExecOptions{
		Detach:       false,
		OutputStream: &buf,
		ErrorStream:  &errBuf,
		RawTerminal:  false,
		Tty:          false,
	}

	err = p.docker.StartExec(exec.ID, startExecOpts)
	if err != nil {
		return err
	}

	inspect, err := p.docker.InspectExec(exec.ID)
	if err != nil {
		return err
	}

	if inspect.ExitCode > 0 {
		return fmt.Errorf("exit code %d: %s", inspect.ExitCode, buf.String())
	}

	return nil
}

// Stop docker container by ID
func (p *Pazuzu) dockerStop(ID string) error {
	err := p.docker.StopContainer(ID, 1)
	if err != nil {
		return err
	}

	rmOpts := docker.RemoveContainerOptions{
		ID: ID,
	}

	err = p.docker.RemoveContainer(rmOpts)
	if err != nil {
		return err
	}

	return nil
}

// RunTestSpec runs the tests against the given image.
func (p *Pazuzu) RunTestSpec(image string) error {
	specs, err := p.readTestSpec()
	if err != nil {
		return err
	}

	container, err := p.dockerStart(image)
	if err != nil {
		return err
	}

	failedTests := 0

	for _, spec := range specs {
		fmt.Printf("Running test spec for feature '%s':\n\t%s\n",
			spec.Feature, spec.Cmd)
		err = p.dockerExec(container.ID, spec.Cmd)
		if err != nil {
			fmt.Println(err)
			failedTests++
		}
	}

	err = p.dockerStop(container.ID)
	if err != nil {
		return err
	}

	if failedTests > 0 {
		return fmt.Errorf("number of failing tests: %d", failedTests)
	}

	return nil
}
