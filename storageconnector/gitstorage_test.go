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

	resolveNonExistingFeatureTest(t, "NotAFeature", storage)

	resolveEmptyFeaturesTest(t, storage)

	resolveSingleFeatureWithoutDependenciesTest(t, "java", storage)

	resolveFeaturesTest(t, "Resolve single feature", []string{"A-java-lein"}, map[string]shared.Feature{
		"A-java-lein": {
			Meta: shared.FeatureMeta{
				Name:         "A-java-lein",
			},
		},
		"java": {
			Meta: shared.FeatureMeta{
				Name:         "java",
			},
		},
		"leiningen": {
			Meta: shared.FeatureMeta{
				Name:   "leiningen",
			},
		}, }, storage)

	resolveFeaturesTest(t, "Resolve features with same dependencies", []string{"A-java-lein", "B-java-node"}, map[string]shared.Feature{
		"A-java-lein": {
			Meta: shared.FeatureMeta{
				Name:         "A-java-lein",
			},
		},
		"java": {
			Meta: shared.FeatureMeta{
				Name:         "java",
			},
		},
		"leiningen": {
			Meta: shared.FeatureMeta{
				Name:   "leiningen",
			},
		},
		"B-java-node": {
			Meta: shared.FeatureMeta{
				Name:         "B-java-node",
			},
		},
		"node": {
			Meta: shared.FeatureMeta{
				Name: "node",
			},
		}, }, storage)
}
