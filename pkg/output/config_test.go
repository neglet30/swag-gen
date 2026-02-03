package output

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig("Test API", "1.0.0", "Test Description")

	assert.NotNil(t, config)
	assert.Equal(t, "Test API", config.Project.Name)
	assert.Equal(t, "1.0.0", config.Project.Version)
	assert.Equal(t, "Test Description", config.Project.Description)
	assert.Equal(t, "./api", config.Parser.Path)
	assert.Equal(t, "./docs", config.Output.Path)
	assert.Equal(t, "json", config.Output.Format)
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid config",
			config:  NewConfig("Test", "1.0.0", ""),
			wantErr: false,
		},
		{
			name: "empty project name",
			config: &Config{
				Project: ProjectConfig{Name: "", Version: "1.0.0"},
				Parser:  ParserConfig{Path: "./api"},
				Output:  OutputConfig{Path: "./docs", Format: "json"},
				Swagger: SwaggerConfig{Title: "Test", Version: "1.0.0"},
			},
			wantErr: true,
			errMsg:  "project name cannot be empty",
		},
		{
			name: "empty project version",
			config: &Config{
				Project: ProjectConfig{Name: "Test", Version: ""},
				Parser:  ParserConfig{Path: "./api"},
				Output:  OutputConfig{Path: "./docs", Format: "json"},
				Swagger: SwaggerConfig{Title: "Test", Version: "1.0.0"},
			},
			wantErr: true,
			errMsg:  "project version cannot be empty",
		},
		{
			name: "empty parser path",
			config: &Config{
				Project: ProjectConfig{Name: "Test", Version: "1.0.0"},
				Parser:  ParserConfig{Path: ""},
				Output:  OutputConfig{Path: "./docs", Format: "json"},
				Swagger: SwaggerConfig{Title: "Test", Version: "1.0.0"},
			},
			wantErr: true,
			errMsg:  "parser path cannot be empty",
		},
		{
			name: "empty output path",
			config: &Config{
				Project: ProjectConfig{Name: "Test", Version: "1.0.0"},
				Parser:  ParserConfig{Path: "./api"},
				Output:  OutputConfig{Path: "", Format: "json"},
				Swagger: SwaggerConfig{Title: "Test", Version: "1.0.0"},
			},
			wantErr: true,
			errMsg:  "output path cannot be empty",
		},
		{
			name: "invalid output format",
			config: &Config{
				Project: ProjectConfig{Name: "Test", Version: "1.0.0"},
				Parser:  ParserConfig{Path: "./api"},
				Output:  OutputConfig{Path: "./docs", Format: "xml"},
				Swagger: SwaggerConfig{Title: "Test", Version: "1.0.0"},
			},
			wantErr: true,
			errMsg:  "output format must be json or yaml",
		},
		{
			name: "empty swagger title",
			config: &Config{
				Project: ProjectConfig{Name: "Test", Version: "1.0.0"},
				Parser:  ParserConfig{Path: "./api"},
				Output:  OutputConfig{Path: "./docs", Format: "json"},
				Swagger: SwaggerConfig{Title: "", Version: "1.0.0"},
			},
			wantErr: true,
			errMsg:  "swagger title cannot be empty",
		},
		{
			name: "empty swagger version",
			config: &Config{
				Project: ProjectConfig{Name: "Test", Version: "1.0.0"},
				Parser:  ParserConfig{Path: "./api"},
				Output:  OutputConfig{Path: "./docs", Format: "json"},
				Swagger: SwaggerConfig{Title: "Test", Version: ""},
			},
			wantErr: true,
			errMsg:  "swagger version cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestConfigToYAML(t *testing.T) {
	config := NewConfig("Test API", "1.0.0", "Test Description")

	data, err := config.ToYAML()
	require.NoError(t, err)

	assert.NotEmpty(t, data)
	assert.Contains(t, string(data), "Test API")
	assert.Contains(t, string(data), "1.0.0")
}

func TestConfigFromYAML(t *testing.T) {
	yamlData := []byte(`project:
  name: Test API
  version: 1.0.0
  description: Test Description
parser:
  path: ./api
  exclude:
    - vendor
    - test
output:
  path: ./docs
  format: json
swagger:
  title: Test API
  version: 1.0.0
  description: Test Description
  basePath: /api/v1
`)

	config, err := FromYAML(yamlData)
	require.NoError(t, err)

	assert.Equal(t, "Test API", config.Project.Name)
	assert.Equal(t, "1.0.0", config.Project.Version)
	assert.Equal(t, "./api", config.Parser.Path)
	assert.Equal(t, "./docs", config.Output.Path)
	assert.Equal(t, "json", config.Output.Format)
}

