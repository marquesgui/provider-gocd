package gocd

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

const (
	elasticAgentProfileSerivcePath = "/go/api/elastic/profiles"
	acceptElasticAgentProfile      = "application/vnd.go.cd.v2+json"
)

type ElasticAgentProfile struct {
	ID               string `json:"id"`
	ClusterProfileID string `json:"cluster_profile_id"`
	Properties       []ConfigProperty
}

type ElasticAgentProfileResponse struct {
	Links HALLinks `json:"_links"`
	ElasticAgentProfile
}

type ElasticAgentProfileService interface {
	Get(ctx context.Context, profileID string) (*ElasticAgentProfileResponse, string, error)
	Create(ctx context.Context, eap ElasticAgentProfile) (*ElasticAgentProfileResponse, string, error)
	Update(ctx context.Context, eap ElasticAgentProfile, eTag string) (*ElasticAgentProfileResponse, string, error)
	Delete(ctx context.Context, id string) error
}

func (c *client) ElasticAgentProfile() ElasticAgentProfileService {
	return &elasticAgentProfileService{c: c}
}

type elasticAgentProfileService struct {
	c *client
}

func (e *elasticAgentProfileService) Get(ctx context.Context, profileID string) (*ElasticAgentProfileResponse, string, error) {
	path := fmt.Sprintf("%s/%s", elasticAgentProfileSerivcePath, url.PathEscape(profileID))
	resp, err := e.c.do(ctx, http.MethodGet, path, acceptElasticAgentProfile, nil, nil)
	if err != nil {
		if IsNotFound(err) {
			return nil, "", nil
		}
		return nil, "", errors.Wrap(err, "gocd: failed to get the elastic agent profile")
	}

	var result ElasticAgentProfileResponse
	if err := decodeJSON(resp, &result); err != nil {
		return nil, "", errors.Wrap(err, "gocd: failed to decode response")
	}
	return &result, resp.Header.Get("ETag"), nil
}

func (e *elasticAgentProfileService) Create(ctx context.Context, eap ElasticAgentProfile) (*ElasticAgentProfileResponse, string, error) {
	resp, err := e.c.do(ctx, http.MethodPost, elasticAgentProfileSerivcePath, acceptElasticAgentProfile, nil, eap)
	if err != nil {
		return nil, "", errors.Wrap(err, "gocd: failed to create the elastic agent profile")
	}

	var result ElasticAgentProfileResponse
	if err := decodeJSON(resp, &result); err != nil {
		return nil, "", errors.Wrap(err, "gocd: failed to decode response")
	}
	return &result, resp.Header.Get("ETag"), nil
}

func (e *elasticAgentProfileService) Update(ctx context.Context, eap ElasticAgentProfile, etag string) (*ElasticAgentProfileResponse, string, error) {
	path := fmt.Sprintf("%s/%s", elasticAgentProfileSerivcePath, url.PathEscape(eap.ID))
	headers := map[string]string{
		"If-Match": etag,
	}
	resp, err := e.c.do(ctx, http.MethodPut, path, acceptElasticAgentProfile, headers, eap)
	if err != nil {
		return nil, "", errors.Wrap(err, fmt.Sprintf("gocd: could not update the elastic agent profile of id %s", eap.ID))
	}

	var newEap ElasticAgentProfileResponse
	if err := decodeJSON(resp, &newEap); err != nil {
		return nil, "", errors.Wrap(err, "gocd: could not decode http body")
	}
	newetag := resp.Header.Get("ETag")

	return &newEap, newetag, nil
}

func (e *elasticAgentProfileService) Delete(ctx context.Context, profileID string) error {
	path := fmt.Sprintf("%s/%s", elasticAgentProfileSerivcePath, url.PathEscape(profileID))
	resp, err := e.c.do(ctx, http.MethodDelete, path, acceptElasticAgentProfile, nil, nil)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("cannot delete the elastic agent profie with id %s", profileID))
	}
	return resp.Body.Close()
}
