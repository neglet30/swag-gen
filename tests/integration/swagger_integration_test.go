package integration

import (
	"encoding/json"
	"testing"

	"github.com/neglet30/swag-gen/pkg/parser"
	"github.com/neglet30/swag-gen/pkg/swagger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSwaggerIntegration_BuildSimpleDocument 测试构建简单的 Swagger 文档
func TestSwaggerIntegration_BuildSimpleDocument(t *testing.T) {
	// 创建构建器
	builder := swagger.NewBuilder("Test API", "1.0.0", "Test API Description")

	// 创建端点
	endpoint := &parser.Endpoint{
		Method:      "GET",
		Path:        "/api/users",
		Summary:     "Get users",
		Description: "Get all users",
		Tags:        []string{"User"},
		Parameters:  []parser.Parameter{},
		Responses: map[string]parser.Response{
			"200": {
				StatusCode:  "200",
				Description: "Success",
			},
		},
		Deprecated: false,
	}

	// 添加端点
	err := builder.AddEndpoint(endpoint)
	require.NoError(t, err)

	// 构建文档
	doc := builder.Build()

	// 验证文档
	assert.NotNil(t, doc)
	assert.Equal(t, "3.0.0", doc.OpenAPI)
	assert.Equal(t, "Test API", doc.Info.Title)
	assert.Equal(t, "1.0.0", doc.Info.Version)
	assert.Equal(t, "Test API Description", doc.Info.Description)

	// 验证路径
	assert.Contains(t, doc.Paths, "/api/users")
	pathItem := doc.Paths["/api/users"]
	assert.NotNil(t, pathItem.Get)
	assert.Equal(t, "Get users", pathItem.Get.Summary)
	assert.Equal(t, "Get all users", pathItem.Get.Description)
}

// TestSwaggerIntegration_BuildMultipleEndpoints 测试构建多个端点的文档
func TestSwaggerIntegration_BuildMultipleEndpoints(t *testing.T) {
	// 创建构建器
	builder := swagger.NewBuilder("Test API", "1.0.0", "")

	// 创建多个端点
	endpoints := []*parser.Endpoint{
		{
			Method:     "GET",
			Path:       "/api/users",
			Summary:    "Get users",
			Tags:       []string{"User"},
			Parameters: []parser.Parameter{},
			Responses:  map[string]parser.Response{"200": {StatusCode: "200", Description: "Success"}},
			Deprecated: false,
		},
		{
			Method:     "POST",
			Path:       "/api/users",
			Summary:    "Create user",
			Tags:       []string{"User"},
			Parameters: []parser.Parameter{},
			Responses:  map[string]parser.Response{"201": {StatusCode: "201", Description: "Created"}},
			Deprecated: false,
		},
		{
			Method:     "GET",
			Path:       "/api/posts",
			Summary:    "Get posts",
			Tags:       []string{"Post"},
			Parameters: []parser.Parameter{},
			Responses:  map[string]parser.Response{"200": {StatusCode: "200", Description: "Success"}},
			Deprecated: false,
		},
	}

	// 添加所有端点
	for _, ep := range endpoints {
		err := builder.AddEndpoint(ep)
		require.NoError(t, err)
	}

	// 构建文档
	doc := builder.Build()

	// 验证文档
	assert.NotNil(t, doc)
	assert.Equal(t, 2, len(doc.Paths), "应该有 2 个路径")

	// 验证路径
	assert.Contains(t, doc.Paths, "/api/users")
	assert.Contains(t, doc.Paths, "/api/posts")

	// 验证 /api/users 路径
	usersPath := doc.Paths["/api/users"]
	assert.NotNil(t, usersPath.Get)
	assert.NotNil(t, usersPath.Post)

	// 验证 /api/posts 路径
	postsPath := doc.Paths["/api/posts"]
	assert.NotNil(t, postsPath.Get)
}

// TestSwaggerIntegration_ToJSON 测试转换为 JSON
func TestSwaggerIntegration_ToJSON(t *testing.T) {
	// 创建构建器
	builder := swagger.NewBuilder("Test API", "1.0.0", "Test Description")

	// 创建端点
	endpoint := &parser.Endpoint{
		Method:     "GET",
		Path:       "/api/users",
		Summary:    "Get users",
		Tags:       []string{"User"},
		Parameters: []parser.Parameter{},
		Responses:  map[string]parser.Response{"200": {StatusCode: "200", Description: "Success"}},
		Deprecated: false,
	}

	// 添加端点
	err := builder.AddEndpoint(endpoint)
	require.NoError(t, err)

	// 转换为 JSON
	jsonData, err := builder.ToJSON()
	require.NoError(t, err)

	// 验证 JSON
	assert.NotEmpty(t, jsonData)

	// 解析 JSON
	var doc map[string]interface{}
	err = json.Unmarshal(jsonData, &doc)
	require.NoError(t, err)

	// 验证内容
	assert.Equal(t, "3.0.0", doc["openapi"])
	assert.NotNil(t, doc["info"])
	assert.NotNil(t, doc["paths"])
}

// TestSwaggerIntegration_ToYAML 测试转换为 YAML
func TestSwaggerIntegration_ToYAML(t *testing.T) {
	// 创建构建器
	builder := swagger.NewBuilder("Test API", "1.0.0", "Test Description")

	// 创建端点
	endpoint := &parser.Endpoint{
		Method:     "GET",
		Path:       "/api/users",
		Summary:    "Get users",
		Tags:       []string{"User"},
		Parameters: []parser.Parameter{},
		Responses:  map[string]parser.Response{"200": {StatusCode: "200", Description: "Success"}},
		Deprecated: false,
	}

	// 添加端点
	err := builder.AddEndpoint(endpoint)
	require.NoError(t, err)

	// 转换为 YAML
	yamlData, err := builder.ToYAML()
	require.NoError(t, err)

	// 验证 YAML
	assert.NotEmpty(t, yamlData)
	assert.Contains(t, string(yamlData), "openapi: 3.0.0")
	assert.Contains(t, string(yamlData), "title: Test API")
}

// TestSwaggerIntegration_WithParameters 测试包含参数的端点
func TestSwaggerIntegration_WithParameters(t *testing.T) {
	// 创建构建器
	builder := swagger.NewBuilder("Test API", "1.0.0", "")

	// 创建带参数的端点
	endpoint := &parser.Endpoint{
		Method:  "GET",
		Path:    "/api/users/{id}",
		Summary: "Get user by ID",
		Tags:    []string{"User"},
		Parameters: []parser.Parameter{
			{
				Name:        "id",
				In:          "path",
				Type:        "integer",
				Required:    true,
				Description: "User ID",
			},
		},
		Responses:  map[string]parser.Response{"200": {StatusCode: "200", Description: "Success"}},
		Deprecated: false,
	}

	// 添加端点
	err := builder.AddEndpoint(endpoint)
	require.NoError(t, err)

	// 构建文档
	doc := builder.Build()

	// 验证文档
	assert.NotNil(t, doc)
	pathItem := doc.Paths["/api/users/{id}"]
	assert.NotNil(t, pathItem.Get)
	assert.Equal(t, 1, len(pathItem.Get.Parameters))

	// 验证参数
	param := pathItem.Get.Parameters[0]
	assert.Equal(t, "id", param.Name)
	assert.Equal(t, "path", param.In)
	assert.True(t, param.Required)
}

// TestSwaggerIntegration_WithMultipleResponses 测试多个响应
func TestSwaggerIntegration_WithMultipleResponses(t *testing.T) {
	// 创建构建器
	builder := swagger.NewBuilder("Test API", "1.0.0", "")

	// 创建带多个响应的端点
	endpoint := &parser.Endpoint{
		Method:     "GET",
		Path:       "/api/users/{id}",
		Summary:    "Get user by ID",
		Tags:       []string{"User"},
		Parameters: []parser.Parameter{},
		Responses: map[string]parser.Response{
			"200": {StatusCode: "200", Description: "Success"},
			"404": {StatusCode: "404", Description: "Not Found"},
			"500": {StatusCode: "500", Description: "Internal Server Error"},
		},
		Deprecated: false,
	}

	// 添加端点
	err := builder.AddEndpoint(endpoint)
	require.NoError(t, err)

	// 构建文档
	doc := builder.Build()

	// 验证文档
	assert.NotNil(t, doc)
	pathItem := doc.Paths["/api/users/{id}"]
	assert.NotNil(t, pathItem.Get)
	assert.Equal(t, 3, len(pathItem.Get.Responses))

	// 验证响应
	assert.Contains(t, pathItem.Get.Responses, "200")
	assert.Contains(t, pathItem.Get.Responses, "404")
	assert.Contains(t, pathItem.Get.Responses, "500")
}

// TestSwaggerIntegration_WithTags 测试标签
func TestSwaggerIntegration_WithTags(t *testing.T) {
	// 创建构建器
	builder := swagger.NewBuilder("Test API", "1.0.0", "")

	// 创建多个端点，使用不同的标签
	endpoints := []*parser.Endpoint{
		{
			Method:     "GET",
			Path:       "/api/users",
			Summary:    "Get users",
			Tags:       []string{"User"},
			Parameters: []parser.Parameter{},
			Responses:  map[string]parser.Response{"200": {StatusCode: "200", Description: "Success"}},
			Deprecated: false,
		},
		{
			Method:     "GET",
			Path:       "/api/posts",
			Summary:    "Get posts",
			Tags:       []string{"Post"},
			Parameters: []parser.Parameter{},
			Responses:  map[string]parser.Response{"200": {StatusCode: "200", Description: "Success"}},
			Deprecated: false,
		},
		{
			Method:     "GET",
			Path:       "/api/comments",
			Summary:    "Get comments",
			Tags:       []string{"Comment"},
			Parameters: []parser.Parameter{},
			Responses:  map[string]parser.Response{"200": {StatusCode: "200", Description: "Success"}},
			Deprecated: false,
		},
	}

	// 添加所有端点
	for _, ep := range endpoints {
		err := builder.AddEndpoint(ep)
		require.NoError(t, err)
	}

	// 构建文档
	doc := builder.Build()

	// 验证文档
	assert.NotNil(t, doc)
	assert.Greater(t, len(doc.Tags), 0, "应该有标签")

	// 验证标签
	tagNames := make(map[string]bool)
	for _, tag := range doc.Tags {
		tagNames[tag.Name] = true
	}

	assert.True(t, tagNames["User"], "应该有 User 标签")
	assert.True(t, tagNames["Post"], "应该有 Post 标签")
	assert.True(t, tagNames["Comment"], "应该有 Comment 标签")
}

// TestSwaggerIntegration_DeprecatedEndpoint 测试已弃用的端点
func TestSwaggerIntegration_DeprecatedEndpoint(t *testing.T) {
	// 创建构建器
	builder := swagger.NewBuilder("Test API", "1.0.0", "")

	// 创建已弃用的端点
	endpoint := &parser.Endpoint{
		Method:     "GET",
		Path:       "/api/old-users",
		Summary:    "Get users (deprecated)",
		Tags:       []string{"User"},
		Parameters: []parser.Parameter{},
		Responses:  map[string]parser.Response{"200": {StatusCode: "200", Description: "Success"}},
		Deprecated: true,
	}

	// 添加端点
	err := builder.AddEndpoint(endpoint)
	require.NoError(t, err)

	// 构建文档
	doc := builder.Build()

	// 验证文档
	assert.NotNil(t, doc)
	pathItem := doc.Paths["/api/old-users"]
	assert.NotNil(t, pathItem.Get)
	assert.True(t, pathItem.Get.Deprecated, "端点应该被标记为已弃用")
}
