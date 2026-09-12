package api

import (
	"context"
	"fmt"
	"net/url"

	"github.com/darthkoax/confluencecli/internal/client"
)

type ContentService struct {
	client *client.Client
}

func NewContentService(c *client.Client) *ContentService {
	return &ContentService{client: c}
}

type Content struct {
	ID        string                 `json:"id,omitempty"`
	Type      string                 `json:"type,omitempty"`
	Status    string                 `json:"status,omitempty"`
	Title     string                 `json:"title,omitempty"`
	Space     *Space                 `json:"space,omitempty"`
	Body      map[string]interface{} `json:"body,omitempty"`
	Version   *Version               `json:"version,omitempty"`
	Ancestors []Content              `json:"ancestors,omitempty"`
	Links     map[string]string      `json:"_links,omitempty"`
}

type Version struct {
	By        *User  `json:"by,omitempty"`
	When      string `json:"when,omitempty"`
	Message   string `json:"message,omitempty"`
	Number    int    `json:"number,omitempty"`
	MinorEdit bool   `json:"minorEdit,omitempty"`
}

type ContentResult struct {
	Results []Content `json:"results"`
	Start   int       `json:"start"`
	Limit   int       `json:"limit"`
	Size    int       `json:"size"`
}

type ContentHistory struct {
	Latest   bool   `json:"latest"`
	CreatedBy *User `json:"createdBy,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
	LastUpdated *Version `json:"lastUpdated,omitempty"`
}

type ContentChildren struct {
	Page    *ContentResult `json:"page,omitempty"`
	Comment *ContentResult `json:"comment,omitempty"`
}

type ContentDescendants struct {
	Child  *ContentResult `json:"child,omitempty"`
	Descendant *ContentResult `json:"descendant,omitempty"`
}

func (s *ContentService) Get(ctx context.Context, contentID string, expand []string) (*Content, error) {
	if err := s.client.CheckEndpoint("content"); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/rest/api/content/%s", url.PathEscape(contentID))
	if len(expand) > 0 {
		path += "?expand=" + joinStrings(expand)
	}
	var content Content
	if err := s.client.Get(ctx, path, &content); err != nil {
		return nil, err
	}
	return &content, nil
}

func (s *ContentService) GetAll(ctx context.Context, contentType, spaceKey, title string, start, limit int, expand []string) (*ContentResult, error) {
	if err := s.client.CheckEndpoint("content"); err != nil {
		return nil, err
	}
	path := "/rest/api/content?"
	params := url.Values{}
	if contentType != "" {
		params.Set("type", contentType)
	}
	if spaceKey != "" {
		params.Set("spaceKey", spaceKey)
	}
	if title != "" {
		params.Set("title", title)
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
	var result ContentResult
	if err := s.client.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *ContentService) Create(ctx context.Context, content *Content) (*Content, error) {
	if err := s.client.CheckEndpoint("content"); err != nil {
		return nil, err
	}
	var result Content
	if err := s.client.Post(ctx, "/rest/api/content", content, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *ContentService) Update(ctx context.Context, contentID string, content *Content) (*Content, error) {
	if err := s.client.CheckEndpoint("content"); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/rest/api/content/%s", url.PathEscape(contentID))
	var result Content
	if err := s.client.Put(ctx, path, content, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *ContentService) Delete(ctx context.Context, contentID string) error {
	if err := s.client.CheckEndpoint("content"); err != nil {
		return err
	}
	path := fmt.Sprintf("/rest/api/content/%s", url.PathEscape(contentID))
	return s.client.Delete(ctx, path)
}

func (s *ContentService) GetHistory(ctx context.Context, contentID string, expand []string) (*ContentHistory, error) {
	if err := s.client.CheckEndpoint("content"); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/rest/api/content/%s/history", url.PathEscape(contentID))
	if len(expand) > 0 {
		path += "?expand=" + joinStrings(expand)
	}
	var history ContentHistory
	if err := s.client.Get(ctx, path, &history); err != nil {
		return nil, err
	}
	return &history, nil
}

func (s *ContentService) GetChildren(ctx context.Context, contentID string, expand []string) (*ContentResult, error) {
	if err := s.client.CheckEndpoint("content"); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/rest/api/content/%s/child/page", url.PathEscape(contentID))
	if len(expand) > 0 {
		path += "?expand=" + joinStrings(expand)
	}
	var children ContentResult
	if err := s.client.Get(ctx, path, &children); err != nil {
		return nil, err
	}
	return &children, nil
}

func (s *ContentService) GetDescendants(ctx context.Context, contentID string, expand []string) (*ContentResult, error) {
	if err := s.client.CheckEndpoint("content"); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/rest/api/content/%s/descendant", url.PathEscape(contentID))
	if len(expand) > 0 {
		path += "?expand=" + joinStrings(expand)
	}
	var descendants ContentResult
	if err := s.client.Get(ctx, path, &descendants); err != nil {
		return nil, err
	}
	return &descendants, nil
}

func (s *ContentService) GetVersions(ctx context.Context, contentID string, start, limit int) (*ContentResult, error) {
	if err := s.client.CheckEndpoint("content"); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/rest/api/content/%s/version", url.PathEscape(contentID))
	params := url.Values{}
	if start > 0 {
		params.Set("start", itoa(start))
	}
	if limit > 0 {
		params.Set("limit", itoa(limit))
	}
	if len(params) > 0 {
		path += "?" + params.Encode()
	}
	var result ContentResult
	if err := s.client.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *ContentService) GetComments(ctx context.Context, contentID string) (*ContentResult, error) {
	if err := s.client.CheckEndpoint("content"); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/rest/api/content/%s/child/comment", url.PathEscape(contentID))
	var result ContentResult
	if err := s.client.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *ContentService) GetAttachments(ctx context.Context, contentID string) (*ContentResult, error) {
	if err := s.client.CheckEndpoint("content"); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/rest/api/content/%s/child/attachment", url.PathEscape(contentID))
	var result ContentResult
	if err := s.client.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func joinStrings(strs []string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += ","
		}
		result += s
	}
	return result
}

func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}
