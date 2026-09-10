package main

import (
	"fmt"
	"os"

	"github.com/darkkoax/confluencecli/internal/config"
)

func handleInit(args []string) {
	if len(args) > 0 && (args[0] == "help" || args[0] == "--help" || args[0] == "-h") {
		printInitHelp()
		return
	}

	var dir string
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--dir":
			if i+1 < len(args) {
				dir = args[i+1]
				i++
			}
		}
	}

	path, err := config.Init(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Configuration file created at: %s\n", path)
	fmt.Println("Edit the file to set your Confluence base_url and api_token.")
}

func printInitHelp() {
	fmt.Print(`confluencecli init - Create a default configuration file

USAGE:
  confluencecli init [--dir <path>]

DESCRIPTION:
  Creates a default configuration file at ~/.config/darthkoax/confluencecli/config.toml
  with all endpoints and methods enabled. The config file contains placeholders
  for your Confluence instance URL and API token.

  The default config directory can be overridden using the CONFLUENCE_CONFIG_DIR
  environment variable.

FLAGS:
  --dir <path>    Directory to create config in
                  Default: ~/.config/darthkoax/confluencecli
                  Can also be set via CONFLUENCE_CONFIG_DIR environment variable

EXAMPLES:
  confluencecli init                          # Create config in default location
  confluencecli init --dir /custom/path       # Create config in custom directory

OUTPUT:
  On success, prints the path to the created config file.
  On failure, prints error message and exits with code 1.

CONFIGURATION FILE STRUCTURE:
  [confluence]
  base_url = "https://confluence.example.com"  # Your Confluence DC URL
  api_token = "YOUR_API_TOKEN_HERE"            # Your API token
  custom_ca_cert = ""                          # Path to custom CA cert (optional)
  timeout = 30                                 # HTTP timeout in seconds

  admin_mode = false                           # Enable admin endpoints

  [methods]
  allow_get = true                             # Allow GET requests
  allow_post = true                            # Allow POST requests
  allow_put = true                             # Allow PUT requests
  allow_delete = true                          # Allow DELETE requests

  [endpoints]
  content = true                               # Enable content endpoint
  spaces = true                                # Enable spaces endpoint
  # ... (all 15 endpoints listed)

NEXT STEPS:
  1. Edit the config file with your Confluence instance details
  2. Run 'confluencecli connect' to verify the connection
  3. Start using the CLI to interact with Confluence
`)
}
