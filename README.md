# confluencecli

A Go wrapper library and CLI for the **Confluence Data Center / Server REST API**. Provides typed, configurable access to all major Confluence DC endpoints with Bearer token authentication and custom CA certificate support.

## Features

- **Full API Coverage**: 15 service wrappers covering Content, Spaces, Search, Users, Groups, Settings, Audit, Templates, Content States, Inline Tasks, Relations, Long Tasks, System, Blueprints, and Health Check
- **Configurable Access Control**: Method gating (GET/POST/PUT/DELETE) and per-endpoint toggles
- **Admin Mode**: Separate admin-only endpoint protection for sensitive operations
- **Custom TLS**: Support for custom CA certificates for self-signed or internal CAs
- **Bearer Token Auth**: Simple API token authentication
- **TOML Configuration**: Human-readable config file format
- **Environment Overrides**: All settings overridable via environment variables

## Installation

### From Source

```bash
git clone https://github.com/darthkoax/confluencecli.git
cd confluencecli
go build -o confluencecli ./cmd/confluencecli
mv ./confluencecli /usr/local/bin/
```

### Read-Only Build

Build a binary that permanently disables all write operations (POST, PUT, DELETE). This is enforced at compile time and cannot be overridden by configuration.

```bash
go build -tags readonly -o confluencecli-readonly ./cmd/confluencecli
```

The read-only binary will display `(read-only)` in the version output:

```
$ confluencecli-readonly version
confluencecli version 0.2.0 (read-only)
```

Any attempt to perform a write operation will fail with an error:

```
POST method is disabled: this binary was compiled in read-only mode (-tags readonly)
```

### Install to GOPATH

```bash
go install github.com/darthkoax/confluencecli/cmd/confluencecli@latest
```

## Quick Start

### 1. Initialize Configuration

```bash
confluencecli init
```

This creates a default config file at `~/.config/confluencecli/config.toml`.

### 2. Edit Configuration

Open the config file and set your Confluence instance details:

```toml
[confluence]
base_url = "https://confluence.example.com"
api_token = "YOUR_API_TOKEN_HERE"
custom_ca_cert = ""
timeout = 30

admin_mode = false

[methods]
allow_get = true
allow_post = true
allow_put = true
allow_delete = true

[endpoints]
content = true
spaces = true
# ... enable/disable endpoints as needed
```

### 3. Test Connection

```bash
confluencecli connect
```

Expected output:
```
Connected to Confluence DC as: John Doe (john.doe)
```

## Configuration

### Config File Location

Default: `~/.config/confluencecli/config.toml`

Override with `CONFLUENCE_CONFIG_DIR` environment variable:
```bash
export CONFLUENCE_CONFIG_DIR="/custom/path"
confluencecli init
```

### Configuration Options

| Key | Description | Default |
|-----|-------------|---------|
| `confluence.base_url` | Confluence DC instance URL (required) | - |
| `confluence.api_token` | Bearer API token (required) | - |
| `confluence.custom_ca_cert` | Path to custom CA certificate bundle | `""` |
| `confluence.timeout` | HTTP timeout in seconds | `30` |
| `admin_mode` | Enable admin-only endpoints | `false` |
| `methods.allow_get` | Enable GET requests | `true` |
| `methods.allow_post` | Enable POST requests | `true` |
| `methods.allow_put` | Enable PUT requests | `true` |
| `methods.allow_delete` | Enable DELETE requests | `true` |
| `endpoints.*` | Enable/disable individual endpoints | `true` |

### Environment Variable Overrides

All configuration options can be overridden via environment variables:

```bash
CONFLUENCE_BASE_URL="https://confluence.example.com"
CONFLUENCE_API_TOKEN="your-token"
CONFLUENCE_CUSTOM_CA_CERT="/path/to/ca.crt"
CONFLUENCE_TIMEOUT="60"
CONFLUENCE_ADMIN_MODE="true"
CONFLUENCE_ALLOW_GET="true"
CONFLUENCE_ALLOW_POST="true"
CONFLUENCE_ALLOW_PUT="true"
CONFLUENCE_ALLOW_DELETE="true"
CONFLUENCE_ENDPOINT_CONTENT="true"
CONFLUENCE_ENDPOINT_SPACES="false"
# ... etc
```

## Agent Setup

### Opencode

#### Skill Config

Copy the included skill to your harness skills directory:
```bash
cp -rp skills/confluence-cli ~/.config/opencode/skills
```

#### Agent Permission

Apply the following to your opencode.json permissions block if you want to manually review write operations of your agent:
```json
"permission": {
  "bash": {
    "confluencecli *": "allow",
    "CONFLUENCE_ALLOW_*=true confluencecli *": "ask"
  }
},
```

