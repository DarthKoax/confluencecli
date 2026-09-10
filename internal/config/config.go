package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Confluence ConfluenceConfig  `toml:"confluence"`
	Methods    MethodsConfig     `toml:"methods"`
	AdminMode  bool              `toml:"admin_mode"`
	Endpoints  EndpointsConfig   `toml:"endpoints"`
}

type ConfluenceConfig struct {
	BaseURL      string `toml:"base_url"`
	APIToken     string `toml:"api_token"`
	CustomCACert string `toml:"custom_ca_cert"`
	Timeout      int    `toml:"timeout"`
}

type MethodsConfig struct {
	AllowGet    bool `toml:"allow_get"`
	AllowPost   bool `toml:"allow_post"`
	AllowPut    bool `toml:"allow_put"`
	AllowDelete bool `toml:"allow_delete"`
}

type EndpointsConfig struct {
	Content      bool `toml:"content"`
	Spaces       bool `toml:"spaces"`
	Search       bool `toml:"search"`
	Users        bool `toml:"users"`
	Groups       bool `toml:"groups"`
	Settings     bool `toml:"settings"`
	Audit        bool `toml:"audit"`
	Templates    bool `toml:"templates"`
	ContentStates bool `toml:"contentstates"`
	InlineTasks  bool `toml:"inlinetasks"`
	Relations    bool `toml:"relations"`
	LongTasks    bool `toml:"longtasks"`
	System       bool `toml:"system"`
	Blueprints   bool `toml:"blueprints"`
	HealthCheck  bool `toml:"healthcheck"`
}

var allEndpointKeys = []string{
	"content", "spaces", "search", "users", "groups",
	"settings", "audit", "templates", "contentstates",
	"inlinetasks", "relations", "longtasks", "system",
	"blueprints", "healthcheck",
}

func DefaultEndpoints() EndpointsConfig {
	return EndpointsConfig{
		Content:       true,
		Spaces:        true,
		Search:        true,
		Users:         true,
		Groups:        true,
		Settings:      true,
		Audit:         true,
		Templates:     true,
		ContentStates: true,
		InlineTasks:   true,
		Relations:     true,
		LongTasks:     true,
		System:        true,
		Blueprints:    true,
		HealthCheck:   true,
	}
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var raw map[string]interface{}
	if _, err := toml.Decode(string(data), &raw); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if cfg.Confluence.BaseURL == "" {
		return nil, fmt.Errorf("confluence.base_url is required")
	}
	if cfg.Confluence.APIToken == "" {
		return nil, fmt.Errorf("confluence.api_token is required")
	}
	if cfg.Confluence.Timeout == 0 {
		cfg.Confluence.Timeout = 30
	}

	endpointsSection, hasEndpoints := raw["endpoints"].(map[string]interface{})
	if !hasEndpoints {
		cfg.Endpoints = DefaultEndpoints()
	} else {
		cfg.Endpoints = DefaultEndpoints()
		for _, key := range allEndpointKeys {
			if val, ok := endpointsSection[key]; ok {
				if b, ok := val.(bool); ok {
					setEndpoint(&cfg.Endpoints, key, b)
				}
			}
		}
	}

	applyEnvOverrides(&cfg)

	return &cfg, nil
}

func LoadFromEnv() (*Config, error) {
	baseURL := os.Getenv("CONFLUENCE_BASE_URL")
	apiToken := os.Getenv("CONFLUENCE_API_TOKEN")

	if baseURL == "" {
		return nil, fmt.Errorf("CONFLUENCE_BASE_URL environment variable is required")
	}
	if apiToken == "" {
		return nil, fmt.Errorf("CONFLUENCE_API_TOKEN environment variable is required")
	}

	cfg := &Config{
		Confluence: ConfluenceConfig{
			BaseURL:      baseURL,
			APIToken:     apiToken,
			CustomCACert: os.Getenv("CONFLUENCE_CUSTOM_CA_CERT"),
			Timeout:      30,
		},
		Methods: MethodsConfig{
			AllowGet:    true,
			AllowPost:   true,
			AllowPut:    true,
			AllowDelete: true,
		},
		AdminMode: false,
		Endpoints: DefaultEndpoints(),
	}

	if timeout := os.Getenv("CONFLUENCE_TIMEOUT"); timeout != "" {
		if t, err := strconv.Atoi(timeout); err == nil {
			cfg.Confluence.Timeout = t
		}
	}

	applyEnvOverrides(cfg)

	return cfg, nil
}

