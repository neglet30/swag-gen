package output

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFormatter(t *testing.T) {
	formatter := NewFormatter()

	assert.NotNil(t, formatter)
	assert.Equal(t, 2, formatter.GetIndentSize())
}

func TestFormatterFormatJSON(t *testing.T) {
	formatter := NewFormatter()

	data := map[string]interface{}{
		"name":    "Test",
		"version": "1.0.0",
	}

	result, err := formatter.FormatJSON(data)
	require.NoError(t, err)

	assert.NotEmpty(t, result)
	assert.Contains(t, string(result), "name")
	assert.Contains(t, string(result), "Test")
}

func TestFormatterFormatJSON_Nil(t *testing.T) {
	formatter := NewFormatter()

	_, err := formatter.FormatJSON(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "data cannot be nil")
}

func TestFormatterFormatYAML(t *testing.T) {
	formatter := NewFormatter()

	data := map[string]interface{}{
		"name":    "Test",
		"version": "1.0.0",
	}

	result, err := formatter.FormatYAML(data)
	require.NoError(t, err)

	assert.NotEmpty(t, result)
	assert.Contains(t, string(result), "name")
	assert.Contains(t, string(result), "Test")
}

func TestFormatterFormatYAML_Nil(t *testing.T) {
	formatter := NewFormatter()

	_, err := formatter.FormatYAML(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "data cannot be nil")
}

func TestFormatterFormatText(t *testing.T) {
	formatter := NewFormatter()

	text := "This is a test"
	result, err := formatter.FormatText(text)
	require.NoError(t, err)

	assert.Equal(t, []byte(text), result)
}

func TestFormatterFormatText_Empty(t *testing.T) {
	formatter := NewFormatter()

	_, err := formatter.FormatText("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "data cannot be empty")
}

func TestFormatterValidateJSON(t *testing.T) {
	formatter := NewFormatter()

	validJSON := []byte(`{"name": "Test", "version": "1.0.0"}`)
	err := formatter.ValidateJSON(validJSON)
	require.NoError(t, err)
}

func TestFormatterValidateJSON_Invalid(t *testing.T) {
	formatter := NewFormatter()

	invalidJSON := []byte(`{invalid json}`)
	err := formatter.ValidateJSON(invalidJSON)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid JSON")
}

func TestFormatterValidateJSON_Empty(t *testing.T) {
	formatter := NewFormatter()

	err := formatter.ValidateJSON([]byte{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "data cannot be empty")
}

func TestFormatterValidateYAML(t *testing.T) {
	formatter := NewFormatter()

	validYAML := []byte(`name: Test
version: 1.0.0`)
	err := formatter.ValidateYAML(validYAML)
	require.NoError(t, err)
}

func TestFormatterValidateYAML_Invalid(t *testing.T) {
	formatter := NewFormatter()

	invalidYAML := []byte(`key: [unclosed`)
	err := formatter.ValidateYAML(invalidYAML)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid YAML")
}

func TestFormatterValidateYAML_Empty(t *testing.T) {
	formatter := NewFormatter()

	err := formatter.ValidateYAML([]byte{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "data cannot be empty")
}

func TestFormatterPrettyPrintJSON(t *testing.T) {
	formatter := NewFormatter()

	jsonData := []byte(`{"name":"Test","version":"1.0.0"}`)
	result, err := formatter.PrettyPrintJSON(jsonData)
	require.NoError(t, err)

	assert.NotEmpty(t, result)
	// Should have indentation
	assert.Contains(t, string(result), "\n")
}

func TestFormatterPrettyPrintJSON_Invalid(t *testing.T) {
	formatter := NewFormatter()

	invalidJSON := []byte(`{invalid}`)
	_, err := formatter.PrettyPrintJSON(invalidJSON)
	assert.Error(t, err)
}

func TestFormatterPrettyPrintJSON_Empty(t *testing.T) {
	formatter := NewFormatter()

	_, err := formatter.PrettyPrintJSON([]byte{})
	assert.Error(t, err)
}

func TestFormatterPrettyPrintYAML(t *testing.T) {
	formatter := NewFormatter()

	yamlData := []byte(`name: Test
version: 1.0.0`)
	result, err := formatter.PrettyPrintYAML(yamlData)
	require.NoError(t, err)

	assert.NotEmpty(t, result)
}

func TestFormatterPrettyPrintYAML_Invalid(t *testing.T) {
	formatter := NewFormatter()

	invalidYAML := []byte(`key: [unclosed`)
	_, err := formatter.PrettyPrintYAML(invalidYAML)
	assert.Error(t, err)
}

func TestFormatterPrettyPrintYAML_Empty(t *testing.T) {
	formatter := NewFormatter()

	_, err := formatter.PrettyPrintYAML([]byte{})
	assert.Error(t, err)
}

func TestFormatterSetIndentSize(t *testing.T) {
	formatter := NewFormatter()

	err := formatter.SetIndentSize(4)
	require.NoError(t, err)
	assert.Equal(t, 4, formatter.GetIndentSize())
}

func TestFormatterSetIndentSize_Invalid(t *testing.T) {
	tests := []struct {
		name    string
		size    int
		wantErr bool
		errMsg  string
	}{
		{"negative size", -1, true, "cannot be negative"},
		{"too large", 9, true, "cannot be greater than 8"},
		{"valid size 0", 0, false, ""},
		{"valid size 8", 8, false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewFormatter()
			err := formatter.SetIndentSize(tt.size)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.size, formatter.GetIndentSize())
			}
		})
	}
}

