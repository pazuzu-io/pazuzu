package storageconnector

import (
	"testing"
	"regexp"
	"os"
)


const (
	// TODO setup and use a dedicated (local?) test repository
	testRepository = "https://github.com/Sangdol/pazuzu-test-repo.git"
)

var (
	storage StorageReader
)

func TestMain(m *testing.M) {
	var err error
	storage, err = NewStorageReader(testRepository)
	if err != nil {
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func TestGitStorage_SearchMeta(t *testing.T) {
	t.Run("FeatureContainingJava", func(t *testing.T) {
		name, _ := regexp.Compile("java")
		expected := 3

		features, err := storage.SearchMeta(SearchParams{Name: name})
		if err != nil {
			t.Error(err)
		}

		if len(features) != expected {
			t.Errorf("Feature count should be %d but was %d", expected, len(features))
		}

		if features[0].Name != "A-java-lein" {
			t.Errorf("Name of feature 0 should be 'A-java-lein' but was '%s'", features[0].Name)
		}
		if features[1].Name != "B-java-node" {
			t.Errorf("Name of feature 1 should be 'B-java-node' but was '%s'", features[1].Name)
		}
		if features[2].Name != "java" {
			t.Errorf("Name of feature 2 should be 'java' but was '%s'", features[2].Name)
		}
	})

	t.Run("FeatureContainingJavaLimit2", func(t *testing.T) {
		name, _ := regexp.Compile("java")
		expected := 2

		features, err := storage.SearchMeta(SearchParams{Name: name, Limit: 2})
		if err != nil {
			t.Error(err)
		}

		if len(features) != expected {
			t.Errorf("Feature count should be %d but was %d", expected, len(features))
		}

		if features[0].Name != "A-java-lein" {
			t.Errorf("Name of feature 0 should be 'A-java-lein' but was '%s'", features[0].Name)
		}
		if features[1].Name != "B-java-node" {
			t.Errorf("Name of feature 1 should be 'B-java-node' but was '%s'", features[1].Name)
		}
	})

	t.Run("FeatureContainingJavaOffset1", func(t *testing.T) {
		name, _ := regexp.Compile("java")
		expected := 2

		features, err := storage.SearchMeta(SearchParams{Name: name, Offset: 1})
		if err != nil {
			t.Error(err)
		}

		if len(features) != expected {
			t.Errorf("Feature count should be %d but was %d", expected, len(features))
		}

		if features[0].Name != "B-java-node" {
			t.Errorf("Name of feature 1 should be 'B-java-node' but was '%s'", features[0].Name)
		}
		if features[1].Name != "java" {
			t.Errorf("Name of feature 2 should be 'java' but was '%s'", features[1].Name)
		}
	})

	t.Run("FeatureContainingJavaOffset1Limit1", func(t *testing.T) {
		name, _ := regexp.Compile("java")
		expected := 1

		features, err := storage.SearchMeta(SearchParams{Name: name, Offset: 1, Limit: 1})
		if err != nil {
			t.Error(err)
		}

		if len(features) != expected {
			t.Errorf("Feature count should be %d but was %d", expected, len(features))
		}

		if features[0].Name != "B-java-node" {
			t.Errorf("Name of feature 2 should be 'B-java-node' but was '%s'", features[0].Name)
		}
	})

}

func TestGitStorage_GetMeta(t *testing.T) {
	t.Run("ExistingFeature", func(t *testing.T) {
		meta, err := storage.GetMeta("java")
		if err != nil {
			t.Error(err)
		}

		if meta.Name != "java" {
			t.Errorf("Feature name shoule be 'java' but was '%s'", meta.Name)
		}
	})

	t.Run("NonExistingFeature", func(t *testing.T) {
		_, err := storage.GetMeta("reallynotafeature")
		if err == nil {
			t.Error("Error expected when getting meatdata for non existing feature")
		}
	})
}

func TestGitStorage_Get(t *testing.T) {
	t.Run("ExistingFeatureWithoutSnippet", func(t *testing.T) {
		feature, err := storage.Get("A-java-lein")
		if err != nil {
			t.Error(err)
		}

		if feature.Meta.Name != "A-java-lein" {
			t.Errorf("Feature name should be 'A-java-lein' but was '%s'", feature.Meta.Name)
		}

		if feature.Snippet != "" {
			t.Errorf("Feature snippet should be empty but was '%s'", feature.Snippet)
		}
	})

	t.Run("ExistingFeatureWithSnippet", func(t *testing.T) {
		feature, err := storage.Get("java")
		if err != nil {
			t.Error(err)
		}

		if feature.Meta.Name != "java" {
			t.Errorf("Feature name should be 'java' but was '%s'", feature.Meta.Name)
		}

		if feature.Snippet != "" {
			t.Error("Feature snippet should not be empty", feature.Snippet)
		}
	})

	t.Run("NonExistingFeature", func(t *testing.T) {
		_, err := storage.Get("reallynotafeature")
		if err == nil {
			t.Error("Error expected when getting non existing feature")
		}
	})
}

func TestGitStorage_Resolve(t *testing.T) {
	// helper to get index of a feature
	idx := func(features []Feature, name string) int {
		for i, f := range features {
			if f.Meta.Name == name {
				return i
			}
		}
		return -1
	}

	t.Run("FeatureWithoutDependencies", func(t *testing.T) {
		expected := 1

		features, err := storage.Resolve("java")
		if err != nil {
			t.Error(err)
		}

		if len(features) != expected {
			t.Errorf("Feature count should be %d but was %d", expected, len(features))
		}
	})

	t.Run("FeatureWithTwoDependencies", func(t *testing.T) {
		expected := 3

		features, err := storage.Resolve("A-java-lein")
		if err != nil {
			t.Error(err)
		}

		if len(features) != expected {
			t.Errorf("Feature count should be %d but was %d", expected, len(features))
		}

		idxA := idx(features, "java")
		idxB := idx(features, "A-java-lein")
		if idxA < 0 || idxB < 0 {
			t.Error("Feature not found")
		}
		if idxA > idxB {
			t.Error("Features sorted incorrectly")
		}
	})
}
