package api

import (
	"context"
	"fmt"
	"net/url"

	"github.com/darthkoax/confluencecli/internal/client"
)

type SpaceService struct {
	client *client.Client
}

func NewSpaceService(c *client.Client) *SpaceService {
	return &SpaceService{client: c}
}

type Space struct {
	ID          int               `json:"id,omitempty"`
	Key         string            `json:"key,omitempty"`
	Name        string            `json:"name,omitempty"`
	Type        string            `json:"type,omitempty"`
	Status      string            `json:"status,omitempty"`
	Description map[string]interface{} `json:"description,omitempty"`
	Homepage    *Content          `json:"homepage,omitempty"`
	Links       map[string]string `json:"_links,omitempty"`
}

type SpaceResult struct {
	Results []Space `json:"results"`
	Start   int     `json:"start"`
	Limit   int     `json:"limit"`
	Size    int     `json:"size"`
}

type SpaceContent struct {
	Page     *ContentResult `json:"page,omitempty"`
	Blogpost *ContentResult `json:"blogpost,omitempty"`
}

func (s *SpaceService) GetAll(ctx context.Context, spaceType, status string, start, limit int, expand []string) (*SpaceResult, error) {
	if err := s.client.CheckEndpoint("spaces"); err != nil {
		return nil, err
	}
	path := "/rest/api/space?"
	params := url.Values{}
	if spaceType != "" {
		params.Set("type", spaceType)
	}
	if status != "" {
		params.Set("status", status)
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
	var result SpaceResult
	if err := s.client.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *SpaceService) Get(ctx context.Context, spaceKey string, expand []string) (*Space, error) {
	if err := s.client.CheckEndpoint("spaces"); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/rest/api/space/%s", url.PathEscape(spaceKey))
	if len(expand) > 0 {
		path += "?expand=" + joinStrings(expand)
	}
	var space Space
	if err := s.client.Get(ctx, path, &space); err != nil {
		return nil, err
	}
	return &space, nil
}

func (s *SpaceService) Create(ctx context.Context, space *Space) (*Space, error) {
	if err := s.client.CheckEndpoint("spaces"); err != nil {
		return nil, err
	}
	var result Space
	if err := s.client.Post(ctx, "/rest/api/space", space, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *SpaceService) Delete(ctx context.Context, spaceKey string) error {
	if err := s.client.CheckEndpoint("spaces"); err != nil {
		return err
	}
	path := fmt.Sprintf("/rest/api/space/%s", url.PathEscape(spaceKey))
	return s.client.Delete(ctx, path)
}

func (s *SpaceService) GetContent(ctx context.Context, spaceKey, contentType string, start, limit int, expand []string) (*SpaceContent, error) {
	if err := s.client.CheckEndpoint("spaces"); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/rest/api/space/%s/content", url.PathEscape(spaceKey))
	if contentType != "" {
		path = fmt.Sprintf("/rest/api/space/%s/content/%s", url.PathEscape(spaceKey), contentType)
	}
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
		path += "?" + params.Encode()
	}
	if contentType != "" {
		var result ContentResult
		if err := s.client.Get(ctx, path, &result); err != nil {
			return nil, err
		}
		sc := &SpaceContent{}
		switch contentType {
		case "page":
			sc.Page = &result
		case "blogpost":
			sc.Blogpost = &result
		}
		return sc, nil
	}
	var result SpaceContent
	if err := s.client.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
