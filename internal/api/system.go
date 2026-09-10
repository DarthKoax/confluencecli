package api

import (
	"context"

	"github.com/darkkoax/confluencecli/internal/client"
)

type SystemService struct {
	client *client.Client
}

func NewSystemService(c *client.Client) *SystemService {
	return &SystemService{client: c}
}

type SystemStatus struct {
	State string `json:"state,omitempty"`
}

func (s *SystemService) GetStatus(ctx context.Context) (*SystemStatus, error) {
	if err := s.client.CheckEndpoint("system"); err != nil {
		return nil, err
	}
	if err := s.client.CheckAdmin("system", "status"); err != nil {
		return nil, err
	}
	var status SystemStatus
	if err := s.client.Get(ctx, "/rest/api/system/status", &status); err != nil {
		return nil, err
	}
	return &status, nil
}
