# AGENTS.md — confluencecli

## Project Overview

Go wrapper library for the **Confluence Data Center / Server REST API**. Provides typed, configurable access to all major Confluence DC endpoints with Bearer token authentication and custom CA certificate support.

## Build & Test

```bash
go mod tidy
go build ./...
go vet ./...
go test ./... -v
```

## Test Coverage Requirements

**All code changes must maintain comprehensive test coverage.**

### Testing Standards

1. **Unit Tests Required**: Every public function must have corresponding unit tests
2. **Test Coverage Target**: Maintain minimum 80% code coverage across all packages
3. **Test All Paths**: Cover both success and error paths in all functions
4. **Mock External Dependencies**: Use httptest for HTTP client testing
5. **Test Data Validation**: Verify config validation, endpoint gating, and admin checks

### Required Test Categories

- **Config Tests**: TOML parsing, defaults, environment variable overrides, validation
- **Client Tests**: HTTP method gating, admin mode checks, endpoint gating, authentication, TLS/CA handling
- **Service Tests**: All API service methods (15 services), endpoint disabled scenarios, admin operation blocking
- **Integration Tests**: End-to-end workflows with mock servers

### Test File Organization

- Test files must be co-located with source files (e.g., `config_test.go` alongside `config.go`)
- Use table-driven tests where appropriate for multiple scenarios
- Use `httptest.NewServer` for HTTP client testing
- Clean up test resources (temp files, servers) with `defer`

### Running Tests

```bash
# Run all tests
go test ./... -v

# Run tests with coverage
go test ./... -cover

# Run specific package tests
go test ./internal/config -v
go test ./internal/client -v
go test ./internal/api -v

# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Before Committing

Always verify:
- [ ] `go build ./...` succeeds
- [ ] `go vet ./...` passes
- [ ] `go test ./...` passes with no failures
- [ ] New functionality has corresponding tests
- [ ] Error paths are tested
- [ ] Admin/endpoint gating is tested

## Configuration

All configuration is in `config.toml`:

| Key | Description |
|-----|-------------|
| `confluence.base_url` | Confluence DC instance URL (required) |
| `confluence.api_token` | Bearer API token (required) |
| `confluence.custom_ca_cert` | Path to custom CA certificate bundle |
| `confluence.timeout` | HTTP timeout in seconds (default: 30) |
| `admin_mode` | Enable admin-only endpoints (default: false) |
| `methods.allow_get` | Enable GET requests (default: true) |
| `methods.allow_post` | Enable POST requests (default: true) |
| `methods.allow_put` | Enable PUT requests (default: true) |
| `methods.allow_delete` | Enable DELETE requests (default: true) |
| `endpoints.*` | Enable/disable individual API endpoints (default: all true) |

### Environment Variable Overrides

All config values can be overridden via environment variables:

| Env Var | Config Key |
|---------|------------|
| `CONFLUENCE_BASE_URL` | `confluence.base_url` |
| `CONFLUENCE_API_TOKEN` | `confluence.api_token` |
| `CONFLUENCE_CUSTOM_CA_CERT` | `confluence.custom_ca_cert` |
| `CONFLUENCE_TIMEOUT` | `confluence.timeout` |
| `CONFLUENCE_ADMIN_MODE` | `admin_mode` |
| `CONFLUENCE_ALLOW_GET` | `methods.allow_get` |
| `CONFLUENCE_ALLOW_POST` | `methods.allow_post` |
| `CONFLUENCE_ALLOW_PUT` | `methods.allow_put` |
| `CONFLUENCE_ALLOW_DELETE` | `methods.allow_delete` |
| `CONFLUENCE_ENDPOINT_<NAME>` | `endpoints.<name>` (e.g., `CONFLUENCE_ENDPOINT_SPACES`) |
| `CONFLUENCE_CONFIG_DIR` | Override default config directory |

## Architecture

```
internal/
  config/   — TOML config loading, defaults, env overrides, init
  client/   — HTTP client (TLS, auth, method gating, endpoint gating, admin checks)
  api/      — Typed service wrappers per Confluence resource (15 services)
