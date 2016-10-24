package storageconnector

import (
	"regexp"
	"time"
)

// StorageReader defines an interface to get Features from data sources
type StorageReader interface {
	// SearchMeta returns an arbitrary ordered list of FeatureMeta records using given SearchParams
	SearchMeta(SearchParams) ([]FeatureMeta, error)

	// GetMeta returns a single FeatureMeta by given Feature ID. Meta is a small piece of data,
	// so it should be indexed by a storage and accessed rather quickly.
	GetMeta(featureID string) (FeatureMeta, error)

	// Get returns a full feature data from a storage. This operation is a way slower than GetMeta, so for
	// quick lookups GetMeta is better to be used.
	Get(featureID string) (Feature, error)

	// Resolve resolves dependencies for a given Feature and returns an **ordered** list of Features.
	// Ordering is critical here, it defines a way how Features should be executed.
	//
	// Example:
	// If you're given a list [FeatureA, FeatureB, FeatureC, FeatureD] that might mean that FeatureC depends on
	// FeatureB which in its turn depends on FeatureA. In given example FeatureD has no dependencies, but it must
	// be executed **after** all dependencies for FeatureC are resolved.
	Resolve(featureID string) ([]Feature, error)
}

// SearchParams define parameters for searching for the Features
type SearchParams struct {
	Name   regexp.Regexp
	Limit  int64
	Offset int64
}

// FeatureMeta provides short information about the Feature.
// This piece of data better to be indexed by a storage
type FeatureMeta struct {
	ID            string
	Name          string
	Author        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DependencyIDs []string
}

// Feature is a definition for a piece of work to be done. Contains meta information as well as
// all necessary data to compose a piece of Dockerfile at the end.
type Feature struct {
	Meta    FeatureMeta
	Snippet string
}
