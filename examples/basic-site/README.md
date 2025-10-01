# GGI Basic Site Example

This example demonstrates how to use the GGI library to create a simple website with CGI-based editing capabilities.

## Features Demonstrated

1. **CGI Detection**: The application automatically detects if it's running as a CGI script or in development mode
2. **Public/Admin Structure**: Shows separation between public site and admin UI
3. **Markdown Editing**: Provides an interface to edit and save markdown content
4. **Form-based Data Management**: Creates forms from Go structs and saves data as JSON
5. **CSS/JS Helpers**: Shows how to load CSS and JS files using helper functions

## Running the Example

### Development Mode
```bash
cd examples/basic-site
go run main.go
```

The development server will start on port 8080.

### As CGI Script
1. Compile the binary:
```bash
go build -o ggi.cgi main.go
```

2. Place the `ggi.cgi` file in your web server's CGI directory and configure the .htaccess as needed.

3. Set up HTTP Basic Authentication for the admin section using the provided .htaccess file.

## Structure

- `_source/public/templates/` - HTML templates for the public site
- `_source/admin/templates/` - HTML templates for the admin UI
- `_source/public/resources/` - Public static assets like CSS, JS, images
- `_source/admin/resources/` - Admin static assets
- `_source/data/` - JSON files for structured data
- `_source/content/` - Markdown content files

The `.htaccess` file shows how to secure the source directories and route requests appropriately.