package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	successResp = `{
  "token_type": "Bearer",
  "access_token": "ad479621-455a-4d8a-9ac8-e1a7e7ce8f01",
  "expires_in": 3600,
  "refresh_token": "0b044aa8-f075-4d10-9629-f067f0928c82",
  "scope": "uid cn"
}`
)

func TestGetToken(t *testing.T) {
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, successResp)
		}),
	)
	defer ts.Close()

	auth := NewOAuth2Authenticator(ts.URL, "jdoe", "password")
	authentication, err := auth.Authenticate()

	if err != nil {
		t.Errorf("should not fail: %s", err)
	}
	if authentication == nil {
		t.Errorf("Failed to acquire token")
	}
}

func TestFailToGetToken(t *testing.T) {
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Unathorized", http.StatusUnauthorized)
		}),
	)
	defer ts.Close()

	auth := NewOAuth2Authenticator(ts.URL, "jdoe", "password")
	_, err := auth.Authenticate()

	if err == nil {
		t.Errorf("should have failed")
	}
}
