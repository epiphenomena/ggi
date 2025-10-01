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
	PublicSourceDir  string // Directory containing public templates and resources
	AdminSourceDir   string // Directory containing admin templates and resources
	OutputDir        string // Directory for generated static HTML
	BaseURL          string // Base URL for the site
}

// Build generates static HTML from templates for both public and admin sections
func Build(config BuildConfig) error {
	// Ensure source directories exist
	if _, err := os.Stat(config.PublicSourceDir); os.IsNotExist(err) {
		return fmt.Errorf("public source directory does not exist: %s", config.PublicSourceDir)
	}
	
	if config.AdminSourceDir != "" {
		if _, err := os.Stat(config.AdminSourceDir); os.IsNotExist(err) {
			return fmt.Errorf("admin source directory does not exist: %s", config.AdminSourceDir)
		}
	}
	
	// Create output directory
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}
	
	// Process public templates
	if err := processTemplates(config.PublicSourceDir, filepath.Join(config.OutputDir, ""), config); err != nil {
		return fmt.Errorf("failed to process public templates: %v", err)
	}
	
	// Process admin templates if admin source directory is specified
	if config.AdminSourceDir != "" {
		adminOutputDir := filepath.Join(config.OutputDir, "admin")
		if err := processTemplates(config.AdminSourceDir, adminOutputDir, config); err != nil {
			return fmt.Errorf("failed to process admin templates: %v", err)
		}
	}
	
	// Copy static resources
	if err := copyResources(config); err != nil {
		return fmt.Errorf("failed to copy resources: %v", err)
	}
	
	return nil
}

// processTemplates processes all HTML templates and generates static HTML
func processTemplates(sourceDir, outputDir string, config BuildConfig) error {
	templatesDir := filepath.Join(sourceDir, "templates")
	
	// Check if templates directory exists
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		// If templates directory doesn't exist, that's fine, just return
		return nil
	}
	
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
			outputPath = filepath.Join(outputDir, relPath)
		} else {
			// For other files, create subdirectory with index.html
			dir := filepath.Dir(relPath)
			filename := strings.TrimSuffix(filepath.Base(relPath), ".html")
			outputPath = filepath.Join(outputDir, dir, filename, "index.html")
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
	// Copy public resources
	publicResourcesDir := filepath.Join(config.PublicSourceDir, "resources")
	
	// Check if public resources directory exists  
	if _, err := os.Stat(publicResourcesDir); err == nil {
		err := filepath.Walk(publicResourcesDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			
			// Skip directories for now - we'll create them as needed
			if info.IsDir() {
				return nil
			}
			
			// Get relative path from resources directory
			relPath, err := filepath.Rel(publicResourcesDir, path)
			if err != nil {
				return err
			}
			
			// Determine output path in public area
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
		
		if err != nil {
			return err
		}
	}
	
	// Copy admin resources if admin directory is specified
	if config.AdminSourceDir != "" {
		adminResourcesDir := filepath.Join(config.AdminSourceDir, "resources")
		
		// Check if admin resources directory exists
		if _, err := os.Stat(adminResourcesDir); err == nil {
			err := filepath.Walk(adminResourcesDir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				
				// Skip directories for now - we'll create them as needed
				if info.IsDir() {
					return nil
				}
				
				// Get relative path from admin resources directory
				relPath, err := filepath.Rel(adminResourcesDir, path)
				if err != nil {
					return err
				}
				
				// Determine output path in admin area
				outputPath := filepath.Join(config.OutputDir, "admin", relPath)
				
				// Create output directory if needed
				if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
					return err
				}
				
				// Read the source file
				content, err := os.ReadFile(path)
				if err != nil {
					return fmt.Errorf("failed to read admin resource %s: %v", path, err)
				}
				
				// Write to output file
				if err := os.WriteFile(outputPath, content, 0644); err != nil {
					return fmt.Errorf("failed to write admin resource %s: %v", outputPath, err)
				}
				
				return nil
			})
			
			if err != nil {
				return err
			}
		}
	}
	
	return nil
}