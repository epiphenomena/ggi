package main

import (
	"fmt"
	"ggi/pkg/ggi"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	if ggi.IsCGI() {
		// Handle as CGI script
		handleCGI()
	} else {
		// Start development server
		handleDevServer()
	}
}

func handleCGI() {
	// Handle both GET and POST requests for admin CGI script
	requestMethod := os.Getenv("REQUEST_METHOD")
	
	// Read request body if it's a POST
	var body []byte
	if requestMethod == "POST" {
		var err error
		body, err = io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println("Status: 400 Bad Request")
			fmt.Println("Content-Type: text/plain")
			fmt.Println()
			fmt.Println("Error reading request body")
			return
		}
	}

	// Create a mock request to process the form
	req, err := http.NewRequest(requestMethod, "/", strings.NewReader(string(body)))
	if err != nil {
		fmt.Println("Status: 500 Internal Server Error")
		fmt.Println("Content-Type: text/plain")
		fmt.Println()
		fmt.Println("Error creating request")
		return
	}
	
	contentType := os.Getenv("CONTENT_TYPE")
	if contentType == "" && requestMethod == "POST" {
		contentType = "application/x-www-form-urlencoded"
	}
	req.Header.Set("Content-Type", contentType)
	
	if requestMethod == "POST" {
		err = req.ParseForm()
		if err != nil {
			fmt.Println("Status: 400 Bad Request")
			fmt.Println("Content-Type: text/plain")
			fmt.Println()
			fmt.Println("Error parsing form data")
			return
		}
	}

	// Process the request based on form action or path
	action := req.FormValue("action")
	
	fmt.Println("Content-Type: text/html")
	fmt.Println()
	
	switch action {
	case "save_markdown":
		content := req.FormValue("content")
		filePath := req.FormValue("file_path")
		if content != "" && filePath != "" {
			// In a real implementation, save the markdown content
			fmt.Println("<html><body><h1>Content Saved</h1><p>Markdown content has been saved.</p></body></html>")
		} else {
			fmt.Println("<html><body><h1>Error</h1><p>Missing content or file path.</p></body></html>")
		}
	default:
		// Default response for CGI script
		fmt.Println("<html><body><h1>CGI Script Running</h1><p>Admin CGI script is working.</p></body></html>")
	}
}

func handleDevServer() {
	http.HandleFunc("/admin/", func(w http.ResponseWriter, r *http.Request) {
		// For development, serve admin pages
		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Error parsing form", http.StatusBadRequest)
				return
			}
		}

		// Serve a simple admin test page
		tmpl := `
{{define "content"}}
<h1>Admin Dashboard</h1>
<p>This is the admin section for site management.</p>
<form method="post" action="/admin/save">
    <input type="hidden" name="action" value="save_markdown">
    <div>
        <label for="file_path">File Path:</label>
        <input type="text" name="file_path" id="file_path" value="content/main.md">
    </div>
    <div>
        <label for="content">Content:</label>
        <textarea name="content" id="content" rows="10" cols="50"># Editable content</textarea>
    </div>
    <button type="submit">Save Content</button>
</form>
{{end}}
`
		parsedTmpl, err := ggi.ParseTemplate(tmpl)
		if err != nil {
			http.Error(w, "Error parsing template", http.StatusInternalServerError)
			return
		}

		context := map[string]interface{}{
			"Title": "Admin Dashboard",
		}

		err = parsedTmpl.ExecuteTemplate(w, "base", context)
		if err != nil {
			http.Error(w, "Error executing template", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Serve public site
		tmpl := `
{{define "content"}}
<h1>Welcome to the Public Site</h1>
<p>This is the public-facing part of the website.</p>
<p><a href="/admin/">Admin Panel</a></p>
{{end}}
`
		parsedTmpl, err := ggi.ParseTemplate(tmpl)
		if err != nil {
			http.Error(w, "Error parsing template", http.StatusInternalServerError)
			return
		}

		context := map[string]interface{}{
			"Title": "Public Site",
		}

		err = parsedTmpl.ExecuteTemplate(w, "base", context)
		if err != nil {
			http.Error(w, "Error executing template", http.StatusInternalServerError)
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Development server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}