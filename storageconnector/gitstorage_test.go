package storageconnector

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/src-d/go-git.v4"
	"github.com/zalando-incubator/pazuzu/shared"
	"github.com/davecgh/go-spew/spew"
)

var (
	testRepository = filepath.Join("fixtures", "git")
	storage StorageReader
)

func TestMain(m *testing.M) {
	spew.Config = spew.ConfigState{
		DisableCapacities: true,
		DisablePointerAddresses: true,
	}
	repo, err := git.NewFilesystemRepository(testRepository)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	storage = &GitStorage{repo: repo}
	os.Exit(m.Run())
}

func TestGitStorage_SearchMeta(t *testing.T) {
	expected := []shared.FeatureMeta{{
		Name: "A-java-lein",
		Description: "Java + Leiningen",
		Author: "",
		Dependencies: []string{"java", "leiningen"},
	}, {
		Name: "B-java-node",
		Description: "Java + Node",
		Author: "",
		Dependencies: []string{"java", "node"},
	}, {
		Name: "java",
		Description: "basic java feature",
		Author: "",
	},
	}
	searchMetaAndFindResultTest(t, "java", expected, storage)
	searchMetaAndFindNothingTest(t, "NotAFeature", storage)
}

func TestGitStorage_GetMeta(t *testing.T) {
	getExistingFeatureMetaTest(t, "java", storage)
	getNonExistingFeatureMetaTest(t, "NotAFeature", storage)
}

func TestGitStorage_Get(t *testing.T) {
	getExistingFeatureTest(t, "java", storage)
	getExistingFeatureWithoutSnippetTest(t, "A-java-lein", storage)
	getNonExistingFeatureTest(t, "NotAFeature", storage)
}

func TestGitStorage_Resolve(t *testing.T) {
	t.Run("FeatureWithoutDependencies", func(t *testing.T) {
		expected := 1

		list, features, err := storage.Resolve("java")
		if err != nil {
			t.Error(err)
		}
		if len(list) != expected {
			t.Errorf("Feature count should be %d but was %d", expected, len(list))
		}

		if len(features) != expected {
			t.Errorf("Feature count should be %d but was %d", expected, len(features))
		}
		if _, ok := features["java"]; !ok {
			t.Error("Missing feature 'java'")
		}
	})

	t.Run("FeatureWithTwoDependencies", func(t *testing.T) {
		expected := 3

		list, features, err := storage.Resolve("A-java-lein")
		if err != nil {
			t.Error(err)
		}

		if len(list) != expected {
			t.Errorf("Feature count should be %d but was %d", expected, len(list))
		}

		if len(features) != expected {
			t.Errorf("Feature count should be %d but was %d", expected, len(features))
		}

		if _, ok := features["A-java-lein"]; !ok {
			t.Error("Missing feature 'A-java-lein")
		}
		if _, ok := features["leiningen"]; !ok {
			t.Error("Missing feature 'leiningen")
		}
		if _, ok := features["java"]; !ok {
			t.Error("Missing feature 'java")
		}
	})
}
