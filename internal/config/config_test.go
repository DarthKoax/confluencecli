package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	content := `admin_mode = true

[confluence]
base_url = "https://confluence.example.com"
api_token = "test-token"
custom_ca_cert = "/path/to/ca.crt"
timeout = 60

[methods]
allow_get = true
allow_post = false
allow_put = true
allow_delete = false

[endpoints]
content = true
spaces = false
`

	tmpfile, err := os.CreateTemp("", "config*.toml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	cfg, err := Load(tmpfile.Name())
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Confluence.BaseURL != "https://confluence.example.com" {
		t.Errorf("BaseURL = %v, want %v", cfg.Confluence.BaseURL, "https://confluence.example.com")
	}
	if cfg.Confluence.APIToken != "test-token" {
		t.Errorf("APIToken = %v, want %v", cfg.Confluence.APIToken, "test-token")
	}
	if cfg.Confluence.CustomCACert != "/path/to/ca.crt" {
		t.Errorf("CustomCACert = %v, want %v", cfg.Confluence.CustomCACert, "/path/to/ca.crt")
	}
	if cfg.Confluence.Timeout != 60 {
		t.Errorf("Timeout = %v, want %v", cfg.Confluence.Timeout, 60)
	}
	if !cfg.AdminMode {
		t.Error("AdminMode = false, want true")
	}
	if !cfg.Methods.AllowGet {
		t.Error("AllowGet = false, want true")
	}
	if cfg.Methods.AllowPost {
		t.Error("AllowPost = true, want false")
	}
	if !cfg.Methods.AllowPut {
		t.Error("AllowPut = false, want true")
	}
	if cfg.Methods.AllowDelete {
		t.Error("AllowDelete = true, want false")
	}
	if !cfg.Endpoints.Content {
		t.Error("Endpoints.Content = false, want true")
	}
	if cfg.Endpoints.Spaces {
		t.Error("Endpoints.Spaces = true, want false")
	}
}