func TestFormatterConvertJSONToYAML(t *testing.T) {
	formatter := NewFormatter()

	jsonData := []byte(`{"name":"Test","version":"1.0.0"}`)
	result, err := formatter.ConvertJSONToYAML(jsonData)
	require.NoError(t, err)

	assert.NotEmpty(t, result)
	assert.Contains(t, string(result), "name")
	assert.Contains(t, string(result), "Test")
}

func TestFormatterConvertJSONToYAML_Invalid(t *testing.T) {
	formatter := NewFormatter()

	invalidJSON := []byte(`{invalid}`)
	_, err := formatter.ConvertJSONToYAML(invalidJSON)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid JSON")
}

func TestFormatterConvertJSONToYAML_Empty(t *testing.T) {
	formatter := NewFormatter()

	_, err := formatter.ConvertJSONToYAML([]byte{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "JSON data cannot be empty")
}

func TestFormatterConvertYAMLToJSON(t *testing.T) {
	formatter := NewFormatter()

	yamlData := []byte(`name: Test
version: 1.0.0`)
	result, err := formatter.ConvertYAMLToJSON(yamlData)
	require.NoError(t, err)

	assert.NotEmpty(t, result)
	assert.Contains(t, string(result), "name")
	assert.Contains(t, string(result), "Test")
}

func TestFormatterConvertYAMLToJSON_Invalid(t *testing.T) {
	formatter := NewFormatter()

	invalidYAML := []byte(`key: [unclosed`)
	_, err := formatter.ConvertYAMLToJSON(invalidYAML)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid YAML")
}

func TestFormatterConvertYAMLToJSON_Empty(t *testing.T) {
	formatter := NewFormatter()

	_, err := formatter.ConvertYAMLToJSON([]byte{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "YAML data cannot be empty")
}

func TestFormatterRoundTrip_JSONToYAMLToJSON(t *testing.T) {
	formatter := NewFormatter()

	originalJSON := []byte(`{"name":"Test","version":"1.0.0","tags":["api","v1"]}`)

	// JSON to YAML
	yaml, err := formatter.ConvertJSONToYAML(originalJSON)
	require.NoError(t, err)

	// YAML to JSON
	jsonData, err := formatter.ConvertYAMLToJSON(yaml)
	require.NoError(t, err)

	// Verify both are valid
	var original, result map[string]interface{}
	err = json.Unmarshal(originalJSON, &original)
	require.NoError(t, err)

	err = json.Unmarshal(jsonData, &result)
	require.NoError(t, err)

	assert.Equal(t, original["name"], result["name"])
	assert.Equal(t, original["version"], result["version"])
}

func TestFormatterComplexData(t *testing.T) {
	formatter := NewFormatter()

	data := map[string]interface{}{
		"project": map[string]interface{}{
			"name":    "Test API",
			"version": "1.0.0",
		},
		"paths": []string{"/users", "/posts"},
		"tags":  []string{"api", "v1"},
	}

	// Format as JSON
	jsonData, err := formatter.FormatJSON(data)
	require.NoError(t, err)

	// Validate JSON
	err = formatter.ValidateJSON(jsonData)
	require.NoError(t, err)

	// Format as YAML
	yamlData, err := formatter.FormatYAML(data)
	require.NoError(t, err)

	// Validate YAML
	err = formatter.ValidateYAML(yamlData)
	require.NoError(t, err)

	// Convert JSON to YAML
	convertedYAML, err := formatter.ConvertJSONToYAML(jsonData)
	require.NoError(t, err)
	assert.NotEmpty(t, convertedYAML)

	// Convert YAML to JSON
	convertedJSON, err := formatter.ConvertYAMLToJSON(yamlData)
	require.NoError(t, err)
	assert.NotEmpty(t, convertedJSON)
}