func setEndpoint(e *EndpointsConfig, key string, val bool) {
	switch key {
	case "content":
		e.Content = val
	case "spaces":
		e.Spaces = val
	case "search":
		e.Search = val
	case "users":
		e.Users = val
	case "groups":
		e.Groups = val
	case "settings":
		e.Settings = val
	case "audit":
		e.Audit = val
	case "templates":
		e.Templates = val
	case "contentstates":
		e.ContentStates = val
	case "inlinetasks":
		e.InlineTasks = val
	case "relations":
		e.Relations = val
	case "longtasks":
		e.LongTasks = val
	case "system":
		e.System = val
	case "blueprints":
		e.Blueprints = val
	case "healthcheck":
		e.HealthCheck = val
	}
}

func applyEnvOverrides(cfg *Config) {
	if v := os.Getenv("CONFLUENCE_BASE_URL"); v != "" {
		cfg.Confluence.BaseURL = v
	}
	if v := os.Getenv("CONFLUENCE_API_TOKEN"); v != "" {
		cfg.Confluence.APIToken = v
	}
	if v := os.Getenv("CONFLUENCE_CUSTOM_CA_CERT"); v != "" {
		cfg.Confluence.CustomCACert = v
	}
	if v := os.Getenv("CONFLUENCE_TIMEOUT"); v != "" {
		if t, err := strconv.Atoi(v); err == nil {
			cfg.Confluence.Timeout = t
		}
	}
	if v := os.Getenv("CONFLUENCE_ADMIN_MODE"); v != "" {
		cfg.AdminMode = v == "true" || v == "1"
	}
	if v := os.Getenv("CONFLUENCE_ALLOW_GET"); v != "" {
		cfg.Methods.AllowGet = v == "true" || v == "1"
	}
	if v := os.Getenv("CONFLUENCE_ALLOW_POST"); v != "" {
		cfg.Methods.AllowPost = v == "true" || v == "1"
	}
	if v := os.Getenv("CONFLUENCE_ALLOW_PUT"); v != "" {
		cfg.Methods.AllowPut = v == "true" || v == "1"
	}
	if v := os.Getenv("CONFLUENCE_ALLOW_DELETE"); v != "" {
		cfg.Methods.AllowDelete = v == "true" || v == "1"
	}

	for _, key := range allEndpointKeys {
		envKey := "CONFLUENCE_ENDPOINT_" + strings.ToUpper(key)
		if v := os.Getenv(envKey); v != "" {
			setEndpoint(&cfg.Endpoints, key, v == "true" || v == "1")
		}
	}
}

func DefaultConfigDir() string {
	if v := os.Getenv("CONFLUENCE_CONFIG_DIR"); v != "" {
		return v
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ".config/darthkoax/confluencecli"
	}
	return filepath.Join(home, ".config", "darthkoax", "confluencecli")
}

func DefaultConfigPath() string {
	return filepath.Join(DefaultConfigDir(), "config.toml")
}

func DefaultConfigContent() string {
	return `[confluence]
base_url = "https://confluence.example.com"
api_token = "YOUR_API_TOKEN_HERE"
custom_ca_cert = ""
timeout = 30

admin_mode = false

[methods]
allow_get = true
allow_post = true
allow_put = true
allow_delete = true

[endpoints]
content = true
spaces = true
search = true
users = true
groups = true
settings = true
audit = true
templates = true
contentstates = true
inlinetasks = true
relations = true
longtasks = true
system = true
blueprints = true
healthcheck = true
`
}

func Init(dir string) (string, error) {
	if dir == "" {
		dir = DefaultConfigDir()
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("creating config directory: %w", err)
	}

	configPath := filepath.Join(dir, "config.toml")

	if _, err := os.Stat(configPath); err == nil {
		return configPath, fmt.Errorf("config file already exists at %s", configPath)
	}

	if err := os.WriteFile(configPath, []byte(DefaultConfigContent()), 0600); err != nil {
		return "", fmt.Errorf("writing config file: %w", err)
	}

	return configPath, nil
}
