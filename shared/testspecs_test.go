package shared

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)


const javaTestSpecFixture = `
	#!/usr/bin/env bats

	@test "Check that Java is installed" {
	    command java -version
	}
`

func TestReadTestSpec(t *testing.T) {
	t.Run("Return an empty string where is nothing to read", func(t *testing.T) {
		reader := strings.NewReader("")
		result, err := ReadTestSpec(reader)

		assert.Equal(t, result, "", "Result should be empty")
		assert.Equal(t, err, nil, "No error should be returned")
	})

	t.Run("Return a proper test spec content", func(t *testing.T) {
		reader := strings.NewReader(javaTestSpecFixture)
		result, err := ReadTestSpec(reader)

		assert.Equal(t, result, javaTestSpecFixture, "Result mismatched")
		assert.Equal(t, err, nil, "No error should be returned")
	})
}

func TestWriteTestSpec(t *testing.T) {
	t.Run("Writes all features test and shebang", func(t *testing.T) {
		var buffer = bytes.NewBufferString("")
		var features = []Feature{
			Feature{
				TestSnippet: "@test \"Check that Java is installed\" {}",
			},
			Feature{
				TestSnippet: "@test \"Check that Node.js is installed\" {}",
			},
		}
		expectedOutput := "#!/usr/bin/env bats\n" +
			"@test \"Check that Java is installed\" {}\n" +
			"@test \"Check that Node.js is installed\" {}\n"

		WriteTestSpec(buffer, features)

		assert.Equal(t, buffer.String(), expectedOutput, "Result should be empty")
	})

	t.Run("Writes shebang when no features are specified", func(t *testing.T) {
		var buffer = bytes.NewBufferString("")
		var features = []Feature{}

		expectedOutput := "#!/usr/bin/env bats\n"
		WriteTestSpec(buffer, features)

		assert.Equal(t, buffer.String(), expectedOutput, "Result should be empty")
	})
}