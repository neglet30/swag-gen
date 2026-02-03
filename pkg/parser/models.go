package parser

// Endpoint 代表一个 API 端点
type Endpoint struct {
	Method      string
	Path        string
	Summary     string
	Description string
	Tags        []string
	Parameters  []Parameter
	Responses   map[string]Response
	Deprecated  bool
	File        string
	Line        int
}

// Parameter 代表一个参数
type Parameter struct {
	Name        string
	In          string // query, path, header, body
	Type        string
	Required    bool
	Description string
}

// Response 代表一个响应
type Response struct {
	StatusCode  string
	Description string
	Schema      *Schema
}

// Schema 代表一个数据模型
type Schema struct {
	Type        string
	Properties  map[string]*Schema
	Items       *Schema
	Required    []string
	Description string
}

// ParseResult 代表解析结果
type ParseResult struct {
	Endpoints []*Endpoint
	Errors    []error
}
