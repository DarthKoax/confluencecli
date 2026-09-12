package api

import (
	"context"
	"fmt"
	"net/url"

	"github.com/darthkoax/confluencecli/internal/client"
)

type ContentStateService struct {
	client *client.Client
}

func NewContentStateService(c *client.Client) *ContentStateService {
	return &ContentStateService{client: c}
}

type ContentState struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Color       string `json:"color,omitempty"`
	FinalState  bool   `json:"isFinalState,omitempty"`
}

type ContentStateResult struct {
	Results []ContentState `json:"results"`
	Start   int            `json:"start"`
	Limit   int            `json:"limit"`
	Size    int            `json:"size"`
}

func (s *ContentStateService) GetAll(ctx context.Context, start, limit int) (*ContentStateResult, error) {
	if err := s.client.CheckEndpoint("contentstates"); err != nil {
		return nil, err
	}
	path := "/rest/api/content-states?"
	params := url.Values{}
	if start > 0 {
		params.Set("start", itoa(start))
	}
	if limit > 0 {
		params.Set("limit", itoa(limit))
	}
	if len(params) > 0 {
		path += params.Encode()
	}
	var result ContentStateResult
	if err := s.client.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *ContentStateService) Get(ctx context.Context, stateID string) (*ContentState, error) {
	if err := s.client.CheckEndpoint("contentstates"); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/rest/api/content-states/%s", url.PathEscape(stateID))
	var state ContentState
	if err := s.client.Get(ctx, path, &state); err != nil {
		return nil, err
	}
	return &state, nil
}

func (s *ContentStateService) Create(ctx context.Context, state *ContentState) (*ContentState, error) {
	if err := s.client.CheckEndpoint("contentstates"); err != nil {
		return nil, err
	}
	var result ContentState
	if err := s.client.Post(ctx, "/rest/api/content-states", state, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *ContentStateService) Update(ctx context.Context, stateID string, state *ContentState) (*ContentState, error) {
	if err := s.client.CheckEndpoint("contentstates"); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/rest/api/content-states/%s", url.PathEscape(stateID))
	var result ContentState
	if err := s.client.Put(ctx, path, state, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *ContentStateService) Delete(ctx context.Context, stateID string) error {
	if err := s.client.CheckEndpoint("contentstates"); err != nil {
		return err
	}
	path := fmt.Sprintf("/rest/api/content-states/%s", url.PathEscape(stateID))
	return s.client.Delete(ctx, path)
}

func (s *ContentStateService) GetContentByState(ctx context.Context, stateID string, start, limit int) (*ContentResult, error) {
	if err := s.client.CheckEndpoint("contentstates"); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/rest/api/content-states/%s/content?start=%d&limit=%d",
		url.PathEscape(stateID), start, limit)
	var result ContentResult
	if err := s.client.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
