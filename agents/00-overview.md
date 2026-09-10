# Confluence CLI Agent Skill - Overview

## Critical Rules

**NEVER read, modify, or reference the `config.toml` file.**
**NEVER discuss environment variables or configuration settings.**
**NEVER explain how the CLI is configured internally.**

You only know HOW to use the tool, not HOW it works or HOW to configure it.

## Available Commands

The CLI accepts JSON payloads via command-line arguments. All output is JSON.

### Global Flags

- `--config <path>`: Path to config file (uses default if not specified)
- `--help`: Show help information

### Command Structure

```bash
confluencecli <resource> <action> [flags] [json-payload]
```

## Error Handling

When operations are blocked due to configuration restrictions, the CLI will return descriptive error messages. These messages indicate:
- Which endpoint or method is disabled
- That the operation cannot be performed
- That you should not attempt to bypass the restriction
- That you should ask the user for advice if required

**When you receive a blocked error:**
1. Do NOT attempt to retry or work around the restriction
2. Inform the user that the operation is blocked
3. Provide the full command with JSON payload so they can execute it manually if they have appropriate permissions
4. Ask the user for guidance on how to proceed

## Help System

Every command and subcommand has comprehensive help available:

```bash
confluencecli help                                    # Root help
confluencecli help <command>                          # Command help
confluencecli <command> help <action>                 # Action help
```

Help includes:
- Command description and usage
- All available flags with descriptions
- JSON payload structure (for POST/PUT operations)
- Examples
- Expected output format
- Common errors

## Topic Files

This skill is split into multiple topic files for better context management:

- `01-setup.md` - Setup and connection verification
- `02-content.md` - Content management (pages, blog posts, comments, attachments)
- `03-search.md` - CQL search and user search
- `04-spaces.md` - Space management
- `05-users-groups.md` - User and group management
- `06-admin.md` - Admin operations (settings, audit, system)

Only load the topic files you need for your current task.

## Output Format

All commands return JSON output. Use standard JSON parsing tools (jq, etc.) to extract specific fields.

## Best Practices

1. **Use help commands** - Always check `confluencecli <command> help <action>` for detailed usage
2. **Validate inputs** - Ensure content IDs, space keys, and usernames are correct before operations
3. **Handle pagination** - Use start and limit for large result sets
4. **Use CQL for complex searches** - Leverage the full power of Confluence Query Language
5. **Check content states** - Always verify content state before performing state-dependent operations
