package storageconnector

import (
	"github.com/zalando-incubator/pazuzu/shared"
	"regexp"
	"strconv"

	"swaggen/client/features"
	"swaggen/client/files"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	httptransport "github.com/go-openapi/runtime/client"
)

type registryStorage struct {
	Hostname string  // localhost
	Port int	 // 8080
	Token string     // OAUTH2 Token

	Features 	*features.Client
	Files		*files.Client
	Transport	runtime.ClientTransport
}

func (store *registryStorage) init(hostname string, port int, formats strfmt.Registry) {
	store.Hostname = hostname
	store.Port = port

	host := hostname +":"+strconv.Itoa(port)
	path := "/"
	schemes := []string{"http"}

	transport := httptransport.New(host, path, schemes)

	store.Transport = transport
	store.Features = features.New(transport, formats)
	store.Files = files.New(transport, formats)
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
	params.Name = name
	apiFeatures,err := store.Features.GetAPIFeatures(params)

	if (err != nil) {
		return nil, err
	}

	// let's check that the name and feature name actually match
	for _,ft := range apiFeatures.Payload {
		if ft.Name == name {
			return shared.NewFeature(ft), err
		}
	}

	return nil, err
}

// Use the given regex to return a list of FeatureMeta.
// name		a regex used to filter out FeatureMeta
// TODO : issue #137 -> implement Features Meta on the registry for speed gain
// TODO : issue registry-#111 -> investigate possibility of regex support server-side to optimize
func (store *registryStorage) SearchMeta(name *regexp.Regexp) ([]shared.FeatureMeta, error) {

	result := []shared.FeatureMeta{}

	apiFeatures,err := store.Features.GetAPIFeatures(nil)
	if (err != nil) {
		return result, err
	}

	for _,ft := range apiFeatures.Payload {
		if name.MatchString(ft.Name) {
			ft2 := shared.NewFeature(ft)
			result = append(result, ft2.Meta)
		}
	}

	return result, err
}

// Return a feature metadata from the storage.
// name:	a value, that must present in feature name
// TODO : issue #137 -> implement Features Meta on the registry for speed gain
func (store *registryStorage) GetMeta(name string) (shared.FeatureMeta, error) {
	ft, err := store.GetFeature(name)
	return ft.Meta, err
}

// Resolve a list of features and their dependencies from the storage. Return non-nil err if at least one feature not found.
// names:	an array of feature names
func (store *registryStorage) Resolve(names ...string) ([]string, map[string]shared.Feature, error) {
	var slice []string
	result := map[string]shared.Feature{}

	for _,name := range names {
		err := store.resolve(name, &slice, result)
		if err != nil {
			return []string{}, map[string]shared.Feature{}, err
		}
	}
	return slice, result, nil
}

// Resolve a single feature and its dependencies from the storage. Return non-nil err if feature not found.
// name:	a value that must present in feature name
// list:	a list of previously visited feature ids
// result:	a result map containing feature ids and data
func (store *registryStorage) resolve(name string, list *[]string, result map[string]shared.Feature) error {
	// if name already in result
	if _, ok := result[name]; ok {
		return nil
	}

	feature, err := store.GetFeature(name)
	if err != nil {
		return err
	}
	for _, depName := range feature.Meta.Dependencies {
		err := store.resolve(depName, list, result)
		if err != nil {
			return err
		}
	}

	result[name] = feature
	*list = append(*list, name)
	return nil
}