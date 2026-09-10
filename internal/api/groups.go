package api

import (
	"context"
	"fmt"
	"net/url"

	"github.com/darkkoax/confluencecli/internal/client"
)

type GroupService struct {
	client *client.Client
}

func NewGroupService(c *client.Client) *GroupService {
	return &GroupService{client: c}
}

type Group struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
	ID   string `json:"id,omitempty"`
}

type GroupMembers struct {
	Results []User `json:"results"`
	Start   int    `json:"start"`
	Limit   int    `json:"limit"`
	Size    int    `json:"size"`
}

func (s *GroupService) Get(ctx context.Context, groupName string, expand []string) (*Group, error) {
	if err := s.client.CheckEndpoint("groups"); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/rest/api/group/%s", url.PathEscape(groupName))
	if len(expand) > 0 {
		path += "?expand=" + joinStrings(expand)
	}
	var group Group
	if err := s.client.Get(ctx, path, &group); err != nil {
		return nil, err
	}
	return &group, nil
}

func (s *GroupService) GetMembers(ctx context.Context, groupName string, start, limit int) (*GroupMembers, error) {
	if err := s.client.CheckEndpoint("groups"); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/rest/api/group/%s/member?start=%d&limit=%d",
		url.PathEscape(groupName), start, limit)
	var members GroupMembers
	if err := s.client.Get(ctx, path, &members); err != nil {
		return nil, err
	}
	return &members, nil
}

func (s *GroupService) AddMember(ctx context.Context, groupName string, user *User) error {
	if err := s.client.CheckEndpoint("groups"); err != nil {
		return err
	}
	path := fmt.Sprintf("/rest/api/group/%s/member", url.PathEscape(groupName))
	return s.client.Post(ctx, path, user, nil)
}

func (s *GroupService) RemoveMember(ctx context.Context, groupName, username string) error {
	if err := s.client.CheckEndpoint("groups"); err != nil {
		return err
	}
	path := fmt.Sprintf("/rest/api/group/%s/member/%s",
		url.PathEscape(groupName), url.PathEscape(username))
	return s.client.Delete(ctx, path)
}
