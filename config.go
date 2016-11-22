package pazuzu

import (
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/zalando-incubator/pazuzu/storageconnector"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
)

const (
	UserConfigFilenamePart = ".pazuzu-cli.yaml"

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
	URL string `yaml:"url" help:"Git Repository URL."`
}

// MemoryConfig : config structure for Memory-storage.
type MemoryConfig struct {
	InitialiseRandom bool `yaml:"random_init" help:"???"`
	RandomSetSize    int  `yaml:"random_size" help:"???"`
}

// Config : actual config data structure.
type Config struct {
	Base        string       `yaml:"base" help:"Base image name and tag (ex: 'ubuntu:14.04')"`
	StorageType string       `yaml:"storage" help:"Storage-type ('git' or 'memory')"`
	Git         GitConfig    `yaml:"git" help:"Git storage configs."`
	Memory      MemoryConfig `yaml:"memory" help:"Memory storage configs."`
}

// SetBase : Setter of "Base".
func (c *Config) SetBase(base string) {
	c.Base = base
}

// SetStorageType : Setter of "StorageType".
func (c *Config) SetStorageType(storageType string) {
	c.StorageType = storageType
}

// SetGit : Setter of Git-Storage specific configuration.
func (c *Config) SetGit(git GitConfig) {
	c.Git = git
}

// SetURL : Setter of GitConfig.URL.
func (g *GitConfig) SetURL(url string) {
	g.URL = url
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

func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func UserConfigFilename() string {
	return filepath.Join(UserHomeDir(), UserConfigFilenamePart)
}

func (c *Config) SaveToWriter(writer io.Writer) error {
	data, err := yaml.Marshal(c)
	_, err = writer.Write(data)
	return err
}

func LoadConfigFromReader(reader io.Reader) (Config, error) {
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return Config{}, err
	}

	c := &Config{}
	err = yaml.Unmarshal(content, c)
	if err != nil {
		return Config{}, err
	}

	return *c, nil
}

func (c *Config) Load() {
	configFn := UserConfigFilename()
	c.LoadFromFile(configFn)
}

func (c *Config) LoadFromFile(configFn string) {
	f, err := os.Open(configFn)
	if err != nil {
		log.Printf("Cannot open config-file [%s], Reason = [%s], SKIP\n",
			configFn, err)
		return
	}
	defer f.Close()

	// replace cfg?
	cfg2, errLoad := LoadConfigFromReader(f)
	if errLoad != nil {
		log.Printf("Cannot load from [%s], Reason = [%s], SKIP\n",
			configFn, errLoad)
		return
	}

	errCopy := copier.Copy(c, &cfg2)
	if errCopy != nil {
		log.Printf("Cannot copy [%v] to [%v], Reason = [%s], SKIP\n",
			cfg2, c, errCopy)
		return
	}
}

func (c *Config) Save() error {
	configFn := UserConfigFilename()
	return c.SaveToFile(configFn)
}

func (c *Config) SaveToFile(configFn string) error {
	f, err := os.Create(configFn)
	if err != nil {
		return err
	}
	defer f.Close()

	errWriter := c.SaveToWriter(f)
	if errWriter != nil {
		return errWriter
	}
	return nil
}

type ConfigTraverseFunc func(field reflect.StructField,
	aVal reflect.Value, aType reflect.Type, ancestors []reflect.StructField) error

func (c *Config) TraverseEachField(cb ConfigTraverseFunc) error {
	aType := reflect.TypeOf(*c)
	aVal := reflect.ValueOf(*c)
	return traverseEachFieldRecur(aVal, aType, []reflect.StructField{}, cb)
}

func traverseEachFieldRecur(aVal reflect.Value, aType reflect.Type,
	ancestors []reflect.StructField, cb ConfigTraverseFunc) error {
	//
	for i := 0; i < aType.NumField(); i++ {
		field := aType.Field(i)
		if field.Type.Kind() == reflect.Struct {
			bType := field.Type
			f := reflect.Indirect(aVal).FieldByName(field.Name)
			//fmt.Printf("\tstruct-val=[%s]\n", f)
			//bVal := reflect.ValueOf(bType)
			err := traverseEachFieldRecur(f, bType, append(ancestors, field), cb)
			if err != nil {
				return err
			}
		} else {
			err := cb(field, aVal, aType, ancestors)
			if err != nil {
				return err
			}
		}
	}
	//
	return nil
}
