package swagger

import (
	"encoding/json"
	"fmt"

	"github.com/neglet30/swag-gen/pkg/parser"
	"gopkg.in/yaml.v3"
)

// Builder is responsible for building Swagger/OpenAPI documentation.
type Builder struct {
	doc *SwaggerDoc
}

// NewBuilder creates a new Swagger builder with the given title, version, and description.
func NewBuilder(title, version, description string) *Builder {
	return &Builder{
		doc: &SwaggerDoc{
			OpenAPI: "3.0.0",
			Info: Info{
				Title:       title,
				Version:     version,
				Description: description,
			},
			Paths: make(map[string]PathItem),
			Components: Components{
				Schemas: make(map[string]*Schema),
			},
			Tags: make([]Tag, 0),
		},
	}
}

// AddEndpoint adds an endpoint to the Swagger documentation.
func (b *Builder) AddEndpoint(endpoint *parser.Endpoint) error {
	if endpoint == nil {
		return fmt.Errorf("endpoint cannot be nil")
	}

	if endpoint.Path == "" {
		return fmt.Errorf("endpoint path cannot be empty")
	}

	if endpoint.Method == "" {
		return fmt.Errorf("endpoint method cannot be empty")
	}

	// Get or create path item
	pathItem, exists := b.doc.Paths[endpoint.Path]
	if !exists {
		pathItem = PathItem{}
	}

	// Create operation
	operation := &Operation{
		Summary:     endpoint.Summary,
		Description: endpoint.Description,
		Tags:        endpoint.Tags,
		Responses:   make(map[string]Response),
		Deprecated:  endpoint.Deprecated,
	}

	// Add parameters
	if len(endpoint.Parameters) > 0 {
		operation.Parameters = make([]Parameter, 0, len(endpoint.Parameters))
		for _, param := range endpoint.Parameters {
			operation.Parameters = append(operation.Parameters, Parameter{
				Name:        param.Name,
				In:          param.In,
				Description: param.Description,
				Required:    param.Required,
				Schema: &Schema{
					Type: param.Type,
				},
			})
		}
	}

	// Add responses
	for statusCode, response := range endpoint.Responses {
		operation.Responses[statusCode] = Response{
			Description: response.Description,
		}
	}

	// Set operation on path item based on method
	method := endpoint.Method
	switch method {
	case "GET":
		pathItem.Get = operation
	case "POST":
		pathItem.Post = operation
	case "PUT":
		pathItem.Put = operation
	case "DELETE":
		pathItem.Delete = operation
	case "PATCH":
		pathItem.Patch = operation
	case "HEAD":
		pathItem.Head = operation
	case "OPTIONS":
		pathItem.Options = operation
	case "TRACE":
		pathItem.Trace = operation
	default:
		return fmt.Errorf("unsupported HTTP method: %s", method)
	}

	// Update path item
	b.doc.Paths[endpoint.Path] = pathItem

	// Add tags if not already present
	for _, tag := range endpoint.Tags {
		if !b.hasTag(tag) {
			b.doc.Tags = append(b.doc.Tags, Tag{Name: tag})
		}
	}

	return nil
}

// hasTag checks if a tag already exists in the document.
func (b *Builder) hasTag(tagName string) bool {
	for _, tag := range b.doc.Tags {
		if tag.Name == tagName {
			return true
		}
	}
	return false
}

// AddSchema adds a schema definition to the components section.
func (b *Builder) AddSchema(name string, schema *Schema) error {
	if name == "" {
		return fmt.Errorf("schema name cannot be empty")
	}

	if schema == nil {
		return fmt.Errorf("schema cannot be nil")
	}

	b.doc.Components.Schemas[name] = schema
	return nil
}

// SetInfo sets the API information.
func (b *Builder) SetInfo(info Info) {
	b.doc.Info = info
}

// AddServer adds a server to the documentation.
func (b *Builder) AddServer(server Server) {
	b.doc.Servers = append(b.doc.Servers, server)
}

// Build returns the built Swagger document.
func (b *Builder) Build() *SwaggerDoc {
	return b.doc
}

// ToJSON converts the Swagger document to JSON format.
func (b *Builder) ToJSON() ([]byte, error) {
	return json.MarshalIndent(b.doc, "", "  ")
}

// ToYAML converts the Swagger document to YAML format.
func (b *Builder) ToYAML() ([]byte, error) {
	return yaml.Marshal(b.doc)
}

// GetDocument returns the underlying Swagger document.
func (b *Builder) GetDocument() *SwaggerDoc {
	return b.doc
}
