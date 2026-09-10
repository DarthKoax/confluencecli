# Confluence CLI Agent Skill

## Overview

This skill teaches you how to use the `confluencecli` binary to interact with Confluence Data Center / Server.

## Critical Rules

**NEVER read, modify, or reference the `config.toml` file.**
**NEVER discuss environment variables or configuration settings.**
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

## Error Handling

When operations are blocked, the CLI returns descriptive errors:
- Endpoint disabled: "X endpoint is disabled in configuration..."
- Method disabled: "X method is disabled in configuration..."
- Admin required: "X requires admin mode to be enabled..."

**When blocked:**
1. Do NOT retry or work around
2. Inform the user the operation is blocked
3. Provide the full command with JSON payload for manual execution
4. Ask for guidance

## Topic Files

Detailed command reference is split into topic files:
- `01-setup.md` - Setup and connection
- `02-content.md` - Content management (pages, blog posts, comments)
- `03-search.md` - CQL search and user search
- `04-spaces.md` - Space management
- `05-users-groups.md` - User and group management
- `06-admin.md` - Admin operations (settings, audit, system, templates, content states, inline tasks, relations, long tasks, blueprints)
- `glossary.md` - Confluence terminology

Load only the topic files you need for your current task.

## Output Format

All commands return JSON output. Use standard JSON parsing tools (jq, etc.) to extract specific fields.

## Best Practices

1. **Use help commands** - Check `confluencecli <command> help <action>` for detailed usage
2. **Validate inputs** - Ensure content IDs, space keys, and usernames are correct
3. **Handle pagination** - Use start and limit for large result sets
4. **Use CQL** - Leverage Confluence Query Language for complex searches
5. **Check content states** - Verify content state before performing state-dependent operations
6. **Use storage format** - Content bodies must be in Confluence storage format (XHTML)
7. **Version numbers** - When updating content, always increment the version number