func TestConfigFromYAML_Invalid(t *testing.T) {
	invalidYAML := []byte(`project:
  name: Test
  version: 1.0.0
  description: Test
parser:
  path: ./api
  exclude:
    - vendor
    - test
output:
  path: ./docs
  format: json
swagger:
  title: Test
  version: 1.0.0
  description: Test
  basePath: /api/v1
invalid_field: [unclosed`)

	_, err := FromYAML(invalidYAML)
	assert.Error(t, err)
}

func TestConfigSetProjectInfo(t *testing.T) {
	config := NewConfig("Old Name", "0.1.0", "Old Description")

	config.SetProjectInfo("New Name", "2.0.0", "New Description")

	assert.Equal(t, "New Name", config.Project.Name)
	assert.Equal(t, "2.0.0", config.Project.Version)
	assert.Equal(t, "New Description", config.Project.Description)
}

func TestConfigSetParserPath(t *testing.T) {
	config := NewConfig("Test", "1.0.0", "")

	err := config.SetParserPath("./new/path")
	require.NoError(t, err)
	assert.Equal(t, "./new/path", config.Parser.Path)
}

func TestConfigSetParserPath_Empty(t *testing.T) {
	config := NewConfig("Test", "1.0.0", "")

	err := config.SetParserPath("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "parser path cannot be empty")
}

func TestConfigSetOutputPath(t *testing.T) {
	config := NewConfig("Test", "1.0.0", "")

	err := config.SetOutputPath("./new/output")
	require.NoError(t, err)
	assert.Equal(t, "./new/output", config.Output.Path)
}

func TestConfigSetOutputPath_Empty(t *testing.T) {
	config := NewConfig("Test", "1.0.0", "")

	err := config.SetOutputPath("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "output path cannot be empty")
}

func TestConfigSetOutputFormat(t *testing.T) {
	tests := []struct {
		name    string
		format  string
		wantErr bool
	}{
		{"json", "json", false},
		{"yaml", "yaml", false},
		{"yml", "yml", false},
		{"invalid", "xml", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := NewConfig("Test", "1.0.0", "")
			err := config.SetOutputFormat(tt.format)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.format, config.Output.Format)
			}
		})
	}
}

func TestConfigSetSwaggerInfo(t *testing.T) {
	config := NewConfig("Test", "1.0.0", "")

	config.SetSwaggerInfo("New Title", "2.0.0", "New Description")

	assert.Equal(t, "New Title", config.Swagger.Title)
	assert.Equal(t, "2.0.0", config.Swagger.Version)
	assert.Equal(t, "New Description", config.Swagger.Description)
}

func TestConfigSetSwaggerBasePath(t *testing.T) {
	config := NewConfig("Test", "1.0.0", "")

	config.SetSwaggerBasePath("/api/v2")

	assert.Equal(t, "/api/v2", config.Swagger.BasePath)
}

func TestConfigAddExcludePath(t *testing.T) {
	config := NewConfig("Test", "1.0.0", "")

	initialCount := len(config.Parser.Exclude)

	config.AddExcludePath("build")
	assert.Equal(t, initialCount+1, len(config.Parser.Exclude))
	assert.Contains(t, config.Parser.Exclude, "build")
}

func TestConfigAddExcludePath_Empty(t *testing.T) {
	config := NewConfig("Test", "1.0.0", "")

	initialCount := len(config.Parser.Exclude)

	config.AddExcludePath("")
	assert.Equal(t, initialCount, len(config.Parser.Exclude))
}

func TestConfigRemoveExcludePath(t *testing.T) {
	config := NewConfig("Test", "1.0.0", "")

	config.AddExcludePath("build")
	assert.Contains(t, config.Parser.Exclude, "build")

	config.RemoveExcludePath("build")
	assert.NotContains(t, config.Parser.Exclude, "build")
}

func TestConfigGetExcludePaths(t *testing.T) {
	config := NewConfig("Test", "1.0.0", "")

	paths := config.GetExcludePaths()
	assert.NotEmpty(t, paths)
	assert.Contains(t, paths, "vendor")
	assert.Contains(t, paths, "test")
}

func TestConfigRoundTrip(t *testing.T) {
	// Create config
	config1 := NewConfig("Test API", "1.0.0", "Test Description")
	config1.SetParserPath("./src")
	config1.SetOutputPath("./output")
	config1.SetOutputFormat("yaml")
	config1.AddExcludePath("build")

	// Convert to YAML
	yamlData, err := config1.ToYAML()
	require.NoError(t, err)

	// Convert back from YAML
	config2, err := FromYAML(yamlData)
	require.NoError(t, err)

	// Verify they match
	assert.Equal(t, config1.Project.Name, config2.Project.Name)
	assert.Equal(t, config1.Project.Version, config2.Project.Version)
	assert.Equal(t, config1.Parser.Path, config2.Parser.Path)
	assert.Equal(t, config1.Output.Path, config2.Output.Path)
	assert.Equal(t, config1.Output.Format, config2.Output.Format)
}
