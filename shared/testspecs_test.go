package shared

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)


const testSpecFixture = `#!/usr/bin/env bats

@test "Check that Java is installed" {  
    command java -version    
}
`

func TestReadTestSpec(t *testing.T) {
	t.Run("Return an empty string where is nothing to read", func(t *testing.T) {
		reader := strings.NewReader("")
		result := ReadTestSpec(reader)

		assert.Equal(t, result, "", "Result should be empty")
	})

	t.Run("Return a proper test spec content without shebang and trailing whitespaces", func(t *testing.T) {
		reader := strings.NewReader(testSpecFixture)
		result := ReadTestSpec(reader)

		expectedTestSpec := "@test \"Check that Java is installed\" {\n" +
			"    command java -version\n" +
			"}"


		assert.Equal(t, result, expectedTestSpec, "Result mismatched")
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
		expectedOutput := "#!/usr/bin/env bats\n\n" +
			"@test \"Check that Java is installed\" {}\n\n" +
			"@test \"Check that Node.js is installed\" {}\n\n"

		WriteTestSpec(buffer, features)

		assert.Equal(t, buffer.String(), expectedOutput, "Result should be empty")
	})

	t.Run("Writes shebang when no features are specified", func(t *testing.T) {
		var buffer = bytes.NewBufferString("")
		var features = []Feature{}

		expectedOutput := "#!/usr/bin/env bats\n\n"
		WriteTestSpec(buffer, features)

		assert.Equal(t, buffer.String(), expectedOutput, "Result should be empty")
	})
}