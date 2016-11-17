package pazuzu

import (
	"fmt"
	"github.com/zalando-incubator/pazuzu/storageconnector"
)

const (
	// URL : default features-repo.
	URL = "https://github.com/Sangdol/pazuzu-test-repo.git"
	// BaseImage : Base feature.
	BaseImage = "ubuntu:14.04"

	// StorageTypeGit : Git storage type.
	StorageTypeGit = "git"
	// StorageTypeMemory : Memory storage type.
	StorageTypeMemory = "memory"
)

var config Config

// GitConfig : config structure for Git-storage.
type GitConfig struct {
	URL string `yaml:"url"`
}

// MemoryConfig : config structure for Memory-storage.
type MemoryConfig struct {
	InitialiseRandom bool `yaml:"random_init"`
	RandomSetSize    int  `yaml:"random_size"`
}

// Config : actual config data structure.
type Config struct {
	Base        string       `yaml:"base"`
	StorageType string       `yaml:"storage"`
	Git         GitConfig    `yaml:"git"`
	Memory      MemoryConfig `yaml:"memory"`
}

// SetBase : Setter of "Base".
func (c *Config) SetBase(base string) {
	c.Base = base
}

// NewConfig : Please call this function before GetConfig and only once in your application.
func NewConfig() error {
	// TODO: add read from $HOME/.pazuzu/config and return error if fail
	// viper library is planned to be used here
	config = Config{
		StorageType: "git",
		Base:        BaseImage,
		Git:         GitConfig{URL: URL},
		Memory: MemoryConfig{
			InitialiseRandom: false,
		},
	}
	return nil
}

// GetConfig : get loaded config.
func GetConfig() Config {
	return config
}

// GetStorageReader : create new StorageReader by StorageType of given config.
func GetStorageReader(config Config) (storageconnector.StorageReader, error) {
	switch config.StorageType {
	case StorageTypeMemory:
		data := []storageconnector.Feature{}
		if config.Memory.InitialiseRandom {
			data = generateRandomFeatures(config.Memory.RandomSetSize)
		}

		return storageconnector.NewMemoryStorage(data), nil // implement a generator of random list of features?
	case StorageTypeGit:
		return storageconnector.NewGitStorage(config.Git.URL)
	}

	return nil, fmt.Errorf("unknown storage type '%s'", config.StorageType)
}

func generateRandomFeatures(setsize int) []storageconnector.Feature {
	// TODO: implement in case of need
	return []storageconnector.Feature{}
}
