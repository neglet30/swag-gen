// Package swagger provides functionality for generating OpenAPI 3.0 Swagger documentation.
package swagger

// SwaggerDoc represents an OpenAPI 3.0 document.
type SwaggerDoc struct {
	OpenAPI    string              `json:"openapi"`
	Info       Info                `json:"info"`
	Paths      map[string]PathItem `json:"paths"`
	Components Components          `json:"components,omitempty"`
	Servers    []Server            `json:"servers,omitempty"`
	Tags       []Tag               `json:"tags,omitempty"`
}

// Info contains metadata about the API.
type Info struct {
	Title       string   `json:"title"`
	Version     string   `json:"version"`
	Description string   `json:"description,omitempty"`
	Contact     *Contact `json:"contact,omitempty"`
	License     *License `json:"license,omitempty"`
}

// Contact information for the exposed API.
type Contact struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

// License information for the exposed API.
type License struct {
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}

// Server represents a server available to the API.
type Server struct {
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
}

// Tag represents a tag for grouping operations.
type Tag struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// PathItem describes the operations available on a single path.
type PathItem struct {
	Get     *Operation `json:"get,omitempty"`
	Post    *Operation `json:"post,omitempty"`
	Put     *Operation `json:"put,omitempty"`
	Delete  *Operation `json:"delete,omitempty"`
	Patch   *Operation `json:"patch,omitempty"`
	Head    *Operation `json:"head,omitempty"`
	Options *Operation `json:"options,omitempty"`
	Trace   *Operation `json:"trace,omitempty"`
}

// Operation describes a single API operation for a path and HTTP method.
type Operation struct {
	Summary     string              `json:"summary,omitempty"`
	Description string              `json:"description,omitempty"`
	Tags        []string            `json:"tags,omitempty"`
	Parameters  []Parameter         `json:"parameters,omitempty"`
	RequestBody *RequestBody        `json:"requestBody,omitempty"`
	Responses   map[string]Response `json:"responses"`
	Deprecated  bool                `json:"deprecated,omitempty"`
	OperationID string              `json:"operationId,omitempty"`
}

// Parameter describes a single operation parameter.
type Parameter struct {
	Name        string  `json:"name"`
	In          string  `json:"in"` // query, path, header, cookie
	Description string  `json:"description,omitempty"`
	Required    bool    `json:"required,omitempty"`
	Schema      *Schema `json:"schema,omitempty"`
	Deprecated  bool    `json:"deprecated,omitempty"`
}

// RequestBody describes a request body that can be used by the operation.
type RequestBody struct {
	Description string               `json:"description,omitempty"`
	Content     map[string]MediaType `json:"content"`
	Required    bool                 `json:"required,omitempty"`
}

// MediaType provides schema and examples for the media type identified by its key.
type MediaType struct {
	Schema  *Schema     `json:"schema,omitempty"`
	Example interface{} `json:"example,omitempty"`
}

// Response describes a single response from an API Operation.
type Response struct {
	Description string               `json:"description"`
	Content     map[string]MediaType `json:"content,omitempty"`
	Headers     map[string]Header    `json:"headers,omitempty"`
}

// Header represents a header parameter.
type Header struct {
	Description string  `json:"description,omitempty"`
	Schema      *Schema `json:"schema,omitempty"`
	Required    bool    `json:"required,omitempty"`
}

// Schema represents a JSON Schema.
type Schema struct {
	Type        string             `json:"type,omitempty"`
	Format      string             `json:"format,omitempty"`
	Title       string             `json:"title,omitempty"`
	Description string             `json:"description,omitempty"`
	Default     interface{}        `json:"default,omitempty"`
	Example     interface{}        `json:"example,omitempty"`
	Items       *Schema            `json:"items,omitempty"`
	Properties  map[string]*Schema `json:"properties,omitempty"`
	Required    []string           `json:"required,omitempty"`
	Enum        []interface{}      `json:"enum,omitempty"`
	Ref         string             `json:"$ref,omitempty"`
	AllOf       []*Schema          `json:"allOf,omitempty"`
	OneOf       []*Schema          `json:"oneOf,omitempty"`
	AnyOf       []*Schema          `json:"anyOf,omitempty"`
	Not         *Schema            `json:"not,omitempty"`
	Minimum     *float64           `json:"minimum,omitempty"`
	Maximum     *float64           `json:"maximum,omitempty"`
	MinLength   *int               `json:"minLength,omitempty"`
	MaxLength   *int               `json:"maxLength,omitempty"`
	Pattern     string             `json:"pattern,omitempty"`
}

// Components holds a set of reusable objects for different aspects of the OAS.
type Components struct {
	Schemas map[string]*Schema `json:"schemas,omitempty"`
}
