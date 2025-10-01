package ggi

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// CSSConfig holds configuration values for CSS generation
type CSSConfig struct {
	PrimaryColor    string `json:"primary_color"`
	SecondaryColor  string `json:"secondary_color"`
	BackgroundColor string `json:"background_color"`
	TextColor       string `json:"text_color"`
	FontSize        string `json:"font_size"`
	FontFamily      string `json:"font_family"`
}

// DefaultCSSConfig provides default styling values
var DefaultCSSConfig = CSSConfig{
	PrimaryColor:    "#3498db",
	SecondaryColor:  "#2ecc71",
	BackgroundColor: "#ffffff",
	TextColor:       "#333333",
	FontSize:        "16px",
	FontFamily:      "Arial, sans-serif",
}

// GenerateCSS creates CSS from a template and configuration
func GenerateCSS(cssTemplate string, config CSSConfig) (string, error) {
	tmpl, err := template.New("css").Parse(cssTemplate)
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	err = tmpl.Execute(&buf, config)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// SaveCSSConfig saves CSS configuration to a JSON file
func SaveCSSConfig(config CSSConfig, filePath string) error {
	// Ensure the file is in a safe location
	if !isSafePath(filePath) {
		return fmt.Errorf("unsafe file path: %s", filePath)
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(config)
}

// LoadCSSConfig loads CSS configuration from a JSON file
func LoadCSSConfig(filePath string) (CSSConfig, error) {
	var config CSSConfig

	// Ensure the file is in a safe location
	if !isSafePath(filePath) {
		return config, fmt.Errorf("unsafe file path: %s", filePath)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(content, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}