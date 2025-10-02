#!/bin/bash
# Build script for GGI

# Build the admin CGI binary
echo "Building admin.cgi..."
go build -o public/admin.cgi .
chmod +x public/admin.cgi

echo "Build completed! The admin.cgi binary is in the public/ directory."
echo "Upload the contents of the public/ directory to your web server."