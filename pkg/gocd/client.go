package gocd

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Config holds parameters to connect to a GoCD server.
type Config struct {
	BaseURL   string // e.g., https://gocd.example.com/go/api
	Username  string
	Password  string
	Token     string // Personal access token (takes precedence over Username/Password if set)
	Insecure  bool   // Skip TLS verification
	UserAgent string // Optional custom user agent
	Timeout   time.Duration
}

// Client is a minimal interface for interacting with GoCD API features used by this provider.
// Extend as needed with more services.
type Client interface {
	AuthorizationConfigurations() AuthorizationConfigurationsService
	Roles() RolesService
	PipelineConfigs() PipelineConfigsService
	ElasticAgentProfile() ElasticAgentProfileService
}

// APIError represents an error returned by the GoCD API.
type APIError struct {
	StatusCode int
	Message    string
	Body       string
}

func (e *APIError) Error() string {
	if e.Body != "" {
		return fmt.Sprintf("gocd: %s (status %d): %s", e.Message, e.StatusCode, e.Body)
	}
	return fmt.Sprintf("gocd: %s (status %d)", e.Message, e.StatusCode)
}

// IsNotFound returns true if the error is a 404 Not Found.
func IsNotFound(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusNotFound
	}
	return false
}

// client is an http-based implementation of Client.
type client struct {
	http  *http.Client
	base  *url.URL
	ua    string
	token string
	basic bool
	user  string
	pass  string
}

// New creates a new GoCD API client.
func New(cfg Config) (Client, error) {
	if cfg.BaseURL == "" {
		return nil, fmt.Errorf("gocd: BaseURL is required")
	}
	// Ensure trailing slash not required; we'll join paths.
	u, err := url.Parse(cfg.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("gocd: invalid BaseURL: %w", err)
	}

	tr := &http.Transport{}
	if u.Scheme == "https" {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: cfg.Insecure} // #nosec G402: intentional, controlled by config
	}
	httpClient := &http.Client{Transport: tr}
	if cfg.Timeout > 0 {
		httpClient.Timeout = cfg.Timeout
	} else {
		httpClient.Timeout = 30 * time.Second
	}

	ua := cfg.UserAgent
	if ua == "" {
		ua = "provider-gocd/0.1"
	}

	c := &client{
		http:  httpClient,
		base:  u,
		ua:    ua,
		token: cfg.Token,
		user:  cfg.Username,
		pass:  cfg.Password,
		basic: cfg.Token == "",
	}
	return c, nil
}

// do builds and executes an HTTP request against the GoCD API.
func (c *client) do(ctx context.Context, method, path, accept string, headers map[string]string, body any) (*http.Response, error) {
	// Build full URL
	rel := &url.URL{Path: strings.TrimSuffix(c.base.Path, "/") + "/" + strings.TrimPrefix(path, "/")}
	u := *c.base
	u.Path = rel.Path

	var rdr io.ReadCloser
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		rdr = io.NopCloser(strings.NewReader(string(b)))
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), rdr)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.ua)
	if accept != "" {
		req.Header.Set("Accept", accept)
	} else {
		req.Header.Set("Accept", "application/vnd.go.cd.v1+json")
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	} else if c.basic {
		req.SetBasicAuth(c.user, c.pass)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    "request failed",
			Body:       string(b),
		}
	}

	return resp, nil
}

// decodeJSON decodes a JSON response and closes the body.
func decodeJSON(resp *http.Response, out any) error {
	defer resp.Body.Close() //nolint:errcheck
	dec := json.NewDecoder(resp.Body)
	return dec.Decode(out)
}
