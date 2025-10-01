# GGI (Go Generated Interfaces)
A foundation library for building custom static sites with auto-generated admin interfaces

## Overview

GGI is a Go library that provides a foundation for creating custom static websites with content management capabilities. Rather than being a generic CMS, GGI serves as a toolkit that professionals can import into their own "generator" packages to build highly customized static sites with personalized admin interfaces.

The core idea is that a professional creates a separate Go package (the "generator") that imports GGI and defines:
- The site's structure and templates
- Which content areas should be editable by end users
- How different content types should be handled and displayed

The generator package can then be run to produce:
- An initial `_source/` directory containing editable content
- Public static HTML pages based on the templates and content
- Admin UI static HTML pages for content editing
- A CGI script that handles content updates from the admin UI

## Architecture

### The Generator Pattern

```
[Generator Package] 
    ├── imports GGI library
    ├── defines content types and editable areas
    ├── specifies templates and styling
    └── runs build to generate site

[Generated Output]
    ├── _source/           # Editable content files
    ├── public/            # Static public site
    ├── admin/             # Static admin UI  
    └── admin.cgi          # CGI script for handling updates
```

### Content Type System

GGI provides a content type registration system where professionals can register different types of editable content:

- **Markdown content**: Text areas for users to write in Markdown format
- **JSON data structures**: Structured data edited through forms
- **Media files**: Images and other media that can be uploaded/replaced
- **Custom content types**: Any specialized content with custom handlers

### Auto-Generated Admin UI

Based on registered content types, GGI automatically generates:

- Admin pages for each content type
- Forms with appropriate input fields (text areas for Markdown, form fields for JSON, upload fields for media)
- Management interfaces for adding, editing, and deleting content

## Key Features

- **Content Type Registration**: Define custom content types with appropriate handlers
- **Auto-Generated Admin UI**: Forms and interfaces created automatically based on content definitions
- **Static Site Generation**: Compile templates and content into static HTML
- **CGI Integration**: Automatic CGI script generation for content updates
- **Extensible Design**: Hooks for custom build processes and content handling
- **Security**: File access limited to appropriate directories only
- **Fast Static Output**: Generated sites are pure static HTML for optimal performance

## How It Works

1. **Professional creates a generator package** that imports GGI and defines their site
2. **Content types are registered** with appropriate handlers and templates
3. **Build process is run** to generate the initial site structure
4. **Generated site is deployed** to a web server with .htaccess for security
5. **End users access admin UI** to edit pre-defined content areas
6. **CGI script processes updates** and regenerates static pages

## Benefits

- **Lightweight**: No heavy CMS overhead - just static HTML
- **Fast**: Pure static output for optimal performance
- **Secure**: Content editing through controlled CGI interface
- **Customizable**: Highly tailored to specific site needs
- **Maintainable**: Professionals can make structural changes and regenerate

## Use Cases

GGI is ideal for creating websites where:
- Static HTML performance is important
- Specific content areas need to be editable by non-technical users
- A custom, lightweight solution is preferred over a generic CMS
- An LLM can be used to generate the specialized code with GGI as a foundation

## Getting Started

Professionals should create their own generator package that imports GGI, defines their content structure, and uses the build system to generate their custom site.
