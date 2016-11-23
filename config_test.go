package pazuzu

import (
	"fmt"
	"github.com/satori/go.uuid"
	"os"
	"path"
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

	defer func() {
		_ = os.Remove(tempFn)
	}()

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

func getConfigMirror(t *testing.T) *ConfigMirror {
	cfg := getConfig(t)
	return cfg.InitConfigFieldMirrors()
}

func TestConfigMirrorGetRepr(t *testing.T) {
	c := getConfigMirror(t)
	strBase, err1 := c.GetRepr("base")
	strGitURL, err2 := c.GetRepr("git.url")
	if err1 != nil || err2 != nil {
		t.Errorf("Error? base=%v or git.url=%v\n", err1, err2)
	}
	if true != (len(strBase) > 0 && len(strGitURL) > 0) {
		t.Errorf("Unexpected base=[%s] or git.url=[%s]\n", strBase, strGitURL)
	}
}

func TestConfigMirrorGetHelp(t *testing.T) {
	c := getConfigMirror(t)
	strBase, err1 := c.GetHelp("base")
	strGitURL, err2 := c.GetHelp("git.url")
	if err1 != nil || err2 != nil {
		t.Errorf("Error? base=%v or git.url=%v\n", err1, err2)
	}
	if true != (len(strBase) > 0 && len(strGitURL) > 0) {
		t.Errorf("Unexpected base=[%s] or git.url=[%s]\n", strBase, strGitURL)
	}

}

func TestConfigMirrorGetKeys(t *testing.T) {
	c := getConfigMirror(t)
	keys := c.GetKeys()
	count := 0
	for _, k := range keys {
		if k == "base" {
			count++
		}
		if k == "git.url" {
			count++
		}
	}
	if count != 2 {
		t.Errorf("Should be 2 == %v\n", count)
	}
}

func TestConfigMirrorSetConfig(t *testing.T) {
	c := getConfigMirror(t)
	cfg := c.C
	_ = c.SetConfig("base", "base123")
	_ = c.SetConfig("git.url", "git.url.456")
	if cfg.Base != "base123" || cfg.Git.URL != "git.url.456" {
		t.Errorf("ConfigMirror SetConfig FAIL! [%v]\n", cfg)
	}
}
