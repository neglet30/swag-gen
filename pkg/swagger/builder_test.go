package swagger

import (
	"encoding/json"
	"testing"

	"github.com/neglet30/swag-gen/pkg/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBuilder(t *testing.T) {
	builder := NewBuilder("Test API", "1.0.0", "Test Description")

	assert.NotNil(t, builder)
	assert.NotNil(t, builder.doc)
	assert.Equal(t, "3.0.0", builder.doc.OpenAPI)
	assert.Equal(t, "Test API", builder.doc.Info.Title)
	assert.Equal(t, "1.0.0", builder.doc.Info.Version)
	assert.Equal(t, "Test Description", builder.doc.Info.Description)
}

func TestBuilderAddEndpoint(t *testing.T) {
	tests := []struct {
		name     string
		endpoint *parser.Endpoint
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid GET endpoint",
			endpoint: &parser.Endpoint{Path: "/users", Method: "GET", Summary: "Get users"},
			wantErr:  false,
		},
		{
			name:     "valid POST endpoint",
			endpoint: &parser.Endpoint{Path: "/users", Method: "POST", Summary: "Create user"},
			wantErr:  false,
		},
		{
			name:     "valid PUT endpoint",
			endpoint: &parser.Endpoint{Path: "/users/{id}", Method: "PUT", Summary: "Update user"},
			wantErr:  false,
		},
		{
			name:     "valid DELETE endpoint",
			endpoint: &parser.Endpoint{Path: "/users/{id}", Method: "DELETE", Summary: "Delete user"},
			wantErr:  false,
		},
		{
			name:     "valid PATCH endpoint",
			endpoint: &parser.Endpoint{Path: "/users/{id}", Method: "PATCH", Summary: "Patch user"},
			wantErr:  false,
		},
		{
			name:     "nil endpoint",
			endpoint: nil,
			wantErr:  true,
			errMsg:   "endpoint cannot be nil",
		},
		{
			name:     "empty path",
			endpoint: &parser.Endpoint{Path: "", Method: "GET"},
			wantErr:  true,
			errMsg:   "endpoint path cannot be empty",
		},
		{
			name:     "empty method",
			endpoint: &parser.Endpoint{Path: "/users", Method: ""},
			wantErr:  true,
			errMsg:   "endpoint method cannot be empty",
		},
		{
			name:     "unsupported method",
			endpoint: &parser.Endpoint{Path: "/users", Method: "INVALID"},
			wantErr:  true,
			errMsg:   "unsupported HTTP method",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewBuilder("Test API", "1.0.0", "")
			err := builder.AddEndpoint(tt.endpoint)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, builder.doc.Paths)
			}
		})
	}
}

func TestBuilderAddEndpointWithParameters(t *testing.T) {
	builder := NewBuilder("Test API", "1.0.0", "")

	endpoint := &parser.Endpoint{
		Path:    "/users/{id}",
		Method:  "GET",
		Summary: "Get user by ID",
		Parameters: []parser.Parameter{
			{Name: "id", In: "path", Type: "string", Required: true, Description: "User ID"},
			{Name: "include", In: "query", Type: "string", Required: false, Description: "Include fields"},
		},
	}

	err := builder.AddEndpoint(endpoint)
	require.NoError(t, err)

	pathItem, exists := builder.doc.Paths["/users/{id}"]
	assert.True(t, exists)
	assert.NotNil(t, pathItem.Get)
	assert.Len(t, pathItem.Get.Parameters, 2)
	assert.Equal(t, "id", pathItem.Get.Parameters[0].Name)
	assert.Equal(t, "path", pathItem.Get.Parameters[0].In)
	assert.True(t, pathItem.Get.Parameters[0].Required)
}

