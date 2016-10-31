package pazuzu

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	https       = "https"
	accessToken = "access_token"
)

type tokenRequest struct {
	user     string
	password string
	scopes   []string
}

// Authentication abstracts the HTTP request authentication, allowing to forward the credentials with the request
type Authentication interface {
	Enrich(*http.Request)
}

// Authenticator abstract the authentication process
type Authenticator interface {
	Authenticate() (Authentication, error)
}

// BearerTokenAuthentication the OAuth2 Bearer token authentication
type bearerTokenAuthentication struct {
	token string
}

// Enrich allows to add the authentication details to the request
func (auth bearerTokenAuthentication) Enrich(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", auth.token))
}

// OAuth2Authenticator hides all of the details of authenticating against authorization server
type oauth2Authenticator struct {
	url      string
	user     string
	password string
}

// NewOAuth2Authenticator instantiates new OAuth2Authenticator
func NewOAuth2Authenticator(url, user, password string) Authenticator {
	return oauth2Authenticator{url: url, user: user, password: password}
}

// Authenticate against the configured OAuth2 authorization server
func (auth oauth2Authenticator) Authenticate() (Authentication, error) {
	req := tokenRequest{user: auth.user, password: string(auth.password), scopes: []string{"uid"}}
	token, err := auth.requestToken(req)
	if err != nil {
		return nil, err
	}
	return &bearerTokenAuthentication{token: *token}, nil
}

func (auth oauth2Authenticator) requestToken(tokReq tokenRequest) (*string, error) {

	url, err := buildTokenURL(auth.url, tokReq)
	if err != nil {
		return nil, err
	}
	info, err := requestTokenInfo(url, tokReq)
	if err != nil {
		return nil, err
	}

	if val, exists := info[accessToken]; exists {
		token := val.(string)
		return &token, nil
	}
	return nil, fmt.Errorf("The access token couldn't be aquired")
}

func requestTokenInfo(url *url.URL, tokReq tokenRequest) (map[string]interface{}, error) {

	client := createClient()

	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(tokReq.user, tokReq.password)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Authentication failed, server has returned %d status code", resp.StatusCode)
	}
	return decodeMap(resp.Body)
}

func buildTokenURL(rawurl string, req tokenRequest) (*url.URL, error) {

	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	u.Scheme = https
	q := u.Query()
	q.Add("json", "true")
	if req.scopes != nil {
		for _, scope := range req.scopes {
			q.Add("scope", scope)
		}
	}
	u.RawQuery = q.Encode()
	return u, nil
}

func decodeMap(body io.ReadCloser) (map[string]interface{}, error) {
	var result map[string]interface{}
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func createClient() *http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: transport}
}
