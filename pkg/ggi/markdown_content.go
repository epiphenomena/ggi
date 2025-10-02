package ggi

import (
	"os"
	"path/filepath"
)

// MarkdownContent represents markdown content with its path
type MarkdownContent struct {
	Path    string
	Content string
}

// NewMarkdownContent creates a new markdown content instance
func NewMarkdownContent(path string) (*MarkdownContent, error) {
	content := &MarkdownContent{Path: path}
	
	// Load existing content if the file exists
	if _, err := os.Stat(path); err == nil {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		content.Content = string(data)
	}
	
	return content, nil
}

// Name returns the content type identifier
func (mc *MarkdownContent) Name() string {
	return "markdown"
}

// AdminForm generates the HTML form for editing this content
func (mc *MarkdownContent) AdminForm() (string, error) {
	form := `
<div class="admin-form">
	<h3>Edit Markdown Content</h3>
	<form method="post" action="/admin/save">
		<input type="hidden" name="content_type" value="markdown">
		<input type="hidden" name="content_path" value="` + mc.Path + `">
		<div>
			<label for="content">Markdown Content:</label>
			<textarea name="content" id="content" rows="15" cols="80">` + mc.Content + `</textarea>
		</div>
		<button type="submit">Save Content</button>
	</form>
</div>
`
	return form, nil
}

// Save saves the content to its file
func (mc *MarkdownContent) Save(formData map[string]string) error {
	content := formData["content"]
	mc.Content = content
	
	// Ensure the directory exists
	dir := filepath.Dir(mc.Path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	return os.WriteFile(mc.Path, []byte(content), 0644)
}

// Load loads content from its file
func (mc *MarkdownContent) Load() (interface{}, error) {
	// Load from file if not already loaded
	if mc.Content == "" {
		data, err := os.ReadFile(mc.Path)
		if err != nil {
			return nil, err
		}
		mc.Content = string(data)
	}
	
	// Convert markdown to HTML
	html, err := ToHTML(mc.Content)
	if err != nil {
		return nil, err
	}
	
	return html, nil
}

// TemplateName returns the name of the template to use for displaying this content
func (mc *MarkdownContent) TemplateName() string {
	return "markdown"
}