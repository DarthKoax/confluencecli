package api

import (
	"context"
	"fmt"
	"net/url"

	"github.com/darkkoax/confluencecli/internal/client"
)

type BlueprintService struct {
	client *client.Client
}

func NewBlueprintService(c *client.Client) *BlueprintService {
	return &BlueprintService{client: c}
}

type ContentBlueprint struct {
	ContentBlueprintID string `json:"contentBlueprintId,omitempty"`
	Name               string `json:"name,omitempty"`
	Description        string `json:"description,omitempty"`
	PluginKey          string `json:"pluginKey,omitempty"`
	TemplateID         string `json:"templateId,omitempty"`
}

type ContentBlueprintResult struct {
	Results []ContentBlueprint `json:"results"`
	Start   int                `json:"start"`
	Limit   int                `json:"limit"`
	Size    int                `json:"size"`
}

func (s *BlueprintService) GetAll(ctx context.Context, start, limit int, expand []string) (*ContentBlueprintResult, error) {
	if err := s.client.CheckEndpoint("blueprints"); err != nil {
		return nil, err
	}
	path := "/rest/api/content-blueprint?"
	params := url.Values{}
	if start > 0 {
		params.Set("start", itoa(start))
	}
	if limit > 0 {
		params.Set("limit", itoa(limit))
	}
	if len(expand) > 0 {
		params.Set("expand", joinStrings(expand))
	}
	if len(params) > 0 {
		path += params.Encode()
	}
	var result ContentBlueprintResult
	if err := s.client.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *BlueprintService) Get(ctx context.Context, blueprintID string, expand []string) (*ContentBlueprint, error) {
	if err := s.client.CheckEndpoint("blueprints"); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/rest/api/content-blueprint/%s", url.PathEscape(blueprintID))
	if len(expand) > 0 {
		path += "?expand=" + joinStrings(expand)
	}
	var bp ContentBlueprint
	if err := s.client.Get(ctx, path, &bp); err != nil {
		return nil, err
	}
	return &bp, nil
}
