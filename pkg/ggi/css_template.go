package ggi

import (
	"text/template"
)

// CSSTemplate is a default CSS template that can be customized
const CSSTemplate = `
<style>
    body { 
        font-family: {{.FontFamily}}; 
        font-size: {{.FontSize}}; 
        margin: 40px; 
        background-color: {{.BackgroundColor}};
        color: {{.TextColor}};
    }
    .container { 
        max-width: 800px; 
        margin: 0 auto; 
    }
    h1, h2, h3, h4, h5, h6 {
        color: {{.PrimaryColor}};
    }
    a {
        color: {{.PrimaryColor}};
    }
    a:hover {
        color: {{.SecondaryColor}};
    }
    .modal { 
        display: none; 
        position: fixed; 
        z-index: 1; 
        left: 0; 
        top: 0; 
        width: 100%; 
        height: 100%; 
        background-color: rgba(0,0,0,0.4); 
    }
    .modal-content { 
        background-color: #fefefe; 
        margin: 15% auto; 
        padding: 20px; 
        border: 1px solid #888; 
        width: 80%; 
    }
    .close { 
        color: #aaa; 
        float: right; 
        font-size: 28px; 
        font-weight: bold; 
    }
    .close:hover, .close:focus { 
        color: black; 
        text-decoration: none; 
        cursor: pointer; 
    }
    .secret-key-form { 
        margin-top: 20px; 
    }
    .secret-key-form input { 
        padding: 5px; 
        margin: 5px; 
        width: 300px; 
    }
    .secret-key-form button { 
        padding: 8px 15px; 
        margin: 5px; 
    }
    .form-field {
        margin: 10px 0;
    }
    .form-field label {
        display: block;
        margin-bottom: 5px;
        font-weight: bold;
    }
    .form-field input, .form-field textarea, .form-field select {
        width: 100%;
        padding: 8px;
        border: 1px solid #ddd;
        border-radius: 4px;
        box-sizing: border-box;
    }
</style>
`

// ParseCSSTemplate parses the CSS template
func ParseCSSTemplate() (*template.Template, error) {
	return template.New("css").Parse(CSSTemplate)
}