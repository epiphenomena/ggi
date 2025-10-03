package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

//go:embed site/templates/*
var siteTemplates embed.FS

//go:embed site/css/*
var siteCSS embed.FS

//go:embed site/js/*
var siteJS embed.FS



// Build generates the static site from templates and data files
func Build() error {
	fmt.Println("Building static site...")

	// Copy all site assets to public directory
	if err := copyEmbeddedAssets(); err != nil {
		return fmt.Errorf("error copying assets: %v", err)
	}

	// Process templates with data files
	if err := processTemplates(); err != nil {
		return fmt.Errorf("error processing templates: %v", err)
	}

	fmt.Println("Site built successfully!")
	return nil
}

// copyEmbeddedAssets copies all embedded assets to the public directory
func copyEmbeddedAssets() error {
	assetPaths := []struct {
		embedFS  embed.FS
		baseDir  string
		destDir  string
	}{
		{siteCSS, "site/css", "public/css"},
		{siteJS, "site/js", "public/js"},
	}

	for _, assetPath := range assetPaths {
		err := copyEmbeddedDir(assetPath.embedFS, assetPath.baseDir, assetPath.destDir)
		if err != nil {
			return err
		}
	}

	return nil
}

// copyEmbeddedDir copies a directory from embedded filesystem to destination
func copyEmbeddedDir(embedFS embed.FS, srcDir, destDir string) error {
	entries, err := fs.ReadDir(embedFS, srcDir)
	if err != nil {
		return err
	}

	err = os.MkdirAll(destDir, 0755)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		destPath := filepath.Join(destDir, entry.Name())

		if entry.IsDir() {
			err := copyEmbeddedDir(embedFS, srcPath, destPath)
			if err != nil {
				return err
			}
		} else {
			data, err := embedFS.ReadFile(srcPath)
			if err != nil {
				return err
			}

			err = os.WriteFile(destPath, data, 0644)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// processTemplates processes all template files with data
func processTemplates() error {
	// Read all data files from public/data
	dataDir := "public/data"
	dataFiles, err := os.ReadDir(dataDir)
	if err != nil {
		if os.IsNotExist(err) {
			// Data directory doesn't exist yet, create it
			err = os.MkdirAll(dataDir, 0755)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	// Load data from all JSON and markdown files
	allData := make(map[string]interface{})
	if dataFiles != nil {
		for _, file := range dataFiles {
			if !file.IsDir() {
				filePath := filepath.Join(dataDir, file.Name())
				ext := strings.ToLower(filepath.Ext(file.Name()))
				
				switch ext {
				case ".json":
					jsonData, err := loadDataFile(filePath)
					if err != nil {
						fmt.Printf("Warning: could not load data file %s: %v\n", filePath, err)
						continue
					}
					key := strings.TrimSuffix(file.Name(), ".json")
					allData[key] = jsonData
				case ".md", ".markdown":
					markdownContent, err := os.ReadFile(filePath)
					if err != nil {
						fmt.Printf("Warning: could not load markdown file %s: %v\n", filePath, err)
						continue
					}
					key := strings.TrimSuffix(file.Name(), ext)
					// Convert markdown to HTML
					htmlContent := processMarkdown(string(markdownContent))
					allData[key+"HTML"] = template.HTML(htmlContent)
					// Also make raw content available
					allData[key] = string(markdownContent)
				}
			}
		}
	}

	// Use data from site.json, with fallback to defaults
	if siteData, ok := allData["site"].(map[string]interface{}); ok {
		allData["SiteTitle"] = getFromData(siteData, "siteTitle", "GGI Sample Site")
		allData["WelcomeText"] = getFromData(siteData, "welcomeText", "Welcome to our website!")
		allData["CurrentYear"] = getFromData(siteData, "currentYear", "2025")
		allData["heroImage"] = getFromData(siteData, "heroImage", "")
		allData["portfolio"] = getFromData(siteData, "portfolio", []interface{}{})
		allData["ContactInfo"] = getFromData(siteData, "contactInfo", nil)
	} else {
		allData["SiteTitle"] = "GGI Sample Site"
		allData["WelcomeText"] = "Welcome to our website!"
		allData["CurrentYear"] = "2025"
		allData["heroImage"] = ""
		allData["portfolio"] = []interface{}{}
	}
	
	// Use about content from about.md if it exists, otherwise fallback to site.json
	if _, exists := allData["about"]; exists {
		allData["aboutText"] = allData["aboutHTML"]
	} else if siteData, ok := allData["site"].(map[string]interface{}); ok {
		allData["aboutText"] = getFromData(siteData, "aboutText", "")
	} else {
		allData["aboutText"] = ""
	}

	// Copy media files from data directory to public/media directory
	if err := copyMediaFiles(); err != nil {
		return fmt.Errorf("error copying media files: %v", err)
	}

	// Find all template files
	entries, err := fs.ReadDir(siteTemplates, "site/templates")
	if err != nil {
		return err
	}

	// Process each template file
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".tmpl") {
			templatePath := filepath.Join("site/templates", entry.Name())
			outputPath := filepath.Join("public", strings.TrimSuffix(entry.Name(), ".tmpl")+".html")
			
			if err := renderTemplate(templatePath, outputPath, allData); err != nil {
				return fmt.Errorf("error rendering template %s: %v", templatePath, err)
			}
		}
	}

	return nil
}

// copyMediaFiles copies media files from public/data to public/media so they can be served
func copyMediaFiles() error {
	dataDir := "public/data"
	mediaDir := "public/media"
	
	// Create media directory if it doesn't exist
	if err := os.MkdirAll(mediaDir, 0755); err != nil {
		return err
	}
	
	// Get all files in data directory
	entries, err := os.ReadDir(dataDir)
	if err != nil {
		return err
	}
	
	// Define media file extensions
	mediaExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, 
		".webp": true, ".svg": true, ".ico": true, ".mp4": true, 
		".webm": true, ".pdf": true, ".mp3": true, ".wav": true,
	}
	
	for _, entry := range entries {
		if !entry.IsDir() {
			ext := strings.ToLower(filepath.Ext(entry.Name()))
			if mediaExts[ext] {
				srcPath := filepath.Join(dataDir, entry.Name())
				dstPath := filepath.Join(mediaDir, entry.Name())
				
				// Copy the file
				if err := copyFile(srcPath, dstPath); err != nil {
					return fmt.Errorf("error copying media file %s: %v", entry.Name(), err)
				}
			}
		}
	}
	
	return nil
}

// copyFile copies a file from source to destination
func copyFile(src, dst string) error {
	content, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	
	return os.WriteFile(dst, content, 0644)
}

// loadDataFile loads data from a JSON file
func loadDataFile(filePath string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	
	err = json.Unmarshal(content, &data)
	if err != nil {
		return nil, err
	}
	
	return data, nil
}

// renderTemplate renders a template with provided data and saves to output file
func renderTemplate(templatePath, outputPath string, data map[string]interface{}) error {
	// Read the template content from embedded filesystem
	tmplContent, err := siteTemplates.ReadFile(templatePath)
	if err != nil {
		return err
	}

	// Parse the template
	tmpl, err := template.New(filepath.Base(templatePath)).Parse(string(tmplContent))
	if err != nil {
		return err
	}

	// Create the output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Execute the template with data
	return tmpl.Execute(outputFile, data)
}

// Clean removes build artifacts from the public folder, keeping only data and admin.cgi
func Clean() error {
	publicDir := "public"
	entries, err := os.ReadDir(publicDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		// Skip the data directory and admin.cgi file
		if entry.Name() == "data" || entry.Name() == "admin.cgi" || entry.Name() == ".htaccess" {
			continue
		}

		// Remove everything else
		entryPath := filepath.Join(publicDir, entry.Name())
		if err := os.RemoveAll(entryPath); err != nil {
			return err
		}
	}

	return nil
}

// processMarkdown performs basic markdown conversion to HTML
func processMarkdown(md string) string {
	// Split content into lines for processing
	lines := strings.Split(md, "\n")
	var htmlParts []string
	
	i := 0
	for i < len(lines) {
		trimmedLine := strings.TrimSpace(lines[i])
		
		// Check for headers
		if strings.HasPrefix(trimmedLine, "# ") {
			// H1 header
			content := strings.TrimPrefix(trimmedLine, "# ")
			htmlParts = append(htmlParts, "<h1>"+content+"</h1>")
			i++
		} else if strings.HasPrefix(trimmedLine, "## ") {
			// H2 header
			content := strings.TrimPrefix(trimmedLine, "## ")
			htmlParts = append(htmlParts, "<h2>"+content+"</h2>")
			i++
		} else if strings.HasPrefix(trimmedLine, "### ") {
			// H3 header
			content := strings.TrimPrefix(trimmedLine, "### ")
			htmlParts = append(htmlParts, "<h3>"+content+"</h3>")
			i++
		} else {
			// Collect all consecutive non-header lines as a paragraph
			var paragraphLines []string
			for i < len(lines) {
				currentLine := lines[i]
				trimmedCurrent := strings.TrimSpace(currentLine)
				
				// Check if this line is a header
				if strings.HasPrefix(trimmedCurrent, "# ") || 
				   strings.HasPrefix(trimmedCurrent, "## ") || 
				   strings.HasPrefix(trimmedCurrent, "### ") {
					break
				}
				
				// Process markdown elements within the line
				processedLine := currentLine
				// Convert bold (**text** -> <strong>text</strong>)
				reBold := regexp.MustCompile(`\*\*(.*?)\*\*`)
				processedLine = reBold.ReplaceAllString(processedLine, "<strong>$1</strong>")
				
				// Convert italic (*text* -> <em>text</em>)
				reItalic := regexp.MustCompile(`\*(.*?)\*`)
				processedLine = reItalic.ReplaceAllString(processedLine, "<em>$1</em>")
				
				// Convert links ([text](url) -> <a href="url">text</a>)
				reLink := regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)
				processedLine = reLink.ReplaceAllString(processedLine, `<a href="$2">$1</a>`)
				
				// Convert image ![](url) -> <img src="url" />
				reImage := regexp.MustCompile(`!\[.*?\]\((.*?)\)`)
				processedLine = reImage.ReplaceAllString(processedLine, `<img src="$1" alt="" />`)
				
				paragraphLines = append(paragraphLines, processedLine)
				i++
				
				// If the next line is empty, we've reached the end of this paragraph
				if i < len(lines) && strings.TrimSpace(lines[i]) == "" {
					break
				}
			}
			
			// Join paragraph lines with <br> tags and wrap in <p>
			if len(paragraphLines) > 0 {
				paragraphContent := strings.Join(paragraphLines, "<br />")
				htmlParts = append(htmlParts, "<p>"+paragraphContent+"</p>")
			}
			
			// Skip the empty line that ends the paragraph
			if i < len(lines) && strings.TrimSpace(lines[i]) == "" {
				i++
			}
		}
	}
	
	return strings.Join(htmlParts, "\n")
}

// getFromData safely gets a value from a data map with a fallback
func getFromData(data map[string]interface{}, key string, fallback interface{}) interface{} {
	if val, exists := data[key]; exists {
		return val
	}
	return fallback
}