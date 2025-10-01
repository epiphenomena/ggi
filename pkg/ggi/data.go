package ggi

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// SaveData saves a list of items to a JSON file
func SaveData(data interface{}, filePath string) error {
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
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	
	return json.Unmarshal(content, &data)
}