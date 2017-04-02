package mock

import (
	"regexp"

	"github.com/zalando-incubator/pazuzu/shared"
)

type TestStorage struct{}

var pythonFeatureMeta = shared.FeatureMeta{
	Name:        "python",
	Description: "Use `python -V`",
}

func GetTestFeatureMeta() shared.FeatureMeta {
	return pythonFeatureMeta
}

func (s *TestStorage) GetFeature(feature string) (shared.Feature, error) {
	return shared.Feature{
		Meta:    pythonFeatureMeta,
		Snippet: "RUN apt-get update && apt-get install python --yes",
	}, nil
}

func (s *TestStorage) GetMeta(name string) (shared.FeatureMeta, error) {
	return pythonFeatureMeta, nil
}

func (s *TestStorage) SearchMeta(name *regexp.Regexp) ([]shared.FeatureMeta, error) {
	return []shared.FeatureMeta{pythonFeatureMeta}, nil
}

func (s *TestStorage) Resolve(names ...string) ([]string, map[string]shared.Feature, error) {
	return []string{}, make(map[string]shared.Feature), nil
}
