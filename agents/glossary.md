# Confluence CLI - Glossary of Terms

This glossary defines all Confluence terms and concepts covered by the `confluencecli` tool.

## Core Concepts

### Content
A unit of information in Confluence. Content can be pages, blog posts, or comments. Each content item has a unique ID and contains a title, body (in storage format), and metadata like space, version, and status.

### Page
The primary content type in Confluence. Pages are organized hierarchically within spaces and can have parent-child relationships. Pages support rich text, images, tables, macros, and other content types.

### Blog Post
A time-stamped content type used for announcements, updates, and news. Blog posts are organized by date within a space and are typically displayed in reverse chronological order.

### Comment
A note attached to a page or blog post. Comments are used for discussion, feedback, and communication about content. Comments can be inline (attached to specific text) or footer comments.

### Space
A container for organizing related content. Spaces have a unique key (e.g., DEV) and contain pages, blog posts, and other content. Spaces can be global (team-wide) or personal (user-specific).

### Space Key
The short code identifying a space (e.g., DEV, TEAM). Space keys are uppercase and used in URLs and CQL queries.

### Content ID
The unique numeric identifier for a content item. Content IDs are unique across all spaces but are less user-friendly than page titles or URLs.

## Content Structure

### Ancestor
A parent page in the content hierarchy. Pages can have multiple ancestors, forming a tree structure within a space.

### Children
Direct child pages of a parent page. Children are one level below their parent in the hierarchy.

### Descendants
All pages below a parent page in the hierarchy, including children, grandchildren, and so on.

### Body
The main content of a page or blog post. The body is stored in Confluence storage format (XHTML-based) and can be rendered in view format for display.

### Storage Format
The XHTML-based format used to store content in Confluence. Storage format includes macros, links, and other Confluence-specific elements encoded as HTML.

### View Format
The rendered HTML format of content, suitable for display in a web browser. View format is generated from storage format.

## Versioning

### Version
A snapshot of content at a specific point in time. Each time content is updated, a new version is created. Versions track the number, author, timestamp, and change message.

### Version Number
The sequential number of a content version. Version numbers start at 1 and increment with each update.

### Minor Edit
A flag indicating that a content update is minor (e.g., typo fix). Minor edits do not trigger notifications to watchers.

### Content History
The complete version history of a content item, including all versions, authors, and change messages.

## Content States

### Content State
A label indicating the lifecycle status of content (e.g., Draft, In Review, Published). Content states help track the progress of content through a workflow.

### Final State
A content state that indicates the content is complete and no further changes are expected. Content in a final state is considered finalized.

## Templates

### Content Template
A pre-defined page layout that users can use as a starting point for new pages. Templates include pre-filled content, macros, and formatting.

### Blueprint Template
A specialized template that creates content with additional logic and configuration. Blueprints can include custom workflows, forms, and automated content generation.

### Template ID
The unique identifier for a template. Template IDs are used with template-related commands.

## Inline Tasks

### Inline Task
A to-do item embedded within page content. Inline tasks can be assigned to users, have due dates, and be marked as complete or incomplete.

### Task Status
The current state of an inline task (complete or incomplete).

## Search

### CQL (Confluence Query Language)
A powerful query language used to search for content in Confluence. CQL supports operators like =, !=, ~, and functions like currentUser(), now(), etc.

### CQL Query
A search expression written in CQL. CQL queries can filter by space, type, title, creator, label, date, and many other fields.

## People & Permissions

### User
A person who can access Confluence. Users have usernames, email addresses, and can create, edit, and comment on content.

### Group
A collection of users. Groups are used to manage permissions and notifications. Users can belong to multiple groups.

### Creator
The user who originally created a content item.

### Contributor
A user who has edited or added content to a page.

### Administrator (Admin)
A user with elevated permissions. Admin mode in confluencecli enables access to administrative operations like managing spaces, groups, settings, and audit logs.

## Communication

### Attachment
A file attached to a content item. Attachments can be documents, images, spreadsheets, or any other file type. (Note: Attachment upload is not yet implemented in confluencecli)

### Watch
A subscription to receive notifications about changes to a content item. Users can watch pages to be notified of updates.

## Relations

### Relation
A connection between two entities in Confluence. Relations define relationships like "depends on", "is related to", etc.

### Relation Name
The type of relation (e.g., depends-on, relates-to). Relation names define the nature of the relationship between source and target entities.

### Source/Target
The two entities connected by a relation. The source is the originating entity and the target is the destination entity.

## System & Configuration

### System Info
Information about the Confluence instance, including build date, commit hash, and cloud ID.

### Theme
A visual style applied to the Confluence instance. Themes control the overall look and feel of the UI.

### Look and Feel
Customization settings for the Confluence UI, including colors, backgrounds, and layout. Look and feel can be set globally or per-space.

### Audit Log
A record of all administrative actions performed in Confluence. Audit logs track who did what and when, including content creation, deletion, and permission changes.

### Audit Retention
The period for which audit log records are kept before being automatically deleted.

### Health Check
A diagnostic check of the Confluence instance's subsystems (database, search index, etc.). Health checks report the status of each component.

### Long Task
A background operation that takes time to complete (e.g., content export, reindex). Long tasks report progress and can be monitored.

### Blueprint
A content generation framework that creates pages with predefined structure and logic. Blueprints extend templates with custom workflows.

## Common Abbreviations

- **CQL**: Confluence Query Language
- **API**: Application Programming Interface
- **REST**: Representational State Transfer (type of web API)
- **DC**: Data Center (Confluence Data Center deployment)
- **TOML**: Tom's Obvious, Minimal Language (configuration file format)
- **TLS**: Transport Layer Security (encryption protocol)
- **CA**: Certificate Authority
- **XHTML**: Extensible HyperText Markup Language

## Command-Specific Terms

### Content ID
The unique numeric identifier for a content item. Content IDs are used with content-related commands.

### Space Key
The short code identifying a space (e.g., DEV). Space keys are used with space-related commands and CQL queries.

### Template ID
The unique identifier for a content or blueprint template. Template IDs are used with template-related commands.

### State ID
The unique identifier for a content state. State IDs are used with content state-related commands.

### Task ID
The unique identifier for an inline task or long task. Task IDs are used with task-related commands.

### Group Name
The name of a user group. Group names are used with group-related commands.

### Username
The login name of a user. Usernames are used with user-related commands.

### Theme Key
The identifier for a theme (e.g., com.atlassian.confluence.theme:default). Theme keys are used with theme-related commands.

### Relation Name
The type of relation between entities. Relation names are used with relation-related commands.

### Entity Type
The type of entity in a relation (e.g., content, user, space). Entity types are used with relation-related commands.

### Entity Key
The identifier of an entity in a relation. Entity keys are used with relation-related commands.
