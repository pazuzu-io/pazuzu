package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
)

func getConfig(t *testing.T) *Config {
	InitDefaultConfig()
	config := GetConfig()
	if len(config.Base) <= 0 {
		t.Error("Please fill 'Base' property of initial config.")
	}

	return config
}

// TestConfigSetBase : tests on its mutability.
func TestConfigSetBase(t *testing.T) {
	config := getConfig(t)

	beforeBase := config.Base
	const someAnotherBase = "foo-bar-zoo-spam-eggs"
	config.SetBase(someAnotherBase)

	if strings.Compare(config.Base, someAnotherBase) != 0 {
		t.Errorf("Unexpected value [%v]", config.Base)
	}

	if strings.Compare(config.Base, beforeBase) == 0 {
		t.Errorf("Not changed 'Base' value!")
	}
}

// TestConfigSetStorageType ...
func TestConfigSetStorageType(t *testing.T) {
	config := getConfig(t)

	beforeStorageType := config.StorageType
	config.SetStorageType("foo")

	if strings.Compare(config.StorageType, "foo") != 0 {
		t.Errorf("SetStorage FAIL! [%v]", config.StorageType)
	}

	if strings.Compare(config.StorageType, beforeStorageType) == 0 {
		t.Error("No changes made.")
	}
}

func TestConfigSetRegistry(t *testing.T) {
	config := getConfig(t)

	beforeHost := config.Registry.Hostname
	config.Registry.SetHostname("betterhost")

	if strings.Compare(config.Registry.Hostname, "betterhost") != 0 {
		t.Errorf("SetHostname FAIL! [%v]", config.Registry.Hostname)
	}
	if strings.Compare(config.Registry.Hostname, beforeHost) == 0 {
		t.Error("No changes made.")
	}

	beforePort := config.Registry.Port
	config.Registry.SetPort(8081)
	if config.Registry.Port-beforePort != 1 {
		t.Errorf("SetPort FAIL! [%v]", config.Registry.Port)
	}
}

func TestConfigSaveAndLoad(t *testing.T) {
	config := getConfig(t)

	tempFile, err := ioutil.TempFile(os.TempDir(), "pazuzu_config_test")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(tempFile.Name())

	const ExpectBase = "MyBase"
	const ExpectStorageType = "registry"

	const UnexpectBase = "NotMyBase"
	const UnexpectStorageType = "memory"

	config.SetBase(ExpectBase)
	config.SetStorageType(ExpectStorageType)

	errSave := config.SaveToFile(tempFile.Name())
	if errSave != nil {
		t.Fatalf("SaveToFile FAIL! filename=[%s], reason=[%v]\n",
			tempFile.Name(), errSave)
	}

	config.SetBase(UnexpectBase)
	config.SetStorageType(UnexpectStorageType)

	config.LoadFromFile(tempFile.Name())

	if config.Base != ExpectBase {
		t.Fatalf("'Base' config-val should be equals with [%s], but [%s]",
			ExpectBase, config.Base)
	}

	if config.StorageType != ExpectStorageType {
		t.Fatalf("'StorageType' config-val should be equals with [%s], but [%s]",
			ExpectStorageType, config.StorageType)
	}

}

func TestConfigUserHomeDir(t *testing.T) {
	home := UserHomeDir()
	fmt.Printf("home = [%v]\n", home)
	if len(home) <= 0 {
		t.Errorf("Too short for home-dir! [%v]", home)
	}
}

func TestConfigUserConfigFilename(t *testing.T) {
	filename := UserConfigFilename()
	fmt.Printf("user-config-filename = [%v]\n", filename)
}

func getConfigMirror(t *testing.T) *ConfigMirror {
	cfg := getConfig(t)
	return cfg.InitConfigFieldMirrors()
}

func dummySetterWithInteger(value int) {}

func TestValueToReflectValue(t *testing.T) {
	setter := reflect.ValueOf(dummySetterWithInteger)
	val, err := valToReflectValue(setter, "10")
	if err != nil {
		t.Error("Couldn't handle integer parameter type")
	}
	if val.Int() != reflect.ValueOf(10).Int() {
		t.Error("Couldn't parse integer correctly.")
	}
}
