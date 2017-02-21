package compose

import (
	"testing"
	"reflect"
)

const checkMark = "\u2713"

func TestGetFeaturesList(t *testing.T) {
	t.Run("Returns empty slice when nothing is specified", func(t *testing.T) {
		var badExamples = []string{"", "    ", ", "}
		for _, example := range badExamples {
			result := getFeaturesList(example)
			if len(result) != 0 {
				t.Errorf("Result should be empty: %s", result)
			}
			t.Logf("%s Example \"%s\"", checkMark, example)
		}
	})

	t.Run("Returns list of features", func(t *testing.T) {
		var featureString = "node,   java,clojure  "
		var expectedFeatures = []string{"node", "java", "clojure"}
		var result = getFeaturesList(featureString)

		if !reflect.DeepEqual(result, expectedFeatures) {
			t.Errorf("Result differs from expected: %s", result)
		}
	})
}
