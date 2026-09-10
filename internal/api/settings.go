package api

import (
	"context"
	"fmt"
	"net/url"

	"github.com/darkkoax/confluencecli/internal/client"
)

type SettingsService struct {
	client *client.Client
}

func NewSettingsService(c *client.Client) *SettingsService {
	return &SettingsService{client: c}
}

type SystemInfo struct {
	CloudID  string `json:"cloudId,omitempty"`
	Commit   string `json:"commitHash,omitempty"`
	BuildDate string `json:"buildDate,omitempty"`
}

type Theme struct {
	ThemeKey    string `json:"themeKey,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Icon        map[string]string `json:"icon,omitempty"`
}

type LookAndFeel struct {
	Head     map[string]interface{} `json:"head,omitempty"`
	Header   map[string]interface{} `json:"header,omitempty"`
	Space    map[string]interface{} `json:"space,omitempty"`
	Content  map[string]interface{} `json:"content,omitempty"`
	Custom   map[string]interface{} `json:"custom,omitempty"`
	Breadcrumbs map[string]interface{} `json:"breadcrumbs,omitempty"`
	Menu     map[string]interface{} `json:"menu,omitempty"`
	ColorScheme map[string]interface{} `json:"colourScheme,omitempty"`
}

func (s *SettingsService) GetSystemInfo(ctx context.Context) (*SystemInfo, error) {
	if err := s.client.CheckEndpoint("settings"); err != nil {
		return nil, err
	}
	if err := s.client.CheckAdmin("settings", "system_info"); err != nil {
		return nil, err
	}
	var info SystemInfo
	if err := s.client.Get(ctx, "/rest/api/settings/systemInfo", &info); err != nil {
		return nil, err
	}
	return &info, nil
}

func (s *SettingsService) GetTheme(ctx context.Context) (*Theme, error) {
	if err := s.client.CheckEndpoint("settings"); err != nil {
		return nil, err
	}
	if err := s.client.CheckAdmin("settings", "theme"); err != nil {
		return nil, err
	}
	var theme Theme
	if err := s.client.Get(ctx, "/rest/api/settings/theme", &theme); err != nil {
		return nil, err
	}
	return &theme, nil
}

func (s *SettingsService) SetTheme(ctx context.Context, themeKey string) error {
	if err := s.client.CheckEndpoint("settings"); err != nil {
		return err
	}
	if err := s.client.CheckAdmin("settings", "theme"); err != nil {
		return err
	}
	body := map[string]string{"themeKey": themeKey}
	return s.client.Put(ctx, "/rest/api/settings/theme", body, nil)
}

func (s *SettingsService) GetLookAndFeel(ctx context.Context, spaceKey string) (*LookAndFeel, error) {
	if err := s.client.CheckEndpoint("settings"); err != nil {
		return nil, err
	}
	if err := s.client.CheckAdmin("settings", "lookandfeel"); err != nil {
		return nil, err
	}
	path := "/rest/api/settings/lookandfeel"
	if spaceKey != "" {
		path = fmt.Sprintf("/rest/api/settings/lookandfeel?spaceKey=%s", url.QueryEscape(spaceKey))
	}
	var laf LookAndFeel
	if err := s.client.Get(ctx, path, &laf); err != nil {
		return nil, err
	}
	return &laf, nil
}

func (s *SettingsService) SetLookAndFeel(ctx context.Context, laf *LookAndFeel) error {
	if err := s.client.CheckEndpoint("settings"); err != nil {
		return err
	}
	if err := s.client.CheckAdmin("settings", "lookandfeel"); err != nil {
		return err
	}
	return s.client.Put(ctx, "/rest/api/settings/lookandfeel", laf, nil)
}

func (s *SettingsService) ResetLookAndFeel(ctx context.Context, spaceKey string) error {
	if err := s.client.CheckEndpoint("settings"); err != nil {
		return err
	}
	if err := s.client.CheckAdmin("settings", "lookandfeel"); err != nil {
		return err
	}
	path := "/rest/api/settings/lookandfeel"
	if spaceKey != "" {
		path = fmt.Sprintf("/rest/api/settings/lookandfeel?spaceKey=%s", url.QueryEscape(spaceKey))
	}
	return s.client.Delete(ctx, path)
}
