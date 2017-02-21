package utils

import (
	"github.com/zalando-incubator/pazuzu"
	"reflect"
	"testing"
)

const checkMark = "\u2713"

func TestAppendIfMissing(t *testing.T) {
	t.Run("Append to an empty slice", func(t *testing.T) {
		var emptySlice []string
		var element = "Test"
		var result = appendIfMissing(emptySlice, element)

		if len(result) != 1 || result[0] != element {
			t.Errorf("Wrong result: %s", result)
		}
	})

	t.Run("Append to a non-empty slice", func(t *testing.T) {
		var nonEmptySlice = []string{"Existing element"}
		var element = "Test"
		var result = appendIfMissing(nonEmptySlice, element)

		if len(result) != 2 || result[1] != element {
			t.Errorf("Wrong result: %s", result)
		}
	})

	t.Run("Does not append duplicates", func(t *testing.T) {
		var nonEmptySlice = []string{"Test"}
		var element = "Test"
		var result = appendIfMissing(nonEmptySlice, element)

		if len(result) != 1 {
			t.Errorf("Wrong result: %s", result)
		}
	})
}

func TestGenerateFeaturesList(t *testing.T) {

	t.Run("Fails when both add and init are specified", func(t *testing.T) {
		var pazuzufileFeatures []string
		var featuresToInit = []string{"a", "b"}
		var featuresToAdd = []string{"c", "d"}

		_, err := GenerateFeaturesList(pazuzufileFeatures, featuresToInit, featuresToAdd)
		if err != pazuzu.ErrInitAndAddAreSpecified {
			t.Error("No error is raised")
		}
	})

	t.Run("Returns features to init if specified", func(t *testing.T) {
		var pazuzufileFeatures []string
		var featuresToInit = []string{"a", "b"}
		var featuresToAdd []string

		result, err := GenerateFeaturesList(pazuzufileFeatures, featuresToInit, featuresToAdd)
		if !reflect.DeepEqual(result, featuresToInit) || err != nil {
			t.Errorf("Result differs from expected: %s", result)
		}
	})

	t.Run("Returns features to add if specified", func(t *testing.T) {
		var pazuzufileFeatures []string
		var featuresToInit []string
		var featuresToAdd = []string{"c", "d"}

		result, err := GenerateFeaturesList(pazuzufileFeatures, featuresToInit, featuresToAdd)
		if !reflect.DeepEqual(result, featuresToAdd) || err != nil {
			t.Errorf("Result differs from expected: %s", result)
		}
	})

	t.Run("Adds features to Pazuzufile features", func(t *testing.T) {
		var pazuzufileFeatures = []string{"a", "b"}
		var featuresToInit []string
		var featuresToAdd = []string{"c", "d"}
		var expectedFeatures = append(pazuzufileFeatures, featuresToAdd...)

		result, err := GenerateFeaturesList(pazuzufileFeatures, featuresToInit, featuresToAdd)
		if !reflect.DeepEqual(result, expectedFeatures) || err != nil {
			t.Errorf("Result differs from expected: %s", result)
		}
	})

	t.Run("Does not append duplicates", func(t *testing.T) {
		var pazuzufileFeatures = []string{"a", "b"}
		var featuresToInit []string
		var featuresToAdd = []string{"a", "c", "d"}
		var expectedFeatures = append(pazuzufileFeatures, "c", "d")

		result, err := GenerateFeaturesList(pazuzufileFeatures, featuresToInit, featuresToAdd)
		if !reflect.DeepEqual(result, expectedFeatures) || err != nil {
			t.Errorf("Result differs from expected: %s", result)
		}
	})
}
