package ggi

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// JSONContent represents JSON content with its path
type JSONContent struct {
	Path string
	Data interface{}
}

// NewJSONContent creates a new JSON content instance
func NewJSONContent(path string, data interface{}) (*JSONContent, error) {
	content := &JSONContent{Path: path, Data: data}
	
	// Load existing content if the file exists
	if _, err := os.Stat(path); err == nil {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&content.Data)
		if err != nil {
			return nil, err
		}
	}
	
	return content, nil
}

// Name returns the content type identifier
func (jc *JSONContent) Name() string {
	return "json"
}

// AdminForm generates the HTML form for editing this content
func (jc *JSONContent) AdminForm() (string, error) {
	// Generate form fields based on the data structure
	formHTML := generateJSONForm(jc.Data)
	
	form := `
<div class="admin-form">
	<h3>Edit JSON Content</h3>
	<form method="post" action="/admin/save">
		<input type="hidden" name="content_type" value="json">
		<input type="hidden" name="content_path" value="` + jc.Path + `">
		` + formHTML + `
		<button type="submit">Save Data</button>
	</form>
</div>
`
	return form, nil
}

// generateJSONForm creates form fields based on the JSON data structure
func generateJSONForm(data interface{}) string {
	v := reflect.ValueOf(data)
	
	// Handle slice of objects (common case for content collections)
	if v.Kind() == reflect.Slice {
		var result string
		for i := 0; i < v.Len(); i++ {
			elem := v.Index(i)
			if elem.Kind() == reflect.Interface {
				elem = elem.Elem()
			}
			
			if elem.Kind() == reflect.Struct {
				result += fmt.Sprintf("<div class=\"json-item\">\n<h4>Item %d</h4>\n", i+1)
				result += generateStructForm(elem, fmt.Sprintf("item_%d_", i))
				result += "</div>\n"
			}
		}
		return result
	}
	
	// Handle single struct/object
	if v.Kind() == reflect.Struct {
		return generateStructForm(v, "")
	}
	
	// Default: simple textarea for primitive values
	return `<div><textarea name="json_data" rows="10" cols="80">` + 
		fmt.Sprintf("%v", data) + 
		`</textarea></div>`
}

// generateStructForm creates form fields for a struct
func generateStructForm(v reflect.Value, prefix string) string {
	t := v.Type()
	var fields []string
	
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)
		
		// Skip unexported fields
		if !field.IsExported() {
			continue
		}
		
		fieldName := field.Name
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" && jsonTag != "-" {
			parts := strings.Split(jsonTag, ",")
			if parts[0] != "" {
				fieldName = parts[0]
			}
		}
		
		// Determine input type based on Go type
		inputType := "text"
		switch field.Type.Name() {
		case "bool":
			inputType = "checkbox"
		case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
			inputType = "number"
		case "float32", "float64":
			inputType = "number"
		}
		
		// Get the current value
		var currentValue string
		switch fieldValue.Kind() {
		case reflect.String:
			currentValue = fieldValue.String()
		case reflect.Bool:
			if fieldValue.Bool() {
				currentValue = "true"
			} else {
				currentValue = "false"
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			currentValue = fmt.Sprintf("%d", fieldValue.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			currentValue = fmt.Sprintf("%d", fieldValue.Uint())
		case reflect.Float32, reflect.Float64:
			currentValue = fmt.Sprintf("%g", fieldValue.Float())
		default:
			currentValue = fmt.Sprintf("%v", fieldValue.Interface())
		}
		
		// Create the form field HTML
		var fieldHTML string
		if inputType == "checkbox" {
			checked := ""
			if currentValue == "true" {
				checked = "checked"
			}
			fieldHTML = fmt.Sprintf(`
<div class="form-field">
    <label>
        <input type="checkbox" name="%s" value="true" %s> %s
    </label>
</div>
`, prefix+fieldName, checked, fieldName)
		} else {
			fieldHTML = fmt.Sprintf(`
<div class="form-field">
    <label for="%s">%s:</label>
    <input type="%s" id="%s" name="%s" value="%s">
</div>
`, prefix+fieldName, fieldName, inputType, prefix+fieldName, prefix+fieldName, currentValue)
		}
		
		fields = append(fields, fieldHTML)
	}
	
	return strings.Join(fields, "\n")
}

// Save saves the content to its file
func (jc *JSONContent) Save(formData map[string]string) error {
	// In a real implementation, this would convert the form data back to the appropriate JSON structure
	// For now, we'll just store the form data and marshal it when needed
	
	// For this example, we'll just save a basic JSON representation
	// In a real implementation, you'd parse the form fields back to your struct
	
	// Ensure the directory exists
	dir := filepath.Dir(jc.Path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	file, err := os.Create(jc.Path)
	if err != nil {
		return err
	}
	defer file.Close()
	
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(formData)
}

// Load loads content from its file
func (jc *JSONContent) Load() (interface{}, error) {
	// Load from file if not already loaded
	if jc.Data == nil {
		content, err := os.ReadFile(jc.Path)
		if err != nil {
			return nil, err
		}
		
		var data interface{}
		if err := json.Unmarshal(content, &data); err != nil {
			return nil, err
		}
		
		jc.Data = data
	}
	
	return jc.Data, nil
}

// TemplateName returns the name of the template to use for displaying this content
func (jc *JSONContent) TemplateName() string {
	return "json"
}