package main

import (
	"fmt"
	"os"
)

func printContentHelp() {
	fmt.Print(`confluencecli content - Manage pages, blog posts, and comments

USAGE:
  confluencecli content <action> [flags]

ACTIONS:
  get <id>              Get content by ID
  list                  List content (pages, blog posts)
  create                Create new content
  update <id>           Update existing content
  delete <id>           Delete content
  history <id>          Get content version history
  children <id>         Get child content
  descendants <id>      Get descendant content
  versions <id>         Get content versions
  comments <id>         Get content comments
  attachments <id>      Get content attachments

FLAGS:
  --config <path>       Path to config file
  --expand <fields>     Comma-separated list of fields to expand
  --type <type>         Content type (page, blogpost, comment)
  --space-key <key>     Space key to filter by
  --title <title>       Title to filter by
  --start <n>           Start index for pagination
  --limit <n>           Maximum results to return
  --json <payload>      JSON payload for create/update

EXAMPLES:
  confluencecli content get 12345
  confluencecli content list --type page --space-key DEV
  confluencecli content create --json '{"type":"page","title":"New Page","space":{"key":"DEV"}}'
  confluencecli content update 12345 --json '{"version":{"number":2},"title":"Updated"}'
  confluencecli content delete 12345
  confluencecli content history 12345
  confluencecli content children 12345 --expand page.body.view
`)
}

func printSpaceHelp() {
	fmt.Print(`confluencecli space - Manage spaces

USAGE:
  confluencecli space <action> [flags]

ACTIONS:
  list                  List all spaces
  get <key>             Get space by key
  create                Create a new space (admin)
  delete <key>          Delete a space (admin)
  content <key>         Get space content

FLAGS:
  --config <path>       Path to config file
  --expand <fields>     Comma-separated list of fields to expand
  --type <type>         Space type (global, personal)
  --status <status>     Space status (current, archived)
  --start <n>           Start index for pagination
  --limit <n>           Maximum results to return
  --json <payload>      JSON payload for create

EXAMPLES:
  confluencecli space list
  confluencecli space get DEV
  confluencecli space create --json '{"key":"NEW","name":"New Space","type":"global"}'
  confluencecli space delete DEV
  confluencecli space content DEV --type page
`)
}

func printSearchHelp() {
	fmt.Print(`confluencecli search - Search using CQL (Confluence Query Language)

USAGE:
  confluencecli search <action> [flags]

ACTIONS:
  cql                   Search content using CQL
  user                  Search users using CQL

FLAGS:
  --config <path>       Path to config file
  --cql <query>         CQL search query
  --expand <fields>     Comma-separated list of fields to expand
  --start <n>           Start index for pagination
  --limit <n>           Maximum results to return

EXAMPLES:
  confluencecli search cql --cql "type=page AND space=DEV"
  confluencecli search cql --cql "title~'Meeting Notes'" --limit 10
  confluencecli search user --cql "user.fullname~'John'"
`)
}

func printUserHelp() {
	fmt.Print(`confluencecli user - Manage users

USAGE:
  confluencecli user <action> [flags]

ACTIONS:
  current               Get current authenticated user
  anonymous             Get anonymous user
  get                   Get user by username, key, or account ID
  unknown               Get unknown user
  email                 Get user by email

FLAGS:
  --config <path>       Path to config file
  --expand <fields>     Comma-separated list of fields to expand
  --username <name>     Username to look up
  --key <key>           User key to look up
  --account-id <id>     Account ID to look up
  --email <email>       Email address to look up

EXAMPLES:
  confluencecli user current
  confluencecli user anonymous
  confluencecli user get --username john.doe
  confluencecli user get --key abc123
  confluencecli user unknown
  confluencecli user email --email john@example.com
`)
}

func printGroupHelp() {
	fmt.Print(`confluencecli group - Manage groups and memberships

USAGE:
  confluencecli group <action> [flags]

ACTIONS:
  get                   Get group by name
  members               List group members
  add-member            Add user to group (admin)
  remove-member         Remove user from group (admin)

FLAGS:
  --config <path>       Path to config file
  --expand <fields>     Comma-separated list of fields to expand
  --groupname <name>    Group name
  --username <name>     Username for add/remove member
  --start <n>           Start index for pagination
  --limit <n>           Maximum results to return

EXAMPLES:
  confluencecli group get --groupname developers
  confluencecli group members --groupname developers
  confluencecli group add-member --groupname developers --username john.doe
  confluencecli group remove-member --groupname developers --username john.doe
`)
}

