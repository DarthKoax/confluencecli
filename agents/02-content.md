# Confluence CLI - Content Management

## Get Content

```bash
confluencecli content get <content-id> [--expand <list>] [--config <path>]
```

**Arguments:**
- `<content-id>`: Content ID (e.g., 12345)

**Flags:**
- `--expand <list>`: Comma-separated list of properties to expand (body.storage, body.view, version, space, ancestors, children, descendants)
- `--config <path>`: Path to config file

**Description:**
Retrieves detailed information about a Confluence page or blog post. Maps to GET /rest/api/content/{id}.

**Examples:**
```bash
confluencecli content get 12345
confluencecli content get 12345 --expand "body.storage,version,space"
```

**Output:**
```json
{
  "id": "12345",
  "type": "page",
  "status": "current",
  "title": "Meeting Notes",
  "space": {
    "key": "DEV",
    "name": "Development"
  },
  "version": {
    "number": 3,
    "when": "2024-01-15T10:00:00.000Z",
    "message": "Updated meeting notes"
  },
  "body": {
    "storage": {
      "value": "<p>Meeting notes content</p>",
      "representation": "storage"
    }
  }
}
```

## List Content

```bash
confluencecli content list [--type <type>] [--space-key <key>] [--title <title>] [--start <num>] [--limit <num>] [--expand <list>] [--config <path>]
```

**Flags:**
- `--type <type>`: Content type (page, blogpost)
- `--space-key <key>`: Space key to filter by
- `--title <title>`: Title to filter by
- `--start <num>`: Start index for pagination (default: 0)
- `--limit <num>`: Maximum results to return (default: 25)
- `--expand <list>`: Comma-separated list of properties to expand
- `--config <path>`: Path to config file

**Description:**
Lists content (pages and blog posts) with optional filtering. Maps to GET /rest/api/content.

**Examples:**
```bash
confluencecli content list
confluencecli content list --type page --space-key DEV
confluencecli content list --title "Meeting Notes" --limit 10
```

**Output:**
```json
{
  "results": [
    {
      "id": "12345",
      "type": "page",
      "status": "current",
      "title": "Meeting Notes"
    }
  ],
  "start": 0,
  "limit": 25,
  "size": 1
}
```

## Create Content

```bash
confluencecli content create --json <json-payload> [--config <path>]
```

**Flags:**
- `--json <payload>`: JSON object with content fields (required)
- `--config <path>`: Path to config file

**Description:**
Creates a new page or blog post in Confluence. Maps to POST /rest/api/content.

**JSON Payload Structure:**
```json
{
  "type": "page",
  "title": "New Page",
  "space": {"key": "DEV"},
  "body": {
    "storage": {
      "value": "<p>Page content in Confluence storage format</p>",
      "representation": "storage"
    }
  },
  "ancestors": [{"id": "12345"}]
}
```

**Examples:**
```bash
# Simple page
confluencecli content create --json '{
  "type": "page",
  "title": "New Page",
  "space": {"key": "DEV"},
  "body": {
    "storage": {
      "value": "<p>Hello World</p>",
      "representation": "storage"
    }
  }
}'

# Blog post
confluencecli content create --json '{
  "type": "blogpost",
  "title": "Weekly Update",
  "space": {"key": "TEAM"},
  "body": {
    "storage": {
      "value": "<p>This week we accomplished...</p>",
      "representation": "storage"
    }
  }
}'

# Child page with parent
confluencecli content create --json '{
  "type": "page",
  "title": "Sub Page",
  "space": {"key": "DEV"},
  "ancestors": [{"id": "12345"}],
  "body": {
    "storage": {
      "value": "<p>Child content</p>",
      "representation": "storage"
    }
  }
}'
```

**Output:**
```json
{
  "id": "12346",
  "type": "page",
  "status": "current",
  "title": "New Page"
}
```

## Update Content

```bash
confluencecli content update <content-id> --json <json-payload> [--config <path>]
```

**Arguments:**
- `<content-id>`: Content ID (e.g., 12345)

**Flags:**
- `--json <payload>`: JSON object with fields to update (required)
- `--config <path>`: Path to config file

**Description:**
Updates an existing page or blog post. You must include the current version number incremented by 1. Maps to PUT /rest/api/content/{id}.

**JSON Payload Structure:**
```json
{
  "type": "page",
  "title": "Updated Title",
  "version": {"number": 4},
  "body": {
    "storage": {
      "value": "<p>Updated content</p>",
      "representation": "storage"
    }
  }
}
```

**Examples:**
```bash
# Update title and content
confluencecli content update 12345 --json '{
  "type": "page",
  "title": "Updated Title",
  "version": {"number": 4, "message": "Updated title and content"},
  "body": {
    "storage": {
      "value": "<p>New content here</p>",
      "representation": "storage"
    }
  }
}'

# Minor edit (no notification)
confluencecli content update 12345 --json '{
  "type": "page",
  "title": "Meeting Notes",
  "version": {"number": 4, "minorEdit": true},
  "body": {
    "storage": {
      "value": "<p>Fixed typo</p>",
      "representation": "storage"
    }
  }
}'
```

