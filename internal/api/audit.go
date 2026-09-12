package api

import (
	"context"
	"net/url"
	"strconv"

	"github.com/darthkoax/confluencecli/internal/client"
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
	CreatedDate string                 `json:"creationDate,omitempty"`
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

type AuditRetention struct {
	Number int    `json:"number,omitempty"`
	Units  string `json:"units,omitempty"`
}

type AuditExport struct {
	ExportID string `json:"exportId,omitempty"`
	Token    string `json:"token,omitempty"`
}

func (s *AuditService) Get(ctx context.Context, startDate, endDate string, searchString string, start, limit int) (*AuditResult, error) {
	if err := s.client.CheckEndpoint("audit"); err != nil {
		return nil, err
	}
	if err := s.client.CheckAdmin("audit", "get"); err != nil {
		return nil, err
	}
	params := url.Values{}
	if startDate != "" {
		params.Set("startDate", startDate)
	}
	if endDate != "" {
		params.Set("endDate", endDate)
	}
	if searchString != "" {
		params.Set("searchString", searchString)
	}
	if start > 0 {
		params.Set("start", strconv.Itoa(start))
	}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	path := "/rest/audit/1.0/audit?" + params.Encode()
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
	params := url.Values{}
	params.Set("number", number)
	if searchString != "" {
		params.Set("searchString", searchString)
	}
	params.Set("start", strconv.Itoa(start))
	params.Set("limit", strconv.Itoa(limit))
	path := "/rest/audit/1.0/audit/since?" + params.Encode()
	var result AuditResult
	if err := s.client.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *AuditService) GetRetentionPeriod(ctx context.Context) (*AuditRetention, error) {
	if err := s.client.CheckEndpoint("audit"); err != nil {
		return nil, err
	}
	if err := s.client.CheckAdmin("audit", "get"); err != nil {
		return nil, err
	}
	var result AuditRetention
	if err := s.client.Get(ctx, "/rest/audit/1.0/audit/retention", &result); err != nil {
		return nil, err
	}
	return &result, nil
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

func (s *AuditService) Export(ctx context.Context, startDate, endDate string) (*AuditExport, error) {
	if err := s.client.CheckEndpoint("audit"); err != nil {
		return nil, err
	}
	if err := s.client.CheckAdmin("audit", "get"); err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Set("startDate", startDate)
	params.Set("endDate", endDate)
	path := "/rest/audit/1.0/audit/export?" + params.Encode()
	var result AuditExport
	if err := s.client.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
