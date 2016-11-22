package storageconnector

import (
	"io/ioutil"
	"path"
	"strings"

	"regexp"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/core"
	"gopkg.in/yaml.v2"
)

const (
	featureDir               = "features"   // name of the directory where features are located.
	featureFile              = "meta.yml"   // name of the file containing all metadata for a feature.
	featureSnippet           = "Dockerfile" // the file containing the actual docker snippet.
	defaultSearchParamsLimit = 100
)

// yamlFeatureMeta is used for unmarshalling of meta.yml files.
type yamlFeatureMeta struct {
	Description  string
	Author       string
	Dependencies []string
}

// GitStorage is an implementation of StorageReader based on
// a git repository as storage back-end.
type GitStorage struct {
	repo *git.Repository
}

// NewGitStorage returns a StorageReader which uses a public git repository
// as data source for pazuzu features.
//
// url:  The URL to the git repository that serves as data source. The
//       repository must be publicly accessible.
//
// If the repository can't be accessed NewStorageReader returns an error.
func NewGitStorage(url string) (*GitStorage, error) {
	// OPTIMIZATION: can be an fs repository which is cached and only pulled when needed
	repo := git.NewMemoryRepository()

	err := repo.Clone(&git.CloneOptions{
		URL:           url,
		ReferenceName: core.HEAD,
		SingleBranch:  true,
		Depth:         1,
	})
	if err != nil {
		return nil, err
	}

	return &GitStorage{repo: repo}, nil
}

func (storage *GitStorage) SearchMeta(name *regexp.Regexp) ([]FeatureMeta, error) {
	commit, err := storage.latestCommit()
	if err != nil {
		return nil, err
	}

	all, err := commit.Files()
	if err != nil {
		return nil, err
	}

	// find matching feature names
	matchedNames := map[string]bool{}
	matchedFeatures := []FeatureMeta{}
	err = all.ForEach(func(file *git.File) error {
		pathComponents := strings.Split(file.Name, "/")

		// check if file is in feature dir
		if pathComponents[0] != featureDir {
			return nil
		}

		// check if feature was already found
		featureName := pathComponents[1]
		if matchedNames[featureName] {
			return nil
		}

		// check if feature matches search params
		if name.MatchString(featureName) {
			meta, err := getMeta(commit, featureName)
			if err != nil {
				return err
			}
			matchedFeatures = append(matchedFeatures, meta)
			matchedNames[featureName] = true
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return matchedFeatures, nil
}

func (storage *GitStorage) GetMeta(name string) (FeatureMeta, error) {
	commit, err := storage.latestCommit()
	if err != nil {
		return FeatureMeta{}, err
	}

	return getMeta(commit, name)
}

// getMeta returns metadata about a feature from a given commit.
//
// commit:  The commit from which to obtain the feature information.
// name:    The exact feature name.
func getMeta(commit *git.Commit, name string) (FeatureMeta, error) {
	file, err := commit.File(path.Join(featureDir, name, featureFile))

	if err != nil {
		return FeatureMeta{}, err
	}

	reader, err := file.Reader()
	if err != nil {
		return FeatureMeta{}, err
	}

	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return FeatureMeta{}, err
	}

	meta := &yamlFeatureMeta{}
	err = yaml.Unmarshal(content, meta)
	if err != nil {
		return FeatureMeta{}, err
	}

	return FeatureMeta{
		Name:         name,
		Dependencies: meta.Dependencies,
		Description:  meta.Description,
		Author:       meta.Author,
		UpdatedAt:    commit.Committer.When,
	}, nil
}

func (storage *GitStorage) GetFeature(name string) (Feature, error) {
	commit, err := storage.latestCommit()
	if err != nil {
		return Feature{}, err
	}

	return getFeature(commit, name)
}

// getFeature returns all data of a feature from a given commit.
//
// commit:  The commit from which to obtain the feature information.
// name:    The exact feature name.
func getFeature(commit *git.Commit, name string) (Feature, error) {
	meta, err := getMeta(commit, name)
	if err != nil {
		return Feature{}, err
	}

	file, err := commit.File(path.Join(featureDir, name, featureSnippet))
	if err != nil {
		if err == git.ErrFileNotFound {
			return Feature{Meta: meta}, nil
		}
		return Feature{}, err
	}

	reader, err := file.Reader()
	if err != nil {
		return Feature{}, err
	}

	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return Feature{}, err
	}

	return Feature{
		Meta:    meta,
		Snippet: string(content),
	}, nil
}

func (storage *GitStorage) Resolve(names ...string) (map[string]Feature, error) {
	commit, err := storage.latestCommit()
	if err != nil {
		return map[string]Feature{}, err
	}

	result := map[string]Feature{}
	for _, name := range names {
		err = resolve(commit, name, result)
		if err != nil {
			return map[string]Feature{}, err
		}
	}

	return result, nil
}

// resolve returns all data for a certain feature and its direct and indirect
// dependencies. All feature data is added to the provided result map.
//
// commit:  The commit from which to obtain the feature information.
// name:    The exact feature name.
// result:  All features collected so far.
func resolve(commit *git.Commit, name string, result map[string]Feature) error {
	if _, ok := result[name]; ok {
		return nil
	}

	feature, err := getFeature(commit, name)
	if err != nil {
		return err
	}

	for _, depName := range feature.Meta.Dependencies {
		err = resolve(commit, depName, result)
		if err != nil {
			return err
		}
	}

	result[name] = feature

	return nil
}

// latestCommit is a helper method which gets the latest commit (HEAD) from
// a the storage git repository.
func (storage *GitStorage) latestCommit() (*git.Commit, error) {
	head, err := storage.repo.Head()
	if err != nil {
		return nil, err
	}

	commit, err := storage.repo.Commit(head.Hash())
	if err != nil {
		return nil, err
	}

	return commit, nil
}
