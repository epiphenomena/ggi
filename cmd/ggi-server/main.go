package main

import (
	"fmt"
	"ggi/pkg/ggi"
	"ggi/pkg/ggi/templates"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Set the secret key from environment or use the default
	if secret := os.Getenv("GGI_SECRET_KEY"); secret != "" {
		ggi.SecretKey = secret
	}

	if ggi.IsCGI() {
		// Handle as CGI script
		handleCGI()
	} else {
		// Start development server
		handleDevServer()
	}
}

func handleCGI() {
	// Only handle POST requests
	if os.Getenv("REQUEST_METHOD") != "POST" {
		fmt.Println("Status: 405 Method Not Allowed")
		fmt.Println("Content-Type: text/plain")
		fmt.Println()
		fmt.Println("Only POST requests are allowed")
		return
	}

	// Read the request body
	body, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println("Status: 400 Bad Request")
		fmt.Println("Content-Type: text/plain")
		fmt.Println()
		fmt.Println("Error reading request body")
		return
	}

	// Create a mock request to parse the form
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
	
	token := req.FormValue("token")
	if token != ggi.SecretKey {
		fmt.Println("Status: 401 Unauthorized")
		fmt.Println("Content-Type: text/plain")
		fmt.Println()
		fmt.Println("Invalid token")
		return
	}

	// Process the request - this is where you'd add specific handlers
	// based on the form action or other parameters
	fmt.Println("Content-Type: text/html")
	fmt.Println()
	fmt.Println("<html><body><h1>CGI Request Processed</h1><p>Token authenticated successfully.</p></body></html>")
}

func handleDevServer() {
	// For development, we'll set the secret key automatically
	ggi.SecretKey = "dev_secret_key_for_testing"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// For development, allow access without token for GET requests
		// but require token for POST requests
		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Error parsing form", http.StatusBadRequest)
				return
			}

			token := r.FormValue("token")
			if token != ggi.SecretKey {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
		}

		// Serve a simple test page
		templates.RenderMarkdownPage(w, "Development Server", "# Welcome to GGI Development Server\n\nThis is a test page.", nil)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Development server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}