package storageconnector

import (
	"regexp"
	"strconv"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/zalando-incubator/pazuzu/shared"
	"swaggen/client/feature_metas"
	"swaggen/client/features"

	"context"
	httptransport "github.com/go-openapi/runtime/client"
	"swaggen/models"
)

type registryStorage struct {
	Hostname string // localhost
	Port     int    // 8080
	Token    string // OAUTH2 Token

	Features  *features.Client
	Metas     *feature_metas.Client
	Transport runtime.ClientTransport
}

func (store *registryStorage) init(hostname string, port int, formats strfmt.Registry) {
	store.Hostname = hostname
	store.Port = port

	host := hostname + ":" + strconv.Itoa(port)
	path := "/"
	schemes := []string{"http"}

	transport := httptransport.New(host, path, schemes)

	store.Transport = transport
	store.Features = features.New(transport, formats)
	store.Metas = feature_metas.New(transport, formats)
}

func NewRegistryStorage(hostname string, port int, formats strfmt.Registry) (*registryStorage, error) {
	if formats == nil {
		formats = strfmt.Default
	}

	var rs registryStorage
	rs.init(hostname, port, formats)
	return &rs, nil
}

// Return a full feature data from the storage.
// For the registry, the filtering is done server-side to reduce result size.
// name:	a value, that must present in feature name (from API doc)
func (store *registryStorage) GetFeature(name string) (shared.Feature, error) {

	// let's get features containing name on the registry
	params := features.NewGetAPIFeaturesParams()
	params.Names = []string{name}
	apiFeatures, err := store.Features.GetAPIFeatures(params)

	if err != nil {
		return shared.Feature{}, err
	}

	// let's check that the name and feature name actually match
	for _, ft := range apiFeatures.Payload {
		if ft.Meta.Name == name {
			return shared.NewFeature(ft), err
		}
	}
	return shared.Feature{}, err
}

// Use the given regex to return a list of FeatureMeta.
// name		a regex used to filter out FeatureMeta
// TODO issue registry-#111 -> investigate regex support server-side to optimize
func (store *registryStorage) SearchMeta(name *regexp.Regexp) ([]shared.FeatureMeta, error) {

	result := []shared.FeatureMeta{}

	params := feature_metas.NewGetAPIFeatureMetasParams()
	params.Name = []string{name.String()}

	metas, err := store.Metas.GetAPIFeatureMetas(params)

	if err == nil {
		for _, meta := range metas.Payload {
			result = append(result, shared.NewMeta(meta))
		}
	}

	return result, err
}

// Return a feature metadata from the storage.
// name:	a value, that must present in feature name
func (store *registryStorage) GetMeta(name string) (shared.FeatureMeta, error) {

	params := feature_metas.NewGetAPIFeatureMetasNameParams()
	params.Name = name

	meta, err := store.Metas.GetAPIFeatureMetasName(params)

	if err != nil {
		return shared.FeatureMeta{}, err
	}
	return shared.NewMeta(meta.Payload), nil
}

// Resolve a list of features and their dependencies from the storage. Return non-nil err if at least one feature not found.
// names:	an array of feature names
func (store *registryStorage) Resolve(names ...string) ([]string, map[string]shared.Feature, error) {
	if len(names) == 0 {
		names = append(names, "")
	}

	params := features.NewGetAPIResolvedFeaturesParams()
	params.Names = names

	features, err := store.Features.GetAPIResolvedFeatures(params)
	if err != nil {
		return []string{}, map[string]shared.Feature{}, err
	}

	var slice []string
	result := map[string]shared.Feature{}

	for _, feature := range features.Payload {

		feature2 := shared.NewFeature(feature)
		name := feature2.Meta.Name
		result[name] = feature2
		slice = append(slice, name)
	}
	return slice, result, nil
}

func (store *registryStorage) AddFeature(feature shared.Feature) error {

	meta := models.FeatureMeta{
		Name:         feature.Meta.Name,
		Description:  feature.Meta.Description,
		Author:       feature.Meta.Author,
		Dependencies: feature.Meta.Dependencies,
		UpdatedAt:    feature.Meta.UpdatedAt.UTC().Format("2006-01-02T15:04:05-0700"),
	}

	newfeature := models.Feature{
		Meta:        &meta,
		Snippet:     feature.Snippet,
		TestSnippet: feature.TestSnippet,
	}

	params := features.PostAPIFeaturesParams{
		Feature: &newfeature,
		Context: context.Background(),
	}

	_, err := store.Features.PostAPIFeatures(&params)
	if err != nil {
		return err
	}
	return nil
}