func printSettingsHelp() {
	fmt.Print(`confluencecli settings - Manage system settings (admin)

USAGE:
  confluencecli settings <action> [flags]

ACTIONS:
  systeminfo            Get system information
  theme-get             Get current theme
  theme-set             Set theme
  laf-get               Get look and feel settings
  laf-set               Set look and feel settings
  laf-reset             Reset look and feel to defaults

FLAGS:
  --config <path>       Path to config file
  --theme-key <key>     Theme key for theme-set
  --space-key <key>     Space key for space-specific settings
  --json <payload>      JSON payload for laf-set

EXAMPLES:
  confluencecli settings systeminfo
  confluencecli settings theme-get
  confluencecli settings theme-set --theme-key com.atlassian.confluence.theme:mytheme
  confluencecli settings laf-get
  confluencecli settings laf-get --space-key DEV
  confluencecli settings laf-set --json '{"head":{"background":"#ffffff"}}'
  confluencecli settings laf-reset
`)
}

func printAuditHelp() {
	fmt.Print(`confluencecli audit - View audit logs (admin)

USAGE:
  confluencecli audit <action> [flags]

ACTIONS:
  get                   Get audit records
  since                 Get audit records since a time period
  retention             Get or set retention period
  export                Export audit records

FLAGS:
  --config <path>       Path to config file
  --start-date <date>   Start date (ISO 8601)
  --end-date <date>     End date (ISO 8601)
  --search <string>     Search string to filter records
  --number <value>      Number of units (for since: e.g., 1m, 7d)
  --units <unit>        Units for retention (DAYS, MONTHS)
  --start <n>           Start index for pagination
  --limit <n>           Maximum results to return

EXAMPLES:
  confluencecli audit get --start-date 2024-01-01 --end-date 2024-12-31
  confluencecli audit since --number 7d
  confluencecli audit retention
  confluencecli audit retention set --number 365 --units DAYS
  confluencecli audit export --start-date 2024-01-01 --end-date 2024-12-31
`)
}

func printTemplateHelp() {
	fmt.Print(`confluencecli template - Manage content and blueprint templates

USAGE:
  confluencecli template <action> [flags]

ACTIONS:
  list-content          List content templates
  get-content <id>      Get content template by ID
  create-content        Create content template
  update-content        Update content template
  delete-content <id>   Delete content template
  list-blueprint        List blueprint templates
  get-blueprint <id>    Get blueprint template by ID

FLAGS:
  --config <path>       Path to config file
  --expand <fields>     Comma-separated list of fields to expand
  --space-key <key>     Space key to filter by
  --start <n>           Start index for pagination
  --limit <n>           Maximum results to return
  --json <payload>      JSON payload for create/update

EXAMPLES:
  confluencecli template list-content --space-key DEV
  confluencecli template get-content tmpl-123
  confluencecli template create-content --json '{"name":"Meeting Notes","templateType":"page"}'
  confluencecli template delete-content tmpl-123
  confluencecli template list-blueprint
  confluencecli template get-blueprint bp-456
`)
}

func printContentStateHelp() {
	fmt.Print(`confluencecli contentstate - Manage content states

USAGE:
  confluencecli contentstate <action> [flags]

ACTIONS:
  list                  List all content states
  get <id>              Get content state by ID
  create                Create content state
  update <id>           Update content state
  delete <id>           Delete content state
  content <id>          Get content with this state

FLAGS:
  --config <path>       Path to config file
  --start <n>           Start index for pagination
  --limit <n>           Maximum results to return
  --json <payload>      JSON payload for create/update

EXAMPLES:
  confluencecli contentstate list
  confluencecli contentstate get state-123
  confluencecli contentstate create --json '{"name":"In Review","color":"#ffcc00"}'
  confluencecli contentstate update state-123 --json '{"name":"Approved","color":"#00cc00"}'
  confluencecli contentstate delete state-123
  confluencecli contentstate content state-123
`)
}

