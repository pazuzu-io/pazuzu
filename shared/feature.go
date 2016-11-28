package shared

import (
	"time"
)

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
	Meta         FeatureMeta
	Snippet      string
	TestSnippet  string
}
