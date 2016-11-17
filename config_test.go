package pazuzu

import (
	"strings"
	"testing"
)

func getConfig(t *testing.T) Config {
	errCnf := NewConfig()
	if errCnf != nil {
		t.Errorf("%v", errCnf)
	}

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

func TestConfigGitSetURL(t *testing.T) {
	config := getConfig(t)

	beforeURL := config.Git.URL
	config.Git.SetURL("foobarzoo")

	if strings.Compare(config.Git.URL, "foobarzoo") != 0 {
		t.Errorf("SetURL FAIL! [%v]", config.Git.URL)
	}

	if strings.Compare(config.Git.URL, beforeURL) == 0 {
		t.Error("No changes made.")
	}
}

// TODO: save
func TestConfigSave(ctx *testing.T) {
}

// TODO: load
func TestConfigLoad(ctx *testing.T) {
}
