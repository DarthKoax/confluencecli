package api

import (
	"context"
	"fmt"
	"net/url"

	"github.com/darthkoax/confluencecli/internal/client"
)

type RelationService struct {
	client *client.Client
}

func NewRelationService(c *client.Client) *RelationService {
	return &RelationService{client: c}
}

type Relation struct {
	Name       string                 `json:"name,omitempty"`
	RelationData map[string]interface{} `json:"relationData,omitempty"`
	Source     map[string]interface{} `json:"source,omitempty"`
	Target     map[string]interface{} `json:"target,omitempty"`
}

type RelationResult struct {
	Results []Relation `json:"results"`
	Start   int        `json:"start"`
	Limit   int        `json:"limit"`
	Size    int        `json:"size"`
}

func (s *RelationService) Get(ctx context.Context, relationName, sourceType, sourceKey, targetType, targetKey string, expand []string) (*Relation, error) {
	if err := s.client.CheckEndpoint("relations"); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/rest/api/relation/%s/%s/%s/to/%s/%s",
		url.PathEscape(relationName),
		url.PathEscape(sourceType), url.PathEscape(sourceKey),
		url.PathEscape(targetType), url.PathEscape(targetKey))
	if len(expand) > 0 {
		path += "?expand=" + joinStrings(expand)
	}
	var relation Relation
	if err := s.client.Get(ctx, path, &relation); err != nil {
		return nil, err
	}
	return &relation, nil
}

func (s *RelationService) Create(ctx context.Context, relationName, sourceType, sourceKey, targetType, targetKey string, body interface{}) (*Relation, error) {
	if err := s.client.CheckEndpoint("relations"); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/rest/api/relation/%s/%s/%s/to/%s/%s",
		url.PathEscape(relationName),
		url.PathEscape(sourceType), url.PathEscape(sourceKey),
		url.PathEscape(targetType), url.PathEscape(targetKey))
	var result Relation
	if err := s.client.Post(ctx, path, body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *RelationService) Delete(ctx context.Context, relationName, sourceType, sourceKey, targetType, targetKey string) error {
	if err := s.client.CheckEndpoint("relations"); err != nil {
		return err
	}
	path := fmt.Sprintf("/rest/api/relation/%s/%s/%s/to/%s/%s",
		url.PathEscape(relationName),
		url.PathEscape(sourceType), url.PathEscape(sourceKey),
		url.PathEscape(targetType), url.PathEscape(targetKey))
	return s.client.Delete(ctx, path)
}
