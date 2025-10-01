package ggi

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

// BuildConfig holds the configuration for the build process
type BuildConfig struct {
	SourceDir  string // Directory containing templates and resources
	OutputDir  string // Directory for generated static HTML
	BaseURL    string // Base URL for the site
}

// Build generates static HTML from templates
func Build(config BuildConfig) error {
	// Ensure source directory exists
	if _, err := os.Stat(config.SourceDir); os.IsNotExist(err) {
		return fmt.Errorf("source directory does not exist: %s", config.SourceDir)
	}
	
	// Create output directory
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}
	
	// Process templates
	if err := processTemplates(config); err != nil {
		return fmt.Errorf("failed to process templates: %v", err)
	}
	
	// Copy static resources
	if err := copyResources(config); err != nil {
		return fmt.Errorf("failed to copy resources: %v", err)
	}
	
	return nil
}

// processTemplates processes all HTML templates and generates static HTML
func processTemplates(config BuildConfig) error {
	templatesDir := filepath.Join(config.SourceDir, "templates")
	
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
		
		// Determine output path
		relPath, err := filepath.Rel(templatesDir, path)
		if err != nil {
			return err
		}
		
		var outputPath string
		if strings.HasSuffix(relPath, "/index.html") || relPath == "index.html" {
			// For index.html files, use the directory structure as is
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
		
		// Execute template with empty data (in a real implementation, 
		// you would pass in actual data from JSON files)
		data := map[string]interface{}{
			"Title":   "Generated Page",
			"Content": template.HTML("<p>This is a generated static page.</p>"),
		}
		
		if err := tmpl.ExecuteTemplate(outputFile, "base", data); err != nil {
			return fmt.Errorf("failed to execute template %s: %v", path, err)
		}
		
		return nil
	})
	
	return err
}

// copyResources copies static resources to the output directory
func copyResources(config BuildConfig) error {
	resourcesDir := filepath.Join(config.SourceDir, "resources")
	
	// Check if resources directory exists
	if _, err := os.Stat(resourcesDir); os.IsNotExist(err) {
		// If resources directory doesn't exist, that's fine
		return nil
	}
	
	err := filepath.Walk(resourcesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Skip directories for now - we'll create them as needed
		if info.IsDir() {
			return nil
		}
		
		// Get relative path from resources directory
		relPath, err := filepath.Rel(resourcesDir, path)
		if err != nil {
			return err
		}
		
		// Determine output path
		outputPath := filepath.Join(config.OutputDir, relPath)
		
		// Create output directory if needed
		if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
			return err
		}
		
		// Read the source file
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read resource %s: %v", path, err)
		}
		
		// Write to output file
		if err := os.WriteFile(outputPath, content, 0644); err != nil {
			return fmt.Errorf("failed to write resource %s: %v", outputPath, err)
		}
		
		return nil
	})
	
	return err
}