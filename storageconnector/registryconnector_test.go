package storageconnector

import (
	"testing"

	"github.com/zalando-incubator/pazuzu/shared"
	"fmt"
)

var (
	setupOk = false
	hostname = "localhost"
	port = 8080
	registry StorageReader

)

func InitRegistryTests() error {
	registry2, err := NewRegistryStorage(hostname, port, nil)

	err = registry2.addFeature(featureA)
	err = registry2.addFeature(featureB)
	err = registry2.addFeature(featureC)
	err = registry2.addFeature(featureD)
	err = registry2.addFeature(featureE)
	err = registry2.addFeature(featureF)

	if err == nil {
		setupOk = true
		registry = registry2
	} else {
		fmt.Printf("Setup for Registry Connector failed: %v\n", err)
		return err
	}
	return nil
}

func TestRegistry_GetMeta(t *testing.T) {
	if !setupOk {
		t.Skipf("No endpoint listening at %v:%v", hostname, port)
	}

	getExistingFeatureMetaTest(t, "A", registry)
	getNonExistingFeatureMetaTest(t, "NotAFeature", registry)
}

func TestRegistry_GetFeature(t *testing.T) {
	if !setupOk {
		t.Skipf("No endpoint listening at %v:%v", hostname, port)
	}

	getExistingFeatureTest(t, "A", registry)
	getNonExistingFeatureTest(t, "NotAFeature", registry)
}

func TestRegistry_ResolveOne(t *testing.T){
	if !setupOk {
		t.Skipf("No endpoint listening at %v:%v", hostname, port)
	}

	resolveNonExistingFeatureTest(t, "NotAFeature", registry)
	resolveEmptyFeaturesTest(t, registry)
	resolveSingleFeatureWithoutDependenciesTest(t, "A", registry)
	resolveFeaturesTest(t, "resolve single feature with deps", []string{"D"},
		map[string]shared.Feature{"A":featureA, "B":featureB, "D":featureD}, registry)

}