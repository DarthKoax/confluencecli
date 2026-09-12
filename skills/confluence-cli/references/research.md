# Confluence CLI - Research Tasks

## Overview

This guide covers how to systematically research multi-page Confluence documents using confluencecli. Research tasks involve discovering document structure, reading content from multiple pages, and compiling summaries with source tracking.

## Research Workflow

### Step 1: Understand the Starting Point

Extract the page ID from the URL or user input:
- URL format: `http://localhost:8090/spaces/ds/pages/131178/Page+Title`
- Page ID: `131178` (the number in the URL path)

### Step 2: Discover Available Commands

```bash
confluencecli help
confluencecli help content
```

This reveals the available subcommands and their syntax.

### Step 3: Read the Top-Level Page

```bash
confluencecli content get <page-id> --expand body.view
```

**Important:** The `--expand body.view` flag is essential. Without it, only metadata (title, type, status) is returned, not the actual content.

**Example:**
```bash
confluencecli content get 131178 --expand body.view
```

**Output structure:**
```json
{
  "id": "131178",
  "type": "page",
  "title": "Document Title",
  "body": {
    "view": {
      "value": "<p>HTML content here...</p>",
      "representation": "view"
    }
  }
}
```

The actual content is in `body.view.value` as HTML.

### Step 4: Discover Child Pages

```bash
confluencecli content children <parent-page-id>
```

**Example:**
```bash
confluencecli content children 131178
```

**Output structure:**
```json
{
  "results": [
    {
      "id": "131179",
      "type": "page",
      "title": "Child Page 1",
      "_links": {
        "webui": "/spaces/ds/pages/131179/Child+Page+1"
      }
    },
    {
      "id": "131180",
      "type": "page",
      "title": "Child Page 2"
    }
  ],
  "start": 0,
  "limit": 25,
  "size": 2
}
```

### Step 5: Read All Child Pages

For each child page discovered, read its content:

```bash
confluencecli content get <child-id> --expand body.view
```

**Tip:** If you have multiple child pages, you can read them in parallel since they are independent operations.

### Step 6: Check for Grandchildren (Optional)

If the document might have deeper nesting, check each child for its own children:

```bash
confluencecli content children <child-id>
```

If `size` is 0, there are no grandchildren for that page.

### Step 7: Compile Summary with Source Tracking

As you read each page, track the source of each piece of information:

**Format:** `[Source: Page Title, p.<page-id>]`

**Example:**
```
The document was accidentally invented in 1698 when a cat named Whiskers knocked over a bottle.
[Source: History & Origins, p.131179]
```

## Complete Research Example

**Scenario:** Research a document starting from page ID 131178.

```bash
# Step 1: Read the top-level page
confluencecli content get 131178 --expand body.view

# Step 2: Discover child pages
confluencecli content children 131178
# Returns: 131179, 131182, 131184, 131185

# Step 3: Read all child pages (can be done in parallel)
confluencecli content get 131179 --expand body.view
confluencecli content get 131182 --expand body.view
confluencecli content get 131184 --expand body.view
confluencecli content get 131185 --expand body.view

# Step 4: Check for grandchildren
confluencecli content children 131179
confluencecli content children 131182
confluencecli content children 131184
confluencecli content children 131185
# All return size: 0, so no deeper nesting
```

## Useful Expand Options for Research

| Expand Option | Description |
|---------------|-------------|
| `body.view` | Rendered HTML content (recommended for reading) |
| `body.storage` | Confluence storage format (XHTML) |
| `version` | Version information (number, author, date) |
| `space` | Space details (key, name) |
| `ancestors` | Parent page chain |
| `children.page` | Direct child pages (alternative to `content children`) |

**Multiple expand fields:**
```bash
confluencecli content get 131178 --expand "body.view,version,space"
```

## Alternative: CQL Search for Finding Related Pages

If you need to find pages related to a topic (not just children), use CQL search:

```bash
# Find all pages in the same space
confluencecli search cql --cql "space=ds"

# Find pages with specific text
confluencecli search cql --cql "text~'stream trains'"

# Find all descendants of a page (recursive)
confluencecli search cql --cql "ancestor=131178"
```

## Best Practices for Research Tasks

1. **Always use `--expand body.view`** when reading page content
2. **Track sources** - Note which page each fact came from
3. **Check for children** - Documents often have multi-level structure
4. **Use parallel reads** - Child pages can be read simultaneously
5. **Handle HTML** - Content is returned as HTML; parse as needed
6. **Start broad, then narrow** - Read the top-level page first to understand structure
7. **Use CQL for cross-page searches** - When you need to find pages by topic, not just hierarchy

## Common Patterns

### Pattern 1: Simple Document (Parent + Children)
```
Parent Page
├── Child Page 1
├── Child Page 2
└── Child Page 3
```
**Approach:** Read parent, get children, read each child.

### Pattern 2: Deep Hierarchy
```
Parent Page
├── Child Page 1
│   ├── Grandchild 1.1
│   └── Grandchild 1.2
└── Child Page 2
    └── Grandchild 2.1
```
**Approach:** Read parent, get children, for each child check for grandchildren, read all pages.

### Pattern 3: Related Pages (Not Hierarchical)
```
Page A (related to Page B via links)
Page B (related to Page C via links)
Page C
```
**Approach:** Use CQL search to find related pages, or follow links in content.

## Error Handling

- **Page not found:** Check the page ID is correct
- **Empty content:** Try `--expand body.storage` instead of `body.view`
- **Permission denied:** The page may be restricted; inform the user
- **No children:** The `size` field will be 0; this is normal for leaf pages