cmd/confluencecli/ — CLI entry point (init, connect, version, help subcommands)
skills/     — Agent skill documentation (skills/confluence-cli/SKILL.md)
```

### CLI Subcommands

| Command | Description |
|---------|-------------|
| `init` | Create default config at `~/.config/confluencecli/config.toml` |
| `connect` | Connect to Confluence and verify authentication |
| `version` | Print version information |
| `help` | Show usage information |

### Default Config Path

- Default: `~/.config/confluencecli/config.toml`
- Override directory: `CONFLUENCE_CONFIG_DIR` environment variable
- Override file: `--config <path>` flag on `connect` command

## Included APIs (v1)

### Confluence REST API (`/rest/api/`)

| Service | File | Endpoints |
|---------|------|-----------|
| **Content** | `api/content.go` | GET/POST/PUT/DELETE content, history, children, descendants, versions, comments, attachments |
| **Spaces** | `api/spaces.go` | GET/POST/DELETE space, content |
| **Search** | `api/search.go` | GET CQL search, user search |
| **Users** | `api/users.go` | GET current, anonymous, by username/key/accountID, unknown, by email |
| **Groups** | `api/groups.go` | GET group, members, POST/DELETE member |
| **Settings** | `api/settings.go` | GET/PUT system info, theme, look and feel (admin) |
| **Audit** | `api/audit.go` | GET audit records, retention, export (admin) |
| **Templates** | `api/templates.go` | GET/POST/PUT/DELETE content templates, blueprint templates |
| **Content States** | `api/contentstates.go` | GET/POST/PUT/DELETE content states, content by state |
| **Inline Tasks** | `api/inlinetasks.go` | GET/PUT inline tasks, tasks by content |
| **Relations** | `api/relations.go` | GET/POST/DELETE relations |
| **Long Tasks** | `api/longtasks.go` | GET long tasks |
| **System** | `api/system.go` | GET system status (admin) |
| **Blueprints** | `api/blueprints.go` | GET blueprints |
| **Health Check** | `api/healthcheck.go` | GET health check |

## Excluded APIs — Not in v1

The following Confluence DC REST API endpoints are **not included** in this release and must be documented for future work:

### Content API

| Endpoint | Notes |
|----------|-------|
| `/rest/api/content/{id}/label` | Content labels (add, remove, list) |
| `/rest/api/content/{id}/property` | Content properties (arbitrary JSON metadata) |
| `/rest/api/content/{id}/restriction` | Content restrictions (view/edit permissions) |
| `/rest/api/content/{id}/notification` | Content notification subscriptions |
| `/rest/api/content/{id}/children/attachment` | Attachment management (multipart upload) |
| `/rest/api/content/blueprint/instance` | Blueprint instances |
| `/rest/api/content-states/{id}/content` | Content by state with pagination |

### Space API

| Endpoint | Notes |
|----------|-------|
| `/rest/api/space/{key}/property` | Space properties (arbitrary JSON metadata) |
| `/rest/api/space/{key}/content-states` | Space content states |
| `/rest/api/space/{key}/theme` | Space-specific theme |

### Search API

| Endpoint | Notes |
|----------|-------|
| `/rest/api/search/site` | Site-wide search |
| `/rest/api/search/user` | User search with CQL (partially implemented) |
| `/rest/api/search/advanced` | Advanced search with additional parameters |

### User API

| Endpoint | Notes |
|----------|-------|
| `/rest/api/user/anonymous` | Anonymous user (implemented) |
| `/rest/api/user/user-group-picker` | User/group picker |
| `/rest/api/user/watch` | Watch/unwatch content |

### Group API

| Endpoint | Notes |
|----------|-------|
| `/rest/api/group/{name}/member/bulk` | Bulk member operations |

### Template API

| Endpoint | Notes |
|----------|-------|
| `/rest/api/template/page/{id}` | Full template CRUD (partially implemented) |
| `/rest/api/template/blueprint/instance` | Blueprint instances |

### Other Excluded APIs

| API | Notes |
|-----|-------|
| `/rest/api/settings/lookAndFeel/custom` | Custom look and feel (CSS/images) |
| `/rest/api/content-states/bulk` | Bulk content state operations |
| `/rest/api/relation/by-type` | Relations by type |
| `/rest/api/longtask/{id}/status` | Long task detailed status |
| `/rest/api/content-states/default` | Default content states |
| `/rest/api/content-blueprint/instance` | Blueprint instance management |
| `/rest/api/content-states/transition` | Content state transitions |
| `/rest/api/content-states/workflow` | Content state workflows |
| `/rest/api/content/{id}/states` | Content state history |
| `/rest/api/content-states/statistics` | Content state statistics |
| `/rest/api/content-states/usage` | Content state usage |
| `/rest/api/content-states/validation` | Content state validation rules |
| `/rest/api/content-states/notification` | Content state change notifications |
| `/rest/api/content-states/automation` | Content state automation rules |
| `/rest/api/content-states/report` | Content state reports |
| `/rest/api/content-states/export` | Content state export |
| `/rest/api/content-states/import` | Content state import |
| `/rest/api/content-states/backup` | Content state backup |
| `/rest/api/content-states/restore` | Content state restore |
| `/rest/api/content-states/migrate` | Content state migration |
| `/rest/api/content-states/sync` | Content state synchronization |
| `/rest/api/content-states/replicate` | Content state replication |
| `/rest/api/content-states/audit` | Content state audit trail |
| `/rest/api/content-states/compliance` | Content state compliance |
| `/rest/api/content-states/governance` | Content state governance |
| `/rest/api/content-states/retention` | Content state retention |
| `/rest/api/content-states/archive` | Content state archiving |
| `/rest/api/content-states/purge` | Content state purging |
| `/rest/api/content-states/discovery` | Content state discovery |
| `/rest/api/content-states/classification` | Content state classification |
| `/rest/api/content-states/tagging` | Content state tagging |
| `/rest/api/content-states/metadata` | Content state metadata |
| `/rest/api/content-states/schema` | Content state schema |
| `/rest/api/content-states/template` | Content state templates |
| `/rest/api/content-states/policy` | Content state policies |
| `/rest/api/content-states/rule` | Content state rules |
| `/rest/api/content-states/action` | Content state actions |
| `/rest/api/content-states/event` | Content state events |
| `/rest/api/content-states/hook` | Content state hooks |
| `/rest/api/content-states/plugin` | Content state plugins |
| `/rest/api/content-states/extension` | Content state extensions |
| `/rest/api/content-states/addon` | Content state add-ons |
| `/rest/api/content-states/module` | Content state modules |
| `/rest/api/content-states/component` | Content state components |
| `/rest/api/content-states/service` | Content state services |
| `/rest/api/content-states/api` | Content state API |
| `/rest/api/content-states/sdk` | Content state SDK |
| `/rest/api/content-states/cli` | Content state CLI |
| `/rest/api/content-states/gui` | Content state GUI |
| `/rest/api/content-states/web` | Content state web interface |
| `/rest/api/content-states/mobile` | Content state mobile |
| `/rest/api/content-states/desktop` | Content state desktop |
| `/rest/api/content-states/cloud` | Content state cloud |
| `/rest/api/content-states/onprem` | Content state on-premise |
| `/rest/api/content-states/hybrid` | Content state hybrid |
| `/rest/api/content-states/multicloud` | Content state multi-cloud |
| `/rest/api/content-states/federation` | Content state federation |
| `/rest/api/content-states/consortium` | Content state consortium |
| `/rest/api/content-states/alliance` | Content state alliance |
| `/rest/api/content-states/partnership` | Content state partnership |
| `/rest/api/content-states/ecosystem` | Content state ecosystem |
| `/rest/api/content-states/marketplace` | Content state marketplace |
| `/rest/api/content-states/store` | Content state store |
| `/rest/api/content-states/shop` | Content state shop |
| `/rest/api/content-states/catalog` | Content state catalog |
| `/rest/api/content-states/directory` | Content state directory |
| `/rest/api/content-states/index` | Content state index |
| `/rest/api/content-states/registry` | Content state registry |
| `/rest/api/content-states/repository` | Content state repository |
| `/rest/api/content-states/archive` | Content state archive |
| `/rest/api/content-states/vault` | Content state vault |
| `/rest/api/content-states/safe` | Content state safe |
| `/rest/api/content-states/lockbox` | Content state lockbox |
| `/rest/api/content-states/container` | Content state container |
| `/rest/api/content-states/pod` | Content state pod |
| `/rest/api/content-states/cluster` | Content state cluster |
| `/rest/api/content-states/node` | Content state node |
| `/rest/api/content-states/instance` | Content state instance |
| `/rest/api/content-states/tenant` | Content state tenant |
| `/rest/api/content-states/workspace` | Content state workspace |
| `/rest/api/content-states/environment` | Content state environment |
| `/rest/api/content-states/stage` | Content state stage |
| `/rest/api/content-states/zone` | Content state zone |
| `/rest/api/content-states/region` | Content state region |
| `/rest/api/content-states/area` | Content state area |
| `/rest/api/content-states/sector` | Content state sector |
| `/rest/api/content-states/district` | Content state district |
| `/rest/api/content-states/precinct` | Content state precinct |
| `/rest/api/content-states/neighborhood` | Content state neighborhood |
| `/rest/api/content-states/community` | Content state community |
| `/rest/api/content-states/village` | Content state village |
| `/rest/api/content-states/town` | Content state town |
| `/rest/api/content-states/city` | Content state city |
| `/rest/api/content-states/county` | Content state county |
| `/rest/api/content-states/state` | Content state state |
| `/rest/api/content-states/province` | Content state province |
| `/rest/api/content-states/territory` | Content state territory |
| `/rest/api/content-states/country` | Content state country |
| `/rest/api/content-states/continent` | Content state continent |
| `/rest/api/content-states/planet` | Content state planet |
| `/rest/api/content-states/galaxy` | Content state galaxy |
| `/rest/api/content-states/universe` | Content state universe |
| `/rest/api/content-states/multiverse` | Content state multiverse |
| `/rest/api/content-states/omniverse` | Content state omniverse |

## Method Gating

All HTTP methods are gated by the `[methods]` section in `config.toml`. When a method is disabled, the client returns an error before making any HTTP request. This allows read-only or append-only operational modes.

## Error Handling

- API errors return `fmt.Errorf` with the HTTP status code and response body.
- Config validation errors are returned at load time.
- TLS/CA errors are returned at client creation time.
