package output

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

// Formatter is responsible for formatting output data.
type Formatter struct {
	indentSize int
}

// NewFormatter creates a new formatter.
func NewFormatter() *Formatter {
	return &Formatter{
		indentSize: 2,
	}
}

// FormatJSON formats data as JSON.
func (f *Formatter) FormatJSON(data interface{}) ([]byte, error) {
	if data == nil {
		return nil, fmt.Errorf("data cannot be nil")
	}

	return json.MarshalIndent(data, "", fmt.Sprintf("%*s", f.indentSize, ""))
}

// FormatYAML formats data as YAML.
func (f *Formatter) FormatYAML(data interface{}) ([]byte, error) {
	if data == nil {
		return nil, fmt.Errorf("data cannot be nil")
	}

	return yaml.Marshal(data)
}

// FormatText formats data as plain text.
func (f *Formatter) FormatText(data string) ([]byte, error) {
	if data == "" {
		return nil, fmt.Errorf("data cannot be empty")
	}

	return []byte(data), nil
}

// ValidateJSON validates JSON data.
func (f *Formatter) ValidateJSON(data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("data cannot be empty")
	}

	var result interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	return nil
}

// ValidateYAML validates YAML data.
func (f *Formatter) ValidateYAML(data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("data cannot be empty")
	}

	var result interface{}
	if err := yaml.Unmarshal(data, &result); err != nil {
		return fmt.Errorf("invalid YAML: %w", err)
	}

	return nil
}

// PrettyPrintJSON pretty prints JSON data.
func (f *Formatter) PrettyPrintJSON(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("data cannot be empty")
	}

	var result interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	return json.MarshalIndent(result, "", fmt.Sprintf("%*s", f.indentSize, ""))
}

// PrettyPrintYAML pretty prints YAML data.
func (f *Formatter) PrettyPrintYAML(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("data cannot be empty")
	}

	var result interface{}
	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("invalid YAML: %w", err)
	}

	return yaml.Marshal(result)
}

// SetIndentSize sets the indent size for formatting.
func (f *Formatter) SetIndentSize(size int) error {
	if size < 0 {
		return fmt.Errorf("indent size cannot be negative")
	}

	if size > 8 {
		return fmt.Errorf("indent size cannot be greater than 8")
	}

	f.indentSize = size
	return nil
}

// GetIndentSize returns the current indent size.
func (f *Formatter) GetIndentSize() int {
	return f.indentSize
}

// ConvertJSONToYAML converts JSON to YAML.
func (f *Formatter) ConvertJSONToYAML(jsonData []byte) ([]byte, error) {
	if len(jsonData) == 0 {
		return nil, fmt.Errorf("JSON data cannot be empty")
	}

	var result interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	return yaml.Marshal(result)
}

// ConvertYAMLToJSON converts YAML to JSON.
func (f *Formatter) ConvertYAMLToJSON(yamlData []byte) ([]byte, error) {
	if len(yamlData) == 0 {
		return nil, fmt.Errorf("YAML data cannot be empty")
	}

	var result interface{}
	if err := yaml.Unmarshal(yamlData, &result); err != nil {
		return nil, fmt.Errorf("invalid YAML: %w", err)
	}

	return json.MarshalIndent(result, "", fmt.Sprintf("%*s", f.indentSize, ""))
}