func TestBuilderAddEndpointWithResponses(t *testing.T) {
	builder := NewBuilder("Test API", "1.0.0", "")

	endpoint := &parser.Endpoint{
		Path:    "/users",
		Method:  "GET",
		Summary: "Get users",
		Responses: map[string]parser.Response{
			"200": {StatusCode: "200", Description: "Success"},
			"400": {StatusCode: "400", Description: "Bad Request"},
			"500": {StatusCode: "500", Description: "Internal Server Error"},
		},
	}

	err := builder.AddEndpoint(endpoint)
	require.NoError(t, err)

	pathItem, exists := builder.doc.Paths["/users"]
	assert.True(t, exists)
	assert.NotNil(t, pathItem.Get)
	assert.Len(t, pathItem.Get.Responses, 3)
	assert.Equal(t, "Success", pathItem.Get.Responses["200"].Description)
}

func TestBuilderAddEndpointWithTags(t *testing.T) {
	builder := NewBuilder("Test API", "1.0.0", "")

	endpoint := &parser.Endpoint{
		Path:    "/users",
		Method:  "GET",
		Summary: "Get users",
		Tags:    []string{"User", "Admin"},
	}

	err := builder.AddEndpoint(endpoint)
	require.NoError(t, err)

	assert.Len(t, builder.doc.Tags, 2)
	assert.Equal(t, "User", builder.doc.Tags[0].Name)
	assert.Equal(t, "Admin", builder.doc.Tags[1].Name)
}

func TestBuilderAddSchema(t *testing.T) {
	builder := NewBuilder("Test API", "1.0.0", "")

	schema := &Schema{
		Type: "object",
		Properties: map[string]*Schema{
			"id":   {Type: "integer"},
			"name": {Type: "string"},
		},
	}

	err := builder.AddSchema("User", schema)
	require.NoError(t, err)

	assert.NotNil(t, builder.doc.Components.Schemas["User"])
	assert.Equal(t, "object", builder.doc.Components.Schemas["User"].Type)
}

