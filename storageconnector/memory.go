package storageconnector

import (
	"fmt"
	"regexp"
	"sort"

	"github.com/zalando-incubator/pazuzu/shared"
)

// MemoryStorage is a simple in-memory storage of features
// usable for writing tests
type MemoryStorage struct {
	features     map[string]shared.Feature
	featureNames []string // sorted list of feature names
}

// NewMemoryStorage is a constructor for in-memory storage
func NewMemoryStorage(features []shared.Feature) *MemoryStorage {
	m := &MemoryStorage{
		features: map[string]shared.Feature{},
	}
	for _, f := range features {
		m.features[f.Meta.Name] = f
		m.featureNames = append(m.featureNames, f.Meta.Name)
	}

	sort.Strings(m.featureNames)
	return m
}

func (m *MemoryStorage) SearchMeta(name *regexp.Regexp) ([]shared.FeatureMeta, error) {
	result := []shared.FeatureMeta{}
	for _, n := range m.featureNames {
		if name.MatchString(n) {
			f := m.features[n]
			result = append(result, f.Meta)
		}
	}
	return result, nil
}

func (m *MemoryStorage) GetMeta(name string) (shared.FeatureMeta, error) {
	f, err := m.GetFeature(name)
	return f.Meta, err
}

func (m *MemoryStorage) GetFeature(name string) (shared.Feature, error) {
	// TODO: make Get case-insensitive
	f, ok := m.features[name]
	if !ok {
		return shared.Feature{}, fmt.Errorf("Feature '%s' was not found", name)
	}

	return f, nil
}

func (m *MemoryStorage) Resolve(names ...string) ([]string, map[string]shared.Feature, error) {
	var slice []string
	result := map[string]shared.Feature{}
	for _, name := range names {
		if err := m.resolve(name, &slice, result); err != nil {
			return []string{}, result, err
		}
	}

	return []string{}, result, nil
}

func (m *MemoryStorage) resolve(name string, list *[]string, resolved map[string]shared.Feature) error {
	f, ok := m.features[name]
	if !ok {
		return fmt.Errorf("Feature '%s' was not found", name)
	}

	if _, ok := resolved[name]; ok {
		return nil
	}

	resolved[f.Meta.Name] = f
	*list = append(*list, f.Meta.Name)
	for _, depName := range f.Meta.Dependencies {
		if err := m.resolve(depName, list, resolved); err != nil {
			return err
		}
	}

	return nil
}
