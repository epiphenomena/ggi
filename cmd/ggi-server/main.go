package main

import (
	"fmt"
	"ggi/pkg/ggi"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	if ggi.IsCGI() {
		// Handle as CGI script
		handleCGI()
	} else {
		// For development, could start a test server
		// but typically, this binary would only run as CGI
		fmt.Println("Status: 400 Bad Request")
		fmt.Println("Content-Type: text/plain")
		fmt.Println()
		fmt.Println("This binary is intended to run as a CGI script")
	}
}

func handleCGI() {
	// Initialize default content types
	ggi.RegisterDefaultContentTypes()

	requestMethod := os.Getenv("REQUEST_METHOD")
	
	if requestMethod != "POST" {
		// For now, only handle POST requests
		fmt.Println("Status: 405 Method Not Allowed")
		fmt.Println("Content-Type: text/plain")
		fmt.Println()
		fmt.Println("Only POST requests are allowed")
		return
	}

	// Read the POST body
	body, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println("Status: 400 Bad Request")
		fmt.Println("Content-Type: text/plain")
		fmt.Println()
		fmt.Println("Error reading request body")
		return
	}

	// Create a mock request to process the form
	req, err := http.NewRequest("POST", "/", strings.NewReader(string(body)))
	if err != nil {
		fmt.Println("Status: 500 Internal Server Error")
		fmt.Println("Content-Type: text/plain")
		fmt.Println()
		fmt.Println("Error creating request")
		return
	}
	
	contentType := os.Getenv("CONTENT_TYPE")
	if contentType == "" {
		contentType = "application/x-www-form-urlencoded"
	}
	req.Header.Set("Content-Type", contentType)
	
	err = req.ParseForm()
	if err != nil {
		fmt.Println("Status: 400 Bad Request")
		fmt.Println("Content-Type: text/plain")
		fmt.Println()
		fmt.Println("Error parsing form data")
		return
	}

	// Get content type and path from form
	contentTypeValue := req.FormValue("content_type")
	contentPath := req.FormValue("content_path")
	
	if contentTypeValue == "" || contentPath == "" {
		fmt.Println("Status: 400 Bad Request")
		fmt.Println("Content-Type: text/plain")
		fmt.Println()
		fmt.Println("Missing content_type or content_path")
		return
	}

	// Get the content type handler
	ct, exists := ggi.GetContentType(contentTypeValue)
	if !exists {
		fmt.Println("Status: 400 Bad Request")
		fmt.Println("Content-Type: text/plain")
		fmt.Println()
		fmt.Printf("Unknown content type: %s\n", contentTypeValue)
		return
	}

	// Prepare form data for saving
	formData := make(map[string]string)
	for key, values := range req.Form {
		if len(values) > 0 {
			formData[key] = values[0] // Take the first value for each key
		}
	}

	// Save the content using the appropriate handler
	err = ct.Save(contentPath, formData)
	if err != nil {
		fmt.Println("Status: 500 Internal Server Error")
		fmt.Println("Content-Type: text/plain")
		fmt.Println()
		fmt.Printf("Error saving content: %v\n", err)
		return
	}

	// Success response
	fmt.Println("Content-Type: text/html")
	fmt.Println()
	fmt.Println("<html><body><h1>Content Saved</h1><p>Content has been successfully saved.</p><p><a href=\"/admin/\">Return to Admin Dashboard</a></p></body></html>")
}