func TestBuilderAddSchemaErrors(t *testing.T) {
	tests := []struct {
		name    string
		schName string
		schema  *Schema
		wantErr bool
		errMsg  string
	}{
		{
			name:    "empty schema name",
			schName: "",
			schema:  &Schema{Type: "object"},
			wantErr: true,
			errMsg:  "schema name cannot be empty",
		},
		{
			name:    "nil schema",
			schName: "User",
			schema:  nil,
			wantErr: true,
			errMsg:  "schema cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewBuilder("Test API", "1.0.0", "")
			err := builder.AddSchema(tt.schName, tt.schema)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}

func TestBuilderSetInfo(t *testing.T) {
	builder := NewBuilder("Test API", "1.0.0", "")

	newInfo := Info{
		Title:       "Updated API",
		Version:     "2.0.0",
		Description: "Updated Description",
	}

	builder.SetInfo(newInfo)

	assert.Equal(t, "Updated API", builder.doc.Info.Title)
	assert.Equal(t, "2.0.0", builder.doc.Info.Version)
	assert.Equal(t, "Updated Description", builder.doc.Info.Description)
}

func TestBuilderAddServer(t *testing.T) {
	builder := NewBuilder("Test API", "1.0.0", "")

	server := Server{
		URL:         "https://api.example.com",
		Description: "Production server",
	}

	builder.AddServer(server)

	assert.Len(t, builder.doc.Servers, 1)
	assert.Equal(t, "https://api.example.com", builder.doc.Servers[0].URL)
	assert.Equal(t, "Production server", builder.doc.Servers[0].Description)
}

func TestBuilderBuild(t *testing.T) {
	builder := NewBuilder("Test API", "1.0.0", "Test Description")

	endpoint := &parser.Endpoint{
		Path:    "/users",
		Method:  "GET",
		Summary: "Get users",
	}

	err := builder.AddEndpoint(endpoint)
	require.NoError(t, err)

	doc := builder.Build()

	assert.NotNil(t, doc)
	assert.Equal(t, "3.0.0", doc.OpenAPI)
	assert.Equal(t, "Test API", doc.Info.Title)
	assert.NotEmpty(t, doc.Paths)
}

func TestBuilderToJSON(t *testing.T) {
	builder := NewBuilder("Test API", "1.0.0", "Test Description")

	endpoint := &parser.Endpoint{
		Path:    "/users",
		Method:  "GET",
		Summary: "Get users",
	}

	err := builder.AddEndpoint(endpoint)
	require.NoError(t, err)

	jsonData, err := builder.ToJSON()
	require.NoError(t, err)

	var doc SwaggerDoc
	err = json.Unmarshal(jsonData, &doc)
	require.NoError(t, err)

	assert.Equal(t, "3.0.0", doc.OpenAPI)
	assert.Equal(t, "Test API", doc.Info.Title)
	assert.NotEmpty(t, doc.Paths)
}

func TestBuilderToYAML(t *testing.T) {
	builder := NewBuilder("Test API", "1.0.0", "Test Description")

	endpoint := &parser.Endpoint{
		Path:    "/users",
		Method:  "GET",
		Summary: "Get users",
	}

	err := builder.AddEndpoint(endpoint)
	require.NoError(t, err)

	yamlData, err := builder.ToYAML()
	require.NoError(t, err)

	assert.NotEmpty(t, yamlData)
	assert.Contains(t, string(yamlData), "openapi: 3.0.0")
	assert.Contains(t, string(yamlData), "title: Test API")
}

func TestBuilderMultipleEndpoints(t *testing.T) {
	builder := NewBuilder("Test API", "1.0.0", "")

	endpoints := []parser.Endpoint{
		{Path: "/users", Method: "GET", Summary: "Get users"},
		{Path: "/users", Method: "POST", Summary: "Create user"},
		{Path: "/users/{id}", Method: "GET", Summary: "Get user"},
		{Path: "/users/{id}", Method: "PUT", Summary: "Update user"},
		{Path: "/users/{id}", Method: "DELETE", Summary: "Delete user"},
	}

	for i := range endpoints {
		err := builder.AddEndpoint(&endpoints[i])
		require.NoError(t, err)
	}

	assert.Len(t, builder.doc.Paths, 2)

	usersPath := builder.doc.Paths["/users"]
	assert.NotNil(t, usersPath.Get)
	assert.NotNil(t, usersPath.Post)

	userIDPath := builder.doc.Paths["/users/{id}"]
	assert.NotNil(t, userIDPath.Get)
	assert.NotNil(t, userIDPath.Put)
	assert.NotNil(t, userIDPath.Delete)
}

func TestBuilderGetDocument(t *testing.T) {
	builder := NewBuilder("Test API", "1.0.0", "")
	doc := builder.GetDocument()

	assert.NotNil(t, doc)
	assert.Equal(t, "3.0.0", doc.OpenAPI)
	assert.Equal(t, "Test API", doc.Info.Title)
}

func TestBuilderHasTag(t *testing.T) {
	builder := NewBuilder("Test API", "1.0.0", "")

	endpoint := &parser.Endpoint{
		Path:    "/users",
		Method:  "GET",
		Summary: "Get users",
		Tags:    []string{"User"},
	}

	err := builder.AddEndpoint(endpoint)
	require.NoError(t, err)

	assert.True(t, builder.hasTag("User"))
	assert.False(t, builder.hasTag("Admin"))
}

func TestBuilderDuplicateTags(t *testing.T) {
	builder := NewBuilder("Test API", "1.0.0", "")

	endpoint1 := &parser.Endpoint{
		Path:    "/users",
		Method:  "GET",
		Summary: "Get users",
		Tags:    []string{"User"},
	}

	endpoint2 := &parser.Endpoint{
		Path:    "/users",
		Method:  "POST",
		Summary: "Create user",
		Tags:    []string{"User"},
	}

	err := builder.AddEndpoint(endpoint1)
	require.NoError(t, err)

	err = builder.AddEndpoint(endpoint2)
	require.NoError(t, err)

	// Should only have one "User" tag
	userTagCount := 0
	for _, tag := range builder.doc.Tags {
		if tag.Name == "User" {
			userTagCount++
		}
	}
	assert.Equal(t, 1, userTagCount)
}

func TestBuilderAddEndpoint_HEAD(t *testing.T) {
	builder := NewBuilder("Test API", "1.0.0", "")

	endpoint := &parser.Endpoint{
		Path:    "/users",
		Method:  "HEAD",
		Summary: "Check users endpoint",
	}

	err := builder.AddEndpoint(endpoint)
	require.NoError(t, err)

	pathItem, exists := builder.doc.Paths["/users"]
	assert.True(t, exists)
	assert.NotNil(t, pathItem.Head)
	assert.Equal(t, "Check users endpoint", pathItem.Head.Summary)
}

func TestBuilderAddEndpoint_OPTIONS(t *testing.T) {
	builder := NewBuilder("Test API", "1.0.0", "")

	endpoint := &parser.Endpoint{
		Path:    "/users",
		Method:  "OPTIONS",
		Summary: "Get options for users",
	}

	err := builder.AddEndpoint(endpoint)
	require.NoError(t, err)

	pathItem, exists := builder.doc.Paths["/users"]
	assert.True(t, exists)
	assert.NotNil(t, pathItem.Options)
	assert.Equal(t, "Get options for users", pathItem.Options.Summary)
}

func TestBuilderAddEndpoint_TRACE(t *testing.T) {
	builder := NewBuilder("Test API", "1.0.0", "")

	endpoint := &parser.Endpoint{
		Path:    "/users",
		Method:  "TRACE",
		Summary: "Trace users endpoint",
	}

	err := builder.AddEndpoint(endpoint)
	require.NoError(t, err)

	pathItem, exists := builder.doc.Paths["/users"]
	assert.True(t, exists)
	assert.NotNil(t, pathItem.Trace)
	assert.Equal(t, "Trace users endpoint", pathItem.Trace.Summary)
}

func TestBuilderAddEndpoint_Deprecated(t *testing.T) {
	builder := NewBuilder("Test API", "1.0.0", "")

	endpoint := &parser.Endpoint{
		Path:       "/users/old",
		Method:     "GET",
		Summary:    "Get users (deprecated)",
		Deprecated: true,
	}

	err := builder.AddEndpoint(endpoint)
	require.NoError(t, err)

	pathItem, exists := builder.doc.Paths["/users/old"]
	assert.True(t, exists)
	assert.NotNil(t, pathItem.Get)
	assert.True(t, pathItem.Get.Deprecated)
}

func TestBuilderAddEndpoint_UpdateExistingPath(t *testing.T) {
	builder := NewBuilder("Test API", "1.0.0", "")

	// Add GET endpoint
	endpoint1 := &parser.Endpoint{
		Path:    "/users",
		Method:  "GET",
		Summary: "Get users",
	}

	err := builder.AddEndpoint(endpoint1)
	require.NoError(t, err)

	// Add POST endpoint to same path
	endpoint2 := &parser.Endpoint{
		Path:    "/users",
		Method:  "POST",
		Summary: "Create user",
	}

	err = builder.AddEndpoint(endpoint2)
	require.NoError(t, err)

	pathItem, exists := builder.doc.Paths["/users"]
	assert.True(t, exists)
	assert.NotNil(t, pathItem.Get)
	assert.NotNil(t, pathItem.Post)
	assert.Equal(t, "Get users", pathItem.Get.Summary)
	assert.Equal(t, "Create user", pathItem.Post.Summary)
}

func TestBuilderAddEndpoint_EmptyParameters(t *testing.T) {
	builder := NewBuilder("Test API", "1.0.0", "")

	endpoint := &parser.Endpoint{
		Path:       "/users",
		Method:     "GET",
		Summary:    "Get users",
		Parameters: []parser.Parameter{},
	}

	err := builder.AddEndpoint(endpoint)
	require.NoError(t, err)

	pathItem, exists := builder.doc.Paths["/users"]
	assert.True(t, exists)
	assert.Empty(t, pathItem.Get.Parameters)
}

func TestBuilderAddEndpoint_EmptyResponses(t *testing.T) {
	builder := NewBuilder("Test API", "1.0.0", "")

	endpoint := &parser.Endpoint{
		Path:      "/users",
		Method:    "GET",
		Summary:   "Get users",
		Responses: map[string]parser.Response{},
	}

	err := builder.AddEndpoint(endpoint)
	require.NoError(t, err)

	pathItem, exists := builder.doc.Paths["/users"]
	assert.True(t, exists)
	assert.Empty(t, pathItem.Get.Responses)
}

func TestBuilderAddEndpoint_EmptyTags(t *testing.T) {
	builder := NewBuilder("Test API", "1.0.0", "")

	endpoint := &parser.Endpoint{
		Path:    "/users",
		Method:  "GET",
		Summary: "Get users",
		Tags:    []string{},
	}

	err := builder.AddEndpoint(endpoint)
	require.NoError(t, err)

	assert.Empty(t, builder.doc.Tags)
}

func TestBuilderAddMultipleSchemas(t *testing.T) {
	builder := NewBuilder("Test API", "1.0.0", "")

	schema1 := &Schema{Type: "object", Properties: map[string]*Schema{"id": {Type: "integer"}}}
	schema2 := &Schema{Type: "object", Properties: map[string]*Schema{"name": {Type: "string"}}}
	schema3 := &Schema{Type: "object", Properties: map[string]*Schema{"email": {Type: "string"}}}

	err := builder.AddSchema("User", schema1)
	require.NoError(t, err)

	err = builder.AddSchema("Post", schema2)
	require.NoError(t, err)

	err = builder.AddSchema("Comment", schema3)
	require.NoError(t, err)

	assert.Len(t, builder.doc.Components.Schemas, 3)
	assert.NotNil(t, builder.doc.Components.Schemas["User"])
	assert.NotNil(t, builder.doc.Components.Schemas["Post"])
	assert.NotNil(t, builder.doc.Components.Schemas["Comment"])
}

func TestBuilderAddMultipleServers(t *testing.T) {
	builder := NewBuilder("Test API", "1.0.0", "")

	server1 := Server{URL: "https://api.example.com", Description: "Production"}
	server2 := Server{URL: "https://staging.example.com", Description: "Staging"}
	server3 := Server{URL: "http://localhost:8080", Description: "Local"}

	builder.AddServer(server1)
	builder.AddServer(server2)
	builder.AddServer(server3)

	assert.Len(t, builder.doc.Servers, 3)
	assert.Equal(t, "https://api.example.com", builder.doc.Servers[0].URL)
	assert.Equal(t, "https://staging.example.com", builder.doc.Servers[1].URL)
	assert.Equal(t, "http://localhost:8080", builder.doc.Servers[2].URL)
}

func TestBuilderToJSON_ValidJSON(t *testing.T) {
	builder := NewBuilder("Test API", "1.0.0", "Test Description")

	endpoint := &parser.Endpoint{
		Path:    "/users",
		Method:  "GET",
		Summary: "Get users",
		Parameters: []parser.Parameter{
			{Name: "id", In: "query", Type: "string", Required: true},
		},
		Responses: map[string]parser.Response{
			"200": {StatusCode: "200", Description: "Success"},
		},
	}

	err := builder.AddEndpoint(endpoint)
	require.NoError(t, err)

	jsonData, err := builder.ToJSON()
	require.NoError(t, err)

	// Verify it's valid JSON
	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	require.NoError(t, err)

	assert.Equal(t, "3.0.0", result["openapi"])
	assert.NotNil(t, result["info"])
	assert.NotNil(t, result["paths"])
}

func TestBuilderToYAML_ValidYAML(t *testing.T) {
	builder := NewBuilder("Test API", "1.0.0", "Test Description")

	endpoint := &parser.Endpoint{
		Path:    "/users",
		Method:  "GET",
		Summary: "Get users",
	}

	err := builder.AddEndpoint(endpoint)
	require.NoError(t, err)

	yamlData, err := builder.ToYAML()
	require.NoError(t, err)

	yamlStr := string(yamlData)
	assert.Contains(t, yamlStr, "openapi: 3.0.0")
	assert.Contains(t, yamlStr, "title: Test API")
	assert.Contains(t, yamlStr, "version: 1.0.0")
	assert.Contains(t, yamlStr, "description: Test Description")
	assert.Contains(t, yamlStr, "/users:")
}

func TestBuilderComplexScenario(t *testing.T) {
	builder := NewBuilder("User Management API", "2.0.0", "API for managing users")

	// Add multiple endpoints
	endpoints := []parser.Endpoint{
		{
			Path:    "/users",
			Method:  "GET",
			Summary: "List all users",
			Tags:    []string{"Users"},
			Parameters: []parser.Parameter{
				{Name: "page", In: "query", Type: "integer", Required: false},
				{Name: "limit", In: "query", Type: "integer", Required: false},
			},
			Responses: map[string]parser.Response{
				"200": {StatusCode: "200", Description: "Success"},
				"400": {StatusCode: "400", Description: "Bad Request"},
			},
		},
		{
			Path:    "/users",
			Method:  "POST",
			Summary: "Create a new user",
			Tags:    []string{"Users"},
			Responses: map[string]parser.Response{
				"201": {StatusCode: "201", Description: "Created"},
				"400": {StatusCode: "400", Description: "Bad Request"},
			},
		},
		{
			Path:    "/users/{id}",
			Method:  "GET",
			Summary: "Get user by ID",
			Tags:    []string{"Users"},
			Parameters: []parser.Parameter{
				{Name: "id", In: "path", Type: "string", Required: true},
			},
			Responses: map[string]parser.Response{
				"200": {StatusCode: "200", Description: "Success"},
				"404": {StatusCode: "404", Description: "Not Found"},
			},
		},
		{
			Path:    "/users/{id}",
			Method:  "PUT",
			Summary: "Update user",
			Tags:    []string{"Users"},
			Parameters: []parser.Parameter{
				{Name: "id", In: "path", Type: "string", Required: true},
			},
			Responses: map[string]parser.Response{
				"200": {StatusCode: "200", Description: "Success"},
				"404": {StatusCode: "404", Description: "Not Found"},
			},
		},
		{
			Path:    "/users/{id}",
			Method:  "DELETE",
			Summary: "Delete user",
			Tags:    []string{"Users"},
			Parameters: []parser.Parameter{
				{Name: "id", In: "path", Type: "string", Required: true},
			},
			Responses: map[string]parser.Response{
				"204": {StatusCode: "204", Description: "No Content"},
				"404": {StatusCode: "404", Description: "Not Found"},
			},
		},
	}

	for i := range endpoints {
		err := builder.AddEndpoint(&endpoints[i])
		require.NoError(t, err)
	}

	// Add schemas
	userSchema := &Schema{
		Type: "object",
		Properties: map[string]*Schema{
			"id":    {Type: "integer"},
			"name":  {Type: "string"},
			"email": {Type: "string"},
		},
	}
	err := builder.AddSchema("User", userSchema)
	require.NoError(t, err)

	// Add servers
	builder.AddServer(Server{URL: "https://api.example.com", Description: "Production"})
	builder.AddServer(Server{URL: "http://localhost:8080", Description: "Local"})

	// Verify the document
	doc := builder.Build()
	assert.Equal(t, "3.0.0", doc.OpenAPI)
	assert.Equal(t, "User Management API", doc.Info.Title)
	assert.Equal(t, "2.0.0", doc.Info.Version)
	assert.Len(t, doc.Paths, 2)
	assert.Len(t, doc.Components.Schemas, 1)
	assert.Len(t, doc.Servers, 2)
	assert.Len(t, doc.Tags, 1)

	// Verify JSON output
	jsonData, err := builder.ToJSON()
	require.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	// Verify YAML output
	yamlData, err := builder.ToYAML()
	require.NoError(t, err)
	assert.NotEmpty(t, yamlData)
}
