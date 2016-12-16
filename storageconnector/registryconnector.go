package storageconnector

import (
	"github.com/zalando-incubator/pazuzu/shared"
	"regexp"
)

type registryStorage struct {
	Hostname string  // localhost
	Port int	 // 8080
	Endpoint string  // /api
	Token string     // OAUTH2 Token
}

func (store *registryStorage) init(hostname string, port int) {
	store.Hostname = hostname
	store.Port = port
	store.Endpoint = "/api"
}

func NewRegistryStorage(hostname string, port int) (*registryStorage, error) {
	var rs registryStorage
	rs.init(hostname, port)
	return &rs, nil
}

func (store *registryStorage) GetFeature(name string) (shared.Feature, error) {
	return shared.Feature{}, nil
}

func (store *registryStorage) SearchMeta(name *regexp.Regexp) ([]shared.FeatureMeta, error) {
	return []shared.FeatureMeta{}, nil
}

func (store *registryStorage) GetMeta(name string) (shared.FeatureMeta, error) {
	return shared.FeatureMeta{}, nil
}

func (store *registryStorage) Resolve(names ...string) ([]string, map[string]shared.Feature, error) {
	m := make(map[string]shared.Feature)
	return []string{}, m, nil
}