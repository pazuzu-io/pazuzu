package pazuzu

import (
	"fmt"
	"github.com/zalando-incubator/pazuzu/storageconnector"
	_ "log"
)

const (
	URL        = "https://github.com/Sangdol/pazuzu-test-repo.git"
	BASE_IMAGE = "ubuntu:14.04"

	StorageTypeGit    = "git"
	StorageTypeMemory = "memory"
)

var config Config

type GitConfig struct {
	Url string `yaml:"url"`
}

type MemoryConfig struct {
	InitialiseRandom bool `yaml:"random_init"`
	RandomSetSize    int  `yaml:"random_size"`
}

type Config struct {
	Base        string       `yaml:"base"`
	StorageType string       `yaml:"storage"`
	Git         GitConfig    `yaml:"git"`
	Memory      MemoryConfig `yaml:"memory"`
}

func NewConfig() error {
	// TODO: add read from $HOME/.pazuzu/config and return error if fail
	// viper library is planned to be used here
	config = Config{
		StorageType: "git",
		Base:        BASE_IMAGE,
		Git:         GitConfig{Url: URL},
		Memory: MemoryConfig{
			InitialiseRandom: false,
		},
	}

	return nil
}

func GetConfig() Config {
	return config
}

func GetStorageReader(config Config) (storageconnector.StorageReader, error) {
	switch config.StorageType {
	case StorageTypeMemory:
		data := []storageconnector.Feature{}
		if config.Memory.InitialiseRandom {
			data = generateRandomFeatures(config.Memory.RandomSetSize)
		}

		return storageconnector.NewMemoryStorage(data), nil // implement a generator of random list of features?
	case StorageTypeGit:
		return storageconnector.NewGitStorage(config.Git.Url)
	}

	return nil, fmt.Errorf("unknown storage type '%s'", config.StorageType)
}

func generateRandomFeatures(setsize int) []storageconnector.Feature {
	// TODO: implement in case of need
	return []storageconnector.Feature{}
}
