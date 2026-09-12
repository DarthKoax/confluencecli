package api

import (
	"context"
	"fmt"
	"net/url"

	"github.com/darthkoax/confluencecli/internal/client"
)

type SearchService struct {
	client *client.Client
}

func NewSearchService(c *client.Client) *SearchService {
	return &SearchService{client: c}
}

type SearchResult struct {
	Results        []SearchResultItem `json:"results"`
	Start          int                `json:"start"`
	Limit          int                `json:"limit"`
	Size           int                `json:"size"`
	TotalSize      int                `json:"totalSize"`
	CQLQuery       string             `json:"cqlQuery"`
	SearchDuration int                `json:"searchDuration"`
}

type SearchResultItem struct {
	Content         *Content           `json:"content,omitempty"`
	User            *User              `json:"user,omitempty"`
	Space           *Space             `json:"space,omitempty"`
	Title           string             `json:"title,omitempty"`
	Excerpt         string             `json:"excerpt,omitempty"`
	URL             string             `json:"url,omitempty"`
	EntityType      string             `json:"entityType,omitempty"`
	LastModified    string             `json:"lastModified,omitempty"`
	FriendlyDate    string             `json:"friendlyLastModified,omitempty"`
}

func (s *SearchService) Search(ctx context.Context, cql string, start, limit int, expand []string) (*SearchResult, error) {
	if err := s.client.CheckEndpoint("search"); err != nil {
		return nil, err
	}
	path := "/rest/api/search?"
	params := url.Values{}
	params.Set("cql", cql)
	if start > 0 {
		params.Set("start", itoa(start))
	}
	if limit > 0 {
		params.Set("limit", itoa(limit))
	}
	if len(expand) > 0 {
		params.Set("expand", joinStrings(expand))
	}
	path += params.Encode()

	var result SearchResult
	if err := s.client.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *SearchService) SearchUser(ctx context.Context, cqlUser string, start, limit int) (*SearchResult, error) {
	if err := s.client.CheckEndpoint("search"); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/rest/api/search/user?cql=%s&start=%d&limit=%d",
		url.QueryEscape(cqlUser), start, limit)
	var result SearchResult
	if err := s.client.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
