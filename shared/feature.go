package shared

import (
	"time"
	"swaggen/models"
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
	Meta        FeatureMeta
	Snippet     string
	TestSnippet string
}

func NewFeature(feature *models.Feature) Feature{
	var m FeatureMeta
	m.Name = feature.Name
	m.Description = feature.Description

	var f Feature
	f.Meta = m
	f.Snippet = feature.DockerData
	f.TestSnippet = feature.TestInstruction

	return f
}
