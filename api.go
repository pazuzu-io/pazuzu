package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Feature defines a feature fetched from pazuzu-registry.
type Feature struct {
	Name            string
	Description     string
	DockerData      string `json:"docker_data"`
	TestInstruction string `json:"test_instruction"`
}

// PazuzuRegistry is an interface for pazuzu-registry.
type PazuzuRegistry interface {
	GetFeatures(features []string) ([]Feature, error)
}

// HttpRegistry is a wrapper for the Pazuzu registry API.
type HttpRegistry string

// APIError defines error response from pazuzu-registry.
type APIError struct {
	Code            string
	Message         string
	DetailedMessage string `json:"detailed_message"`
}

// GetFeatures gets features from the pazuzu-registry.
func (r HttpRegistry) GetFeatures(features []string) ([]Feature, error) {
	return getFeatures(fmt.Sprintf("%s/features?name=%s", r,
		strings.Join(features, ",")))
}

// SearchFeatures searches for features based on name.
func (r HttpRegistry) SearchFeatures(query string) ([]Feature, error) {
	return getFeatures(fmt.Sprintf("%s/features/search/%s", r, query))
}

// ListFeatures lists all features from registry.
func (r HttpRegistry) ListFeatures() ([]Feature, error) {
	return getFeatures(fmt.Sprintf("%s/features", r))
}

// Makes HTTP request to pazuzu registry and decodes the json response.
func getFeatures(url string) ([]Feature, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp APIError

		dec := json.NewDecoder(resp.Body)
		err = dec.Decode(&errResp)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf(errResp.Message)
	}

	var res []Feature

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
