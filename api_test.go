package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	featureResp = `
[{
  "name": "python",
  "docker_data": "RUN apt-get update && apt-get install python --yes",
  "test_instruction": "python -V",
  "dependencies": []
}]`
	featureRespError = `
{
  "code": "feature_not_found",
  "message": "Feature was not found",
  "detailed_message": null
}`
)

// Test getting features response from API.
func TestGetFeatures(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, featureResp)
		}),
	)
	defer ts.Close()

	registry := HttpRegistry{URL: ts.URL}
	fs, err := registry.getFeatures(ts.URL)
	if err != nil {
		t.Errorf("should not fail: %s", err)
	}

	if len(fs) != 1 {
		t.Errorf("expected 1 feature, got %d", 1, len(fs))
	}
}

// Test getting error response from API.
func TestGetFeaturesError(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, featureRespError, http.StatusNotFound)
		}),
	)
	defer ts.Close()

	registry := HttpRegistry{URL: ts.URL}
	_, err := registry.getFeatures(ts.URL)
	if err == nil {
		t.Errorf("should fail")
	}

	msg := "Feature was not found"

	if err.Error() != msg {
		t.Errorf("expected %s, got %s", msg, err.Error())
	}
}
