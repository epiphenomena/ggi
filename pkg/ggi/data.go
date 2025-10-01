package ggi

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// SaveData saves a list of items to a JSON file
func SaveData(data interface{}, filePath string) error {
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
	return encoder.Encode(data)
}

// LoadData loads a list of items from a JSON file
func LoadData(data interface{}, filePath string) error {
	// Ensure the file is in a safe location
	if !isSafePath(filePath) {
		return fmt.Errorf("unsafe file path: %s", filePath)
	}
	
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	
	return json.Unmarshal(content, data)
}

// GenerateFormFromStruct generates an HTML form based on a struct
func GenerateFormFromStruct(data interface{}, action string) (string, error) {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return "", fmt.Errorf("data must be a struct or pointer to struct")
	}
	
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
		fieldType := field.Type.Name()
		
		// Determine input type based on Go type
		inputType := "text"
		switch fieldType {
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
		
		fieldLabel := fieldName
		if jsonTag := field.Tag.Get("json"); jsonTag != "" {
			parts := strings.Split(jsonTag, ",")
			if parts[0] != "" {
				fieldLabel = parts[0]
			}
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
`, fieldName, checked, fieldLabel)
		} else {
			fieldHTML = fmt.Sprintf(`
<div class="form-field">
    <label for="%s">%s:</label>
    <input type="%s" id="%s" name="%s" value="%s">
</div>
`, fieldName, fieldLabel, inputType, fieldName, fieldName, currentValue)
		}
		
		fields = append(fields, fieldHTML)
	}
	
	formHTML := fmt.Sprintf(`
<form method="post" action="%s">
    %s
    <input type="hidden" name="token" value="">
    <input type="submit" value="Submit">
</form>
`, action, strings.Join(fields, "\n"))
	
	return formHTML, nil
}

// RenderDataFormPage renders a page with a form for editing structured data
func RenderDataFormPage(w io.Writer, title, action string, data interface{}) error {
	formHTML, err := GenerateFormFromStruct(data, action)
	if err != nil {
		return err
	}
	
	tmpl := `
{{define "content"}}
<h2>` + title + `</h2>
<div class="form-container">
	` + formHTML + `
</div>
{{end}}
`
	
	parsedTmpl, err := ParseTemplate(tmpl)
	if err != nil {
		return err
	}
	
	context := map[string]interface{}{
		"Title": title,
		"Content": template.HTML(formHTML),
	}
	
	return parsedTmpl.ExecuteTemplate(w, "base", context)
}