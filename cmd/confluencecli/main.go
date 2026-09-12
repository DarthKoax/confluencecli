package main

import (
	"fmt"
	"os"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		printRootHelp()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "help", "--help", "-h":
		if len(os.Args) > 2 {
			printCommandHelp(os.Args[2])
		} else {
			printRootHelp()
		}
	case "version", "--version", "-v":
		printVersion()
	case "init":
		handleInit(os.Args[2:])
	case "connect":
		handleConnect(os.Args[2:])
	case "content":
		handleContent(os.Args[2:])
	case "space":
		handleSpace(os.Args[2:])
	case "search":
		handleSearch(os.Args[2:])
	case "user":
		handleUser(os.Args[2:])
	case "group":
		handleGroup(os.Args[2:])
	case "settings":
		handleSettings(os.Args[2:])
	case "audit":
		handleAudit(os.Args[2:])
	case "template":
		handleTemplate(os.Args[2:])
	case "contentstate":
		handleContentState(os.Args[2:])
	case "inlinetask":
		handleInlineTask(os.Args[2:])
	case "relation":
		handleRelation(os.Args[2:])
	case "longtask":
		handleLongTask(os.Args[2:])
	case "system":
		handleSystem(os.Args[2:])
	case "blueprint":
		handleBlueprint(os.Args[2:])
	case "healthcheck":
		handleHealthCheck(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		printRootHelp()
		os.Exit(1)
	}
}

func printVersion() {
	fmt.Printf("confluencecli version %s\n", version)
	fmt.Printf("GitHub: https://github.com/darthkoax/confluencecli\n")
}

func printRootHelp() {
	fmt.Printf(`confluencecli %s - Confluence Data Center / Server CLI

USAGE:
  confluencecli <command> [flags]
  confluencecli <command> <action> [flags]
  confluencecli help <command>

COMMANDS:
  Setup & Connection:
    init              Create a default configuration file
    connect           Connect to Confluence and verify authentication

  Content Management:
    content           Manage pages, blog posts, and comments
    space             Manage spaces
    search            Search using CQL (Confluence Query Language)
    template          Manage content and blueprint templates
    contentstate      Manage content states
    blueprint         Manage content blueprints
    inlinetask        Manage inline tasks

  User & Group Management:
    user              Manage users
    group             Manage groups and memberships

  Relations:
    relation          Manage content relations

  System & Admin:
    settings          Manage system settings (admin)
    audit             View audit logs (admin)
    system            View system status (admin)
    longtask          Monitor long-running tasks
    healthcheck       Check system health

  Other:
    version           Print version information
    help              Show help for a command

EXAMPLES:
  confluencecli init                                          # Initialize config
  confluencecli connect                                       # Test connection
  confluencecli content get 12345                             # Get page by ID
  confluencecli content create --json '{"type":"page",...}'   # Create page
  confluencecli search cql --cql "type=page AND space=DEV"    # Search content
  confluencecli help content                                  # Get help for content

For more information about a command, run:
  confluencecli help <command>
`, version)
}


