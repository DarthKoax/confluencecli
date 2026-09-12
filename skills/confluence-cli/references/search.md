# Confluence CLI - Search

## CQL Search

```bash
confluencecli search cql --cql <cql-query> [--start <num>] [--limit <num>] [--expand <list>] [--config <path>]
```

**Flags:**
- `--cql <query>`: CQL query string (required)
- `--start <num>`: Start index for pagination (default: 0)
- `--limit <num>`: Maximum results to return (default: 25)
- `--expand <list>`: Comma-separated list of properties to expand
- `--config <path>`: Path to config file

**Description:**
Searches for content using Confluence Query Language (CQL). Maps to GET /rest/api/search.

**Examples:**
```bash
# Simple search by space
confluencecli search cql --cql "space=DEV"

# Search by content type
confluencecli search cql --cql "type=page AND space=DEV"

# Search by title
confluencecli search cql --cql "title~'Meeting Notes'"

# Search with pagination
confluencecli search cql --cql "space=DEV" --start 0 --limit 10

# Search by creator
confluencecli search cql --cql "creator=currentUser()"

# Search by label
confluencecli search cql --cql "label=important"

# Search with date range
confluencecli search cql --cql "space=DEV AND lastModified>startOfMonth()"

# Search by content type blogpost
confluencecli search cql --cql "type=blogpost AND space=TEAM"

# Complex CQL query
confluencecli search cql --cql "type=page AND space=DEV AND title~'API' AND lastModified>startOfMonth(-1)"
```

**Output:**
```json
{
  "results": [
    {
      "content": {
        "id": "12345",
        "type": "page",
        "status": "current",
        "title": "API Documentation"
      },
      "title": "API Documentation",
      "excerpt": "...REST API endpoints for the project...",
      "url": "/display/DEV/API+Documentation",
      "entityType": "content",
      "lastModified": "2024-01-15T10:00:00.000Z",
      "friendlyLastModified": "2 days ago"
    }
  ],
  "start": 0,
  "limit": 25,
  "size": 1,
  "totalSize": 1,
  "cqlQuery": "type=page AND space=DEV AND title~'API'"
}
```

**CQL Reference:**

Common operators:
- `=` (equals)
- `!=` (not equals)
- `~` (contains, for text fields)
- `>` `<` `>=` `<=` (comparison, for dates and numbers)
- `IN` (in list)
- `NOT IN` (not in list)

Common fields:
- `type` - Content type (page, blogpost, comment)
- `space` - Space key
- `title` - Content title
- `creator` - Content creator username
- `contributor` - Content contributor username
- `label` - Content label
- `lastModified` - Last modification date
- `created` - Creation date
- `content` - Content body text (for full-text search)
- `ancestor` - Parent page ID
- `space.type` - Space type (global, personal)

Common functions:
- `currentUser()` - Current authenticated user
- `now()` - Current timestamp
- `startOfDay()`, `endOfDay()` - Day boundaries
- `startOfWeek()`, `endOfWeek()` - Week boundaries
- `startOfMonth()`, `endOfMonth()` - Month boundaries
- `startOfYear()`, `endOfYear()` - Year boundaries
- `myContent()` - Content created by current user

## User Search

```bash
confluencecli search user --cql <cql-query> [--start <num>] [--limit <num>] [--config <path>]
```

**Flags:**
- `--cql <query>`: CQL query string for user search (required)
- `--start <num>`: Start index for pagination (default: 0)
- `--limit <num>`: Maximum results to return (default: 25)
- `--config <path>`: Path to config file

**Description:**
Searches for users using CQL. Maps to GET /rest/api/search/user.

**Examples:**
```bash
confluencecli search user --cql "user.fullname~'John'"
confluencecli search user --cql "user.type=known"
```

**Output:**
```json
{
  "results": [
    {
      "user": {
        "type": "known",
        "username": "john.doe",
        "displayName": "John Doe",
        "email": "john.doe@example.com"
      },
      "entityType": "user"
    }
  ],
  "start": 0,
  "limit": 25,
  "size": 1
}
```
