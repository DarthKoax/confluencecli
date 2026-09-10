package api

import (
	"context"
	"fmt"

	"github.com/darkkoax/confluencecli/internal/client"
)

type AuditService struct {
	client *client.Client
}

func NewAuditService(c *client.Client) *AuditService {
	return &AuditService{client: c}
}

type AuditRecord struct {
	Author      *User                  `json:"author,omitempty"`
	RemoteAddr  string                 `json:"remoteAddress,omitempty"`
	CreatedDate int64                  `json:"creationDate,omitempty"`
	Summary     string                 `json:"summary,omitempty"`
	Description string                 `json:"description,omitempty"`
	Category    string                 `json:"category,omitempty"`
	SysAdmin    bool                   `json:"sysAdmin,omitempty"`
	ChangedValues []AuditChangedValue  `json:"changedValues,omitempty"`
	AssociatedObjects []AuditAssociatedObject `json:"associatedObjects,omitempty"`
}

type AuditChangedValue struct {
	Name     string `json:"name,omitempty"`
	OldValue string `json:"oldValue,omitempty"`
	NewValue string `json:"newValue,omitempty"`
}

type AuditAssociatedObject struct {
	Name       string `json:"name,omitempty"`
	ObjectType string `json:"objectType,omitempty"`
}

type AuditResult struct {
	Results []AuditRecord `json:"results"`
	Start   int           `json:"start"`
	Limit   int           `json:"limit"`
	Size    int           `json:"size"`
}

func (s *AuditService) Get(ctx context.Context, startDate, endDate string, searchString string, start, limit int) (*AuditResult, error) {
	if err := s.client.CheckEndpoint("audit"); err != nil {
		return nil, err
	}
	if err := s.client.CheckAdmin("audit", "get"); err != nil {
		return nil, err
	}
	path := "/rest/audit/1.0/audit?"
	params := ""
	if startDate != "" {
		params += "startDate=" + startDate + "&"
	}
	if endDate != "" {
		params += "endDate=" + endDate + "&"
	}
	if searchString != "" {
		params += "searchString=" + searchString + "&"
	}
	if start > 0 {
		params += "start=" + itoa(start) + "&"
	}
	if limit > 0 {
		params += "limit=" + itoa(limit) + "&"
	}
	path += params
	var result AuditResult
	if err := s.client.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *AuditService) GetSince(ctx context.Context, number string, searchString string, start, limit int) (*AuditResult, error) {
	if err := s.client.CheckEndpoint("audit"); err != nil {
		return nil, err
	}
	if err := s.client.CheckAdmin("audit", "get"); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/rest/audit/1.0/audit/since?number=%s&searchString=%s&start=%d&limit=%d",
		number, searchString, start, limit)
	var result AuditResult
	if err := s.client.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *AuditService) GetRetentionPeriod(ctx context.Context) (map[string]interface{}, error) {
	if err := s.client.CheckEndpoint("audit"); err != nil {
		return nil, err
	}
	if err := s.client.CheckAdmin("audit", "get"); err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := s.client.Get(ctx, "/rest/audit/1.0/audit/retention", &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *AuditService) SetRetentionPeriod(ctx context.Context, number int, units string) error {
	if err := s.client.CheckEndpoint("audit"); err != nil {
		return err
	}
	if err := s.client.CheckAdmin("audit", "get"); err != nil {
		return err
	}
	body := map[string]interface{}{
		"number": number,
		"units":  units,
	}
	return s.client.Put(ctx, "/rest/audit/1.0/audit/retention", body, nil)
}

func (s *AuditService) Export(ctx context.Context, startDate, endDate string) (map[string]interface{}, error) {
	if err := s.client.CheckEndpoint("audit"); err != nil {
		return nil, err
	}
	if err := s.client.CheckAdmin("audit", "get"); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/rest/audit/1.0/audit/export?startDate=%s&endDate=%s", startDate, endDate)
	var result map[string]interface{}
	if err := s.client.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result, nil
}
