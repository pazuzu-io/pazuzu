package storageconnector

import (
	"fmt"
	"regexp"
	"sort"
)

// Memory is a simple in-memory storage of features
// usable for writing tests
type Memory struct {
	features     map[string]Feature
	featureNames []string // sorted list of feature names
}

// NewMemory is a constructor for in-memory storage
func NewMemory(features []Feature) *Memory {
	m := &Memory{
		features: map[string]Feature{},
	}
	for _, f := range features {
		m.features[f.Meta.Name] = f
		m.featureNames = append(m.featureNames, f.Meta.Name)
	}

	sort.Strings(m.featureNames)
	return m
}

func (m *Memory) SearchMeta(name *regexp.Regexp) ([]FeatureMeta, error) {
	result := []FeatureMeta{}
	for _, n := range m.featureNames {
		if name.MatchString(n) {
			f := m.features[n]
			result = append(result, f.Meta)
		}
	}
	return result, nil
}

func (m *Memory) GetMeta(name string) (FeatureMeta, error) {
	f, err := m.GetFeature(name)
	return f.Meta, err
}

func (m *Memory) GetFeature(name string) (Feature, error) {
	// TODO: make Get case-insensitive
	f, ok := m.features[name]
	if !ok {
		return Feature{}, fmt.Errorf("Feature '%s' was not found", name)
	}

	return f, nil
}

func (m *Memory) Resolve(names ...string) (map[string]Feature, error) {
	result := map[string]Feature{}
	for _, name := range names {
		if err := m.resolve(name, result); err != nil {
			return result, err
		}
	}

	return result, nil
}

func (m *Memory) resolve(name string, resolved map[string]Feature) error {
	f, ok := m.features[name]
	if !ok {
		return fmt.Errorf("Feature '%s' was not found", name)
	}

	if _, ok := resolved[name]; ok {
		return nil
	}

	resolved[f.Meta.Name] = f
	for _, depName := range f.Meta.Dependencies {
		if err := m.resolve(depName, resolved); err != nil {
			return err
		}
	}

	return nil
}
