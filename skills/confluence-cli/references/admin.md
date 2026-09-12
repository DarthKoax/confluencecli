# Confluence CLI - Admin Operations

**Note:** All admin operations require `admin_mode = true` in config.toml.

## Settings

### Get System Info
```bash
confluencecli settings systeminfo [--config <path>]
```

**Description:**
Retrieves system information about the Confluence instance. Maps to GET /rest/api/settings/systemInfo.

**Example:**
```bash
confluencecli settings systeminfo
```

**Output:**
```json
{
  "cloudId": "abc-123-def-456",
  "commitHash": "a1b2c3d4e5f6",
  "buildDate": "2024-01-10T00:00:00.000Z"
}
```

### Get Theme
```bash
confluencecli settings theme-get [--config <path>]
```

**Description:**
Retrieves the current global theme. Maps to GET /rest/api/settings/theme.

**Example:**
```bash
confluencecli settings theme-get
```

**Output:**
```json
{
  "themeKey": "com.atlassian.confluence.theme:default",
  "name": "Default Theme",
  "description": "The default Confluence theme"
}
```

### Set Theme
```bash
confluencecli settings theme-set --theme-key <theme-key> [--config <path>]
```

**Flags:**
- `--theme-key <key>`: Theme key to set (required)

**Example:**
```bash
confluencecli settings theme-set --theme-key "com.atlassian.confluence.theme:mytheme"
```

**Output:**
```
Theme set to com.atlassian.confluence.theme:mytheme successfully
```

### Get Look and Feel
```bash
confluencecli settings laf-get [--space-key <key>] [--config <path>]
```

**Flags:**
- `--space-key <key>`: Space key for space-specific settings (omit for global)

**Description:**
Retrieves the look and feel settings. Maps to GET /rest/api/settings/lookandfeel.

**Example:**
```bash
confluencecli settings laf-get
confluencecli settings laf-get --space-key DEV
```

**Output:**
```json
{
  "head": {
    "background": {"value": "#ffffff"}
  },
  "header": {
    "background": {"value": "#205081"}
  },
  "content": {
    "background": {"value": "#ffffff"},
    "header": {"value": "#205081"},
    "body": {"value": "#333333"},
    "links": {"value": "#205081"},
    "borders": {"value": "#cccccc"}
  }
}
```

### Set Look and Feel
```bash
confluencecli settings laf-set --json <json-payload> [--config <path>]
```

**Flags:**
- `--json <payload>`: JSON object with look and feel settings (required)

**Example:**
```bash
confluencecli settings laf-set --json '{
  "content": {
    "background": {"value": "#ffffff"},
    "links": {"value": "#0066cc"}
  }
}'
```

**Output:**
```
Look and feel updated successfully
```

### Reset Look and Feel
```bash
confluencecli settings laf-reset [--space-key <key>] [--config <path>]
```

**Flags:**
- `--space-key <key>`: Space key for space-specific reset (omit for global)

**Example:**
```bash
confluencecli settings laf-reset
confluencecli settings laf-reset --space-key DEV
```

**Output:**
```
Look and feel reset successfully
```

## Audit

### Get Audit Records
```bash
confluencecli audit get [--start-date <date>] [--end-date <date>] [--search <string>] [--start <num>] [--limit <num>] [--config <path>]
```

**Flags:**
- `--start-date <date>`: Start date (ISO 8601 format)
- `--end-date <date>`: End date (ISO 8601 format)
- `--search <string>`: Search string to filter records
- `--start <num>`: Start index for pagination (default: 0)
- `--limit <num>`: Maximum results to return (default: 25)

**Description:**
Retrieves audit log records. Maps to GET /rest/audit/1.0/audit.

**Examples:**
```bash
confluencecli audit get
confluencecli audit get --start-date 2024-01-01 --end-date 2024-12-31
confluencecli audit get --search "page created"
```

**Output:**
```json
{
  "results": [
    {
      "author": {
        "username": "john.doe",
        "displayName": "John Doe"
      },
      "remoteAddress": "192.168.1.100",
      "creationDate": 1705312800000,
      "summary": "Page created",
      "description": "Created page 'Meeting Notes' in space DEV",
      "category": "content",
      "sysAdmin": false
    }
  ],
  "start": 0,
  "limit": 25,
  "size": 1
}
```

### Get Audit Records Since
```bash
confluencecli audit since --number <value> [--search <string>] [--start <num>] [--limit <num>] [--config <path>]
```

