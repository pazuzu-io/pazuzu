package storageconnector

import (
	"github.com/zalando-incubator/pazuzu/shared"
	"testing"
)

// TODO issue #159 -> method does not test regex contrary to specs
func TestMemory_SearchMeta(t *testing.T) {
	storage := NewMemoryStorage([]shared.Feature{
		{
			Meta: shared.FeatureMeta{
				Name:         "FeatureA",
				Author:       "SomeAuthor",
				Dependencies: []string{"FeatureB", "FeatureC"},
			},
			Snippet: "",
		},
		{
			Meta: shared.FeatureMeta{
				Name:   "FeatureB",
				Author: "SomeAuthor",
			},
			Snippet: "",
		},
	})
	expected := []shared.FeatureMeta{{
		Name:         "FeatureA",
		Author:       "SomeAuthor",
		Dependencies: []string{"FeatureB", "FeatureC"},
	},
	}
	searchMetaAndFindResultTest(t, "FeatureA", expected, storage)
	searchMetaAndFindNothingTest(t, "NotAFeature", storage)
}

func TestMemory_GetMeta(t *testing.T) {
	storage := NewMemoryStorage([]shared.Feature{
		{
			Meta: shared.FeatureMeta{
				Name:         "FeatureA",
				Author:       "SomeAuthor",
				Dependencies: []string{"FeatureB", "FeatureC"},
			},
			Snippet: "",
		},
		{
			Meta: shared.FeatureMeta{
				Name:   "FeatureB",
				Author: "SomeAuthor",
			},
			Snippet: "",
		},
	})
	getExistingFeatureMetaTest(t, "FeatureA", storage)
	getNonExistingFeatureMetaTest(t, "NotAFeature", storage)
}

func TestMemory_Get(t *testing.T) {
	storage := NewMemoryStorage([]shared.Feature{
		{
			Meta: shared.FeatureMeta{
				Name:         "FeatureA",
				Author:       "SomeAuthor",
				Dependencies: []string{"FeatureB", "FeatureC"},
			},
			Snippet: "RUN something",
		},
		{
			Meta: shared.FeatureMeta{
				Name:   "FeatureB",
				Author: "SomeAuthor",
			},
			Snippet: "",
		},
	})
	getExistingFeatureTest(t, "FeatureA", storage)
	getExistingFeatureWithoutSnippetTest(t, "FeatureB", storage)
	getNonExistingFeatureMetaTest(t, "NotAFeature", storage)
}

func TestMemoryResolve(t *testing.T) {
	storage := NewMemoryStorage([]shared.Feature{
		{
			Meta: shared.FeatureMeta{
				Name:         "FeatureA",
				Author:       "SomeAuthor",
				Dependencies: []string{"FeatureB", "FeatureC"},
			},
			Snippet: "",
		},
		{
			Meta: shared.FeatureMeta{
				Name:   "FeatureB",
				Author: "SomeAuthor",
			},
			Snippet: "",
		},
		{
			Meta: shared.FeatureMeta{
				Name:         "FeatureC",
				Author:       "SomeAuthor",
				Dependencies: []string{"FeatureD"},
			},
			Snippet: "",
		},
		{
			Meta: shared.FeatureMeta{
				Name:   "FeatureD",
				Author: "SomeAuthor",
			},
			Snippet: "",
		},
		{
			Meta: shared.FeatureMeta{
				Name:         "FeatureE",
				Author:       "SomeAuthor",
				Dependencies: []string{"FeatureD"},
			},
			Snippet: "",
		},
		{
			Meta: shared.FeatureMeta{
				Name:         "FeatureF",
				Author:       "SomeAuthor",
				Dependencies: []string{"FeatureG"},
			},
			Snippet: "",
		},
		{
			Meta: shared.FeatureMeta{
				Name:         "FeatureG",
				Author:       "SomeAuthor",
				Dependencies: []string{"FeatureF"},
			},
			Snippet: "",
		},
	})

	resolveEmptyFeaturesTest(t, storage)

	resolveNonExistingFeatureTest(t, "NotAFeatuer", storage)

	resolveFeaturesTest(t, "Resolve single feature", []string{"FeatureA"}, map[string]shared.Feature{
		"FeatureA": {
			Meta: shared.FeatureMeta{
				Name: "FeatureA",
			},
		},
		"FeatureB": {
			Meta: shared.FeatureMeta{
				Name: "FeatureB",
			},
		},
		"FeatureC": {
			Meta: shared.FeatureMeta{
				Name: "FeatureC",
			},
		},
		"FeatureD": {
			Meta: shared.FeatureMeta{
				Name: "FeatureD",
			},
		},
	}, storage)

	resolveFeaturesTest(t, "Resolve multiple features with the same dependencies", []string{"FeatureC", "FeatureE"}, map[string]shared.Feature{
		"FeatureC": {
			Meta: shared.FeatureMeta{
				Name: "FeatureC",
			},
		},
		"FeatureD": {
			Meta: shared.FeatureMeta{
				Name: "FeatureD",
			},
		},
		"FeatureE": {
			Meta: shared.FeatureMeta{
				Name: "FeatureE",
			},
		},
	}, storage)

	resolveFeaturesTest(t, "Resolve multiple features with the circular dependencies", []string{"FeatureF", "FeatureG"}, map[string]shared.Feature{
		"FeatureF": {
			Meta: shared.FeatureMeta{
				Name: "FeatureF",
			},
		},
		"FeatureG": {
			Meta: shared.FeatureMeta{
				Name: "FeatureG",
			},
		},
	}, storage)
}
