package shared

import (
	"fmt"
	"io"
	"io/ioutil"
)

const shebang = "#!/usr/bin/env bats"


func ReadTestSpec(reader io.Reader) (string, error) {
	// TODO: Skip shebang line
	snippet, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(snippet), nil
}

func WriteTestSpec(writer io.Writer, features []Feature) error {
	var lines = []string{shebang}

	for _, feature := range features {
		lines = append(lines, feature.TestSnippet)
	}

	for _, line := range lines {
		_, err := fmt.Fprintln(writer, line)
		if err != nil {
			return err
		}
	}

	return nil
}