**Flags:**
- `--number <value>`: Time period (e.g., 1m, 7d, 24h)
- `--search <string>`: Search string to filter records
- `--start <num>`: Start index for pagination (default: 0)
- `--limit <num>`: Maximum results to return (default: 25)

**Example:**
```bash
confluencecli audit since --number 7d
confluencecli audit since --number 24h --search "deleted"
```

### Get Audit Retention Period
```bash
confluencecli audit retention [--config <path>]
```

**Description:**
Retrieves the current audit log retention period. Maps to GET /rest/audit/1.0/audit/retention.

**Example:**
```bash
confluencecli audit retention
```

**Output:**
```json
{
  "number": 365,
  "units": "DAYS"
}
```

### Set Audit Retention Period
```bash
confluencecli audit retention set --number <num> --units <unit> [--config <path>]
```

**Flags:**
- `--number <num>`: Number of units (required)
- `--units <unit>`: Units (DAYS, MONTHS) (required)

**Example:**
```bash
confluencecli audit retention set --number 365 --units DAYS
```

**Output:**
```
Audit retention period updated successfully
```

### Export Audit Records
```bash
confluencecli audit export --start-date <date> --end-date <date> [--config <path>]
```

**Flags:**
- `--start-date <date>`: Start date (ISO 8601 format) (required)
- `--end-date <date>`: End date (ISO 8601 format) (required)

**Example:**
```bash
confluencecli audit export --start-date 2024-01-01 --end-date 2024-12-31
```

## System

### Get System Status
```bash
confluencecli system status [--config <path>]
```

**Description:**
Retrieves the system status. Maps to GET /rest/api/system/status.

**Example:**
```bash
confluencecli system status
```

**Output:**
```json
{
  "state": "RUNNING"
}
```

## Templates

### List Content Templates
```bash
confluencecli template list-content [--space-key <key>] [--start <num>] [--limit <num>] [--expand <list>] [--config <path>]
```

**Flags:**
- `--space-key <key>`: Space key to filter by
- `--start <num>`: Start index for pagination (default: 0)
- `--limit <num>`: Maximum results to return (default: 25)
- `--expand <list>`: Comma-separated list of properties to expand

**Example:**
```bash
confluencecli template list-content
confluencecli template list-content --space-key DEV
```

### Get Content Template
```bash
confluencecli template get-content <template-id> [--expand <list>] [--config <path>]
```

**Arguments:**
- `<template-id>`: Template ID

**Example:**
```bash
confluencecli template get-content tmpl-123
```

### Create Content Template
```bash
confluencecli template create-content --json <json-payload> [--config <path>]
```

**JSON Payload:**
```json
{
  "name": "Meeting Notes",
  "description": "Template for meeting notes",
  "templateType": "page",
  "body": {
    "storage": {
      "value": "<h1>Meeting Notes</h1><p>Date:</p><p>Attendees:</p><p>Agenda:</p>",
      "representation": "storage"
    }
  }
}
```

**Example:**
```bash
confluencecli template create-content --json '{
  "name": "Meeting Notes",
  "templateType": "page",
  "body": {
    "storage": {
      "value": "<h1>Meeting Notes</h1>",
      "representation": "storage"
    }
  }
}'
```

### Delete Content Template
```bash
confluencecli template delete-content <template-id> [--config <path>]
```

**Example:**
```bash
confluencecli template delete-content tmpl-123
```

### List Blueprint Templates
```bash
confluencecli template list-blueprint [--space-key <key>] [--start <num>] [--limit <num>] [--expand <list>] [--config <path>]
```

**Example:**
```bash
confluencecli template list-blueprint
```

### Get Blueprint Template
```bash
confluencecli template get-blueprint <blueprint-id> [--expand <list>] [--config <path>]
```

**Example:**
```bash
confluencecli template get-blueprint bp-456
```

## Content States

### List Content States
```bash
confluencecli contentstate list [--start <num>] [--limit <num>] [--config <path>]
```

**Example:**
```bash
confluencecli contentstate list
```

**Output:**
```json
{
  "results": [
    {"id": "state-1", "name": "Draft", "color": "#999999"},
    {"id": "state-2", "name": "In Review", "color": "#ffcc00"},
    {"id": "state-3", "name": "Published", "color": "#00cc00"}
  ]
}
```

### Get Content State
```bash
confluencecli contentstate get <state-id> [--config <path>]
```

