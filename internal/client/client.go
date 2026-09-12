package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/darthkoax/confluencecli/internal/config"
)

var adminEndpoints = map[string]bool{
	"settings": true,
	"audit":    true,
	"system":   true,
}

var adminOperations = map[string]bool{}

type Client struct {
	httpClient *http.Client
	baseURL    string
	apiToken   string
	methods    config.MethodsConfig
	adminMode  bool
	endpoints  config.EndpointsConfig
}

func New(cfg *config.Config) (*Client, error) {
	tlsConfig := &tls.Config{}

	if cfg.Confluence.CustomCACert != "" {
		caCert, err := os.ReadFile(cfg.Confluence.CustomCACert)
		if err != nil {
			return nil, fmt.Errorf("reading custom CA cert: %w", err)
		}
		pool, err := x509.SystemCertPool()
		if err != nil {
			pool = x509.NewCertPool()
		}
		if !pool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to append custom CA cert")
		}
		tlsConfig.RootCAs = pool
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(cfg.Confluence.Timeout) * time.Second,
	}

	return &Client{
		httpClient: httpClient,
		baseURL:    strings.TrimRight(cfg.Confluence.BaseURL, "/"),
		apiToken:   cfg.Confluence.APIToken,
		methods:    cfg.Methods,
		adminMode:  cfg.AdminMode,
		endpoints:  cfg.Endpoints,
	}, nil
}

func (c *Client) CheckEndpoint(name string) error {
	switch name {
	case "content":
		if !c.endpoints.Content {
			return fmt.Errorf("content endpoint is disabled in configuration. This operation cannot be performed. Please ask the user to enable the content endpoint in config.toml or contact an administrator")
		}
	case "spaces":
		if !c.endpoints.Spaces {
			return fmt.Errorf("spaces endpoint is disabled in configuration. This operation cannot be performed. Please ask the user to enable the spaces endpoint in config.toml or contact an administrator")
		}
	case "search":
		if !c.endpoints.Search {
			return fmt.Errorf("search endpoint is disabled in configuration. This operation cannot be performed. Please ask the user to enable the search endpoint in config.toml or contact an administrator")
		}
	case "users":
		if !c.endpoints.Users {
			return fmt.Errorf("users endpoint is disabled in configuration. This operation cannot be performed. Please ask the user to enable the users endpoint in config.toml or contact an administrator")
		}
	case "groups":
		if !c.endpoints.Groups {
			return fmt.Errorf("groups endpoint is disabled in configuration. This operation cannot be performed. Please ask the user to enable the groups endpoint in config.toml or contact an administrator")
		}
	case "settings":
		if !c.endpoints.Settings {
			return fmt.Errorf("settings endpoint is disabled in configuration. This operation cannot be performed. Please ask the user to enable the settings endpoint in config.toml or contact an administrator")
		}
	case "audit":
		if !c.endpoints.Audit {
			return fmt.Errorf("audit endpoint is disabled in configuration. This operation cannot be performed. Please ask the user to enable the audit endpoint in config.toml or contact an administrator")
		}
	case "templates":
		if !c.endpoints.Templates {
			return fmt.Errorf("templates endpoint is disabled in configuration. This operation cannot be performed. Please ask the user to enable the templates endpoint in config.toml or contact an administrator")
		}
	case "contentstates":
		if !c.endpoints.ContentStates {
			return fmt.Errorf("contentstates endpoint is disabled in configuration. This operation cannot be performed. Please ask the user to enable the contentstates endpoint in config.toml or contact an administrator")
		}
	case "inlinetasks":
		if !c.endpoints.InlineTasks {
			return fmt.Errorf("inlinetasks endpoint is disabled in configuration. This operation cannot be performed. Please ask the user to enable the inlinetasks endpoint in config.toml or contact an administrator")
		}
	case "relations":
		if !c.endpoints.Relations {
			return fmt.Errorf("relations endpoint is disabled in configuration. This operation cannot be performed. Please ask the user to enable the relations endpoint in config.toml or contact an administrator")
		}
	case "longtasks":
		if !c.endpoints.LongTasks {
			return fmt.Errorf("longtasks endpoint is disabled in configuration. This operation cannot be performed. Please ask the user to enable the longtasks endpoint in config.toml or contact an administrator")
		}
	case "system":
		if !c.endpoints.System {
			return fmt.Errorf("system endpoint is disabled in configuration. This operation cannot be performed. Please ask the user to enable the system endpoint in config.toml or contact an administrator")
		}
	case "blueprints":
		if !c.endpoints.Blueprints {
			return fmt.Errorf("blueprints endpoint is disabled in configuration. This operation cannot be performed. Please ask the user to enable the blueprints endpoint in config.toml or contact an administrator")
		}
	case "healthcheck":
		if !c.endpoints.HealthCheck {
			return fmt.Errorf("healthcheck endpoint is disabled in configuration. This operation cannot be performed. Please ask the user to enable the healthcheck endpoint in config.toml or contact an administrator")
		}
	default:
		return fmt.Errorf("unknown endpoint: %s", name)
	}
	return nil
}

func (c *Client) CheckAdmin(endpoint, operation string) error {
	key := endpoint + "." + operation
	if adminEndpoints[endpoint] || adminOperations[key] {
		if !c.adminMode {
			return fmt.Errorf("%s.%s requires admin mode to be enabled. This operation cannot be performed. You should not attempt to bypass this restriction. Please ask the user to enable admin_mode in config.toml or contact an administrator", endpoint, operation)
		}
	}
	return nil
}

func (c *Client) Do(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	switch strings.ToUpper(method) {
	case http.MethodGet:
		if !c.methods.AllowGet {
			return fmt.Errorf("GET method is disabled in configuration. This operation cannot be performed. You should not attempt to bypass this restriction. Please ask the user to enable GET requests in config.toml or contact an administrator")
		}
	case http.MethodPost:
		if !c.methods.AllowPost {
			return fmt.Errorf("POST method is disabled in configuration. This operation cannot be performed. You should not attempt to bypass this restriction. Please ask the user to enable POST requests in config.toml or contact an administrator")
		}
	case http.MethodPut:
		if !c.methods.AllowPut {
			return fmt.Errorf("PUT method is disabled in configuration. This operation cannot be performed. You should not attempt to bypass this restriction. Please ask the user to enable PUT requests in config.toml or contact an administrator")
		}
	case http.MethodDelete:
		if !c.methods.AllowDelete {
			return fmt.Errorf("DELETE method is disabled in configuration. This operation cannot be performed. You should not attempt to bypass this restriction. Please ask the user to enable DELETE requests in config.toml or contact an administrator")
		}
	default:
		return fmt.Errorf("unsupported HTTP method: %s", method)
	}

	url := fmt.Sprintf("%s%s", c.baseURL, path)

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshaling request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("unmarshaling response: %w", err)
		}
	}

	return nil
}

func (c *Client) Get(ctx context.Context, path string, result interface{}) error {
	return c.Do(ctx, http.MethodGet, path, nil, result)
}

func (c *Client) Post(ctx context.Context, path string, body interface{}, result interface{}) error {
	return c.Do(ctx, http.MethodPost, path, body, result)
}

func (c *Client) Put(ctx context.Context, path string, body interface{}, result interface{}) error {
	return c.Do(ctx, http.MethodPut, path, body, result)
}

func (c *Client) Delete(ctx context.Context, path string) error {
	return c.Do(ctx, http.MethodDelete, path, nil, nil)
}
