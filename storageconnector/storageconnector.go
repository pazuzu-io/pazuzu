package storageconnector

import (
	"regexp"

	"github.com/zalando-incubator/pazuzu/shared"
)

// StorageReader defines an interface to get Features from data sources
type StorageReader interface {
	// SearchMeta returns an arbitrary ordered list of FeatureMeta records using given expression
	SearchMeta(name *regexp.Regexp) ([]shared.FeatureMeta, error)

	// GetMeta returns a single FeatureMeta by given Name. Meta is a small piece of data,
	// so it should be indexed by a storage and accessed rather quickly.
	GetMeta(name string) (shared.FeatureMeta, error)

	// Get returns a full feature data from a storage. This operation is a way slower than GetMeta, so for
	// quick lookups GetMeta is better to be used.
	GetFeature(name string) (shared.Feature, error)

	// Resolve finds all dependencies for a given list of Feature names and returns them as a map of
	// Features. The returned map will contain the Feature information for all listed names as well as
	// the Feature information of all their direct or indirect dependencies.
	//
	// names:  The names of the features which dependencies should be resolved.
	//
	// If a feature can't be found or a dependency can't be resolved an error is returned.
	Resolve(names ...string) ([]string, map[string]shared.Feature, error)
}
