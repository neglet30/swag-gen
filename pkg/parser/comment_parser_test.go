package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewCommentParser(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	cp := NewCommentParser(logger)

	assert.NotNil(t, cp)
	assert.Equal(t, logger, cp.logger)
}

func TestCommentParserParseRouter(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cp := NewCommentParser(logger)

	tests := []struct {
		name    string
		text    string
		wantErr bool
		method  string
		path    string
	}{
		{
			name:    "valid router",
			text:    "// @Router /api/users [GET]",
			wantErr: false,
			method:  "GET",
			path:    "/api/users",
		},
		{
			name:    "valid router POST",
			text:    "// @Router /api/users [POST]",
			wantErr: false,
			method:  "POST",
			path:    "/api/users",
		},
		{
			name:    "invalid router",
			text:    "// @Router /api/users",
			wantErr: true,
		},
		{
			name:    "no router tag",
			text:    "// @Summary Get users",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := cp.parseRouter(tt.text)

			if tt.wantErr {
				assert.Nil(t, router)
			} else {
				assert.NotNil(t, router)
				assert.Equal(t, tt.method, router.Method)
				assert.Equal(t, tt.path, router.Path)
			}
		})
	}
}

func TestCommentParserParseSimpleTag(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cp := NewCommentParser(logger)

	tests := []struct {
		name     string
		text     string
		tag      string
		expected string
	}{
		{
			name:     "summary tag",
			text:     "// @Summary Get all users",
			tag:      "@Summary",
			expected: "Get all users",
		},
		{
			name:     "description tag",
			text:     "// @Description Get all users from database",
			tag:      "@Description",
			expected: "Get all users from database",
		},
		{
			name:     "tags tag",
			text:     "// @Tags User",
			tag:      "@Tags",
			expected: "User",
		},
		{
			name:     "no tag",
			text:     "// Some comment",
			tag:      "@Summary",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cp.parseSimpleTag(tt.text, tt.tag)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCommentParserParseParam(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cp := NewCommentParser(logger)

	tests := []struct {
		name    string
		text    string
		wantErr bool
		param   *Parameter
	}{
		{
			name:    "valid param",
			text:    `// @Param page query int false "Page number"`,
			wantErr: false,
			param: &Parameter{
				Name:        "page",
				In:          "query",
				Type:        "int",
				Required:    false,
				Description: "Page number",
			},
		},
		{
			name:    "valid param required",
			text:    `// @Param id path int true "User ID"`,
			wantErr: false,
			param: &Parameter{
				Name:        "id",
				In:          "path",
				Type:        "int",
				Required:    true,
				Description: "User ID",
			},
		},
		{
			name:    "invalid param",
			text:    `// @Param page query`,
			wantErr: true,
		},
		{
			name:    "no param tag",
			text:    "// @Summary Get users",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			param := cp.parseParam(tt.text)

			if tt.wantErr {
				assert.Nil(t, param)
			} else {
				assert.NotNil(t, param)
				assert.Equal(t, tt.param.Name, param.Name)
				assert.Equal(t, tt.param.In, param.In)
				assert.Equal(t, tt.param.Type, param.Type)
				assert.Equal(t, tt.param.Required, param.Required)
				assert.Equal(t, tt.param.Description, param.Description)
			}
		})
	}
}

func TestCommentParserParseResponse(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cp := NewCommentParser(logger)

	tests := []struct {
		name       string
		text       string
		tag        string
		wantErr    bool
		statusCode string
	}{
		{
			name:       "valid success response",
			text:       "// @Success 200 {object} User",
			tag:        "@Success",
			wantErr:    false,
			statusCode: "200",
		},
		{
			name:       "valid failure response",
			text:       "// @Failure 400 {object} ErrorResponse",
			tag:        "@Failure",
			wantErr:    false,
			statusCode: "400",
		},
		{
			name:    "invalid response",
			text:    "// @Success 200",
			tag:     "@Success",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := cp.parseResponse(tt.text, tt.tag)

			if tt.wantErr {
				assert.Nil(t, resp)
			} else {
				assert.NotNil(t, resp)
				assert.Equal(t, tt.statusCode, resp.StatusCode)
			}
		})
	}
}

func TestCommentParserParseEndpoint(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cp := NewCommentParser(logger)

	comments := []string{
		"// @Router /api/users [GET]",
		"// @Summary Get all users",
		"// @Description Get all users from database",
		"// @Tags User",
		"// @Deprecated",
	}

	endpoint := cp.ParseEndpoint(comments, "test.go", 10)

	assert.NotNil(t, endpoint)
	assert.Equal(t, "GET", endpoint.Method)
	assert.Equal(t, "/api/users", endpoint.Path)
	assert.Equal(t, "Get all users", endpoint.Summary)
	assert.Equal(t, "Get all users from database", endpoint.Description)
	assert.True(t, endpoint.Deprecated)
	assert.Equal(t, "test.go", endpoint.File)
	assert.Equal(t, 10, endpoint.Line)
}

func TestCommentParserParseEndpointNoRouter(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cp := NewCommentParser(logger)

	comments := []string{
		"// @Summary Get all users",
		"// @Description Get all users from database",
	}

	endpoint := cp.ParseEndpoint(comments, "test.go", 10)

	assert.Nil(t, endpoint)
}

func TestCommentParserExtractAllTags(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cp := NewCommentParser(logger)

	comments := []string{
		"// @Router /api/users [GET]",
		"// @Summary Get all users",
		"// @Tags User",
		"// @Tags Admin",
	}

	tags := cp.ExtractAllTags(comments)

	assert.Len(t, tags, 3)
	assert.Contains(t, tags, "@Router")
	assert.Contains(t, tags, "@Summary")
	assert.Contains(t, tags, "@Tags")
	assert.Len(t, tags["@Tags"], 2)
}

func TestCommentParserSupportedTags(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cp := NewCommentParser(logger)

	tags := cp.SupportedTags()

	assert.Len(t, tags, 8)
	assert.Contains(t, tags, "@Router")
	assert.Contains(t, tags, "@Summary")
	assert.Contains(t, tags, "@Description")
	assert.Contains(t, tags, "@Tags")
	assert.Contains(t, tags, "@Param")
	assert.Contains(t, tags, "@Success")
	assert.Contains(t, tags, "@Failure")
	assert.Contains(t, tags, "@Deprecated")
}

func TestCommentParserValidateTag(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cp := NewCommentParser(logger)

	tests := []struct {
		name    string
		text    string
		wantErr bool
	}{
		{
			name:    "valid tag",
			text:    "@Router /api/users [GET]",
			wantErr: false,
		},
		{
			name:    "invalid tag",
			text:    "Some comment",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cp.ValidateTag(tt.text)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
