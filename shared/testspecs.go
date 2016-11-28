package shared

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

const shebang = "#!/usr/bin/env bats"


func ReadTestSpec(reader io.Reader) string {
	var scanner = bufio.NewScanner(reader)
	var buffer = bytes.NewBufferString("")
	var line string

	for scanner.Scan() {
		line = strings.TrimRight(scanner.Text(), " \t\r\n")
		if line != shebang {
			buffer.WriteString(line + "\n")
		}
	}

	return strings.TrimSpace(buffer.String())
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
