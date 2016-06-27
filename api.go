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
type HttpRegistry struct {
	URL           string
	Authenticator Authenticator
}

// APIError defines error response from pazuzu-registry.
type APIError struct {
	Code            string
	Message         string
	DetailedMessage string `json:"detailed_message"`
}

// GetFeatures gets features from the pazuzu-registry.
func (r HttpRegistry) GetFeatures(features []string) ([]Feature, error) {
	return r.getFeatures(fmt.Sprintf("%s/features?name=%s", r.URL,
		strings.Join(features, ",")))
}

// SearchFeatures searches for features based on name.
func (r HttpRegistry) SearchFeatures(query string) ([]Feature, error) {
	return r.getFeatures(fmt.Sprintf("%s/features/search/%s", r.URL, query))
}

// ListFeatures lists all features from registry.
func (r HttpRegistry) ListFeatures() ([]Feature, error) {
	return r.getFeatures(fmt.Sprintf("%s/features", r.URL))
}

// Makes HTTP request to pazuzu registry and decodes the json response.
func (r HttpRegistry) getFeatures(url string) ([]Feature, error) {

	auth, err := r.authenticate()
	if err != nil {
		return nil, err
	}
	resp, err := makeRequest(url, auth)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, decodeError(resp)
	}
	return decodeFeatures(resp)
}

func (r HttpRegistry) authenticate() (Authentication, error) {
	if r.Authenticator != nil {
		return r.Authenticator.Authenticate()
	}
	return nil, nil
}

func makeRequest(url string, auth Authentication) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if auth != nil {
		auth.Enrich(req)
	}
	return client.Do(req)
}

func decodeError(resp *http.Response) error {
	var errResp APIError
	dec := json.NewDecoder(resp.Body)
	err := dec.Decode(&errResp)
	if err != nil {
		return err
	}
	return fmt.Errorf(errResp.Message)
}

func decodeFeatures(resp *http.Response) ([]Feature, error) {
	var res []Feature
	dec := json.NewDecoder(resp.Body)
	err := dec.Decode(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
