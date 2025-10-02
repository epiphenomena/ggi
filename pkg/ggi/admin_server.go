package ggi

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// AdminServer handles admin functionality including CGI and development server
type AdminServer struct {
	// Holds registered content instances
	contentInstances map[string]ContentType
}

// NewAdminServer creates a new admin server
func NewAdminServer() *AdminServer {
	return &AdminServer{
		contentInstances: make(map[string]ContentType),
	}
}

// RegisterContent registers a content instance for management
func (as *AdminServer) RegisterContent(content ContentType) {
	as.contentInstances[content.Name()] = content
}

// HandleCGI handles the CGI script functionality
func (as *AdminServer) HandleCGI() {
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

	// Get the content instance
	ct, exists := as.contentInstances[contentTypeValue]
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
	err = ct.Save(formData)
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

// HandleDevServer starts a development server for testing
func (as *AdminServer) HandleDevServer(port string) {
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/admin/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Error parsing form", http.StatusBadRequest)
				return
			}

			// Handle form submission
			contentTypeValue := r.FormValue("content_type")
			contentPath := r.FormValue("content_path")
			
			if contentTypeValue == "" || contentPath == "" {
				http.Error(w, "Missing content_type or content_path", http.StatusBadRequest)
				return
			}

			// Get the content instance
			ct, exists := as.contentInstances[contentTypeValue]
			if !exists {
				http.Error(w, "Unknown content type", http.StatusBadRequest)
				return
			}

			// Prepare form data for saving
			formData := make(map[string]string)
			for key, values := range r.Form {
				if len(values) > 0 {
					formData[key] = values[0] // Take the first value for each key
				}
			}

			// Save the content
			err = ct.Save(formData)
			if err != nil {
				http.Error(w, "Error saving content: "+err.Error(), http.StatusInternalServerError)
				return
			}

			// Success page
			tmpl := `
{{define "content"}}
<h1>Content Saved</h1>
<p>Content has been successfully saved.</p>
<p><a href="/admin/">Return to Admin Dashboard</a></p>
{{end}}
`
			parsedTmpl, err := ParseTemplate(tmpl)
			if err != nil {
				http.Error(w, "Error parsing template", http.StatusInternalServerError)
				return
			}

			context := map[string]interface{}{
				"Title": "Success",
			}

			err = parsedTmpl.ExecuteTemplate(w, "base", context)
			if err != nil {
				http.Error(w, "Error executing template", http.StatusInternalServerError)
			}
		} else {
			// Show admin dashboard
			adminContent := `
{{define "content"}}
<h1>Site Administration</h1>
<div class="admin-nav">
	<h2>Content Management</h2>
	<ul>
`
			// Add links for each registered content type
			for _, ct := range as.contentInstances {
				adminContent += fmt.Sprintf(`<li><a href="/admin/manage?type=%s">Manage %s Content</a></li>`, ct.Name(), ct.Name())
			}
			
			adminContent += `
	</ul>
</div>
{{end}}
`

			parsedTmpl, err := ParseTemplate(adminContent)
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
		}
	})

	fmt.Printf("Development server starting on :%s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}

// RunServer runs either as CGI or dev server based on environment
func (as *AdminServer) RunServer() {
	if IsCGI() {
		as.HandleCGI()
	} else {
		as.HandleDevServer("")
	}
}

// RunServerOnPort runs the dev server on a specific port (for non-CGI mode)
func (as *AdminServer) RunServerOnPort(port string) {
	if IsCGI() {
		as.HandleCGI()
	} else {
		as.HandleDevServer(port)
	}
}