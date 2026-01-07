package gocd

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

const (
	acceptRoles     = "application/vnd.go.cd.v3+json" // API media type for Roles; version can be adjusted as needed.
	roleServicePath = "/go/api/admin/security/roles"
)

// RolesService defines methods for GoCD Roles API.
// See: https://api.gocd.org/current/#roles
// It supports both GoCD (config) roles and plugin roles via the Type and Attributes fields.
// All methods return the ETag from the server when applicable.
//
// Accepted status codes:
// - Get: 200, 404 (returns nil, "", nil)
// - Create: 200 or 201
// - Update: 200
// - Delete: 200 or 204
//
// ETag handling:
// - Update/Delete can include If-Match header when provided to handle concurrency.
// - The returned ETag (if any) is the value from the response header.
type RolesService interface {
	Get(ctx context.Context, name string) (*Role, string, error)
	Create(ctx context.Context, role Role) (*Role, string, error)
	Update(ctx context.Context, name string, role Role, etag string) (*Role, string, error)
	Delete(ctx context.Context, name string, etag string) error
}

// Attributes fields vary by Type:
// - For Type "gocd": use Users to list usernames.
// - For Type "plugin": set AuthConfigID and Properties according to the plugin.
// The server also returns standard HAL links.
//
// Docs: https://api.gocd.org/current/#roles

type Role struct {
	Name       string          `json:"name,omitempty"`
	Type       string          `json:"type,omitempty"`
	Attributes *RoleAttributes `json:"attributes,omitempty"`
	Policy     []Policy        `json:"policy,omitempty"`
	Links      *HALLinks       `json:"_links,omitempty"`
}

// RoleAttributes contains either built-in GoCD role data (Users) or plugin role attributes.
// Only relevant fields should be populated based on the role Type.

type RoleAttributes struct {
	// For GoCD roles
	Users []string `json:"users,omitempty"`
	// For plugin roles
	AuthConfigID string           `json:"auth_config_id,omitempty"`
	Properties   []ConfigProperty `json:"properties,omitempty"`
}

type Policy struct {
	Permission string `json:"permission"`
	Action     string `json:"action"`
	Type       string `json:"type"`
	Resource   string `json:"resource"`
}

type rolesService struct{ c *client }

func (c *client) Roles() RolesService { return &rolesService{c: c} }

func (s *rolesService) Get(ctx context.Context, name string) (*Role, string, error) {
	resp, err := s.c.do(ctx, http.MethodGet, fmt.Sprintf("%s/%s", roleServicePath, url.PathEscape(name)), acceptRoles, nil, nil)
	if err != nil {
		if IsNotFound(err) {
			return nil, "", nil
		}
		return nil, "", errors.Wrap(err, "gocd: failed to get role")
	}
	var out Role
	if err := decodeJSON(resp, &out); err != nil {
		return nil, "", err
	}
	etag := resp.Header.Get("ETag")
	return &out, etag, nil
}

func (s *rolesService) Create(ctx context.Context, role Role) (*Role, string, error) {
	resp, err := s.c.do(ctx, http.MethodPost, roleServicePath, acceptRoles, nil, role)
	if err != nil {
		return nil, "", err
	}
	var out Role
	if err := decodeJSON(resp, &out); err != nil {
		return nil, "", err
	}
	etag := resp.Header.Get("ETag")
	return &out, etag, nil
}

func (s *rolesService) Update(ctx context.Context, name string, role Role, etag string) (*Role, string, error) {
	path := fmt.Sprintf("%s/%s", roleServicePath, url.PathEscape(name))
	headers := map[string]string{
		"If-Match": etag,
	}
	resp, err := s.c.do(ctx, http.MethodPut, path, acceptRoles, headers, role)
	if err != nil {
		return nil, "", errors.Wrap(err, "gocd: failed to update role")
	}
	var out Role
	if err := decodeJSON(resp, &out); err != nil {
		return nil, "", err
	}
	newETag := resp.Header.Get("ETag")
	return &out, newETag, nil
}

func (s *rolesService) Delete(ctx context.Context, name string, etag string) error {
	path := fmt.Sprintf("%s/%s", roleServicePath, url.PathEscape(name))
	headers := make(map[string]string)
	if etag != "" {
		headers["If-Match"] = etag
	}
	resp, err := s.c.do(ctx, http.MethodDelete, path, acceptRoles, headers, nil)
	if err != nil {
		return err
	}
	return resp.Body.Close()
}
