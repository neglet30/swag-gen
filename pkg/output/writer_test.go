package output

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/neglet30/swag-gen/pkg/swagger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewWriter(t *testing.T) {
	writer := NewWriter("/tmp/output")

	assert.NotNil(t, writer)
	assert.Equal(t, "/tmp/output", writer.GetOutputPath())
}

func TestWriterWriteSwagger_JSON(t *testing.T) {
	tmpDir := t.TempDir()
	writer := NewWriter(tmpDir)

	doc := &swagger.SwaggerDoc{
		OpenAPI: "3.0.0",
		Info: swagger.Info{
			Title:       "Test API",
			Version:     "1.0.0",
			Description: "Test Description",
		},
		Paths:      make(map[string]swagger.PathItem),
		Components: swagger.Components{Schemas: make(map[string]*swagger.Schema)},
	}

	err := writer.WriteSwagger(doc, "swagger", "json")
	require.NoError(t, err)

	// Verify file exists
	filePath := filepath.Join(tmpDir, "swagger.json")
	assert.True(t, FileExists(filePath))

	// Verify file content
	data, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.NotEmpty(t, data)
	assert.Contains(t, string(data), "3.0.0")
	assert.Contains(t, string(data), "Test API")
}

func TestWriterWriteSwagger_YAML(t *testing.T) {
	tmpDir := t.TempDir()
	writer := NewWriter(tmpDir)

	doc := &swagger.SwaggerDoc{
		OpenAPI: "3.0.0",
		Info: swagger.Info{
			Title:       "Test API",
			Version:     "1.0.0",
			Description: "Test Description",
		},
		Paths:      make(map[string]swagger.PathItem),
		Components: swagger.Components{Schemas: make(map[string]*swagger.Schema)},
	}

	err := writer.WriteSwagger(doc, "swagger", "yaml")
	require.NoError(t, err)

	// Verify file exists
	filePath := filepath.Join(tmpDir, "swagger.yaml")
	assert.True(t, FileExists(filePath))

	// Verify file content
	data, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.NotEmpty(t, data)
	assert.Contains(t, string(data), "openapi: 3.0.0")
	assert.Contains(t, string(data), "title: Test API")
}

func TestWriterWriteSwagger_DefaultFormat(t *testing.T) {
	tmpDir := t.TempDir()
	writer := NewWriter(tmpDir)

	doc := &swagger.SwaggerDoc{
		OpenAPI: "3.0.0",
		Info: swagger.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
		Paths:      make(map[string]swagger.PathItem),
		Components: swagger.Components{Schemas: make(map[string]*swagger.Schema)},
	}

	err := writer.WriteSwagger(doc, "swagger", "")
	require.NoError(t, err)

	// Should default to JSON
	filePath := filepath.Join(tmpDir, "swagger.json")
	assert.True(t, FileExists(filePath))
}

