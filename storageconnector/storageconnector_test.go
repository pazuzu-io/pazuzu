package storageconnector

import (
	"testing"
	"regexp"
	"github.com/stretchr/testify/assert"
	"github.com/zalando-incubator/pazuzu/shared"
)

func searchMetaAndFindResultTest(t *testing.T, name string, expected []shared.FeatureMeta, storage StorageReader) {
	t.Run("Search meta and find results", func(t *testing.T) {
		assert := assert.New(t)
		term, _ := regexp.Compile(name)
		result, err := storage.SearchMeta(term)
		assert.Nil(err)

		for i, meta := range result {
			assert.Equal(meta.Name, expected[i].Name)
			assert.Equal(meta.Description, expected[i].Description)
			assert.Equal(meta.Author, expected[i].Author)
			assert.Equal(meta.Dependencies, expected[i].Dependencies)
		}
	})
}

func searchMetaAndFindNothingTest(t *testing.T, name string, storage StorageReader) {
	t.Run("Search meta and find nothing", func(t *testing.T) {
		assert := assert.New(t)
		term, _ := regexp.Compile(name)
		result, err := storage.SearchMeta(term)
		assert.Nil(err)
		assert.Empty(result)
	})
}

func getExistingFeatureMetaTest(t *testing.T, name string, storage StorageReader) {
	t.Run("Get existing feature meta", func(t *testing.T) {
		assert := assert.New(t)
		meta, err := storage.GetMeta(name)
		assert.Nil(err)
		assert.Equal(name, meta.Name)
	})
}

func getNonExistingFeatureMetaTest(t *testing.T, name string, storage StorageReader) {
	t.Run("Get non existing feature meta", func(t *testing.T) {
		_, err := storage.GetMeta(name)
		assert.NotNil(t, err)
	})
}

func getExistingFeatureTest(t *testing.T, name string, storage StorageReader) {
	t.Run("Get existing feature meta", func(t *testing.T) {
		assert := assert.New(t)
		feature, err := storage.GetFeature(name)
		assert.Nil(err)
		assert.Equal(name, feature.Meta.Name)
		assert.NotEmpty(feature.Snippet)
	})
}

func getExistingFeatureWithoutSnippetTest(t *testing.T, name string, storage StorageReader) {
	t.Run("Get existing feature meta", func(t *testing.T) {
		assert := assert.New(t)
		feature, err := storage.GetFeature(name)
		assert.Nil(err)
		assert.Equal(name, feature.Meta.Name)
		assert.Empty(feature.Snippet)
	})
}

func getNonExistingFeatureTest(t *testing.T, name string, storage StorageReader) {
	t.Run("Get non existing feature meta", func(t *testing.T) {
		_, err := storage.GetFeature(name)
		assert.NotNil(t, err)
	})
}

func resolveFeaturesTest(t *testing.T, message string, names []string, expected map[string]shared.Feature, storage StorageReader) {
	t.Run(message, func(t *testing.T) {
		assert := assert.New(t)
		_, result, err := storage.Resolve(names...)
		assert.Nil(err)
		assert.Equal(len(result), len(expected))
		for k, v := range result {
			assert.NotNil(expected[k])
			assert.Equal(v.Meta.Name, expected[k].Meta.Name)
		}
	})
}

func resolveSingleFeatureWithoutDependenciesTest(t *testing.T, name string, storage StorageReader) {
	t.Run("Resolve single feature without dependencies", func(t *testing.T) {
		assert := assert.New(t)
		_, result, err := storage.Resolve(name)
		assert.Nil(err)
		assert.Equal(len(result), 1)
		assert.NotNil(result[name])
		assert.Equal(name, result[name].Meta.Name)
	})
}

func resolveEmptyFeaturesTest(t *testing.T, storage StorageReader) {
	t.Run("Resolve empty feature list", func(t *testing.T) {
		assert := assert.New(t)
		_, result, err := storage.Resolve()
		assert.Nil(err)
		assert.Equal(map[string]shared.Feature{}, result)
	})
}

func resolveNonExistingFeatureTest(t *testing.T, name string,  storage StorageReader) {
	t.Run("Resolve non-existing feature", func(t *testing.T) {
		_, _, err := storage.Resolve(name)
		assert.NotNil(t, err)
	})
}