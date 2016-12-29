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

	_ = registry2.addFeature(featureA)
	_ = registry2.addFeature(featureB)
	_ = registry2.addFeature(featureC)
	_ = registry2.addFeature(featureD)
	_ = registry2.addFeature(featureE)
	_ = registry2.addFeature(featureF)

	_,err = registry2.GetMeta("F")

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

// TODO issue #159 -> method does not test regex contrary to specs
func TestRegistry_SearchMeta(t *testing.T) {
	if !setupOk {
		t.Skipf("No endpoint listening at %v:%v", hostname, port)
	}

	searchMetaAndFindResultTest(t, "D", []shared.FeatureMeta{featureA.Meta, featureB.Meta, featureD.Meta}, storage)
	searchMetaAndFindNothingTest(t, "NotAFeature", storage)
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

	resolveFeaturesTest(t, "resolve features with intersecting deps", []string{"E", "F"},
		map[string]shared.Feature{"A":featureA, "B":featureB, "C":featureC, "D":featureD, "E":featureE, "F":featureF}, registry)

}