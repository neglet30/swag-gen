# Phase 2: 核心功能 - 设计文档

## 架构设计

### 系统架构

```
┌─────────────────────────────────────────────────────────┐
│                    CLI 工具 (Cobra)                      │
│  ┌──────────────────────────────────────────────────┐   │
│  │ init 命令 │ server 命令 │ help                   │   │
│  └──────────────────────────────────────────────────┘   │
└────────────────┬────────────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────────────┐
│              代码解析模块 (Parser)                       │
│  ┌──────────────────────────────────────────────────┐   │
│  │ AST 解析 │ 注释解析 │ API 信息提取              │   │
│  └──────────────────────────────────────────────────┘   │
└────────────────┬────────────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────────────┐
│           Swagger 生成模块 (Swagger)                     │
│  ┌──────────────────────────────────────────────────┐   │
│  │ 文档构建 │ Schema 生成 │ 格式转换                │   │
│  └──────────────────────────────────────────────────┘   │
└────────────────┬────────────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────────────┐
│              文件输出模块 (Output)                       │
│  ┌──────────────────────────────────────────────────┐   │
│  │ 文件写入 │ 格式转换 │ 配置生成                  │   │
│  └──────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
```

### 模块设计

#### 1. 代码解析模块 (pkg/parser/)

**文件结构**:
```
pkg/parser/
├── parser.go           # 主解析器
├── ast_parser.go       # AST 解析
├── comment_parser.go   # 注释解析
├── models.go           # 数据模型
├── parser_test.go
├── ast_parser_test.go
├── comment_parser_test.go
└── testdata/           # 测试数据
```

**主要类型**:
```go
// 解析器
type Parser struct {
    config *config.Config
    logger *zap.Logger
}

// API 端点信息
type Endpoint struct {
    Method      string
    Path        string
    Summary     string
    Description string
    Tags        []string
    Parameters  []Parameter
    Responses   map[string]Response
    Deprecated  bool
}

// 参数信息
type Parameter struct {
    Name        string
    In          string // query, path, header, body
    Type        string
    Required    bool
    Description string
}

// 响应信息
type Response struct {
    StatusCode  string
    Description string
    Schema      *Schema
}

// Schema 定义
type Schema struct {
    Type        string
    Properties  map[string]*Schema
    Items       *Schema
    Required    []string
    Description string
}
```

**主要函数**:
```go
// 创建解析器
func NewParser(cfg *config.Config, logger *zap.Logger) *Parser

// 解析项目
func (p *Parser) ParseProject(projectPath string) ([]*Endpoint, error)

// 解析单个文件
func (p *Parser) ParseFile(filePath string) ([]*Endpoint, error)

// 解析注释
func (p *Parser) parseComments(file *ast.File) ([]*Endpoint, error)

// 提取 API 信息
func (p *Parser) extractAPIInfo(comments []*ast.CommentGroup) *Endpoint
```

#### 2. Swagger 生成模块 (pkg/swagger/)

**文件结构**:
```
pkg/swagger/
├── builder.go          # Swagger 构建器
├── models.go           # OpenAPI 模型
├── schema_builder.go   # Schema 构建器
├── builder_test.go
├── schema_builder_test.go
└── testdata/           # 测试数据
```

**主要类型**:
```go
// Swagger 文档
type SwaggerDoc struct {
    OpenAPI    string                 `json:"openapi"`
    Info       Info                   `json:"info"`
    Paths      map[string]PathItem    `json:"paths"`
    Components Components             `json:"components"`
}

// 信息
type Info struct {
    Title       string `json:"title"`
    Version     string `json:"version"`
    Description string `json:"description,omitempty"`
}

// 路径项
type PathItem struct {
    Get     *Operation `json:"get,omitempty"`
    Post    *Operation `json:"post,omitempty"`
    Put     *Operation `json:"put,omitempty"`
    Delete  *Operation `json:"delete,omitempty"`
    Patch   *Operation `json:"patch,omitempty"`
}

// 操作
type Operation struct {
    Summary     string                 `json:"summary,omitempty"`
    Description string                 `json:"description,omitempty"`
    Tags        []string               `json:"tags,omitempty"`
    Parameters  []Parameter            `json:"parameters,omitempty"`
    RequestBody *RequestBody           `json:"requestBody,omitempty"`
    Responses   map[string]Response    `json:"responses"`
    Deprecated  bool                   `json:"deprecated,omitempty"`
}

// 组件
type Components struct {
    Schemas map[string]Schema `json:"schemas,omitempty"`
}
```

**主要函数**:
```go
// 创建构建器
func NewBuilder(title, version, description string) *Builder

// 添加端点
func (b *Builder) AddEndpoint(endpoint *Endpoint) error

// 构建文档
func (b *Builder) Build() *SwaggerDoc

// 转换为 JSON
func (b *Builder) ToJSON() ([]byte, error)

// 转换为 YAML
func (b *Builder) ToYAML() ([]byte, error)
```

#### 3. 文件输出模块 (pkg/output/)

**文件结构**:
```
pkg/output/
├── writer.go           # 文件写入器
├── formatter.go        # 格式转换器
├── writer_test.go
└── formatter_test.go
```

**主要函数**:
```go
// 写入 Swagger 文档
func WriteSwagger(doc *swagger.SwaggerDoc, outputPath string, format string) error

// 写入配置文件
func WriteConfig(config *Config, outputPath string) error

// 写入 README
func WriteREADME(outputPath string, title string) error
```

### 数据流