func TestLoad_MissingRequired(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name: "missing base_url",
			content: `[confluence]
api_token = "test"
`,
		},
		{
			name: "missing api_token",
			content: `[confluence]
base_url = "https://confluence.example.com"
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpfile, err := os.CreateTemp("", "config*.toml")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpfile.Name())

			if _, err := tmpfile.Write([]byte(tt.content)); err != nil {
				t.Fatal(err)
			}
			tmpfile.Close()

			_, err = Load(tmpfile.Name())
			if err == nil {
				t.Error("Load() expected error, got nil")
			}
		})
	}
}

func TestLoad_Defaults(t *testing.T) {
	content := `[confluence]
base_url = "https://confluence.example.com"
api_token = "test-token"
`

	tmpfile, err := os.CreateTemp("", "config*.toml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	cfg, err := Load(tmpfile.Name())
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Confluence.Timeout != 30 {
		t.Errorf("Timeout = %v, want %v", cfg.Confluence.Timeout, 30)
	}

	if !cfg.Endpoints.Content {
		t.Error("Endpoints.Content = false, want true (default)")
	}
	if !cfg.Endpoints.Spaces {
		t.Error("Endpoints.Spaces = false, want true (default)")
	}
}

func TestLoad_EnvOverrides(t *testing.T) {
	content := `[confluence]
base_url = "https://confluence.example.com"
api_token = "test-token"
timeout = 30

[methods]
allow_get = true
allow_post = true

[endpoints]
content = true
spaces = true
`

	tmpfile, err := os.CreateTemp("", "config*.toml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	os.Setenv("CONFLUENCE_BASE_URL", "https://override.example.com")
	os.Setenv("CONFLUENCE_API_TOKEN", "override-token")
	os.Setenv("CONFLUENCE_TIMEOUT", "90")
	os.Setenv("CONFLUENCE_ADMIN_MODE", "true")
	os.Setenv("CONFLUENCE_ALLOW_POST", "false")
	os.Setenv("CONFLUENCE_ENDPOINT_SPACES", "false")
	defer func() {
		os.Unsetenv("CONFLUENCE_BASE_URL")
		os.Unsetenv("CONFLUENCE_API_TOKEN")
		os.Unsetenv("CONFLUENCE_TIMEOUT")
		os.Unsetenv("CONFLUENCE_ADMIN_MODE")
		os.Unsetenv("CONFLUENCE_ALLOW_POST")
		os.Unsetenv("CONFLUENCE_ENDPOINT_SPACES")
	}()

	cfg, err := Load(tmpfile.Name())
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Confluence.BaseURL != "https://override.example.com" {
		t.Errorf("BaseURL = %v, want %v", cfg.Confluence.BaseURL, "https://override.example.com")
	}
	if cfg.Confluence.APIToken != "override-token" {
		t.Errorf("APIToken = %v, want %v", cfg.Confluence.APIToken, "override-token")
	}
	if cfg.Confluence.Timeout != 90 {
		t.Errorf("Timeout = %v, want %v", cfg.Confluence.Timeout, 90)
	}
	if !cfg.AdminMode {
		t.Error("AdminMode = false, want true")
	}
	if cfg.Methods.AllowPost {
		t.Error("AllowPost = true, want false")
	}
	if cfg.Endpoints.Spaces {
		t.Error("Endpoints.Spaces = true, want false")
	}
}

func TestLoadFromEnv(t *testing.T) {
	os.Setenv("CONFLUENCE_BASE_URL", "https://env.example.com")
	os.Setenv("CONFLUENCE_API_TOKEN", "env-token")
	os.Setenv("CONFLUENCE_TIMEOUT", "45")
	os.Setenv("CONFLUENCE_ADMIN_MODE", "1")
	os.Setenv("CONFLUENCE_ENDPOINT_CONTENT", "false")
	defer func() {
		os.Unsetenv("CONFLUENCE_BASE_URL")
		os.Unsetenv("CONFLUENCE_API_TOKEN")
		os.Unsetenv("CONFLUENCE_TIMEOUT")
		os.Unsetenv("CONFLUENCE_ADMIN_MODE")
		os.Unsetenv("CONFLUENCE_ENDPOINT_CONTENT")
	}()

	cfg, err := LoadFromEnv()
	if err != nil {
		t.Fatalf("LoadFromEnv() error = %v", err)
	}

	if cfg.Confluence.BaseURL != "https://env.example.com" {
		t.Errorf("BaseURL = %v, want %v", cfg.Confluence.BaseURL, "https://env.example.com")
	}
	if cfg.Confluence.APIToken != "env-token" {
		t.Errorf("APIToken = %v, want %v", cfg.Confluence.APIToken, "env-token")
	}
	if cfg.Confluence.Timeout != 45 {
		t.Errorf("Timeout = %v, want %v", cfg.Confluence.Timeout, 45)
	}
	if !cfg.AdminMode {
		t.Error("AdminMode = false, want true")
	}
	if cfg.Endpoints.Content {
		t.Error("Endpoints.Content = true, want false")
	}
	if !cfg.Endpoints.Spaces {
		t.Error("Endpoints.Spaces = false, want true (default)")
	}
}

func TestLoadFromEnv_MissingRequired(t *testing.T) {
	os.Unsetenv("CONFLUENCE_BASE_URL")
	os.Unsetenv("CONFLUENCE_API_TOKEN")

	_, err := LoadFromEnv()
	if err == nil {
		t.Error("LoadFromEnv() expected error, got nil")
	}

	os.Setenv("CONFLUENCE_BASE_URL", "https://example.com")
	defer os.Unsetenv("CONFLUENCE_BASE_URL")

	_, err = LoadFromEnv()
	if err == nil {
		t.Error("LoadFromEnv() expected error for missing API token, got nil")
	}
}

func TestDefaultConfigDir(t *testing.T) {
	os.Unsetenv("CONFLUENCE_CONFIG_DIR")
	dir := DefaultConfigDir()
	if dir == "" {
		t.Error("DefaultConfigDir() returned empty string")
	}

	os.Setenv("CONFLUENCE_CONFIG_DIR", "/custom/path")
	defer os.Unsetenv("CONFLUENCE_CONFIG_DIR")

	dir = DefaultConfigDir()
	if dir != "/custom/path" {
		t.Errorf("DefaultConfigDir() = %v, want /custom/path", dir)
	}
}

func TestDefaultConfigPath(t *testing.T) {
	os.Unsetenv("CONFLUENCE_CONFIG_DIR")
	path := DefaultConfigPath()
	if path == "" {
		t.Error("DefaultConfigPath() returned empty string")
	}
	if len(path) < len("config.toml") {
		t.Errorf("DefaultConfigPath() = %v, too short", path)
	}
}

func TestDefaultConfigContent(t *testing.T) {
	content := DefaultConfigContent()
	if content == "" {
		t.Error("DefaultConfigContent() returned empty string")
	}
	if len(content) < 100 {
		t.Errorf("DefaultConfigContent() = %v chars, expected more", len(content))
	}
}

func TestInit(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "confluencecli-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	path, err := Init(tmpdir)
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if path == "" {
		t.Error("Init() returned empty path")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("Init() did not create file at %s", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read created config: %v", err)
	}
	if len(data) == 0 {
		t.Error("Created config file is empty")
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("Config file permissions = %v, want 0600", info.Mode().Perm())
	}
}

func TestInit_AlreadyExists(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "confluencecli-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	_, err = Init(tmpdir)
	if err != nil {
		t.Fatalf("First Init() error = %v", err)
	}

	_, err = Init(tmpdir)
	if err == nil {
		t.Error("Second Init() expected error for existing file, got nil")
	}
}

func TestInit_DefaultDir(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "confluencecli-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	os.Setenv("CONFLUENCE_CONFIG_DIR", tmpdir)
	defer os.Unsetenv("CONFLUENCE_CONFIG_DIR")

	path, err := Init("")
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if path == "" {
		t.Error("Init() returned empty path")
	}
}
