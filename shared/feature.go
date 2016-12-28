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
	var f Feature
	f.Meta = NewMeta(feature.Meta)
	f.Snippet = feature.Snippet
	f.TestSnippet = feature.TestSnippet
	return f
}

func NewFeature_str(name string, desc string, auth string, updated time.Time, dependencies []string, snippet string, testSnippet string) Feature{
	m := NewMeta_str(name, desc, auth, updated, dependencies)
	return Feature{Meta:m, Snippet:snippet, TestSnippet:testSnippet}
}

func NewMeta(meta *models.FeatureMeta) FeatureMeta{
	var m FeatureMeta
	m.Name = meta.Name
	m.Description = meta.Description
	m.Author = meta.Author
	m.UpdatedAt,_ = time.Parse(meta.UpdatedAt, "2006-01-02T15:04:05-0700")
	m.Dependencies = meta.Dependencies

	return m
}

func NewMeta_str(name string, desc string, auth string, updated time.Time, dependencies []string) FeatureMeta{
	var m FeatureMeta
	m.Name = name
	m.Description = desc
	m.Author = auth
	m.UpdatedAt = updated
	m.Dependencies = dependencies

	return m
}
