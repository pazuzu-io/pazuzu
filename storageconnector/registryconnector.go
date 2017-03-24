package storageconnector

import (
	"regexp"
	"strconv"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/zalando-incubator/pazuzu/shared"
	"github.com/zalando-incubator/pazuzu/swagger/client/features"

	httptransport "github.com/go-openapi/runtime/client"
)

const (
	DefaultApiPath = "/api"
)


type registryStorage struct {
	Hostname string // localhost
	Port     int    // 8080
	Scheme   string // http
	Token    string // OAUTH2 Token

	Features  *features.Client
	Transport runtime.ClientTransport
}

func (store *registryStorage) init(hostname string, port int, scheme string, formats strfmt.Registry) {
	store.Hostname = hostname
	store.Port = port
	store.Scheme = scheme

	host := hostname + ":" + strconv.Itoa(port)
	path := DefaultApiPath
	schemes := []string{scheme}

	transport := httptransport.New(host, path, schemes)

	store.Transport = transport
	store.Features = features.New(transport, formats)
}

func NewRegistryStorage(hostname string, port int, scheme string, formats strfmt.Registry) (*registryStorage, error) {
	if formats == nil {
		formats = strfmt.Default
	}

	var rs registryStorage
	rs.init(hostname, port, scheme, formats)
	return &rs, nil
}

// Return a full feature data from the storage.
// For the registry, the filtering is done server-side to reduce result size.
// name:	a value, that must present in feature name (from API doc)
func (store *registryStorage) GetFeature(name string) (shared.Feature, error) {

	// let's get features containing name on the registry
	// params := features.NewGetAPIFeaturesParams()
	// params.Names = []string{name}
	// apiFeatures, err := store.Features.GetAPIFeatures(params)

	params := features.NewGetFeaturesNameParams().WithName(name)
	feature, err := store.Features.GetFeaturesName(params)
	if err != nil {
		return shared.Feature{}, err
	}
	return shared.NewFeature(feature.Payload), err
}

// Use the given regex to return a list of FeatureMeta.
// name		a regex used to filter out FeatureMeta
// TODO issue registry-#111 -> investigate regex support server-side to optimize
func (store *registryStorage) SearchMeta(name *regexp.Regexp) ([]shared.FeatureMeta, error) {

	result := []shared.FeatureMeta{}
	nameStr := name.String()
	params := features.NewGetFeaturesParams().WithQ(&nameStr)

	features, err := store.Features.GetFeatures(params)
	if err == nil {
		for _, feature := range features.Payload.Features {
			result = append(result, shared.NewMeta(feature.Meta))
		}
	}

	return result, err
}

// Return a feature metadata from the storage.
// name:	a value, that must present in feature name
func (store *registryStorage) GetMeta(name string) (shared.FeatureMeta, error) {

	feature, err := store.GetFeature(name)

	if err != nil {
		return shared.FeatureMeta{}, err
	}
	return feature.Meta, nil
}

// Resolve a list of features and their dependencies from the storage. Return non-nil err if at least one feature not found.
// names:	an array of feature names
func (store *registryStorage) Resolve(names ...string) ([]string, map[string]shared.Feature, error) {

	params := features.NewGetDependenciesParams()
	params.Names = names

	features, err := store.Features.GetDependencies(params)
	if err != nil {
		return []string{}, map[string]shared.Feature{}, err
	}

	var slice []string
	result := map[string]shared.Feature{}

	for _, feature := range features.Payload.Dependencies {
		feature2 := shared.NewFeature(feature)
		name := feature2.Meta.Name
		result[name] = feature2
		slice = append(slice, name)
	}
	return slice, result, nil
}