**Output:**
```json
{
  "id": "12345",
  "type": "page",
  "status": "current",
  "title": "Updated Title",
  "version": {"number": 4}
}
```

## Delete Content

```bash
confluencecli content delete <content-id> [--config <path>]
```

**Arguments:**
- `<content-id>`: Content ID (e.g., 12345)

**Flags:**
- `--config <path>`: Path to config file

**Description:**
Deletes a page or blog post from Confluence. This action cannot be undone. Maps to DELETE /rest/api/content/{id}.

**Examples:**
```bash
confluencecli content delete 12345
```

**Output:**
```
Content 12345 deleted successfully
```

## Content History

```bash
confluencecli content history <content-id> [--expand <list>] [--config <path>]
```

**Arguments:**
- `<content-id>`: Content ID (e.g., 12345)

**Flags:**
- `--expand <list>`: Comma-separated list of properties to expand
- `--config <path>`: Path to config file

**Description:**
Retrieves the version history of a page or blog post. Maps to GET /rest/api/content/{id}/history.

**Example:**
```bash
confluencecli content history 12345
```

**Output:**
```json
{
  "latest": true,
  "createdBy": {
    "username": "john.doe",
    "displayName": "John Doe"
  },
  "createdAt": "2024-01-10T08:00:00.000Z",
  "lastUpdated": {
    "by": {"username": "jane.smith", "displayName": "Jane Smith"},
    "when": "2024-01-15T10:00:00.000Z",
    "message": "Updated content",
    "number": 3
  }
}
```

## Content Children

```bash
confluencecli content children <content-id> [--expand <list>] [--config <path>]
```

**Arguments:**
- `<content-id>`: Content ID (e.g., 12345)

**Flags:**
- `--expand <list>`: Comma-separated list of properties to expand
- `--config <path>`: Path to config file

**Description:**
Retrieves the child pages and comments of a content item. Maps to GET /rest/api/content/{id}/child.

**Example:**
```bash
confluencecli content children 12345
```

**Output:**
```json
{
  "page": {
    "results": [
      {"id": "12346", "title": "Child Page 1"},
      {"id": "12347", "title": "Child Page 2"}
    ]
  },
  "comment": {
    "results": []
  }
}
```

## Content Descendants

```bash
confluencecli content descendants <content-id> [--expand <list>] [--config <path>]
```

**Arguments:**
- `<content-id>`: Content ID (e.g., 12345)

**Flags:**
- `--expand <list>`: Comma-separated list of properties to expand
- `--config <path>`: Path to config file

**Description:**
Retrieves all descendant pages of a content item (recursive children). Maps to GET /rest/api/content/{id}/descendant.

**Example:**
```bash
confluencecli content descendants 12345
```

## Content Versions

```bash
confluencecli content versions <content-id> [--start <num>] [--limit <num>] [--config <path>]
```

**Arguments:**
- `<content-id>`: Content ID (e.g., 12345)

**Flags:**
- `--start <num>`: Start index for pagination (default: 0)
- `--limit <num>`: Maximum results to return (default: 25)
- `--config <path>`: Path to config file

**Description:**
Retrieves the version history of a content item. Maps to GET /rest/api/content/{id}/version.

**Example:**
```bash
confluencecli content versions 12345
```

## Content Comments

```bash
confluencecli content comments <content-id> [--config <path>]
```

**Arguments:**
- `<content-id>`: Content ID (e.g., 12345)

**Flags:**
- `--config <path>`: Path to config file

**Description:**
Retrieves the comments on a content item. Maps to GET /rest/api/content/{id}/child/comment.

**Example:**
```bash
confluencecli content comments 12345
```

**Output:**
```json
{
  "results": [
    {
      "id": "100",
      "type": "comment",
      "title": "Re: Meeting Notes",
      "body": {
        "storage": {
          "value": "<p>Great meeting!</p>",
          "representation": "storage"
        }
      }
    }
  ]
}
```

## Content Attachments

```bash
confluencecli content attachments <content-id> [--config <path>]
```

**Arguments:**
- `<content-id>`: Content ID (e.g., 12345)

**Flags:**
- `--config <path>`: Path to config file

**Description:**
Retrieves the attachments on a content item. Maps to GET /rest/api/content/{id}/child/attachment.

**Example:**
```bash
confluencecli content attachments 12345
```

**Output:**
```json
{
  "results": [
    {
      "id": "att-100",
      "type": "attachment",
      "title": "document.pdf",
      "metadata": {
        "mediaType": "application/pdf"
      }
    }
  ]
}
```
