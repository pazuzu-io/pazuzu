package storageconnector

import (
	"strings"
	"path"
	"sort"
	"gopkg.in/yaml.v2"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/core"
	"io/ioutil"
)

const (
	featureDir = "features"
	featureFile = "meta.yml"
	featureSnippet = "Dockerfile"
	defaultSearchParamsLimit = 100
)

type yamlFeatureMeta struct {
	Description string
	Dependencies []string
}

type gitStorage struct {
	repo 	*git.Repository
}

func NewStorageReader(url string) (StorageReader, error) {
	// OPTIMIZATION: can be an fs repository which is cached and only pulled when needed
	repo := git.NewMemoryRepository()

	err  := repo.Clone(&git.CloneOptions{
		URL: url,
		ReferenceName: core.HEAD,
		SingleBranch: true,
		Depth: 1,
	})
	if err != nil {
		return nil, err
	}

	return &gitStorage{repo: repo}, nil
}

func (storage *gitStorage) SearchMeta(params SearchParams) ([]FeatureMeta, error) {
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
	matchedNamesList := []string{}
	all.ForEach(func(file *git.File) error {
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
		if params.Name.MatchString(featureName) {
			matchedNames[featureName] = true
			matchedNamesList = append(matchedNamesList, featureName)
		}

		return nil
	})

	if params.Limit == 0 {
		params.Limit = defaultSearchParamsLimit
	}

	// prepare resulting feature metadata list
	// OPTIMIZATION: if the above ForEach call was based on some kind of reliable ordering
	//               the following Sort call could be omitted.
	sort.Sort(sort.StringSlice(matchedNamesList))
	matchedFeatures := []FeatureMeta{}
	for i := params.Offset; i < params.Offset + params.Limit && i < int64(len(matchedNamesList)); i++ {
		meta, _ := getMeta(commit, matchedNamesList[i])
		matchedFeatures = append(matchedFeatures, meta)
	}

	return matchedFeatures, nil
}

func (storage *gitStorage) GetMeta(name string) (FeatureMeta, error) {
	commit, err := storage.latestCommit()
	if err != nil {
		return FeatureMeta{}, err
	}

	return getMeta(commit, name)
}

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
		Name: name,
		Dependencies: meta.Dependencies,
		UpdatedAt: commit.Committer.When,
		// FIX: CreatedAt is missing
		// FIX: Description from meta.yml is ignored
	}, nil
}

func (storage *gitStorage) Get(name string) (Feature, error) {
	commit, err := storage.latestCommit()
	if err != nil {
		return Feature{}, err
	}

	return getFeature(commit, name)
}

func getFeature(commit *git.Commit, name string) (Feature, error) {
	meta, err := getMeta(commit, name)
	if err != nil {
		return Feature{}, err
	}

	file, err := commit.File(path.Join(featureFile, name, featureSnippet))
	if err != nil {
		if err == git.ErrFileNotFound {
			return Feature{Meta: meta}, nil
		} else {
			return Feature{}, err
		}
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
		Meta: meta,
		Snippet: string(content),
	}, nil
}

func (storage *gitStorage) Resolve(name string) ([]Feature, error) {
	commit, err := storage.latestCommit()
	if err != nil {
		return []Feature{}, err
	}

	return resolve(commit, name, []Feature{})
}

func resolve(commit *git.Commit, name string, result []Feature) ([]Feature, error)  {
	if containsFeatureWithName(result, name) {
		return result, nil
	}

	feature, err := getFeature(commit, name)
	if err != nil {
		return []Feature{}, err
	}

	for _, depName := range feature.Meta.Dependencies {
		result, err = resolve(commit, depName, result)
		if err != nil {
			return []Feature{}, err
		}
	}

	result = append(result, feature)

	return result, nil
}

func containsFeatureWithName(list []Feature, name string) bool {
	for _, f := range list {
		if f.Meta.Name == name {
			return true
		}
	}
	return false
}

func (storage *gitStorage) latestCommit() (*git.Commit, error) {
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
