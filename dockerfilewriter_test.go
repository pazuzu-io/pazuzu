package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestDockerfileWriter(t *testing.T) {
	writer := NewDockerfileWriter()

	writer.AppendRaw(`
  FROM debian:jessie
  RUN apt-get update
  `)

	features := []struct {
		name  string
		data  string
		files map[string]string
	}{
		{
			name:  "feature-1",
			data:  "COPY a/b/c /\nCOPY c /home",
			files: map[string]string{"a/b/c": "/", "c": "/"},
		},
		{
			name:  "feature-2",
			data:  "RUN /bin/bash\nCOPY lib.jar /var/lib",
			files: map[string]string{"lib.jar": "/var/lib"},
		},
	}

	for _, f := range features {
		err := writer.AppendFeature(Feature{Name: f.name, DockerData: f.data})
		if err != nil {
			t.Fatalf("Feature should be appended: %s", err)
		}
	}

	dockerfile := string(writer.Bytes())

	for _, f := range features {
		for srcF, dstF := range f.files {
			former := fmt.Sprintf("COPY %s %s", srcF, dstF)
			fixed := fmt.Sprintf("COPY %s/%s %s", f.name, srcF, dstF)

			if strings.Contains(dockerfile, former) || !strings.Contains(dockerfile, fixed) {
				t.Fatalf("All 'COPY' commands should be fixed to include feature-name as prefix")
			}
		}
	}
}
