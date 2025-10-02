package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

//go:embed site/templates/*
var siteTemplates embed.FS

//go:embed site/css/*
var siteCSS embed.FS



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

	// Load data from all JSON files
	allData := make(map[string]interface{})
	if dataFiles != nil {
		for _, file := range dataFiles {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
				filePath := filepath.Join(dataDir, file.Name())
				jsonData, err := loadDataFile(filePath)
				if err != nil {
					fmt.Printf("Warning: could not load data file %s: %v\n", filePath, err)
					continue
				}
				key := strings.TrimSuffix(file.Name(), ".json")
				allData[key] = jsonData
			}
		}
	}

	// Use data from JSON files, with fallback to defaults
	siteData, hasSiteData := allData["site"].(map[string]interface{})
	exampleData, hasExampleData := allData["example"].(map[string]interface{})
	
	// Use site.json data if available, otherwise example.json, otherwise defaults
	if hasSiteData {
		allData["SiteTitle"] = getFromData(siteData, "siteTitle", "GGI Sample Site")
		allData["WelcomeText"] = getFromData(siteData, "welcomeText", "Welcome to our website!")
		allData["CurrentYear"] = getFromData(siteData, "currentYear", "2025")
		allData["aboutText"] = getFromData(siteData, "aboutText", "")
		allData["portfolio"] = getFromData(siteData, "portfolio", []interface{}{})
		allData["ContactInfo"] = getFromData(siteData, "contactInfo", nil)
	} else if hasExampleData {
		allData["SiteTitle"] = getFromData(exampleData, "siteTitle", "GGI Sample Site")
		allData["WelcomeText"] = getFromData(exampleData, "welcomeText", "Welcome to our website!")
		allData["CurrentYear"] = getFromData(exampleData, "currentYear", "2025")
		allData["aboutText"] = getFromData(exampleData, "aboutText", "")
		allData["portfolio"] = getFromData(exampleData, "portfolio", []interface{}{})
		allData["ContactInfo"] = getFromData(exampleData, "contactInfo", nil)
	} else {
		allData["SiteTitle"] = "GGI Sample Site"
		allData["WelcomeText"] = "Welcome to our website!"
		allData["CurrentYear"] = "2025"
		allData["aboutText"] = ""
		allData["portfolio"] = []interface{}{}
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

// getFromData safely gets a value from a data map with a fallback
func getFromData(data map[string]interface{}, key string, fallback interface{}) interface{} {
	if val, exists := data[key]; exists {
		return val
	}
	return fallback
}