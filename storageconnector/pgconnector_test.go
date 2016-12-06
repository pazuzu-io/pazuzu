package storageconnector

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zalando-incubator/pazuzu/shared"
)

var beginningOf2016 = time.Date(2016, 1, 1, 11, 11, 11, 0, time.UTC)
var randomError = errors.New("Random error")

func fakeScanner(dest ...interface{}) error {
	*dest[0].(*int) = 123
	*dest[1].(*string) = "Java"
	*dest[2].(*string) = "Java8 SDK"
	*dest[3].(*string) = "John Smith"
	*dest[4].(*time.Time) = beginningOf2016
	*dest[5].(*string) = "java node"
	*dest[6].(*string) = "RUN apt-get java"
	*dest[7].(*string) = "#!/usr/bin/env bats\n\n@test \"Check that Java is installed\" {command java -version}"
	return nil
}

func errorScanner(dest ...interface{}) error {
	return randomError
}

func TestReadFeature(t *testing.T) {

	t.Run("Test error", func(t *testing.T) {
		result, err := readFeature(errorScanner)

		assert.Equal(t, err, randomError)
		assert.Equal(t, result, shared.Feature{})
	})

	t.Run("Test success", func(t *testing.T) {
		result, err := readFeature(fakeScanner)

		assert.Equal(t, result.Meta.Name, "Java")
		assert.Equal(t, result.Meta.Description, "Java8 SDK")
		assert.Equal(t, result.Meta.Author, "John Smith")
		assert.Equal(t, result.Meta.UpdatedAt, beginningOf2016)
		assert.Equal(t, result.Meta.Dependencies, []string{"java", "node"})
		assert.Equal(t, result.Snippet, "RUN apt-get java")
		assert.Equal(
			t,
			result.TestSnippet,
			"@test \"Check that Java is installed\" {command java -version}")

		assert.Equal(t, err, nil)
	})
}
