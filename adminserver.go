package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

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

// handlePostRequest handles POST requests to the admin interface
func (a *AdminServer) handlePostRequest(queryString string) {
	// For this implementation, we'll just show a message since proper form parsing
	// in a CGI context would require more complex handling
	fmt.Println("Content-Type: text/html")
	fmt.Println("")
	fmt.Println("<html><body>")
	fmt.Println("<h1>Form Submitted</h1>")
	fmt.Printf("<p>Received POST request with query string: %s</p>", queryString)
	fmt.Println("<p>In a full implementation, this would process the form data and update the file.</p>")
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
    <script src="/adminui/js/admin.js"></script>
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
    <script src="/adminui/js/admin.js"></script>
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
		// Generate form inputs for JSON data
		var jsonData map[string]interface{}
		if err := json.Unmarshal(content, &jsonData); err == nil {
			fmt.Println("<div class=\"json-form\">")
			generateJSONForm(jsonData, "")
			fmt.Println("</div>")
		} else {
			// If JSON parsing fails, fall back to text area
			fmt.Printf("<textarea name=\"content\">%s</textarea>", template.HTMLEscapeString(string(content)))
		}
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

// generateJSONForm generates HTML form inputs for JSON data
func generateJSONForm(data map[string]interface{}, prefix string) {
	for key, value := range data {
		fieldName := key
		if prefix != "" {
			fieldName = prefix + "." + key
		}
		
		switch v := value.(type) {
		case string:
			fmt.Printf("<div class=\"form-field\"><label for=\"%s\">%s:</label><input type=\"text\" id=\"%s\" name=\"%s\" value=\"%s\" /></div>\n", 
				template.HTMLEscapeString(fieldName), 
				template.HTMLEscapeString(key), 
				template.HTMLEscapeString(fieldName), 
				template.HTMLEscapeString(fieldName), 
				template.HTMLEscapeString(v))
		case float64: // JSON numbers are unmarshaled as float64
			fmt.Printf("<div class=\"form-field\"><label for=\"%s\">%s:</label><input type=\"number\" id=\"%s\" name=\"%s\" value=\"%g\" /></div>\n", 
				template.HTMLEscapeString(fieldName), 
				template.HTMLEscapeString(key), 
				template.HTMLEscapeString(fieldName), 
				template.HTMLEscapeString(fieldName), 
				v)
		case bool:
			checked := ""
			if v {
				checked = "checked"
			}
			fmt.Printf("<div class=\"form-field\"><label for=\"%s\">%s:</label><input type=\"checkbox\" id=\"%s\" name=\"%s\" %s /></div>\n", 
				template.HTMLEscapeString(fieldName), 
				template.HTMLEscapeString(key), 
				template.HTMLEscapeString(fieldName), 
				template.HTMLEscapeString(fieldName), 
				checked)
		case map[string]interface{}:
			fmt.Printf("<div class=\"form-section\"><h4>%s</h4>", template.HTMLEscapeString(key))
			generateJSONForm(v, fieldName)
			fmt.Println("</div>")
		case []interface{}:
			fmt.Printf("<div class=\"form-section array-section\" data-array-name=\"%s\"><h4>%s (Array)</h4>", template.HTMLEscapeString(fieldName), template.HTMLEscapeString(key))
			
			// Add button to add new items
			fmt.Printf("<button type=\"button\" class=\"array-add-btn\" onclick=\"addArrayItem('%s')\">Add Item</button>\n", template.HTMLEscapeString(fieldName))
			
			for i, arrItem := range v {
				itemName := fmt.Sprintf("%s[%d]", fieldName, i)
				
				// Container for array item with controls
				fmt.Printf("<div class=\"array-item\" data-index=\"%d\">\n", i)
				fmt.Printf("<div class=\"array-item-controls\">\n")
				fmt.Printf("<button type=\"button\" class=\"array-move-up-btn\" onclick=\"moveArrayItemUp('%s', %d)\">&#8593;</button>\n", template.HTMLEscapeString(fieldName), i)
				fmt.Printf("<button type=\"button\" class=\"array-move-down-btn\" onclick=\"moveArrayItemDown('%s', %d)\">&#8595;</button>\n", template.HTMLEscapeString(fieldName), i)
				fmt.Printf("<button type=\"button\" class=\"array-remove-btn\" onclick=\"removeArrayItem('%s', %d)\">Remove</button>\n", template.HTMLEscapeString(fieldName), i)
				fmt.Println("</div>")
				
				if str, ok := arrItem.(string); ok {
					fmt.Printf("<div class=\"form-field\"><label>Item %d:</label><input type=\"text\" name=\"%s\" value=\"%s\" /></div>\n", 
						i, template.HTMLEscapeString(itemName), template.HTMLEscapeString(str))
				} else if m, ok := arrItem.(map[string]interface{}); ok {
					fmt.Printf("<div class=\"form-field\"><h5>Item %d:</h5>", i)
					generateJSONForm(m, itemName)
					fmt.Println("</div>")
				} else {
					// For other types in arrays, convert to string
					itemStr := fmt.Sprintf("%v", arrItem)
					fmt.Printf("<div class=\"form-field\"><label>Item %d:</label><input type=\"text\" name=\"%s\" value=\"%s\" /></div>\n", 
						i, template.HTMLEscapeString(itemName), template.HTMLEscapeString(itemStr))
				}
				
				fmt.Println("</div>") // Close array-item div
			}
			fmt.Println("</div>")
		default:
			// For other types, convert to string
			str := fmt.Sprintf("%v", v)
			fmt.Printf("<div class=\"form-field\"><label for=\"%s\">%s:</label><input type=\"text\" id=\"%s\" name=\"%s\" value=\"%s\" /></div>\n", 
				template.HTMLEscapeString(fieldName), 
				template.HTMLEscapeString(key), 
				template.HTMLEscapeString(fieldName), 
				template.HTMLEscapeString(fieldName), 
				template.HTMLEscapeString(str))
		}
	}
}