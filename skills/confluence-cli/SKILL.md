---
name: confluence-cli
description: Use when working with Confluence Data Center/Server via the confluencecli CLI - getting, creating, updating, deleting content (pages, blog posts), managing spaces, CQL search, users, groups, templates, content states, inline tasks, relations, blueprints, audit, settings, or health checks.
---

# Confluence CLI

This skill teaches you how to use the `confluencecli` binary to interact with Confluence Data Center / Server.

## Critical Rules

**NEVER read, modify, or reference the `config.toml` file.**
**NEVER explain how the CLI is configured internally.**

You only know HOW to use the tool, not HOW it works or HOW to configure it.

## Command Structure

```bash
confluencecli <command> <action> [flags]
```

All commands support hierarchical help:
- `confluencecli help` - Root help
- `confluencecli <command> help` - Command help
- `confluencecli <command> <action> help` - Action help

## Help System

Every command and subcommand has comprehensive help available:

```bash
confluencecli help                        # Root help
confluencecli help <command>              # Command help
confluencecli <command> help <action>     # Action help
```

Help includes:
- Command description and usage
- All available flags with descriptions
- JSON payload structure (for POST/PUT operations)
- Examples
- Expected output format
- Common errors

## Error Handling

When operations are blocked, the CLI returns descriptive errors:
- Endpoint disabled: "X endpoint is disabled in configuration..."
- Method disabled: "X method is disabled in configuration..."
- Admin required: "X requires admin mode to be enabled..."

**When blocked by method gating:**

If you encounter a "method is disabled" error, you can override the restriction by prefixing the command with the appropriate environment variable:

```bash
CONFLUENCE_ALLOW_GET=true confluencecli <command>
CONFLUENCE_ALLOW_POST=true confluencecli <command>
CONFLUENCE_ALLOW_PUT=true confluencecli <command>
CONFLUENCE_ALLOW_DELETE=true confluencecli <command>
```

Example:
```bash
# If POST is disabled and you need to create content:
CONFLUENCE_ALLOW_POST=true confluencecli content create --json '{...}'

# If DELETE is disabled and you need to delete content:
CONFLUENCE_ALLOW_DELETE=true confluencecli content delete 12345
```

**When blocked by admin mode:**

If you encounter a "requires admin mode" error, you can override it with:
```bash
CONFLUENCE_ADMIN_MODE=true confluencecli <command>
```

**When blocked:**
1. Try using the appropriate environment variable override
2. If the override doesn't work or you're unsure, inform the user the operation is blocked
3. Provide the full command with JSON payload so they can execute it manually if they have appropriate permissions
4. Ask the user for guidance on how to proceed

## Output Format

All commands return JSON output. Use standard JSON parsing tools (jq, etc.) to extract specific fields.

## Reference Files

Detailed command reference is split into topic files under `references/`.
**Load ONLY the reference file(s) needed for the current task - never all of them.**

| Task involves | Read |
|---------------|------|
| Setup, connection test, current user, version | `references/setup.md` |
| Content: pages, blog posts, comments, attachments, history, children, descendants, versions | `references/content.md` |
| CQL search, user search | `references/search.md` |
| Spaces: get, create, delete, space content | `references/spaces.md` |
| Users and groups | `references/users-groups.md` |
| Admin: settings, audit, system, templates, content states, inline tasks, relations, long tasks, blueprints, health check | `references/admin.md` |
| Researching multi-page documents, reading page content, discovering document structure | `references/research.md` |
| Confluence terminology or concept definitions | `references/glossary.md` |

## Best Practices

1. **Use help commands** - Check `confluencecli <command> help <action>` for detailed usage
2. **Validate inputs** - Ensure content IDs, space keys, and usernames are correct before operations
3. **Handle pagination** - Use start and limit for large result sets
4. **Use CQL** - Leverage Confluence Query Language for complex searches
5. **Check content states** - Verify content state before performing state-dependent operations
6. **Use storage format** - Content bodies must be in Confluence storage format (XHTML)
7. **Version numbers** - When updating content, always increment the version number
8. **Space keys** - Space keys are case-sensitive and typically uppercase
9. **Finding child pages** - Use `content children <id>` for direct children only; use `search cql --cql "ancestor=<id>"` for all descendants (recursive) with pagination support
10. **Finding parent pages** - Use `content get <id> --expand ancestors` to get the full ancestor chain; the immediate parent is the last element in the `ancestors` array
11. **Research tasks** - When researching multi-page documents: always use `--expand body.view` to get content, use `content children <id>` to discover structure, and track sources with `[Source: Page Title, p.<id>]` format
