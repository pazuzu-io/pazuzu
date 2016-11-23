package main

import (
	"reflect"
	"testing"
)


func TestAppendIfMissingWithEmptySlice(t *testing.T) {
	var emptySlice []string
	var element = "Test"
	var result = appendIfMissing(emptySlice, element)

	if len(result) != 1 || result[0] != element {
		t.Errorf("Wrong result: %s", result)
	}
}

func TestAppendIfMissingWithNonEmptySlice(t *testing.T) {
	var nonEmptySlice = []string{"Existing element"}
	var element = "Test"
	var result = appendIfMissing(nonEmptySlice, element)

	if len(result) != 2 || result[1] != element {
		t.Errorf("Wrong result: %s", result)
	}
}

func TestAppendIfMissingDoesNotAppendDuplicate(t *testing.T) {
	var nonEmptySlice = []string{"Test"}
	var element = "Test"
	var result = appendIfMissing(nonEmptySlice, element)

	if len(result) != 1 {
		t.Errorf("Wrong result: %s", result)
	}
}

func TestGenerateFeaturesListFailsWhenBothAddAndInitAreSpecified(t *testing.T) {
	var pazuzufileFeatures []string
	var featuresToInit = []string{"a", "b"}
	var featuresToAdd = []string{"c", "d"}

	_, err := generateFeaturesList(pazuzufileFeatures, featuresToInit, featuresToAdd)
	if err != ErrInitAndAddAreSpecified {
		t.Error("No error is raised")
	}
}

func TestGenerateFeaturesListReturnsFeaturesToInitIfSpecified(t *testing.T) {
	var pazuzufileFeatures []string
	var featuresToInit = []string{"a", "b"}
	var featuresToAdd []string

	result, err := generateFeaturesList(pazuzufileFeatures, featuresToInit, featuresToAdd)
	if !reflect.DeepEqual(result, featuresToInit) || err != nil {
		t.Errorf("Result differs from expected: %s", result)
	}
}

func TestGenerateFeaturesListReturnsFeaturesToAddIfSpecified(t *testing.T) {
	var pazuzufileFeatures []string
	var featuresToInit []string
	var featuresToAdd = []string{"c", "d"}

	result, err := generateFeaturesList(pazuzufileFeatures, featuresToInit, featuresToAdd)
	if !reflect.DeepEqual(result, featuresToAdd) || err != nil {
		t.Errorf("Result differs from expected: %s", result)
	}
}

func TestGenerateFeaturesListAddsFeaturesToPazuzuFeatures(t *testing.T) {
	var pazuzufileFeatures = []string{"a", "b"}
	var featuresToInit []string
	var featuresToAdd = []string{"c", "d"}
	var expectedFeatures = append(pazuzufileFeatures, featuresToAdd...)

	result, err := generateFeaturesList(pazuzufileFeatures, featuresToInit, featuresToAdd)
	if !reflect.DeepEqual(result, expectedFeatures) || err != nil {
		t.Errorf("Result differs from expected: %s", result)
	}
}

func TestGenerateFeaturesListDoesntAppendDuplicates(t *testing.T) {
	var pazuzufileFeatures = []string{"a", "b"}
	var featuresToInit []string
	var featuresToAdd = []string{"a", "c", "d"}
	var expectedFeatures = append(pazuzufileFeatures, "c", "d")

	result, err := generateFeaturesList(pazuzufileFeatures, featuresToInit, featuresToAdd)
	if !reflect.DeepEqual(result, expectedFeatures) || err != nil {
		t.Errorf("Result differs from expected: %s", result)
	}
}

func TestGetFeaturesListReturnsEmptySliceWhenNothingSpecified(t *testing.T) {
	var badExamples = []string{"", "    ", ", "}
	for _, example := range badExamples {
		result := getFeaturesList(example)
		if len(result) != 0 {
			t.Errorf("Result should be empty: %s", result)
		}
	}
}

func TestGetFeaturesListReturnsListOfFeatures(t *testing.T) {
	var featureString = "node,   java,clojure  "
	var expectedFeatures = []string{"node", "java", "clojure"}
	var result = getFeaturesList(featureString)

	if !reflect.DeepEqual(result, expectedFeatures) {
		t.Errorf("Result differs from expected: %s", result)
	}
}
