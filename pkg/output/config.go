package output

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// Config represents the swag-gen configuration.
type Config struct {
	Project ProjectConfig `yaml:"project"`
	Parser  ParserConfig  `yaml:"parser"`
	Output  OutputConfig  `yaml:"output"`
	Swagger SwaggerConfig `yaml:"swagger"`
}

// ProjectConfig represents project configuration.
type ProjectConfig struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
}

// ParserConfig represents parser configuration.
type ParserConfig struct {
	Path    string   `yaml:"path"`
	Exclude []string `yaml:"exclude,omitempty"`
}

// OutputConfig represents output configuration.
type OutputConfig struct {
	Path   string `yaml:"path"`
	Format string `yaml:"format"`
}

// SwaggerConfig represents Swagger configuration.
type SwaggerConfig struct {
	Title       string `yaml:"title"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
	BasePath    string `yaml:"basePath,omitempty"`
}

// NewConfig creates a new configuration.
func NewConfig(projectName, projectVersion, projectDescription string) *Config {
	return &Config{
		Project: ProjectConfig{
			Name:        projectName,
			Version:     projectVersion,
			Description: projectDescription,
		},
		Parser: ParserConfig{
			Path:    "./api",
			Exclude: []string{"vendor", "test"},
		},
		Output: OutputConfig{
			Path:   "./docs",
			Format: "json",
		},
		Swagger: SwaggerConfig{
			Title:       projectName,
			Version:     projectVersion,
			Description: projectDescription,
			BasePath:    "/api/v1",
		},
	}
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	if c.Project.Name == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	if c.Project.Version == "" {
		return fmt.Errorf("project version cannot be empty")
	}

	if c.Parser.Path == "" {
		return fmt.Errorf("parser path cannot be empty")
	}

	if c.Output.Path == "" {
		return fmt.Errorf("output path cannot be empty")
	}

	if c.Output.Format != "json" && c.Output.Format != "yaml" && c.Output.Format != "yml" {
		return fmt.Errorf("output format must be json or yaml")
	}

	if c.Swagger.Title == "" {
		return fmt.Errorf("swagger title cannot be empty")
	}

	if c.Swagger.Version == "" {
		return fmt.Errorf("swagger version cannot be empty")
	}

	return nil
}

// ToYAML converts the configuration to YAML.
func (c *Config) ToYAML() ([]byte, error) {
	return yaml.Marshal(c)
}

// FromYAML creates a configuration from YAML data.
func FromYAML(data []byte) (*Config, error) {
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return &config, nil
}

// SetProjectInfo sets project information.
func (c *Config) SetProjectInfo(name, version, description string) {
	c.Project.Name = name
	c.Project.Version = version
	c.Project.Description = description
}

// SetParserPath sets the parser path.
func (c *Config) SetParserPath(path string) error {
	if path == "" {
		return fmt.Errorf("parser path cannot be empty")
	}

	c.Parser.Path = path
	return nil
}

// SetOutputPath sets the output path.
func (c *Config) SetOutputPath(path string) error {
	if path == "" {
		return fmt.Errorf("output path cannot be empty")
	}

	c.Output.Path = path
	return nil
}

// SetOutputFormat sets the output format.
func (c *Config) SetOutputFormat(format string) error {
	if format != "json" && format != "yaml" && format != "yml" {
		return fmt.Errorf("output format must be json or yaml")
	}

	c.Output.Format = format
	return nil
}

// SetSwaggerInfo sets Swagger information.
func (c *Config) SetSwaggerInfo(title, version, description string) {
	c.Swagger.Title = title
	c.Swagger.Version = version
	c.Swagger.Description = description
}

// SetSwaggerBasePath sets the Swagger base path.
func (c *Config) SetSwaggerBasePath(basePath string) {
	c.Swagger.BasePath = basePath
}

// AddExcludePath adds a path to exclude.
func (c *Config) AddExcludePath(path string) {
	if path != "" {
		c.Parser.Exclude = append(c.Parser.Exclude, path)
	}
}

// RemoveExcludePath removes a path from exclude.
func (c *Config) RemoveExcludePath(path string) {
	for i, p := range c.Parser.Exclude {
		if p == path {
			c.Parser.Exclude = append(c.Parser.Exclude[:i], c.Parser.Exclude[i+1:]...)
			break
		}
	}
}

// GetExcludePaths returns the exclude paths.
func (c *Config) GetExcludePaths() []string {
	return c.Parser.Exclude
}
