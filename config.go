package pazuzu

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/cevaris/ordered_map"
	"github.com/jinzhu/copier"
	"gopkg.in/yaml.v2"

	"github.com/zalando-incubator/pazuzu/storageconnector"
)

const (
	UserConfigFilenamePart = ".pazuzu-cli.yaml"

	// BaseImage : Base feature.
	BaseImage = "ubuntu:14.04"

	// StorageTypeRegistry: pazuzu-registry storage
	StorageTypeRegistry = "registry"
	// Default hostnamefor the registry
	DefaultRegistryHostname = "localhost"
	// StorageTypeRegistry: pazuzu-registry storage
	DefaultRegistryPort = 8080
	// Default hostnamefor the registry
	DefaultRegistryScheme = "http"
)

var config Config

// registryConfig : config structure for Registry-storage
type RegistryConfig struct {
	Hostname string `yaml:"hostname" setter:"SetHostname" help:"Hostname String"`
	Port     int    `yaml:"port" setter:"SetPort" help:"Port Integer"`
	Scheme   string `yaml:"scheme" setter:"SetScheme" help:"Scheme String"`
}

// Config : actual config data structure.
type Config struct {
	Base        string         `yaml:"base" setter:"SetBase" help:"Base image name and tag (ex: 'ubuntu:14.04')"`
	StorageType string         `yaml:"storage" setter:"SetStorageType" help:"Storage-type(registry) "`
	Registry    RegistryConfig `yaml:"registry" help:"Pazuzu-registry configs"`
}

// SetBase : Setter of "Base".
func (c *Config) SetBase(base string) {
	c.Base = base
}

// SetStorageType : Setter of "StorageType".
func (c *Config) SetStorageType(storageType string) {
	c.StorageType = storageType
}

// SetRegistryHostname : Setter of RegistryConfig.Hostname.
func (r *RegistryConfig) SetHostname(hostname string) {
	r.Hostname = hostname
}

// SetRegistryPort : Setter of RegistryConfig.Port.
func (r *RegistryConfig) SetPort(port int) {
	r.Port = port
}

// SetRegistryScheme : Setter of RegistryConfig.Scheme.
func (r *RegistryConfig) SetScheme(scheme string) {
	r.Scheme = scheme
}

// InitDefaultConfig : Initialize config variable with defaults. (Does not loading configuration file)
func InitDefaultConfig() {
	config = Config{
		StorageType: "registry",
		Base:        BaseImage,
		Registry:    RegistryConfig{DefaultRegistryHostname, DefaultRegistryPort, DefaultRegistryScheme},
	}
}

// NewConfig : Please call this function before GetConfig and only once in your application.
// Attempts load config file, but when it fails just use default configuration.
func NewConfig() error {
	InitDefaultConfig()
	config.Load()
	configMirror = config.InitConfigFieldMirrors()
	return nil
}

// GetConfig : get loaded config.
func GetConfig() *Config {
	return &config
}

// GetStorageReader : create new StorageReader by StorageType of given config.
func GetStorageReader(config Config) (storageconnector.StorageReader, error) {
	switch config.StorageType {
	case StorageTypeRegistry:
		return storageconnector.NewRegistryStorage(config.Registry.Hostname, config.Registry.Port, config.Registry.Scheme, nil)
	}

	return nil, fmt.Errorf("unknown storage type '%s'", config.StorageType)
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
	aVal reflect.Value, aType reflect.Type,
	addressableVal reflect.Value,
	ancestors []reflect.StructField) error

func (c *Config) TraverseEachField(cb ConfigTraverseFunc) error {
	aType := reflect.TypeOf(*c)
	aVal := reflect.ValueOf(*c)
	addressableVal := reflect.ValueOf(c)
	return traverseEachFieldRecur(aVal, aType, addressableVal, []reflect.StructField{}, cb)
}

