# Confluence CLI - Users & Groups

## Users

### Get Current User
```bash
confluencecli user current [--expand <list>] [--config <path>]
```

**Flags:**
- `--expand <list>`: Comma-separated list of properties to expand
- `--config <path>`: Path to config file

**Example:**
```bash
confluencecli user current
```

**Output:**
```json
{
  "type": "known",
  "username": "john.doe",
  "userKey": "abc123",
  "accountId": "5b10ac8d82e05b22cc7d4ef5",
  "displayName": "John Doe",
  "email": "john.doe@example.com"
}
```

### Get Anonymous User
```bash
confluencecli user anonymous [--expand <list>] [--config <path>]
```

**Example:**
```bash
confluencecli user anonymous
```

**Output:**
```json
{
  "type": "anonymous",
  "displayName": "Anonymous"
}
```

### Get User
```bash
confluencecli user get --username <username> [--expand <list>] [--config <path>]
```

**Flags:**
- `--username <username>`: Username to look up (one of --username, --key, or --account-id required)
- `--key <key>`: User key to look up
- `--account-id <id>`: Account ID to look up
- `--expand <list>`: Comma-separated list of properties to expand

**Example:**
```bash
confluencecli user get --username john.doe
confluencecli user get --account-id 5b10ac8d82e05b22cc7d4ef5
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

### Get Unknown User
```bash
confluencecli user unknown [--expand <list>] [--config <path>]
```

**Description:**
Retrieves the unknown user representation (used when a user cannot be found).

**Example:**
```bash
confluencecli user unknown
```

**Output:**
```json
{
  "type": "unknown",
  "displayName": "Unknown"
}
```

### Get User by Email
```bash
confluencecli user email --email <email> [--expand <list>] [--config <path>]
```

**Flags:**
- `--email <email>`: Email address to look up (required)
- `--expand <list>`: Comma-separated list of properties to expand

**Example:**
```bash
confluencecli user email --email john.doe@example.com
```

## Groups

### Get Group
```bash
confluencecli group get --groupname <groupname> [--expand <list>] [--config <path>]
```

**Flags:**
- `--groupname <groupname>`: Group name (required)
- `--expand <list>`: Comma-separated list of properties to expand

**Example:**
```bash
confluencecli group get --groupname developers
```

**Output:**
```json
{
  "type": "group",
  "name": "developers",
  "id": "group-123"
}
```

### Get Group Members
```bash
confluencecli group members --groupname <groupname> [--start <num>] [--limit <num>] [--config <path>]
```

**Flags:**
- `--groupname <groupname>`: Group name (required)
- `--start <num>`: Start index for pagination (default: 0)
- `--limit <num>`: Maximum results to return (default: 25)

**Example:**
```bash
confluencecli group members --groupname developers
confluencecli group members --groupname developers --limit 100
```

**Output:**
```json
{
  "results": [
    {
      "type": "known",
      "username": "john.doe",
      "displayName": "John Doe",
      "email": "john.doe@example.com"
    },
    {
      "type": "known",
      "username": "jane.smith",
      "displayName": "Jane Smith",
      "email": "jane.smith@example.com"
    }
  ],
  "start": 0,
  "limit": 25,
  "size": 2
}
```

### Add Member to Group (Admin)
```bash
confluencecli group add-member --groupname <groupname> --username <username> [--config <path>]
```

**Flags:**
- `--groupname <groupname>`: Group name (required)
- `--username <username>`: Username to add (required)

**Example:**
```bash
confluencecli group add-member --groupname developers --username john.doe
```

**Output:**
```
User john.doe added to group developers successfully
```

**Note:** Requires admin_mode = true in config.

### Remove Member from Group (Admin)
```bash
confluencecli group remove-member --groupname <groupname> --username <username> [--config <path>]
```

**Flags:**
- `--groupname <groupname>`: Group name (required)
- `--username <username>`: Username to remove (required)

**Example:**
```bash
confluencecli group remove-member --groupname developers --username john.doe
```

**Output:**
```
User john.doe removed from group developers successfully
```

**Note:** Requires admin_mode = true in config.
