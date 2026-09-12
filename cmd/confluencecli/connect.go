package main

import (
	"context"
	"fmt"
	"os"

	"github.com/darthkoax/confluencecli/internal/api"
	"github.com/darthkoax/confluencecli/internal/client"
	"github.com/darthkoax/confluencecli/internal/config"
)

func handleConnect(args []string) {
	if len(args) > 0 && (args[0] == "help" || args[0] == "--help" || args[0] == "-h") {
		printConnectHelp()
		return
	}

	var configPath string
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--config":
			if i+1 < len(args) {
				configPath = args[i+1]
				i++
			}
		}
	}

	if configPath == "" {
		configPath = config.DefaultConfigPath()
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	c, err := client.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
		os.Exit(1)
	}

	services := api.NewServices(c)
	ctx := context.Background()

	user, err := services.Users.GetCurrent(ctx, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to Confluence: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Connected to Confluence DC as: %s (%s)\n", user.DisplayName, user.Username)
	fmt.Printf("Confluence instance: %s\n", cfg.Confluence.BaseURL)
}

func printConnectHelp() {
	fmt.Print(`confluencecli connect - Connect to Confluence and verify authentication

USAGE:
  confluencecli connect [--config <path>]

DESCRIPTION:
  Connects to the Confluence Data Center instance using the configuration file and
  verifies authentication by fetching the current user's information.

  This command is useful for testing your configuration before using other
  commands. It will fail if:
  - The config file doesn't exist or is invalid
  - The base_url is unreachable
  - The API token is invalid or expired
  - Custom CA certificate is invalid (if specified)

FLAGS:
  --config <path>    Path to config file
                     Default: ~/.config/confluencecli/config.toml

EXAMPLES:
  confluencecli connect                              # Use default config location
  confluencecli connect --config /path/to/config     # Use custom config file

OUTPUT:
  On success:
    Connected to Confluence DC as: John Doe (john.doe)
    Confluence instance: https://confluence.example.com

  On failure:
    Error connecting to Confluence: <error message>
    Exit code: 1

COMMON ERRORS:
  - "reading config file: no such file or directory"
    -> Run 'confluencecli init' to create a config file

  - "confluence.base_url is required"
    -> Edit config.toml and set base_url

  - "confluence.api_token is required"
    -> Edit config.toml and set api_token

  - "API error (status 401)"
    -> API token is invalid or expired. Generate a new token in Confluence.

  - "x509: certificate signed by unknown authority"
    -> Set custom_ca_cert in config.toml to your CA certificate bundle

TROUBLESHOOTING:
  1. Verify config file exists: confluencecli init
  2. Check base_url is correct and accessible
  3. Verify API token has appropriate permissions
  4. If using custom CA, ensure cert path is correct
  5. Check network connectivity to Confluence instance
`)
}
