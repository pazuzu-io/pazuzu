package pazuzu

import (
	"archive/tar"
	"bytes"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"time"

	"github.com/zalando-incubator/pazuzu/shared"
	"github.com/zalando-incubator/pazuzu/storageconnector"
	"os"
	"os/exec"
)

const (
	tempDir    = "/tmp/pazuzu/"
	mountPoint = "/pazuzu/"
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

	if err := p.generateTestSpec(featuresWithDep); err != nil {
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
		err = fmt.Errorf("Error: %s", err2)
		return err
	}

	err = p.testDockerImage(name)

	return err
}

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
		Tty: true,
	}
	exec, err := p.docker.CreateExec(execOpts)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	var errBuf bytes.Buffer

	startExecOpts := docker.StartExecOptions{
		Detach:       false,
		OutputStream: os.Stdout,
		ErrorStream:  &errBuf,
		RawTerminal:  true,
		Tty:          true,
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

func (p *Pazuzu) dockerStart(image string) (*docker.Container, error) {
	var err error
	p.docker, err = docker.NewClient(p.DockerEndpoint)
	if err != nil {
		return nil, err
	}
	opts := docker.CreateContainerOptions{
		Config: &docker.Config{
			Image: image,
			Tty:   true,
			Cmd: []string{
				"/bin/sh",
			},
		},
		HostConfig: &docker.HostConfig{
			Binds: []string{
				tempDir + ":" + mountPoint,
			},
		},
	}

	container, err := p.docker.CreateContainer(opts)
	if err != nil {
		return nil, err
	}

	if err := p.docker.StartContainer(container.ID, nil); err != nil {
		return nil, err
	}

	return container, nil
}

func (p *Pazuzu) dockerStop(ID string) error {
	if err := p.docker.StopContainer(ID, 1); err != nil {
		return err
	}

	if err := p.docker.RemoveContainer(docker.RemoveContainerOptions{
		ID: ID,
	}); err != nil {
		return err
	}

	return nil
}

func (p *Pazuzu) generateTestSpec(features []shared.Feature) error {
	var buffer = bytes.NewBufferString("")
	if err := shared.WriteTestSpec(buffer, features); err != nil {
		return err
	}
	p.TestSpec = buffer.Bytes()
	return nil
}

func (p *Pazuzu) testDockerImage(image string) error {
	os.MkdirAll(tempDir, 0777)

	batsZip := tempDir + "master.zip"

	if err := exec.Command(
		"wget",
		"https://github.com/sstephenson/bats/archive/master.zip",
		"-O", batsZip).Run(); err != nil {
		fmt.Println("Couldn't download bats")
		return err
	}
	if err := exec.Command("unzip", "-o", batsZip, "-d", tempDir).Run(); err != nil {
		fmt.Println("Couldn't unzip bats to " + tempDir)
		return err
	}
	if err := exec.Command("rm", batsZip).Run(); err != nil {
		fmt.Println("Couldn't delete master.zip")
		return err
	}
	if err := exec.Command("cp", shared.TestSpecFilename, tempDir).Run(); err != nil {
		fmt.Println("Couldn't copy test.bats file to " + tempDir)
		return err
	}

	container, err := p.dockerStart(image)
	if err != nil {
		fmt.Println("Couldn't start docker container")
		fmt.Println(err)
		return err
	}

	if err := p.dockerExec(
		container.ID,
		fmt.Sprintf("%sbats-master/install.sh /usr/local && /usr/local/bin/bats -p %s%s", mountPoint, mountPoint, shared.TestSpecFilename)); err != nil {
		fmt.Println("Couldn't exec test commands on container")
		fmt.Println(err)
		return err
	}

	if err := p.dockerStop(container.ID); err != nil {
		fmt.Println("Couldn't stop container")
		return err
	}

	if err := os.RemoveAll(tempDir); err != nil {
		fmt.Println("Couldn't delete " + tempDir)
		return err
	}

	return nil
}
