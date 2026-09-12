# Confluence CLI - Space Management

## List Spaces

```bash
confluencecli space list [--type <type>] [--status <status>] [--start <num>] [--limit <num>] [--expand <list>] [--config <path>]
```

**Flags:**
- `--type <type>`: Space type (global, personal)
- `--status <status>`: Space status (current, archived)
- `--start <num>`: Start index for pagination (default: 0)
- `--limit <num>`: Maximum results to return (default: 25)
- `--expand <list>`: Comma-separated list of properties to expand (description, homepage, icon)
- `--config <path>`: Path to config file

**Description:**
Lists all spaces in the Confluence instance. Maps to GET /rest/api/space.

**Examples:**
```bash
confluencecli space list
confluencecli space list --type global
confluencecli space list --status current --limit 10
```

**Output:**
```json
{
  "results": [
    {
      "id": 10000,
      "key": "DEV",
      "name": "Development",
      "type": "global",
      "status": "current"
    },
    {
      "id": 10001,
      "key": "TEAM",
      "name": "Team Space",
      "type": "global",
      "status": "current"
    }
  ],
  "start": 0,
  "limit": 25,
  "size": 2
}
```

## Get Space

```bash
confluencecli space get <space-key> [--expand <list>] [--config <path>]
```

**Arguments:**
- `<space-key>`: Space key (e.g., DEV)

**Flags:**
- `--expand <list>`: Comma-separated list of properties to expand (description, homepage, icon)
- `--config <path>`: Path to config file

**Description:**
Retrieves detailed information about a space. Maps to GET /rest/api/space/{spaceKey}.

**Example:**
```bash
confluencecli space get DEV
confluencecli space get DEV --expand description,homepage
```

**Output:**
```json
{
  "id": 10000,
  "key": "DEV",
  "name": "Development",
  "type": "global",
  "status": "current",
  "description": {
    "plain": {
      "value": "Development team documentation",
      "representation": "plain"
    }
  },
  "homepage": {
    "id": "12345",
    "type": "page",
    "title": "Development Home"
  }
}
```

## Create Space (Admin)

```bash
confluencecli space create --json <json-payload> [--config <path>]
```

**Flags:**
- `--json <payload>`: JSON object with space fields (required)
- `--config <path>`: Path to config file

**Description:**
Creates a new space in Confluence. Requires admin mode. Maps to POST /rest/api/space.

**JSON Payload Structure:**
```json
{
  "key": "NEW",
  "name": "New Space",
  "type": "global",
  "description": {
    "plain": {
      "value": "Space description",
      "representation": "plain"
    }
  }
}
```

**Examples:**
```bash
confluencecli space create --json '{
  "key": "NEW",
  "name": "New Space",
  "type": "global",
  "description": {
    "plain": {
      "value": "A new space for the team",
      "representation": "plain"
    }
  }
}'
```

**Output:**
```json
{
  "id": 10002,
  "key": "NEW",
  "name": "New Space",
  "type": "global",
  "status": "current"
}
```

**Note:** Requires admin_mode = true in config.

## Delete Space (Admin)

```bash
confluencecli space delete <space-key> [--config <path>]
```

**Arguments:**
- `<space-key>`: Space key (e.g., DEV)

**Flags:**
- `--config <path>`: Path to config file

**Description:**
Deletes a space from Confluence. This action cannot be undone. Requires admin mode. Maps to DELETE /rest/api/space/{spaceKey}.

**Examples:**
```bash
confluencecli space delete OLD
```

**Output:**
```
Space OLD deleted successfully
```

**Note:** Requires admin_mode = true in config.

## Get Space Content

```bash
confluencecli space content <space-key> [--type <type>] [--start <num>] [--limit <num>] [--expand <list>] [--config <path>]
```

**Arguments:**
- `<space-key>`: Space key (e.g., DEV)

**Flags:**
- `--type <type>`: Content type (page, blogpost)
- `--start <num>`: Start index for pagination (default: 0)
- `--limit <num>`: Maximum results to return (default: 25)
- `--expand <list>`: Comma-separated list of properties to expand
- `--config <path>`: Path to config file

**Description:**
Retrieves content within a space. Maps to GET /rest/api/space/{spaceKey}/content.

**Examples:**
```bash
confluencecli space content DEV
confluencecli space content DEV --type page
confluencecli space content DEV --type blogpost --limit 10
```

**Output:**
```json
{
  "results": [
    {
      "id": "12345",
      "type": "page",
      "status": "current",
      "title": "Development Home"
    },
    {
      "id": "12346",
      "type": "page",
      "status": "current",
      "title": "API Documentation"
    }
  ],
  "start": 0,
  "limit": 25,
  "size": 2
}
```
