package ggi

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// MarkdownContentType handles markdown content
type MarkdownContentType struct{}

func (mct *MarkdownContentType) Name() string {
	return "markdown"
}

func (mct *MarkdownContentType) AdminForm(contentPath string) (string, error) {
	// Read existing content if it exists
	content := ""
	if _, err := os.Stat(contentPath); err == nil {
		data, err := os.ReadFile(contentPath)
		if err != nil {
			return "", err
		}
		content = string(data)
	}

	form := `
<div class="admin-form">
	<h3>Edit Markdown Content</h3>
	<form method="post" action="/admin/save">
		<input type="hidden" name="content_type" value="markdown">
		<input type="hidden" name="content_path" value="` + contentPath + `">
		<div>
			<label for="content">Markdown Content:</label>
			<textarea name="content" id="content" rows="15" cols="80">` + content + `</textarea>
		</div>
		<button type="submit">Save Content</button>
	</form>
</div>
`
	return form, nil
}

func (mct *MarkdownContentType) Save(contentPath string, formData map[string]string) error {
	content := formData["content"]
	
	// Ensure the directory exists
	dir := filepath.Dir(contentPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	return os.WriteFile(contentPath, []byte(content), 0644)
}

func (mct *MarkdownContentType) Load(contentPath string) (interface{}, error) {
	content, err := os.ReadFile(contentPath)
	if err != nil {
		return "", err
	}
	
	// Convert markdown to HTML
	html, err := ToHTML(string(content))
	if err != nil {
		return "", err
	}
	
	return html, nil
}

func (mct *MarkdownContentType) TemplateName() string {
	return "markdown"
}

// JSONContentType handles structured JSON data
type JSONContentType struct {
	ItemType interface{} // The struct type for items in this collection
}

func (jct *JSONContentType) Name() string {
	return "json"
}

func (jct *JSONContentType) AdminForm(contentPath string) (string, error) {
	// For now, implement a basic text editor for JSON
	jsonContent := "[]"
	if _, err := os.Stat(contentPath); err == nil {
		data, err := os.ReadFile(contentPath)
		if err != nil {
			return "", err
		}
		jsonContent = string(data)
	}

	form := `
<div class="admin-form">
	<h3>Edit JSON Data</h3>
	<form method="post" action="/admin/save">
		<input type="hidden" name="content_type" value="json">
		<input type="hidden" name="content_path" value="` + contentPath + `">
		<div>
			<label for="content">JSON Data:</label>
			<textarea name="content" id="content" rows="15" cols="80">` + jsonContent + `</textarea>
		</div>
		<button type="submit">Save Data</button>
	</form>
</div>
`
	return form, nil
}

func (jct *JSONContentType) Save(contentPath string, formData map[string]string) error {
	content := formData["content"]
	
	// Validate that it's valid JSON
	var temp interface{}
	if err := json.Unmarshal([]byte(content), &temp); err != nil {
		return err
	}
	
	// Ensure the directory exists
	dir := filepath.Dir(contentPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	return os.WriteFile(contentPath, []byte(content), 0644)
}

func (jct *JSONContentType) Load(contentPath string) (interface{}, error) {
	content, err := os.ReadFile(contentPath)
	if err != nil {
		return nil, err
	}
	
	var data interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		return nil, err
	}
	
	return data, nil
}

func (jct *JSONContentType) TemplateName() string {
	return "json"
}

// MediaContentType handles media files like images
type MediaContentType struct{}

func (mct *MediaContentType) Name() string {
	return "media"
}

func (mct *MediaContentType) AdminForm(contentPath string) (string, error) {
	// Check if file exists and get current URL
	currentURL := ""
	if _, err := os.Stat(contentPath); err == nil {
		// Convert file path to URL path (this is a simplified approach)
		relPath, err := filepath.Rel("_source", contentPath)
		if err == nil {
			currentURL = "/" + strings.Replace(relPath, string(filepath.Separator), "/", -1)
		}
	}

	form := `
<div class="admin-form">
	<h3>Manage Media</h3>
	<form method="post" action="/admin/save" enctype="multipart/form-data">
		<input type="hidden" name="content_type" value="media">
		<input type="hidden" name="content_path" value="` + contentPath + `">
		<div>
			<label for="file">Upload New File:</label>
			<input type="file" name="file" id="file">
		</div>
		<div>
			<label>Current File:</label>
			` + func() string {
				if currentURL != "" {
					return `<img src="` + currentURL + `" alt="Current image" style="max-width: 300px;">`
				}
				return "<p>No file uploaded yet</p>"
			}() + `
		</div>
		<button type="submit">Upload File</button>
	</form>
</div>
`
	return form, nil
}

func (mct *MediaContentType) Save(contentPath string, formData map[string]string) error {
	// For now, this would handle file upload - in a real implementation
	// we'd need to handle multipart form data properly
	// This is a placeholder implementation
	return nil
}

func (mct *MediaContentType) Load(contentPath string) (interface{}, error) {
	// Return the relative path that can be used in templates
	relPath, err := filepath.Rel("_source", contentPath)
	if err != nil {
		return "", err
	}
	
	return "/" + strings.Replace(relPath, string(filepath.Separator), "/", -1), nil
}

func (mct *MediaContentType) TemplateName() string {
	return "media"
}

// RegisterDefaultContentTypes registers the built-in content types
func RegisterDefaultContentTypes() {
	RegisterContentType(&MarkdownContentType{})
	RegisterContentType(&JSONContentType{})
	RegisterContentType(&MediaContentType{})
}