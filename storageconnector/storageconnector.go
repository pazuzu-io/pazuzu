package storageconnector

import (
	"regexp"
	"time"
)

// StorageReader defines an interface to get Features from data sources
type StorageReader interface {
	// SearchMeta returns an arbitrary ordered list of FeatureMeta records using given expression
	SearchMeta(name *regexp.Regexp) ([]FeatureMeta, error)

	// GetMeta returns a single FeatureMeta by given Name. Meta is a small piece of data,
	// so it should be indexed by a storage and accessed rather quickly.
	GetMeta(name string) (FeatureMeta, error)

	// Get returns a full feature data from a storage. This operation is a way slower than GetMeta, so for
	// quick lookups GetMeta is better to be used.
	GetFeature(name string) (Feature, error)

	// Resolve finds all dependencies for a given list of Feature names and returns them as a map of
	// Features. The returned map will contain the Feature information for all listed names as well as
	// the Feature information of all their direct or indirect dependencies.
	//
	// names:  The names of the features which dependencies should be resolved.
	//
	// If a feature can't be found or a dependency can't be resolved an error is returned.
	Resolve(names ...string) ([]string, map[string]Feature, error)
}

type StorageWriter interface {
	// storeFeature is used to store a new feature
	storeFeature(f *Feature)

	// modify FeatureMeta of a feature that currently exists
	modifyFeatureMeta (name string, meta *FeatureMeta)

	// modify Snippet of a already existing feature
	modifyFeatureSnippet (name string, snippet string)

	// delete feature
	deleteFeature (name string)

}


// FeatureMeta provides short information about the Feature.
// This piece of data better to be indexed by a storage.
type FeatureMeta struct {
	Name         string
	Description  string
	Author       string
	UpdatedAt    time.Time
	Dependencies []string
}

// Feature is a definition for a piece of work to be done. Contains meta information as well as
// all necessary data to compose a piece of Dockerfile at the end.
type Feature struct {
	Meta    FeatureMeta
	Snippet string
}