#### 初始化流程
```
用户执行 swag-gen init -p ./api -o ./docs
    ↓
验证参数
    ↓
创建解析器
    ↓
扫描项目目录
    ↓
解析 Go 源文件 (并发)
    ↓
提取 API 信息
    ↓
创建 Swagger 构建器
    ↓
添加所有端点
    ↓
构建 Swagger 文档
    ↓
写入输出文件
    ↓
生成配置文件
    ↓
返回成功信息
```

#### 解析流程
```
读取源文件
    ↓
使用 go/parser 解析 AST
    ↓
遍历 AST 节点
    ↓
查找注释
    ↓
解析 Swagger 标签
    ↓
提取 API 信息
    ↓
返回端点列表
```

## 详细设计

### 1. 代码解析器设计

#### 1.1 AST 解析
- 使用 `go/parser` 解析 Go 源文件
- 支持递归扫描目录
- 使用 goroutine 并发解析
- 使用 sync.WaitGroup 管理并发

**实现步骤**:
1. 获取目录中的所有 Go 文件
2. 为每个文件创建 goroutine
3. 使用 go/parser 解析文件
4. 收集解析结果
5. 处理错误

#### 1.2 注释解析
- 遍历 AST 中的注释
- 识别 Swagger 标签
- 提取标签值
- 验证标签格式

**支持的标签格式**:
```go
// @Router /api/users [GET]
// @Summary Get all users
// @Description Get all users from database
// @Tags User
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} []User
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
```

#### 1.3 API 信息提取
- 从注释中提取路由信息
- 从注释中提取参数信息
- 从注释中提取响应信息
- 从函数签名中提取额外信息

### 2. Swagger 生成器设计

#### 2.1 文档构建
- 创建 OpenAPI 3.0 文档结构
- 添加基本信息
- 添加路径和操作
- 添加 Schema 定义

#### 2.2 Schema 生成
- 从 Go 类型生成 Schema
- 支持基本类型
- 支持结构体
- 支持数组和指针
- 支持嵌套结构

#### 2.3 格式转换
- 支持 JSON 格式
- 支持 YAML 格式
- 使用标准库进行转换

### 3. CLI init 命令设计

#### 3.1 参数处理
```go
type InitOptions struct {
    Path        string // 项目路径
    Output      string // 输出路径
    Title       string // API 标题
    Version     string // API 版本
    Description string // API 描述
    Format      string // 输出格式
}
```

#### 3.2 执行流程
1. 解析命令行参数
2. 验证参数有效性
3. 创建解析器
4. 解析项目
5. 生成 Swagger 文档
6. 写入输出文件
7. 返回结果

#### 3.3 错误处理
- 参数验证错误
- 文件读取错误
- 解析错误
- 写入错误
- 提供清晰的错误信息

### 4. 并发处理设计

#### 4.1 文件解析并发
- 使用 goroutine 并发解析文件
- 使用 sync.WaitGroup 等待所有 goroutine
- 使用 channel 收集结果
- 使用 sync.Mutex 保护共享资源

#### 4.2 错误处理
- 使用 error channel 收集错误
- 支持部分失败继续处理
- 记录所有错误
- 返回汇总错误

### 5. 性能优化

#### 5.1 缓存
- 缓存已解析的文件
- 缓存生成的 Schema
- 支持增量解析

#### 5.2 内存管理
- 及时释放资源
- 避免不必要的复制
- 使用对象池

#### 5.3 并发优化
- 合理设置 goroutine 数量
- 避免过度并发
- 使用 buffered channel

## 接口设计

### Parser 接口
```go
type Parser interface {
    ParseProject(projectPath string) ([]*Endpoint, error)
    ParseFile(filePath string) ([]*Endpoint, error)
}
```

### Builder 接口
```go
type Builder interface {
    AddEndpoint(endpoint *Endpoint) error
    Build() *SwaggerDoc
    ToJSON() ([]byte, error)
    ToYAML() ([]byte, error)
}
```

### Writer 接口
```go
type Writer interface {
    WriteSwagger(doc *SwaggerDoc, outputPath string, format string) error
    WriteConfig(config *Config, outputPath string) error
}
```

## 测试策略

### 单元测试
- 测试 AST 解析
- 测试注释解析
- 测试 API 信息提取
- 测试 Schema 生成
- 测试文件写入

### 集成测试
- 测试完整的初始化流程
- 测试大型项目处理
- 测试错误处理
- 测试性能

### 测试数据
- 创建示例 Go 文件
- 创建示例项目
- 创建边界情况测试

## 配置设计

### 项目配置文件 (swag-gen.yaml)
```yaml
project:
  name: "My API"
  version: "1.0.0"
  description: "API Documentation"

parser:
  path: "./api"
  exclude:
    - "vendor"
    - "test"

output:
  path: "./docs"
  format: "json"

swagger:
  title: "My API"
  version: "1.0.0"
  description: "API Documentation"
  basePath: "/api/v1"
```

## 错误处理

### 错误类型
```go
type ParseError struct {
    File    string
    Line    int
    Message string
}

type ValidationError struct {
    Field   string
    Message string
}

type GenerationError struct {
    Message string
}
```

### 错误处理策略
- 验证输入参数
- 处理文件读取错误
- 处理解析错误
- 处理写入错误
- 提供清晰的错误信息

## 日志设计

### 日志级别
- DEBUG: 详细的调试信息
- INFO: 一般信息
- WARN: 警告信息
- ERROR: 错误信息

### 日志内容
- 解析进度
- 生成进度
- 错误和警告
- 性能指标

## 相关文档

- [需求文档](./requirements.md)
- [代码规范](../../steering/code-standards.md)
- [API 设计](../../steering/api-design.md)
- [技术栈](../../steering/tech.md)
