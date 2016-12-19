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

func (store *registryStorage) GetFeature(name string) (shared.Feature, error) {
	// TODO : unimplemented method
	resp,err := store.Features.GetAPIFeatures(nil)

	if (err != nil) {
		return shared.NewEmptyFeature(), err
	}

	resp2 := shared.NewFeature(resp.Payload[0])

	return resp2, err
}

func (store *registryStorage) SearchMeta(name *regexp.Regexp) ([]shared.FeatureMeta, error) {
	// TODO : unimplemented method
	return []shared.FeatureMeta{}, nil
}

func (store *registryStorage) GetMeta(name string) (shared.FeatureMeta, error) {
	// TODO : unimplemented method
	return shared.FeatureMeta{}, nil
}

func (store *registryStorage) Resolve(names ...string) ([]string, map[string]shared.Feature, error) {
	// TODO : unimplemented method
	m := make(map[string]shared.Feature)
	return []string{}, m, nil
}