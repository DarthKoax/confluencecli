package api

import (
	"context"
	"fmt"
	"net/url"

	"github.com/darthkoax/confluencecli/internal/client"
)

type InlineTaskService struct {
	client *client.Client
}

func NewInlineTaskService(c *client.Client) *InlineTaskService {
	return &InlineTaskService{client: c}
}

type InlineTask struct {
	ID            interface{} `json:"id,omitempty"`
	ContentID     interface{} `json:"contentId,omitempty"`
	Status        string      `json:"status,omitempty"`
	Creator       string      `json:"creator,omitempty"`
	Assignee      string      `json:"assignee,omitempty"`
	CompleteUser  string      `json:"completeUser,omitempty"`
	CompleteDate  int64       `json:"completeDate,omitempty"`
	DueDate       int64       `json:"dueDate,omitempty"`
	TaskDetailURL string      `json:"taskDetailUrl,omitempty"`
}

type InlineTaskResult struct {
	Results []InlineTask `json:"results"`
	Start   int          `json:"start"`
	Limit   int          `json:"limit"`
	Size    int          `json:"size"`
}

func (s *InlineTaskService) Get(ctx context.Context, taskID string) (*InlineTask, error) {
	if err := s.client.CheckEndpoint("inlinetasks"); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/rest/api/inlinetasks/%s", url.PathEscape(taskID))
	var task InlineTask
	if err := s.client.Get(ctx, path, &task); err != nil {
		return nil, err
	}
	return &task, nil
}

func (s *InlineTaskService) GetAll(ctx context.Context, start, limit int) (*InlineTaskResult, error) {
	if err := s.client.CheckEndpoint("inlinetasks"); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/rest/api/inlinetasks?start=%d&limit=%d", start, limit)
	var result InlineTaskResult
	if err := s.client.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *InlineTaskService) GetByContent(ctx context.Context, contentID string, start, limit int) (*InlineTaskResult, error) {
	if err := s.client.CheckEndpoint("inlinetasks"); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/rest/api/content/%s/inline-task?start=%d&limit=%d",
		url.PathEscape(contentID), start, limit)
	var result InlineTaskResult
	if err := s.client.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *InlineTaskService) Update(ctx context.Context, contentID, taskID string, status string) (*InlineTask, error) {
	if err := s.client.CheckEndpoint("inlinetasks"); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/rest/api/content/%s/inline-task/%s", url.PathEscape(contentID), url.PathEscape(taskID))
	body := map[string]string{"status": status}
	var result InlineTask
	if err := s.client.Put(ctx, path, body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
