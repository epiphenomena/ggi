package ggi

import (
	"os"
	"path/filepath"
	"strings"
)

// MediaContent represents media content with its path
type MediaContent struct {
	Path string
	URL  string
}

// NewMediaContent creates a new media content instance
func NewMediaContent(path string) (*MediaContent, error) {
	content := &MediaContent{Path: path}
	
	// Generate URL from path (simplified)
	if _, err := os.Stat(path); err == nil {
		relPath, err := filepath.Rel("_source", path)
		if err == nil {
			content.URL = "/" + strings.Replace(relPath, string(filepath.Separator), "/", -1)
		}
	}
	
	return content, nil
}

// Name returns the content type identifier
func (mc *MediaContent) Name() string {
	return "media"
}

// AdminForm generates the HTML form for uploading/replacing media
func (mc *MediaContent) AdminForm() (string, error) {
	form := `
<div class="admin-form">
	<h3>Manage Media</h3>
	<form method="post" action="/admin/save" enctype="multipart/form-data">
		<input type="hidden" name="content_type" value="media">
		<input type="hidden" name="content_path" value="` + mc.Path + `">
		<div>
			<label for="file">Upload New File:</label>
			<input type="file" name="file" id="file">
		</div>
		<div>
			<label>Current File:</label>
			` + func() string {
				if mc.URL != "" {
					return `<img src="` + mc.URL + `" alt="Current media" style="max-width: 300px;">`
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

// Save saves the media file (handles multipart upload in a real implementation)
func (mc *MediaContent) Save(formData map[string]string) error {
	// This would handle file upload in a real implementation
	// For now, we just return nil
	return nil
}

// Load loads the media URL
func (mc *MediaContent) Load() (interface{}, error) {
	if mc.URL == "" {
		// Generate URL from path if not set
		relPath, err := filepath.Rel("_source", mc.Path)
		if err != nil {
			return "", err
		}
		
		mc.URL = "/" + strings.Replace(relPath, string(filepath.Separator), "/", -1)
	}
	
	return mc.URL, nil
}

// TemplateName returns the name of the template to use for displaying this content
func (mc *MediaContent) TemplateName() string {
	return "media"
}