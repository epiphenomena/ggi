# GGI Basic Site Example

This example demonstrates how to use the GGI library to create a simple website with CGI-based editing capabilities.

## Features Demonstrated

1. **CGI Detection**: The application automatically detects if it's running as a CGI script or in development mode
2. **Authentication**: Uses a secret key system for form submissions
3. **Markdown Editing**: Provides an interface to edit and save markdown content
4. **Form-based Data Management**: Creates forms from Go structs and saves data as JSON
5. **CSS Customization**: Shows how CSS can be customized through configuration

## Running the Example

### Development Mode
```bash
cd examples/basic-site
go run main.go
```

The development server will start on port 8080 with an automatically set secret key.

### As CGI Script
1. Set your `GGI_SECRET_KEY` environment variable:
```bash
export GGI_SECRET_KEY="your_secret_key"
```

2. Compile the binary:
```bash
go build -o ggi.cgi main.go
```

3. Place the `ggi.cgi` file in your web server's CGI directory and configure the .htaccess as needed.

## Structure

- `_source/templates/` - HTML templates for the site
- `_source/resources/` - Static assets like CSS, JS, images
- `_source/data/` - JSON files for structured data
- `_source/markdown/` - Markdown content files

The `.htaccess` file shows how to secure the source directories and route requests appropriately.