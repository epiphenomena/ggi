package ggi

import (
	"html/template"
	"os"
)

// IsCGI checks if the program is running as a CGI script
func IsCGI() bool {
	// Check for common CGI environment variables
	_, hasRequestMethod := os.LookupEnv("REQUEST_METHOD")
	_, hasScriptName := os.LookupEnv("SCRIPT_NAME")
	
	return hasRequestMethod && hasScriptName
}

// CSSLink generates an HTML link tag for a CSS file
func CSSLink(path string) template.HTML {
	return template.HTML(`<link rel="stylesheet" type="text/css" href="` + path + `" />`)
}

// JSScript generates an HTML script tag for a JavaScript file
func JSScript(path string) template.HTML {
	return template.HTML(`<script src="` + path + `"></script>`)
}

// CSSStyleFile generates an HTML style tag by loading CSS content from a file
func CSSStyleFile(filePath string) template.HTML {
	content, err := os.ReadFile(filePath)
	if err != nil {
		// If file reading fails, return an empty style tag
		return template.HTML(`<style></style>`)
	}
	return template.HTML(`<style>` + string(content) + `</style>`)
}

// JSScriptContentFile generates an HTML script tag by loading JS content from a file
func JSScriptContentFile(filePath string) template.HTML {
	content, err := os.ReadFile(filePath)
	if err != nil {
		// If file reading fails, return an empty script tag
		return template.HTML(`<script></script>`)
	}
	return template.HTML(`<script>` + string(content) + `</script>`)
}



// BaseTemplate defines the base HTML template
const BaseTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>{{.Title}}</title>
    {{block "css" .}}{{end}}
    {{block "head" .}}{{end}}
    {{block "js" .}}{{end}}
</head>
<body>
    <div class="container">
        {{block "header" .}}
        <header>
            <h1>{{.Title}}</h1>
        </header>
        {{end}}
        
        {{block "content" .}}
        <main>
            {{.Content}}
        </main>
        {{end}}
        
        {{block "footer" .}}
        <footer>
            <p>&copy; 2025 GGI Site</p>
        </footer>
        {{end}}
    </div>
</body>
</html>
`

// ParseBaseTemplate parses the base template
func ParseBaseTemplate() (*template.Template, error) {
	return template.New("base").Parse(BaseTemplate)
}



// ParseTemplate parses a template with the base layout
func ParseTemplate(tmpl string) (*template.Template, error) {
	baseTmpl, err := ParseBaseTemplate()
	if err != nil {
		return nil, err
	}
	
	// Create a new template based on the base
	t := template.Must(baseTmpl.Clone())
	return template.Must(t.New("page").Parse(tmpl)), nil
}