func TestWriterWriteSwagger_Errors(t *testing.T) {
	tests := []struct {
		name     string
		doc      *swagger.SwaggerDoc
		filename string
		format   string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "nil document",
			doc:      nil,
			filename: "swagger",
			format:   "json",
			wantErr:  true,
			errMsg:   "swagger document cannot be nil",
		},
		{
			name:     "empty filename",
			doc:      &swagger.SwaggerDoc{},
			filename: "",
			format:   "json",
			wantErr:  true,
			errMsg:   "filename cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			writer := NewWriter(tmpDir)
			err := writer.WriteSwagger(tt.doc, tt.filename, tt.format)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}

func TestWriterWriteConfig(t *testing.T) {
	tmpDir := t.TempDir()
	writer := NewWriter(tmpDir)

	config := NewConfig("Test API", "1.0.0", "Test Description")

	err := writer.WriteConfig(config, "swag-gen.yaml")
	require.NoError(t, err)

	// Verify file exists
	filePath := filepath.Join(tmpDir, "swag-gen.yaml")
	assert.True(t, FileExists(filePath))

	// Verify file content
	data, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.NotEmpty(t, data)
	assert.Contains(t, string(data), "Test API")
}

func TestWriterWriteConfig_Errors(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		filename string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "nil config",
			config:   nil,
			filename: "config.yaml",
			wantErr:  true,
			errMsg:   "config cannot be nil",
		},
		{
			name:     "empty filename",
			config:   NewConfig("Test", "1.0.0", ""),
			filename: "",
			wantErr:  true,
			errMsg:   "filename cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			writer := NewWriter(tmpDir)
			err := writer.WriteConfig(tt.config, tt.filename)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}

func TestWriterWriteREADME(t *testing.T) {
	tmpDir := t.TempDir()
	writer := NewWriter(tmpDir)

	err := writer.WriteREADME("README.md", "Test API", "This is a test API")
	require.NoError(t, err)

	// Verify file exists
	filePath := filepath.Join(tmpDir, "README.md")
	assert.True(t, FileExists(filePath))

	// Verify file content
	data, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.NotEmpty(t, data)
	assert.Contains(t, string(data), "# Test API")
	assert.Contains(t, string(data), "This is a test API")
}

func TestWriterWriteREADME_DefaultTitle(t *testing.T) {
	tmpDir := t.TempDir()
	writer := NewWriter(tmpDir)

	err := writer.WriteREADME("README.md", "", "")
	require.NoError(t, err)

	// Verify file exists
	filePath := filepath.Join(tmpDir, "README.md")
	assert.True(t, FileExists(filePath))

	// Verify file content
	data, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.Contains(t, string(data), "# API Documentation")
}

func TestWriterWriteREADME_Error(t *testing.T) {
	tmpDir := t.TempDir()
	writer := NewWriter(tmpDir)

	err := writer.WriteREADME("", "Test API", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "filename cannot be empty")
}

func TestFileExists(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")

	// File doesn't exist
	assert.False(t, FileExists(filePath))

	// Create file
	err := os.WriteFile(filePath, []byte("test"), 0644)
	require.NoError(t, err)

	// File exists
	assert.True(t, FileExists(filePath))
}

func TestDirectoryExists(t *testing.T) {
	tmpDir := t.TempDir()
	dirPath := filepath.Join(tmpDir, "testdir")

	// Directory doesn't exist
	assert.False(t, DirectoryExists(dirPath))

	// Create directory
	err := os.Mkdir(dirPath, 0755)
	require.NoError(t, err)

	// Directory exists
	assert.True(t, DirectoryExists(dirPath))
}

func TestCreateDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	dirPath := filepath.Join(tmpDir, "newdir")

	// Create directory
	err := CreateDirectory(dirPath)
	require.NoError(t, err)

	// Verify directory exists
	assert.True(t, DirectoryExists(dirPath))

	// Create again (should not error)
	err = CreateDirectory(dirPath)
	require.NoError(t, err)
}

func TestRemoveFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")

	// Create file
	err := os.WriteFile(filePath, []byte("test"), 0644)
	require.NoError(t, err)

	// Verify file exists
	assert.True(t, FileExists(filePath))

	// Remove file
	err = RemoveFile(filePath)
	require.NoError(t, err)

	// Verify file doesn't exist
	assert.False(t, FileExists(filePath))

	// Remove non-existent file (should not error)
	err = RemoveFile(filePath)
	require.NoError(t, err)
}

func TestWriterCreateMultipleFiles(t *testing.T) {
	tmpDir := t.TempDir()
	writer := NewWriter(tmpDir)

	// Create Swagger document
	doc := &swagger.SwaggerDoc{
		OpenAPI: "3.0.0",
		Info: swagger.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
		Paths:      make(map[string]swagger.PathItem),
		Components: swagger.Components{Schemas: make(map[string]*swagger.Schema)},
	}

	// Write Swagger JSON
	err := writer.WriteSwagger(doc, "swagger", "json")
	require.NoError(t, err)

	// Write Swagger YAML
	err = writer.WriteSwagger(doc, "swagger", "yaml")
	require.NoError(t, err)

	// Write Config
	config := NewConfig("Test API", "1.0.0", "")
	err = writer.WriteConfig(config, "swag-gen.yaml")
	require.NoError(t, err)

	// Write README
	err = writer.WriteREADME("README.md", "Test API", "")
	require.NoError(t, err)

	// Verify all files exist
	assert.True(t, FileExists(filepath.Join(tmpDir, "swagger.json")))
	assert.True(t, FileExists(filepath.Join(tmpDir, "swagger.yaml")))
	assert.True(t, FileExists(filepath.Join(tmpDir, "swag-gen.yaml")))
	assert.True(t, FileExists(filepath.Join(tmpDir, "README.md")))
}
