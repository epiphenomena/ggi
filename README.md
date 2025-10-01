# GGI (Go Generated Interfaces)
Go lang library that combines a simple static site generator with a CGI-based editor

The fastest and easiest site to serve is static html. But it's complicated for normal users to edit and does not support templating elements that are repeated within a page or across pages. 

SSGs are a solution to the templating problem, but require learning and conforming to their idiosyncrasies. 

There are admin UI's (notably Wordpress) that attempt to make it easier for users to edit content without needing to understand the tech stack.

However, both SSGs and the admin UIs necessarily develop a great deal of complexity in order to support a wide range of use cases.

LLMs make creating customized websites and customized admin UIs easy. The goal of this project is to create a simple go lang based library to support an LLM driven website creation and maintenance.

The idea is to import this library into a new website project, add the needed customizations for that project, and then compile to single binary that acts as a CGI script.

The CGI script supports editing source files and then generating the resulting static html.

## Features

- **CGI Detection**: Automatically detects if running as CGI script or in development mode
- **Authentication**: Uses a secret token system for form submissions
- **Template System**: Uses Go templates with base layout including header, content, footer, CSS and JS blocks
- **Modal Interface**: Built-in modal for settings and secret key management with localStorage integration
- **Markdown Support**: Edit and render markdown content
- **Form-based Data Management**: Generate forms from Go structs and save as JSON
- **CSS Customization**: Dynamic CSS generation from configuration
- **Static Site Generation**: Compile templates to static HTML
- **Security**: File access limited to current directory and subdirectories only

## Getting Started

### Installation

```bash
go get ggi
```

### Basic Usage

```go
package main

import (
    "ggi/pkg/ggi"
)

func main() {
    // Set your secret key
    ggi.SecretKey = "your_secret_key_here"

    // The library will automatically detect if it's running as CGI or in development mode
    if ggi.IsCGI() {
        // Handle as CGI script
        handleCGI()
    } else {
        // Start development server
        handleDevServer()
    }
}
```

### Example

See the [examples/basic-site](examples/basic-site) directory for a complete working example.

## Structure

- `_source/templates/` - HTML templates for the site
- `_source/resources/` - Static assets like CSS, JS, images  
- `_source/data/` - JSON files for structured data
- `_source/markdown/` - Markdown content files

## Security

- All file operations are restricted to the current directory and subdirectories
- Form submissions require a secret token for authentication
- The .htaccess file generated blocks access to source directories

## Development vs Production

The library automatically detects if it's running as a CGI script or in development mode:
- In development mode: runs as a web server for easy testing
- As CGI: responds to POST requests with token authentication
