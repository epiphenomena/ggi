package ggi

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// MarkdownToHTML converts markdown content to HTML
func MarkdownToHTML(mdContent string) (string, error) {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	var buf strings.Builder
	if err := md.Convert([]byte(mdContent), &buf); err != nil {
		return "", err
	}
	
	return buf.String(), nil
}

// SaveMarkdown saves markdown content to a file
func SaveMarkdown(content, filePath string) error {
	// Ensure the file is in a safe location
	if !isSafePath(filePath) {
		return fmt.Errorf("unsafe file path: %s", filePath)
	}
	
	return os.WriteFile(filePath, []byte(content), 0644)
}

// LoadMarkdown loads markdown content from a file
func LoadMarkdown(filePath string) (string, error) {
	// Ensure the file is in a safe location
	if !isSafePath(filePath) {
		return "", fmt.Errorf("unsafe file path: %s", filePath)
	}
	
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	
	return string(content), nil
}

// isSafePath ensures the path is in the current directory or a subdirectory
func isSafePath(path string) bool {
	// Resolve to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return false
	}
	
	// Check that the absolute path starts with current directory
	return strings.HasPrefix(absPath, cwd)
}

// RenderMarkdownPage renders a page with markdown content
func RenderMarkdownPage(w io.Writer, title, mdContent string, data interface{}) error {
	htmlContent, err := MarkdownToHTML(mdContent)
	if err != nil {
		return err
	}
	
	tmpl := `
{{define "content"}}
<div class="markdown-content">
	` + htmlContent + `
</div>
{{end}}
`
	
	parsedTmpl, err := ParseTemplate(tmpl)
	if err != nil {
		return err
	}
	
	context := map[string]interface{}{
		"Title":   title,
		"Content": template.HTML(htmlContent),
		"Data":    data,
	}
	
	return parsedTmpl.ExecuteTemplate(w, "base", context)
}