func printInlineTaskHelp() {
	fmt.Print(`confluencecli inlinetask - Manage inline tasks

USAGE:
  confluencecli inlinetask <action> [flags]

ACTIONS:
  get <id>              Get inline task by ID
  list                  List all inline tasks
  content <id>          Get inline tasks for a content page
  update <content> <id> Update inline task status

FLAGS:
  --config <path>       Path to config file
  --status <status>     Task status (complete, incomplete)
  --start <n>           Start index for pagination
  --limit <n>           Maximum results to return

EXAMPLES:
  confluencecli inlinetask get task-123
  confluencecli inlinetask list
  confluencecli inlinetask content 456
  confluencecli inlinetask update 456 task-123 --status complete
`)
}

func printRelationHelp() {
	fmt.Print(`confluencecli relation - Manage content relations

USAGE:
  confluencecli relation <action> [flags]

ACTIONS:
  get                   Get a relation
  create                Create a relation
  delete                Delete a relation

FLAGS:
  --config <path>       Path to config file
  --expand <fields>     Comma-separated list of fields to expand
  --relation <name>     Relation name
  --source-type <type>  Source entity type
  --source-key <key>    Source entity key
  --target-type <type>  Target entity type
  --target-key <key>    Target entity key
  --json <payload>      JSON payload for create

EXAMPLES:
  confluencecli relation get --relation depends-on --source-type content --source-key 123 --target-type content --target-key 456
  confluencecli relation create --relation depends-on --source-type content --source-key 123 --target-type content --target-key 456
  confluencecli relation delete --relation depends-on --source-type content --source-key 123 --target-type content --target-key 456
`)
}

func printLongTaskHelp() {
	fmt.Print(`confluencecli longtask - Monitor long-running tasks

USAGE:
  confluencecli longtask <action> [flags]

ACTIONS:
  get <id>              Get long task by ID
  list                  List all long tasks

FLAGS:
  --config <path>       Path to config file
  --start <n>           Start index for pagination
  --limit <n>           Maximum results to return

EXAMPLES:
  confluencecli longtask get task-abc123
  confluencecli longtask list
`)
}

func printSystemHelp() {
	fmt.Print(`confluencecli system - View system status (admin)

USAGE:
  confluencecli system <action> [flags]

ACTIONS:
  status                Get system status

FLAGS:
  --config <path>       Path to config file

EXAMPLES:
  confluencecli system status
`)
}

func printBlueprintHelp() {
	fmt.Print(`confluencecli blueprint - Manage content blueprints

USAGE:
  confluencecli blueprint <action> [flags]

ACTIONS:
  list                  List all content blueprints
  get <id>              Get content blueprint by ID

FLAGS:
  --config <path>       Path to config file
  --expand <fields>     Comma-separated list of fields to expand
  --start <n>           Start index for pagination
  --limit <n>           Maximum results to return

EXAMPLES:
  confluencecli blueprint list
  confluencecli blueprint get bp-123
`)
}

func printHealthCheckHelp() {
	fmt.Print(`confluencecli healthcheck - Check system health

USAGE:
  confluencecli healthcheck <action> [flags]

ACTIONS:
  get                   Get system health check results

FLAGS:
  --config <path>       Path to config file

EXAMPLES:
  confluencecli healthcheck get
`)
}

func printCommandHelp(command string) {
	switch command {
	case "init":
		printInitHelp()
	case "connect":
		printConnectHelp()
	case "content":
		printContentHelp()
	case "space":
		printSpaceHelp()
	case "search":
		printSearchHelp()
	case "user":
		printUserHelp()
	case "group":
		printGroupHelp()
	case "settings":
		printSettingsHelp()
	case "audit":
		printAuditHelp()
	case "template":
		printTemplateHelp()
	case "contentstate":
		printContentStateHelp()
	case "inlinetask":
		printInlineTaskHelp()
	case "relation":
		printRelationHelp()
	case "longtask":
		printLongTaskHelp()
	case "system":
		printSystemHelp()
	case "blueprint":
		printBlueprintHelp()
	case "healthcheck":
		printHealthCheckHelp()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		printRootHelp()
		os.Exit(1)
	}
}
