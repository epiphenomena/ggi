package ggi

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

// Build generates static HTML for public site, admin UI, and initializes the source directory
func Build(config BuildConfig) error {
	// Validate config
	if config.OutputDir == "" {
		return fmt.Errorf("output directory must be specified")
	}
	
	if config.ContentDir == "" {
		return fmt.Errorf("content directory must be specified")
	}

	// Create output directory
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// Create content directory
	if err := os.MkdirAll(config.ContentDir, 0755); err != nil {
		return fmt.Errorf("failed to create content directory: %v", err)
	}

	// Generate initial content structure
	if err := generateInitialContent(config); err != nil {
		return fmt.Errorf("failed to generate initial content: %v", err)
	}

	// Generate admin UI if enabled
	if config.EnableAdmin {
		if err := generateAdminUI(config); err != nil {
			return fmt.Errorf("failed to generate admin UI: %v", err)
		}
	}

	// Generate public site
	if err := generatePublicSite(config); err != nil {
		return fmt.Errorf("failed to generate public site: %v", err)
	}

	// Generate CGI script if enabled
	if config.EnableCGI {
		if err := generateCGIScript(config); err != nil {
			return fmt.Errorf("failed to generate CGI script: %v", err)
		}
	}

	// Copy resources
	if err := copyResources(config); err != nil {
		return fmt.Errorf("failed to copy resources: %v", err)
	}

	return nil
}

// generateInitialContent creates the initial _source directory structure with example content
func generateInitialContent(config BuildConfig) error {
	// Create content subdirectories
	contentSubdirs := []string{
		"markdown",
		"data", 
		"media",
		"pages",
	}
	
	for _, subdir := range contentSubdirs {
		fullPath := filepath.Join(config.ContentDir, subdir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return err
		}
	}
	
	// Create example markdown content
	exampleMdPath := filepath.Join(config.ContentDir, "markdown", "home.md")
	exampleMdContent := "# Welcome\n\nThis is your new website. Edit this content in the admin panel."
	if err := os.WriteFile(exampleMdPath, []byte(exampleMdContent), 0644); err != nil {
		return err
	}
	
	// Create example JSON data
	exampleJsonPath := filepath.Join(config.ContentDir, "data", "cards.json")
	exampleJsonContent := `[
  {
    "title": "Example Card",
    "description": "This is an example card that can be managed in the admin panel.",
    "image": "/media/example.jpg"
  }
]`
	if err := os.WriteFile(exampleJsonPath, []byte(exampleJsonContent), 0644); err != nil {
		return err
	}
	
	return nil
}

// generateAdminUI generates the admin interface with forms for each content type
func generateAdminUI(config BuildConfig) error {
	adminDir := filepath.Join(config.OutputDir, "admin")
	if err := os.MkdirAll(adminDir, 0755); err != nil {
		return err
	}

	// For static site generation, we would pre-generate admin UI pages
	// based on the registered content types. For now, we'll create a
	// basic admin interface as an example.
	
	// Create a basic admin dashboard
	adminDashboard := `
{{define "content"}}
<h1>Site Administration</h1>
<div class="admin-nav">
	<h2>Content Management</h2>
	<ul>
		<li><a href="/admin/markdown">Manage Markdown Content</a></li>
		<li><a href="/admin/json">Manage JSON Data</a></li>
		<li><a href="/admin/media">Manage Media Files</a></li>
	</ul>
</div>
{{end}}
`

	// Parse and write the admin dashboard
	tmpl, err := ParseTemplate(adminDashboard)
	if err != nil {
		return err
	}

	outputFile, err := os.Create(filepath.Join(adminDir, "index.html"))
	if err != nil {
		return err
	}
	defer outputFile.Close()

	data := map[string]interface{}{
		"Title": "Admin Dashboard",
	}

	if err := tmpl.ExecuteTemplate(outputFile, "base", data); err != nil {
		return err
	}

	return nil
}

// generatePublicSite generates the public-facing static site from templates and content
func generatePublicSite(config BuildConfig) error {
	templatesDir := config.PublicTemplatesDir

	// Check if templates directory exists
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		// If templates directory doesn't exist, create a basic index.html
		templatesDir = filepath.Join(config.OutputDir, "_templates")
		if err := os.MkdirAll(templatesDir, 0755); err != nil {
			return err
		}

		// Create a basic index template
		basicTemplate := `{{define "content"}}<h1>Welcome to Your Site</h1><p>Your custom site content goes here.</p>{{end}}`
		if err := os.WriteFile(filepath.Join(templatesDir, "index.html"), []byte(basicTemplate), 0644); err != nil {
			return err
		}
	}

	// Process all templates in the directory
	err := filepath.Walk(templatesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Only process .html files
		if !strings.HasSuffix(path, ".html") {
			return nil
		}

		// Read the template file
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %v", path, err)
		}

		// Parse the template
		tmpl, err := ParseTemplate(string(content))
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %v", path, err)
		}

		// Determine output path in the public directory
		relPath, err := filepath.Rel(templatesDir, path)
		if err != nil {
			return err
		}

		var outputPath string
		if strings.HasSuffix(relPath, "/index.html") || relPath == "index.html" {
			// For index.html files, use the directory structure in public dir
			outputPath = filepath.Join(config.OutputDir, relPath)
		} else {
			// For other files, create subdirectory with index.html
			dir := filepath.Dir(relPath)
			filename := strings.TrimSuffix(filepath.Base(relPath), ".html")
			outputPath = filepath.Join(config.OutputDir, dir, filename, "index.html")
		}

		// Create output directory if needed
		if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
			return err
		}

		// Create output file
		outputFile, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("failed to create output file %s: %v", outputPath, err)
		}
		defer outputFile.Close()

		// Execute template with basic data
		data := map[string]interface{}{
			"Title":   "Site Page",
			"Content": template.HTML(""),
		}

		if err := tmpl.ExecuteTemplate(outputFile, "base", data); err != nil {
			return fmt.Errorf("failed to execute template %s: %v", path, err)
		}

		return nil
	})

	return err
}

// generateCGIScript creates the CGI script for handling admin form submissions
func generateCGIScript(config BuildConfig) error {
	// For now, we'll create a placeholder - in a real implementation this would
	// be a Go binary compiled as a CGI script
	cgiContent := `#!/usr/bin/env bash
echo "Content-Type: text/html"
echo ""
echo "<html><body><h1>GGI CGI Script Placeholder</h1>"
echo "<p>This is where the compiled CGI script would handle form submissions.</p>"
echo "</body></html>"
`

	cgiPath := filepath.Join(config.OutputDir, "admin.cgi")
	return os.WriteFile(cgiPath, []byte(cgiContent), 0755)
}

// copyResources copies static resources to the output directory
func copyResources(config BuildConfig) error {
	// This is a simplified implementation - in a real system, 
	// resources would be copied from appropriate directories
	return nil
}