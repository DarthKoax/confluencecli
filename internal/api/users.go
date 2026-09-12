package api

import (
	"context"
	"fmt"
	"net/url"

	"github.com/darthkoax/confluencecli/internal/client"
)

type UserService struct {
	client *client.Client
}

func NewUserService(c *client.Client) *UserService {
	return &UserService{client: c}
}

type User struct {
	Type         string            `json:"type,omitempty"`
	Username     string            `json:"username,omitempty"`
	UserKey      string            `json:"userKey,omitempty"`
	AccountID    string            `json:"accountId,omitempty"`
	DisplayName  string            `json:"displayName,omitempty"`
	Email        string            `json:"email,omitempty"`
	ProfilePicture map[string]interface{} `json:"profilePicture,omitempty"`
	IsExternalCollaborator bool   `json:"isExternalCollaborator,omitempty"`
	Links        map[string]string `json:"_links,omitempty"`
}

func (s *UserService) GetCurrent(ctx context.Context, expand []string) (*User, error) {
	if err := s.client.CheckEndpoint("users"); err != nil {
		return nil, err
	}
	path := "/rest/api/user/current"
	if len(expand) > 0 {
		path += "?expand=" + joinStrings(expand)
	}
	var user User
	if err := s.client.Get(ctx, path, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetAnonymous(ctx context.Context, expand []string) (*User, error) {
	if err := s.client.CheckEndpoint("users"); err != nil {
		return nil, err
	}
	path := "/rest/api/user/anonymous"
	if len(expand) > 0 {
		path += "?expand=" + joinStrings(expand)
	}
	var user User
	if err := s.client.Get(ctx, path, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) Get(ctx context.Context, username, key, accountID string, expand []string) (*User, error) {
	if err := s.client.CheckEndpoint("users"); err != nil {
		return nil, err
	}
	path := "/rest/api/user?"
	params := url.Values{}
	if username != "" {
		params.Set("username", username)
	}
	if key != "" {
		params.Set("key", key)
	}
	if accountID != "" {
		params.Set("accountId", accountID)
	}
	if len(expand) > 0 {
		params.Set("expand", joinStrings(expand))
	}
	path += params.Encode()
	var user User
	if err := s.client.Get(ctx, path, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetUnknown(ctx context.Context, expand []string) (*User, error) {
	if err := s.client.CheckEndpoint("users"); err != nil {
		return nil, err
	}
	path := "/rest/api/user/unknown"
	if len(expand) > 0 {
		path += "?expand=" + joinStrings(expand)
	}
	var user User
	if err := s.client.Get(ctx, path, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string, expand []string) (*User, error) {
	if err := s.client.CheckEndpoint("users"); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/rest/api/user/email?email=%s", url.QueryEscape(email))
	if len(expand) > 0 {
		path += "&expand=" + joinStrings(expand)
	}
	var user User
	if err := s.client.Get(ctx, path, &user); err != nil {
		return nil, err
	}
	return &user, nil
}
