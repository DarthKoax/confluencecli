package api

import (
	"context"
	"fmt"

	"github.com/darkkoax/confluencecli/internal/client"
)

type LongTaskService struct {
	client *client.Client
}

func NewLongTaskService(c *client.Client) *LongTaskService {
	return &LongTaskService{client: c}
}

type LongTask struct {
	ID                 string `json:"id,omitempty"`
	Name               map[string]string `json:"name,omitempty"`
	ElapsedTime        int64  `json:"elapsedTime,omitempty"`
	PercentageComplete int    `json:"percentageComplete,omitempty"`
	Successful         bool   `json:"successful,omitempty"`
	Finished           bool   `json:"finished,omitempty"`
	Messages           []map[string]string `json:"messages,omitempty"`
	Status             map[string]string `json:"status,omitempty"`
	Errors             []map[string]string `json:"errors,omitempty"`
}

type LongTaskResult struct {
	Results []LongTask `json:"results"`
	Start   int        `json:"start"`
	Limit   int        `json:"limit"`
	Size    int        `json:"size"`
}

func (s *LongTaskService) Get(ctx context.Context, taskID string) (*LongTask, error) {
	if err := s.client.CheckEndpoint("longtasks"); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/rest/api/longtask/%s", taskID)
	var task LongTask
	if err := s.client.Get(ctx, path, &task); err != nil {
		return nil, err
	}
	return &task, nil
}

func (s *LongTaskService) GetAll(ctx context.Context, start, limit int) (*LongTaskResult, error) {
	if err := s.client.CheckEndpoint("longtasks"); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/rest/api/longtask?start=%d&limit=%d", start, limit)
	var result LongTaskResult
	if err := s.client.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