**Example:**
```bash
confluencecli contentstate get state-1
```

### Create Content State
```bash
confluencecli contentstate create --json <json-payload> [--config <path>]
```

**JSON Payload:**
```json
{
  "name": "In Review",
  "description": "Content is under review",
  "color": "#ffcc00",
  "isFinalState": false
}
```

**Example:**
```bash
confluencecli contentstate create --json '{
  "name": "In Review",
  "color": "#ffcc00"
}'
```

### Update Content State
```bash
confluencecli contentstate update <state-id> --json <json-payload> [--config <path>]
```

**Example:**
```bash
confluencecli contentstate update state-1 --json '{
  "name": "Draft - Updated",
  "color": "#666666"
}'
```

### Delete Content State
```bash
confluencecli contentstate delete <state-id> [--config <path>]
```

**Example:**
```bash
confluencecli contentstate delete state-1
```

### Get Content by State
```bash
confluencecli contentstate content <state-id> [--start <num>] [--limit <num>] [--config <path>]
```

**Example:**
```bash
confluencecli contentstate content state-1
```

## Inline Tasks

### Get Inline Task
```bash
confluencecli inlinetask get <task-id> [--config <path>]
```

**Example:**
```bash
confluencecli inlinetask get task-123
```

**Output:**
```json
{
  "id": "task-123",
  "contentId": "12345",
  "status": "incomplete",
  "creator": "john.doe",
  "dueDate": "2024-02-01"
}
```

### List Inline Tasks
```bash
confluencecli inlinetask list [--start <num>] [--limit <num>] [--config <path>]
```

**Example:**
```bash
confluencecli inlinetask list
```

### Get Tasks by Content
```bash
confluencecli inlinetask content <content-id> [--start <num>] [--limit <num>] [--config <path>]
```

**Example:**
```bash
confluencecli inlinetask content 12345
```

### Update Inline Task
```bash
confluencecli inlinetask update <content-id> <task-id> --status <status> [--config <path>]
```

**Flags:**
- `--status <status>`: Task status (complete, incomplete)

**Example:**
```bash
confluencecli inlinetask update 12345 task-123 --status complete
```

## Relations

### Get Relation
```bash
confluencecli relation get --relation <name> --source-type <type> --source-key <key> --target-type <type> --target-key <key> [--expand <list>] [--config <path>]
```

**Flags:**
- `--relation <name>`: Relation name (required)
- `--source-type <type>`: Source entity type (required)
- `--source-key <key>`: Source entity key (required)
- `--target-type <type>`: Target entity type (required)
- `--target-key <key>`: Target entity key (required)
- `--expand <list>`: Comma-separated list of properties to expand

**Example:**
```bash
confluencecli relation get --relation depends-on --source-type content --source-key 12345 --target-type content --target-key 67890
```

### Create Relation
```bash
confluencecli relation create --relation <name> --source-type <type> --source-key <key> --target-type <type> --target-key <key> [--json <payload>] [--config <path>]
```

**Example:**
```bash
confluencecli relation create --relation depends-on --source-type content --source-key 12345 --target-type content --target-key 67890
```

### Delete Relation
```bash
confluencecli relation delete --relation <name> --source-type <type> --source-key <key> --target-type <type> --target-key <key> [--config <path>]
```

**Example:**
```bash
confluencecli relation delete --relation depends-on --source-type content --source-key 12345 --target-type content --target-key 67890
```

## Long Tasks

### Get Long Task
```bash
confluencecli longtask get <task-id> [--config <path>]
```

**Example:**
```bash
confluencecli longtask get lt-123
```

**Output:**
```json
{
  "id": "lt-123",
  "name": {"key": "task.name", "value": "Content Export"},
  "elapsedTime": 30000,
  "percentageComplete": 75,
  "successful": false,
  "finished": false
}
```

### List Long Tasks
```bash
confluencecli longtask list [--start <num>] [--limit <num>] [--config <path>]
```

**Example:**
```bash
confluencecli longtask list
```

## Blueprints

### List Blueprints
```bash
confluencecli blueprint list [--start <num>] [--limit <num>] [--expand <list>] [--config <path>]
```

**Example:**
```bash
confluencecli blueprint list
```

### Get Blueprint
```bash
confluencecli blueprint get <blueprint-id> [--expand <list>] [--config <path>]
```

**Example:**
```bash
confluencecli blueprint get bp-123
```
