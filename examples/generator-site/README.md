# GGI Generator Site Example

This example demonstrates the GGI generator pattern where a professional creates a custom package that imports GGI to build a specialized static site with content management.

## Overview

In this pattern:
- The professional creates a "generator" package (this example)
- The generator imports GGI and configures it for their specific site needs
- Running the generator creates all necessary files for deployment
- The end result is a static site with a custom admin UI

## Files Generated

When this generator runs, it creates:

### Source Directory (`_source/`)
- `markdown/` - Editable markdown content
- `data/` - JSON data files that can be edited through forms
- `media/` - Media files

### Public Site (`_output/`)
- Static HTML pages for the public site
- CSS, JS, and other static assets

### Admin UI (`_output/admin/`)
- Administrative interface for editing content
- Auto-generated forms for different content types
- Management pages for different content categories

### CGI Script (`_output/admin.cgi`)
- Handles form submissions from the admin UI
- Saves updated content to source files
- Updates static pages as needed

## Usage

1. Run the generator: `go run main.go`
2. Upload the `_output` directory to your web server
3. Set up .htaccess to password-protect the admin section
4. Users can now access the admin UI to edit content

## Customization

Professionals can customize:
- Content types and how they're handled
- Templates for public pages
- How content is saved and processed
- The build process itself