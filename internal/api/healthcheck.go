package api

import (
	"context"

	"github.com/darkkoax/confluencecli/internal/client"
)

type HealthCheckService struct {
	client *client.Client
}

func NewHealthCheckService(c *client.Client) *HealthCheckService {
	return &HealthCheckService{client: c}
}

type HealthCheckResult struct {
	Name    string `json:"name"`
	Description string `json:"description"`
	Passed  bool   `json:"passed"`
	Failed  bool   `json:"failed,omitempty"`
	Severity string `json:"severity,omitempty"`
}

type HealthCheckResponse struct {
	State  string              `json:"state"`
	Checks []HealthCheckResult `json:"checks"`
}

func (s *HealthCheckService) Get(ctx context.Context) (*HealthCheckResponse, error) {
	if err := s.client.CheckEndpoint("healthcheck"); err != nil {
		return nil, err
	}
	var result HealthCheckResponse
	if err := s.client.Get(ctx, "/status", &result); err != nil {
		return nil, err
	}
	return &result, nil
}
