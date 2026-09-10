# Confluence CLI - Setup & Connection

## Initialize Configuration

```bash
confluencecli init [--dir <path>]
```

**Flags:**
- `--dir <path>`: Directory to create config in (default: ~/.config/darthkoax/confluencecli)

**Description:**
Creates a default configuration file at the specified location. The config file contains placeholders for your Confluence instance URL and API token.

**Example:**
```bash
confluencecli init
confluencecli init --dir /custom/path
```

**Output:**
```
Configuration file created at: ~/.config/darthkoax/confluencecli/config.toml
Edit the file to set your Confluence base_url and api_token.
```

## Test Connection

```bash
confluencecli connect [--config <path>]
```

**Flags:**
- `--config <path>`: Path to config file (default: ~/.config/darthkoax/confluencecli/config.toml)

**Description:**
Connects to the Confluence Data Center instance using the configuration file and verifies authentication by fetching the current user's information.

**Example:**
```bash
confluencecli connect
confluencecli connect --config /path/to/config.toml
```

**Output:**
```
Connected to Confluence DC as: John Doe (john.doe)
Confluence instance: https://confluence.example.com
```

**Errors:**
- Config file not found: Run `confluencecli init` first
- Invalid credentials: Check API token
- Network error: Verify Confluence instance is accessible

## Get Current User

```bash
confluencecli user current [--expand <list>] [--config <path>]
```

**Flags:**
- `--expand <list>`: Comma-separated list of properties to expand
- `--config <path>`: Path to config file

**Description:**
Retrieves information about the currently authenticated user. Maps to GET /rest/api/user/current.

**Example:**
```bash
confluencecli user current
confluencecli user current --expand groups
```

**Output:**
```json
{
  "type": "known",
  "username": "john.doe",
  "userKey": "abc123",
  "accountId": "5b10ac8d82e05b22cc7d4ef5",
  "displayName": "John Doe",
  "email": "john.doe@example.com",
  "profilePicture": {
    "path": "https://confluence.example.com/avatars/123"
  }
}
```

## Version Information

```bash
confluencecli version
```

**Description:**
Prints the CLI version and GitHub repository information.

**Output:**
```
confluencecli version 0.1.0
GitHub: https://github.com/darkkoax/confluencecli
```

## Health Check

```bash
confluencecli healthcheck get [--config <path>]
```

**Description:**
Returns the health status of the Confluence instance. Maps to GET /rest/health.

**Example:**
```bash
confluencecli healthcheck get
```

**Output:**
```json
{
  "state": "HEALTHY",
  "checks": [
    {
      "name": "database",
      "description": "Database connection",
      "passed": true
    },
    {
      "name": "search",
      "description": "Search index",
      "passed": true
    }
  ]
}
```