## CLI Commands

### `confluencecli init`

Create a default configuration file.

```bash
confluencecli init [--dir <path>]
```

| Flag | Description |
|------|-------------|
| `--dir <path>` | Directory to create config in (default: `~/.config/confluencecli`) |

### `confluencecli connect`

Connect to Confluence and verify authentication.

```bash
confluencecli connect [--config <path>]
```

| Flag | Description |
|------|-------------|
| `--config <path>` | Path to config file (default: `~/.config/confluencecli/config.toml`) |

### `confluencecli version`

Print version information.

```bash
confluencecli version
```

### `confluencecli help`

Show help message.

```bash
confluencecli help
```

## API Coverage

### Confluence REST API (`/rest/api/`)

| Service | Endpoints |
|---------|-----------|
| **Content** | GET/POST/PUT/DELETE content, history, children, descendants, versions, comments, attachments |
| **Spaces** | GET/POST/DELETE space, content |
| **Search** | GET CQL search, user search |
| **Users** | GET current, anonymous, by username/key/accountID, unknown, by email |
| **Groups** | GET group, members, POST/DELETE member |
| **Settings** | GET/PUT system info, theme, look and feel (admin) |
| **Audit** | GET audit records, retention, export (admin) |
| **Templates** | GET/POST/PUT/DELETE content templates, blueprint templates |
| **Content States** | GET/POST/PUT/DELETE content states, content by state |
| **Inline Tasks** | GET/PUT inline tasks, tasks by content |
| **Relations** | GET/POST/DELETE relations |
| **Long Tasks** | GET long tasks |
| **System** | GET system status (admin) |
| **Blueprints** | GET blueprints |
| **Health Check** | GET health check |

## Admin Mode

Some endpoints require admin privileges. Enable `admin_mode = true` in config to access:

- **Settings**: System info, theme, look and feel management
- **Audit**: Audit log retrieval, retention management, export
- **System**: System status monitoring
- **Space Management**: Create/delete spaces
- **Group Management**: Add/remove group members

## Access Control

### Compile-Time Read-Only Mode

Build a binary with all write operations permanently disabled:

```bash
go build -tags readonly -o confluencecli-readonly ./cmd/confluencecli
```

This enforces read-only access at compile time. POST, PUT, and DELETE methods cannot be enabled via configuration. See [Installation](#read-only-build) for details.

### Method Gating

Disable HTTP methods to create read-only or append-only modes:

```toml
[methods]
allow_get = true
allow_post = false    # Read-only mode
allow_put = false
allow_delete = false
```

### Endpoint Gating

Disable specific endpoints to restrict functionality:

```toml
[endpoints]
content = true
spaces = true
audit = false         # Disable audit operations
settings = false      # Disable settings operations
```

## Development

### Build

```bash
go build ./...
```

### Test

```bash
go test ./... -v
```

### Test with Coverage

```bash
go test ./... -cover
```

### Lint

```bash
go vet ./...
```

## Project Structure

```
.
├── cmd/confluencecli/    # CLI entry point
├── internal/
│   ├── api/              # API service wrappers (15 services)
│   ├── client/           # HTTP client with TLS, auth, method gating
│   └── config/           # TOML config loading
├── skills/               # Agent skill documentation (skills/confluence-cli/)
├── AGENTS.md             # Developer documentation
└── README.md             # This file
```

## Architecture

### Client

The HTTP client (`internal/client/`) handles:
- Bearer token authentication
- Custom CA certificate loading
- Method gating (GET/POST/PUT/DELETE)
- Endpoint gating
- Admin mode checks
- TLS configuration

### Services

Each API resource has a dedicated service (`internal/api/`):
- Type-safe request/response structures
- Endpoint and admin checks before each operation
- Consistent error handling

### Config

Configuration (`internal/config/`) supports:
- TOML file parsing
- Environment variable overrides
- Default values
- Validation

## Error Handling

- **Config errors**: Returned at load time with descriptive messages
- **TLS/CA errors**: Returned at client creation time
- **API errors**: Return HTTP status code and response body
- **Method/endpoint disabled**: Return error before HTTP request

## Limitations

- Multipart file uploads (attachments) not yet implemented
- OAuth authentication not supported (Bearer token only)
- Some advanced endpoints excluded (see AGENTS.md for full list)

## License

See LICENSE file for details.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Write tests for new functionality
4. Ensure all tests pass: `go test ./... -v`
5. Ensure code passes lint: `go vet ./...`
6. Submit a pull request

See AGENTS.md for detailed development guidelines and test coverage requirements.
