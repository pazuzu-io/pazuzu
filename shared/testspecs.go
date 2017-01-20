package shared

import (
	"fmt"
	"io"
)

const (
	shebang          = "#!/usr/bin/env bats"
	TestSpecFilename = "test.bats"
)

func WriteTestSpec(writer io.Writer, features []Feature) error {
	var lines = []string{shebang}

	for _, feature := range features {
		lines = append(lines, feature.TestSnippet)
	}

	for _, line := range lines {
		_, err := fmt.Fprintf(writer, "%s\n\n", line)
		if err != nil {
			return err
		}
	}

	return nil
}
