package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

//go:embed adminui/templates/*
var adminTemplates embed.FS

//go:embed adminui/css/*
var adminCSS embed.FS



// AdminServer handles the admin interface
type AdminServer struct{}

// NewAdminServer creates a new admin server instance
func NewAdminServer() *AdminServer {
	return &AdminServer{}
}

// HandleCGIRequest processes admin requests in CGI mode
func (a *AdminServer) HandleCGIRequest(method, queryString string) {
	if method == "POST" {
		a.handlePostRequest(queryString)
	} else {
		a.handleGetRequest(queryString)
	}
}

// handleGetRequest handles GET requests to the admin interface
func (a *AdminServer) handleGetRequest(queryString string) {
	if queryString == "" {
		// Show admin home page with links to edit data files
		a.showAdminHome()
	} else {
		// Show form to edit specific data file
		a.showEditForm(queryString)
	}
}

// handlePostRequest handles POST requests to update data files
func (a *AdminServer) handlePostRequest(queryString string) {
	// For POST requests, we need to read from stdin
	contentLengthStr := os.Getenv("CONTENT_LENGTH")
	if contentLengthStr != "" {
		// Parse content length and read form data
		// This is a simplified implementation
		// In a real CGI script, we would parse the form data properly
		fmt.Println("Content-Type: text/html")
		fmt.Println("")
		fmt.Println("<html><body>")
		fmt.Println("<h1>Form Processing</h1>")
		fmt.Printf("<p>Processing form data for: %s</p>", queryString)
		fmt.Println("<p>Form processing not fully implemented in this example.</p>")
		fmt.Println("</body></html>")
		return
	}

	// Show success message and rebuild site
	fmt.Println("Content-Type: text/html")
	fmt.Println("")
	fmt.Println("<html><body>")
	fmt.Println("<h1>Success</h1>")
	fmt.Println("<p>Data updated successfully!</p>")
	
	// Trigger site rebuild
	if err := Build(); err != nil {
		fmt.Printf("<p>Error rebuilding site: %v</p>", err)
	} else {
		fmt.Println("<p>Site rebuilt successfully!</p>")
	}
	
	fmt.Println("<a href='/admin.cgi'>Back to Admin Home</a>")
	fmt.Println("</body></html>")
}

// showAdminHome displays the admin home page
func (a *AdminServer) showAdminHome() {
	fmt.Println("Content-Type: text/html")
	fmt.Println("")
	
	fmt.Println(`
<!DOCTYPE html>
<html>
<head>
    <title>GGI Admin Interface</title>
    <link rel="stylesheet" href="/adminui/css/admin.css">
</head>
<body>
    <header>
        <nav>
            <a href="/">&larr; Site Home</a>
            <a href="/admin.cgi">Admin Home</a>
        </nav>
    </header>
    
    <h1>GGI Admin Interface</h1>
    <p>Welcome to the GGI Admin Interface. You can edit your site content here.</p>
    
    <div class="section">
        <h2>Data Files</h2>
        <ul class="file-list">
    `)
    
	// List all data files in public/data
	dataDir := "public/data"
	files, err := os.ReadDir(dataDir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("<li>No data files found. Add files to public/data to get started.</li>")
		} else {
			fmt.Printf("<li>Error reading data directory: %v</li>", err)
		}
	} else {
		for _, file := range files {
			if !file.IsDir() {
				fileName := file.Name()
				fmt.Printf("<li><a class=\"file-link\" href=\"/admin.cgi?file=%s\">%s</a></li>", fileName, fileName)
			}
		}
	}
	
	fmt.Println(`
        </ul>
    </div>
    
    <div class="section">
        <h2>Add New Content</h2>
        <p>Use the filesystem to add new content files to public/data</p>
    </div>
</body>
</html>
	`)
}

// showEditForm displays a form to edit a specific data file
func (a *AdminServer) showEditForm(queryString string) {
	// Parse the file parameter from query string
	var fileName string
	if strings.HasPrefix(queryString, "file=") {
		fileName = strings.TrimPrefix(queryString, "file=")
	} else {
		fileName = queryString
	}
	
	filePath := filepath.Join("public", "data", fileName)
	
	// Check if file exists
	fileInfo, err := os.Stat(filePath)
	if err != nil || fileInfo.IsDir() {
		// File doesn't exist, show error
		fmt.Println("Content-Type: text/html")
		fmt.Println("")
		fmt.Printf("<html><body><h1>File Not Found</h1><p>The file '%s' does not exist.</p><a href='/admin.cgi'>Back to Admin Home</a></body></html>", fileName)
		return
	}
	
	fmt.Println("Content-Type: text/html")
	fmt.Println("")
	
	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("<html><body><h1>Error Reading File</h1><p>Could not read file '%s': %v</p><a href='/admin.cgi'>Back to Admin Home</a></body></html>", fileName, err)
		return
	}
	
	// Determine appropriate editor based on file type
	fileType := getFileType(fileName)
	fmt.Printf(`
<!DOCTYPE html>
<html>
<head>
    <title>Edit %s - GGI Admin</title>
    <link rel="stylesheet" href="/adminui/css/admin.css">
</head>
<body>
    <header>
        <nav>
            <a href="/">&larr; Site Home</a>
            <a href="/admin.cgi">Admin Home</a>
        </nav>
    </header>
    
    <div class="form-container">
        <h1>Edit %s</h1>
        <form method="POST" action="/admin.cgi?%s">
    `, fileName, fileName, queryString)
	
	if fileType == "json" {
		// Pretty-print JSON for editing
		var prettyJSON map[string]interface{}
		if err := json.Unmarshal(content, &prettyJSON); err == nil {
			prettyContent, err := json.MarshalIndent(prettyJSON, "", "  ")
			if err == nil {
				content = prettyContent
			}
		}
		
		fmt.Printf("<textarea name=\"content\">%s</textarea>", template.HTMLEscapeString(string(content)))
	} else if fileType == "media" {
		// For media files, provide upload form
		fmt.Printf(`
        <p>Current file: <a href="/data/%s" target="_blank">%s</a></p>
        <p>Replace with new file:</p>
        <input type="file" name="file" />
        `, fileName, fileName)
	} else {
		// For other files (markdown, text), show text area
		fmt.Printf("<textarea name=\"content\">%s</textarea>", template.HTMLEscapeString(string(content)))
	}
	
	fmt.Println(`
            <br><br>
            <input type="submit" value="Save Changes" />
        </form>
    </div>
</body>
</html>
	`)
}

// getFileType determines the type of file based on extension
func getFileType(fileName string) string {
	ext := strings.ToLower(filepath.Ext(fileName))
	switch ext {
	case ".json":
		return "json"
	case ".md", ".markdown":
		return "markdown"
	case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg", ".ico", ".mp4", ".webm":
		return "media"
	default:
		return "text"
	}
}