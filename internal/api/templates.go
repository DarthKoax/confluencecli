package api

import (
	"context"
	"fmt"
	"net/url"

	"github.com/darkkoax/confluencecli/internal/client"
)

type TemplateService struct {
	client *client.Client
}

func NewTemplateService(c *client.Client) *TemplateService {
	return &TemplateService{client: c}
}

type ContentTemplate struct {
	TemplateID    string                 `json:"templateId,omitempty"`
	Name          string                 `json:"name,omitempty"`
	Description   string                 `json:"description,omitempty"`
	SpaceKey      string                 `json:"spaceKey,omitempty"`
	TemplateType  string                 `json:"templateType,omitempty"`
	Body          map[string]interface{} `json:"body,omitempty"`
	Labels        []map[string]string    `json:"labels,omitempty"`
}

type ContentTemplateResult struct {
	Results []ContentTemplate `json:"results"`
	Start   int               `json:"start"`
	Limit   int               `json:"limit"`
	Size    int               `json:"size"`
}

type BlueprintTemplate struct {
	ContentBlueprintID string `json:"contentBlueprintId,omitempty"`
	Name               string `json:"name,omitempty"`
	Description        string `json:"description,omitempty"`
	PluginKey          string `json:"pluginKey,omitempty"`
	TemplateID         string `json:"templateId,omitempty"`
}

type BlueprintTemplateResult struct {
	Results []BlueprintTemplate `json:"results"`
	Start   int                 `json:"start"`
	Limit   int                 `json:"limit"`
	Size    int                 `json:"size"`
}

func (s *TemplateService) GetContentTemplates(ctx context.Context, spaceKey string, start, limit int, expand []string) (*ContentTemplateResult, error) {
	if err := s.client.CheckEndpoint("templates"); err != nil {
		return nil, err
	}
	path := "/rest/api/template/page?"
	params := url.Values{}
	if spaceKey != "" {
		params.Set("spaceKey", spaceKey)
	}
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
	var result ContentTemplateResult
	if err := s.client.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *TemplateService) GetContentTemplate(ctx context.Context, contentTemplateID string, expand []string) (*ContentTemplate, error) {
	if err := s.client.CheckEndpoint("templates"); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/rest/api/template/page/%s", url.PathEscape(contentTemplateID))
	if len(expand) > 0 {
		path += "?expand=" + joinStrings(expand)
	}
	var tmpl ContentTemplate
	if err := s.client.Get(ctx, path, &tmpl); err != nil {
		return nil, err
	}
	return &tmpl, nil
}

func (s *TemplateService) CreateContentTemplate(ctx context.Context, tmpl *ContentTemplate) (*ContentTemplate, error) {
	if err := s.client.CheckEndpoint("templates"); err != nil {
		return nil, err
	}
	var result ContentTemplate
	if err := s.client.Post(ctx, "/rest/api/template/page", tmpl, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *TemplateService) UpdateContentTemplate(ctx context.Context, tmpl *ContentTemplate) (*ContentTemplate, error) {
	if err := s.client.CheckEndpoint("templates"); err != nil {
		return nil, err
	}
	var result ContentTemplate
	if err := s.client.Put(ctx, "/rest/api/template/page", tmpl, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *TemplateService) RemoveContentTemplate(ctx context.Context, contentTemplateID string) error {
	if err := s.client.CheckEndpoint("templates"); err != nil {
		return err
	}
	path := fmt.Sprintf("/rest/api/template/page/%s", url.PathEscape(contentTemplateID))
	return s.client.Delete(ctx, path)
}

func (s *TemplateService) GetBlueprintTemplates(ctx context.Context, spaceKey string, start, limit int, expand []string) (*BlueprintTemplateResult, error) {
	if err := s.client.CheckEndpoint("templates"); err != nil {
		return nil, err
	}
	path := "/rest/api/template/blueprint?"
	params := url.Values{}
	if spaceKey != "" {
		params.Set("spaceKey", spaceKey)
	}
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
	var result BlueprintTemplateResult
	if err := s.client.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *TemplateService) GetBlueprintTemplate(ctx context.Context, contentBlueprintID string, expand []string) (*BlueprintTemplate, error) {
	if err := s.client.CheckEndpoint("templates"); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/rest/api/template/blueprint/%s", url.PathEscape(contentBlueprintID))
	if len(expand) > 0 {
		path += "?expand=" + joinStrings(expand)
	}
	var tmpl BlueprintTemplate
	if err := s.client.Get(ctx, path, &tmpl); err != nil {
		return nil, err
	}
	return &tmpl, nil
}
