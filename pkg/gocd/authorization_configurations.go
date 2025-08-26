package gocd

import (
  "context"
  "fmt"
  "io"
  "net/http"
  "net/url"

  "github.com/pkg/errors"
)

// AuthorizationConfigurationsService defines methods for GoCD Authorization Configurations.
type AuthorizationConfigurationsService interface {
  // Get returns an authorization configuration by id.
  Get(ctx context.Context, id string) (*AuthorizationConfiguration, string, error)
  // Create creates a new authorization configuration.
  Create(ctx context.Context, cfg AuthorizationConfiguration) (*AuthorizationConfiguration, string, error)
  // Update updates an existing authorization configuration using If-Match etag.
  Update(ctx context.Context, id string, cfg AuthorizationConfiguration, etag string) (*AuthorizationConfiguration, string, error)
  // Delete deletes an authorization configuration by id.
  Delete(ctx context.Context, id string, etag string) error
}

// authorizationConfigurationsService implements AuthorizationConfigurationsService.
type authorizationConfigurationsService struct{ c *client }

func (c *client) AuthorizationConfigurations() AuthorizationConfigurationsService {
  return &authorizationConfigurationsService{c: c}
}

const (
  acceptAuthzCfg = "application/vnd.go.cd.v2+json" // Based on GoCD API; adjust if needed.
  servicePath    = "/go/api/admin/security/auth_configs"
)

// AuthorizationConfiguration represents GoCD Authorization Configuration payload.
type AuthorizationConfiguration struct {
  ID                         string           `json:"id,omitempty"`
  PluginID                   string           `json:"plugin_id,omitempty"`
  AllowOnlyKnownUsersToLogin bool             `json:"allow_only_known_users_to_login"`
  Properties                 []ConfigProperty `json:"properties,omitempty"`
  Links                      *HALLinks        `json:"_links,omitempty"`
}

func (s *authorizationConfigurationsService) Get(ctx context.Context, id string) (*AuthorizationConfiguration, string, error) {
  resp, err := s.c.do(ctx, http.MethodGet, fmt.Sprintf("%s/%s", servicePath, url.PathEscape(id)), acceptAuthzCfg, nil, nil)
  if err != nil {
    return nil, "", err
  }
  if resp.StatusCode == http.StatusNotFound {
    _ = resp.Body.Close()
    return nil, "", nil
  }
  if resp.StatusCode != http.StatusOK {
    b, _ := io.ReadAll(resp.Body)
    _ = resp.Body.Close()
    return nil, "", fmt.Errorf("gocd: unexpected status %d: %s", resp.StatusCode, string(b))
  }
  var out AuthorizationConfiguration
  if err := decodeJSON(resp, &out); err != nil {
    return nil, "", err
  }
  etag := resp.Header.Get("ETag")
  return &out, etag, nil
}

func (s *authorizationConfigurationsService) Create(ctx context.Context, cfg AuthorizationConfiguration) (*AuthorizationConfiguration, string, error) {
  resp, err := s.c.do(ctx, http.MethodPost, servicePath, acceptAuthzCfg, nil, cfg)
  if err != nil {
    return nil, "", err
  }
  if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
    b, _ := io.ReadAll(resp.Body)
    _ = resp.Body.Close()
    return nil, "", fmt.Errorf("gocd: unexpected status %d: %s", resp.StatusCode, string(b))
  }
  var out AuthorizationConfiguration
  if err := decodeJSON(resp, &out); err != nil {
    return nil, "", err
  }
  etag := resp.Header.Get("ETag")
  return &out, etag, nil
}

func (s *authorizationConfigurationsService) Update(ctx context.Context, id string, cfg AuthorizationConfiguration, etag string) (*AuthorizationConfiguration, string, error) {
  path := fmt.Sprintf("%s/%s", servicePath, url.PathEscape(id))
  headers := map[string]string{
    "If-Match": etag,
  }
  resp, err := s.c.do(ctx, http.MethodPut, path, acceptAuthzCfg, headers, cfg)
  if err != nil {
    return nil, "", errors.Wrap(err, "cannot update authorization configuration")
  }
  if resp.StatusCode != http.StatusOK {
    b, _ := io.ReadAll(resp.Body)
    _ = resp.Body.Close()
    return nil, "", fmt.Errorf("gocd: unexpected status %d: %s", resp.StatusCode, string(b))
  }
  var out AuthorizationConfiguration
  if err := decodeJSON(resp, &out); err != nil {
    return nil, "", fmt.Errorf("gocd: failed to decode response: %w", err)
  }
  newETag := resp.Header.Get("ETag")
  return &out, newETag, nil
}

func (s *authorizationConfigurationsService) Delete(ctx context.Context, id string, etag string) error {
  path := fmt.Sprintf("%s/%s", servicePath, url.PathEscape(id))
  headers := map[string]string{
    "If-Match": etag,
  }
  resp, err := s.c.do(ctx, http.MethodDelete, path, acceptAuthzCfg, headers, nil)
  if err != nil {
    return errors.Wrap(err, "cannot delete authorization configuration")
  }
  if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
    b, _ := io.ReadAll(resp.Body)
    _ = resp.Body.Close()
    return fmt.Errorf("gocd: unexpected status %d: %s", resp.StatusCode, string(b))
  }
  return nil
}