func traverseEachFieldRecur(aVal reflect.Value, aType reflect.Type, addressableVal reflect.Value,
	ancestors []reflect.StructField, cb ConfigTraverseFunc) error {
	//
	for i := 0; i < aType.NumField(); i++ {
		field := aType.Field(i)
		if field.Type.Kind() == reflect.Struct {
			bType := field.Type
			f := reflect.Indirect(aVal).FieldByName(field.Name)
			f2 := reflect.Indirect(addressableVal).FieldByName(field.Name)
			err := traverseEachFieldRecur(f, bType, f2.Addr(), append(ancestors, field), cb)
			if err != nil {
				return err
			}
		} else {
			err := cb(field, aVal, aType, addressableVal, ancestors)
			if err != nil {
				return err
			}
		}
	}
	//
	return nil
}

type ConfigFieldMirror struct {
	Help   string
	Repr   string
	Setter reflect.Value
}

type ConfigMirror struct {
	M *ordered_map.OrderedMap
	C *Config
}

var configMirror *ConfigMirror

func (c *Config) InitConfigFieldMirrors() *ConfigMirror {
	m := ordered_map.NewOrderedMap()
	_ = c.TraverseEachField(func(field reflect.StructField,
		aVal reflect.Value, aType reflect.Type, addressableVal reflect.Value,
		ancestors []reflect.StructField) error {
		//
		configPath := makeConfigPathString(ancestors, field)
		tag := field.Tag
		help := ""
		repr := ""
		setter := reflect.ValueOf(nil)
		// setter.
		setterName := field.Tag.Get("setter")
		if len(setterName) >= 1 {
			setter = addressableVal.MethodByName(setterName)
		}
		// help
		help = tag.Get("help")
		// repr
		f := reflect.Indirect(aVal).FieldByName(field.Name)
		repr = toReprFromReflectValue(f)
		//
		m.Set(configPath, &ConfigFieldMirror{
			Help:   help,
			Repr:   repr,
			Setter: setter,
		})
		//
		return nil
	})
	//
	return &ConfigMirror{M: m, C: c}
}

func toReprFromReflectValue(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Bool:
		b := v.Bool()
		return fmt.Sprintf("%v", b)
	case reflect.Int:
		n := v.Int()
		return fmt.Sprintf("%v", n)
	case reflect.String:
		return v.String()
	default:
		return v.String()
	}
}

func joinConfigPath(path []reflect.StructField) string {
	yamlNames := []string{}
	for _, field := range path {
		yamlNames = append(yamlNames, field.Tag.Get("yaml"))
	}
	return strings.Join(yamlNames, ".")
}

func makeConfigPathString(ancestors []reflect.StructField, field reflect.StructField) string {
	return joinConfigPath(append(ancestors, field))
}

func GetConfigMirror() *ConfigMirror {
	return configMirror
}

func (c *ConfigMirror) GetKeys() []string {
	iter := c.M.IterFunc()
	result := []string{}
	for kv, ok := iter(); ok; kv, ok = iter() {
		result = append(result, kv.Key.(string))
	}
	return result
}

func (c *ConfigMirror) GetHelp(key string) (string, error) {
	v, ok := c.M.Get(key)
	if ok {
		return v.(*ConfigFieldMirror).Help, nil
	}
	return "", ErrNotFound
}

func (c *ConfigMirror) GetRepr(key string) (string, error) {
	v, ok := c.M.Get(key)
	if ok {
		return v.(*ConfigFieldMirror).Repr, nil
	}
	return "", ErrNotFound
}

func valToReflectValue(setter reflect.Value, val string) (reflect.Value, error) {
	switch setter.Type().In(0).Kind() {
	case reflect.Int:
		integerArg, err := strconv.Atoi(val)
		return reflect.ValueOf(integerArg), err

	default:
		return reflect.ValueOf(val), nil
	}
}

func (c *ConfigMirror) SetConfig(key string, val string) error {
	v, ok := c.M.Get(key)
	if ok {
		setter := v.(*ConfigFieldMirror).Setter
		if !setter.IsValid() {
			fmt.Println("INVALID SETTER!!!")
			return ErrNotImplemented
		}
		arg, err := valToReflectValue(setter, val)
		if err != nil {
			return ErrInvalidConfigValue
		}
		setter.Call([]reflect.Value{arg})
		return nil
	}
	return ErrNotFound
}
