package pazuzu

import (
	"fmt"
	"github.com/satori/go.uuid"
	"os"
	"path"
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

func TestConfigSaveAndLoad(t *testing.T) {
	config := getConfig(t)

	tempFn := path.Join(os.TempDir(), uuid.NewV4().String())
	t.Logf("TEMP-FN = [%s]\n", tempFn)

	defer os.Remove(tempFn)

	const ExpectBase = "MyBase"
	const ExpectStorageType = "git"
	const ExpectGitURL = "some-git-url"

	const UnexpectBase = "NotMyBase"
	const UnexpectStorageType = "memory"
	const UnexpectGitURL = "not-git-url"

	config.SetBase(ExpectBase)
	config.SetStorageType(ExpectStorageType)
	config.Git.SetURL(ExpectGitURL)

	errSave := config.SaveToFile(tempFn)
	if errSave != nil {
		t.Fatalf("SaveToFile FAIL! filename=[%s], reason=[%v]\n",
			tempFn, errSave)
	}

	config.SetBase(UnexpectBase)
	config.SetStorageType(UnexpectStorageType)
	config.Git.SetURL(UnexpectGitURL)

	config.LoadFromFile(tempFn)

	if config.Base != ExpectBase {
		t.Fatalf("'Base' config-val should be equals with [%s], but [%s]",
			ExpectBase, config.Base)
	}

	if config.StorageType != ExpectStorageType {
		t.Fatalf("'StorageType' config-val should be equals with [%s], but [%s]",
			ExpectStorageType, config.StorageType)
	}

	if config.Git.URL != ExpectGitURL {
		t.Fatalf("'Git.URL' config-val should be equals with [%s], but [%s]",
			ExpectGitURL, config.Git.URL)
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
