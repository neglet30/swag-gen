# swag-gen 代码规范

## Go 代码规范

### 1. 命名规范

#### 包名
- 使用小写单词
- 避免使用下划线
- 简洁明了，通常为单个单词
- 示例: `parser`, `swagger`, `server`, `config`

#### 文件名
- 使用小写字母
- 多单词使用下划线分隔
- 示例: `ast_parser.go`, `comment_parser.go`, `swagger_builder.go`

#### 函数名
- 导出函数使用 PascalCase
- 未导出函数使用 camelCase
- 示例:
  ```go
  func ParseFile(path string) error { }      // 导出
  func parseComments(ast *ast.File) { }      // 未导出
  ```

#### 变量名
- 使用 camelCase
- 避免单字母变量（除了循环计数器）
- 示例:
  ```go
  var projectName string
  var apiEndpoints []string
  for i := 0; i < len(items); i++ { }
  ```

#### 常量名
- 使用 UPPER_SNAKE_CASE
- 示例:
  ```go
  const (
      DEFAULT_TIMEOUT = 30 * time.Second
      MAX_RETRY_COUNT = 3
  )
  ```

#### 接口名
- 使用 PascalCase
- 通常以 "er" 结尾
- 示例:
  ```go
  type Parser interface { }
  type Reader interface { }
  type Writer interface { }
  ```

### 2. 代码组织

#### 文件结构
```go
package parser

import (
    "fmt"
    "go/ast"
    
    "github.com/your-org/swag-gen/pkg/config"
)

// 常量定义
const (
    DefaultTimeout = 30 * time.Second
)

// 类型定义
type Parser struct {
    config *config.Config
}

// 接口定义
type Reader interface {
    Read() ([]byte, error)
}

// 函数实现
func NewParser(cfg *config.Config) *Parser {
    return &Parser{config: cfg}
}

func (p *Parser) Parse(path string) error {
    // 实现
}
```

#### 导入组织
```go
import (
    // 标准库
    "fmt"
    "os"
    "path/filepath"
    
    // 第三方库
    "github.com/gin-gonic/gin"
    "github.com/spf13/cobra"
    
    // 本项目
    "github.com/your-org/swag-gen/pkg/config"
    "github.com/your-org/swag-gen/internal/models"
)
```

### 3. 注释规范

#### 包注释
```go
// Package parser provides functionality for parsing Go source files
// and extracting API information from comments.
package parser
```

#### 函数注释
```go
// ParseFile parses a Go source file and returns the AST.
// It returns an error if the file cannot be parsed.
func ParseFile(path string) (*ast.File, error) {
    // 实现
}
```

#### 变量/常量注释
```go
// DefaultTimeout is the default timeout for API requests.
const DefaultTimeout = 30 * time.Second

// projectCache stores parsed projects in memory.
var projectCache = make(map[string]*Project)
```

#### 内联注释
```go
// 解释复杂逻辑
if err != nil {
    // 记录错误并继续处理其他文件
    logger.Error("failed to parse file", err)
    continue
}
```

### 4. 错误处理

#### 错误检查
```go
// 不好
file, _ := os.Open("file.txt")

// 好
file, err := os.Open("file.txt")
if err != nil {
    return fmt.Errorf("failed to open file: %w", err)
}
defer file.Close()
```

#### 自定义错误
```go
type ParseError struct {
    File    string
    Line    int
    Message string
}

func (e *ParseError) Error() string {
    return fmt.Sprintf("%s:%d: %s", e.File, e.Line, e.Message)
}
```

### 5. 并发编程

#### Goroutine 使用
```go
// 使用 WaitGroup 管理 goroutine
var wg sync.WaitGroup
for _, file := range files {
    wg.Add(1)
    go func(f string) {
        defer wg.Done()
        parseFile(f)
    }(file)
}
wg.Wait()
```

#### 通道使用
```go
// 使用通道传递数据
results := make(chan *ParseResult, len(files))
for _, file := range files {
    go func(f string) {
        result, err := parseFile(f)
        results <- result
    }(file)
}
```

### 6. 测试规范

#### 测试文件命名
- 使用 `_test.go` 后缀
- 示例: `parser_test.go`, `swagger_builder_test.go`

#### 测试函数命名
```go
func TestParseFile(t *testing.T) { }
func TestParseFile_InvalidPath(t *testing.T) { }
func TestParseFile_EmptyFile(t *testing.T) { }
```

#### 测试结构
```go
func TestParseFile(t *testing.T) {
    // 准备测试数据
    testFile := "testdata/sample.go"
    
    // 执行测试
    result, err := ParseFile(testFile)
    
    // 验证结果
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    
    if result == nil {
        t.Error("expected non-nil result")
    }
}
```

#### 表驱动测试
```go
func TestParseFile(t *testing.T) {
    tests := []struct {
        name    string
        path    string
        wantErr bool
    }{
        {"valid file", "testdata/valid.go", false},
        {"invalid file", "testdata/invalid.go", true},
        {"empty file", "testdata/empty.go", false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := ParseFile(tt.path)
            if (err != nil) != tt.wantErr {
                t.Errorf("ParseFile() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### 7. 代码格式化

#### 使用 go fmt
```bash
go fmt ./...
```

#### 行长度
- 建议不超过 100 个字符
- 必须不超过 120 个字符

#### 缩进
- 使用 Tab 缩进（Go 标准）

### 8. 最佳实践

#### 接口设计
```go
// 好：接口小而专
type Reader interface {
    Read() ([]byte, error)
}

// 不好：接口过大
type FileHandler interface {
    Read() ([]byte, error)
    Write([]byte) error
    Delete() error
    Rename(string) error
}
```

#### 依赖注入
```go
// 好：通过构造函数注入依赖
func NewParser(logger Logger, config *Config) *Parser {
    return &Parser{
        logger: logger,
        config: config,
    }
}

// 不好：全局变量
var globalLogger Logger
```

#### 避免 nil 指针
```go
// 好：检查 nil
if user != nil {
    fmt.Println(user.Name)
}

// 不好：直接访问
fmt.Println(user.Name)  // 可能 panic
```

---

## React/JavaScript 代码规范

### 1. 命名规范

#### 文件名
- 使用 kebab-case
- 示例: `swagger-viewer.jsx`, `api-tester.jsx`, `code-editor.jsx`

#### 变量名
- 使用 camelCase
- 示例: `projectName`, `apiEndpoints`, `isLoading`

#### 常量名
- 使用 UPPER_SNAKE_CASE
- 示例: `DEFAULT_TIMEOUT`, `MAX_RETRY_COUNT`

#### 函数名
- 使用 camelCase
- 示例: `parseProject()`, `generateSwagger()`, `testAPI()`

#### 组件名
- 使用 PascalCase
- 示例: `SwaggerViewer`, `APITester`, `CodeEditor`

### 2. 代码组织

#### 文件结构
```javascript
// 导入
import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useDispatch, useSelector } from 'react-redux';

// 常量
const DEFAULT_TIMEOUT = 30000;
const API_BASE_URL = 'http://localhost:8080';

// 类型定义（TypeScript）
interface Project {
    id: string;
    name: string;
    description: string;
}

// 组件
function SwaggerViewer() {
    const [swagger, setSwagger] = useState(null);
    const [loading, setLoading] = useState(false);
    
    // 逻辑
    useEffect(() => {
        loadSwagger();
    }, []);
    
    const loadSwagger = async () => {
        setLoading(true);
        try {
            const response = await axios.get(`${API_BASE_URL}/swagger`);
            setSwagger(response.data.data);
        } catch (error) {
            console.error('Failed to load swagger:', error);
        } finally {
            setLoading(false);
        }
    };
    
    // 渲染
    return (
        <div>
            {loading ? <Spinner /> : <SwaggerUI spec={swagger} />}
        </div>
    );
}

// 导出
export default SwaggerViewer;
```

#### 导入组织
```javascript
// React 和第三方库
import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useDispatch, useSelector } from 'react-redux';

// Material-UI 组件
import { Box, Button, Card, Container } from '@mui/material';

// 本项目组件
import SwaggerUI from './components/swagger-ui';
import APITester from './components/api-tester';

// 本项目工具和常量
import { API_BASE_URL } from './config';
import { fetchSwagger } from './services/api';
```

### 3. 注释规范

#### JSDoc 注释
```javascript
/**
 * Parses a project and generates Swagger documentation.
 * @param {string} projectId - The project ID
 * @param {Object} options - Configuration options
 * @param {string} options.path - Project path
 * @returns {Promise<Object>} The generated Swagger document
 * @throws {Error} If parsing fails
 */
function parseProject(projectId, options) {
    // 实现
}
```

#### 内联注释
```javascript
// 解释复杂逻辑
if (response.status === 200) {
    // 解析响应数据并更新状态
    setSwagger(response.data);
}
```

### 4. 错误处理

#### 异步错误处理
```javascript
// 好
async function loadSwagger() {
    try {
        const response = await axios.get('/swagger');
        setSwagger(response.data);
    } catch (error) {
        console.error('Failed to load swagger:', error);
        setError(error.message);
    }
}

// 不好
async function loadSwagger() {
    const response = await axios.get('/swagger');
    setSwagger(response.data);
}
```

### 5. React 最佳实践

#### 函数组件
```javascript
function SwaggerViewer() {
    const [swagger, setSwagger] = useState(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    
    useEffect(() => {
        loadSwagger();
    }, []);
    
    const loadSwagger = async () => {
        setLoading(true);
        try {
            const response = await axios.get('/swagger');
            setSwagger(response.data.data);
        } catch (err) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    };
    
    if (loading) return <Spinner />;
    if (error) return <ErrorMessage message={error} />;
    
    return <SwaggerUI spec={swagger} />;
}
```

#### Props 验证
```javascript
import PropTypes from 'prop-types';

function SwaggerCard({ swagger, onTest }) {
    return (
        <Card>
            <h3>{swagger.info.title}</h3>
            <button onClick={() => onTest(swagger)}>Test</button>
        </Card>
    );
}

SwaggerCard.propTypes = {
    swagger: PropTypes.shape({
        info: PropTypes.shape({
            title: PropTypes.string.required,
            version: PropTypes.string.required,
        }).required,
    }).required,
    onTest: PropTypes.func.required,
};
```

#### 自定义 Hooks
```javascript
// 好：提取可复用逻辑
function useSwagger() {
    const [swagger, setSwagger] = useState(null);
    const [loading, setLoading] = useState(false);
    
    const load = async () => {
        setLoading(true);
        try {
            const response = await axios.get('/swagger');
            setSwagger(response.data.data);
        } finally {
            setLoading(false);
        }
    };
    
    return { swagger, loading, load };
}

// 使用
function SwaggerViewer() {
    const { swagger, loading, load } = useSwagger();
    
    useEffect(() => {
        load();
    }, []);
    
    return <div>{loading ? <Spinner /> : <SwaggerUI spec={swagger} />}</div>;
}
```

### 6. 测试规范

#### 测试文件命名
- 使用 `.test.jsx` 或 `.spec.jsx` 后缀
- 示例: `swagger-viewer.test.jsx`, `api-tester.test.jsx`

#### 测试结构
```javascript
import { render, screen, waitFor } from '@testing-library/react';
import SwaggerViewer from './swagger-viewer';

describe('SwaggerViewer', () => {
    it('should render swagger viewer', () => {
        const { getByText } = render(<SwaggerViewer />);
        expect(getByText('Swagger')).toBeInTheDocument();
    });
    
    it('should load swagger on mount', async () => {
        const { getByText } = render(<SwaggerViewer />);
        await waitFor(() => {
            expect(getByText('API Documentation')).toBeInTheDocument();
        });
    });
});
```

### 7. 代码格式化

#### 使用 Prettier
```bash
npm run format
```

#### ESLint 检查
```bash
npm run lint
```

---

## 通用规范

### 1. 提交规范

#### 提交信息格式
```
<type>(<scope>): <subject>

<body>

<footer>
```

#### 类型
- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `style`: 代码格式化
- `refactor`: 代码重构
- `test`: 添加测试
- `chore`: 构建工具或依赖更新

#### 示例
```
feat(parser): add support for parsing swagger tags

- Add comment parser for swagger tags
- Support @Router, @Param, @Success tags
- Add unit tests

Closes #123
```

### 2. 代码审查

#### 审查清单
- [ ] 代码遵循规范
- [ ] 有适当的注释
- [ ] 有单元测试
- [ ] 没有明显的 bug
- [ ] 性能可接受
- [ ] 安全性考虑

### 3. 文档规范

#### README 结构
- 项目描述
- 功能特性
- 安装说明
- 使用示例
- API 文档
- 贡献指南

#### 代码文档
- 每个模块有 README
- 每个公共 API 有文档
- 复杂逻辑有详细注释

### 4. 性能考虑

#### Go 性能
- 避免不必要的内存分配
- 使用对象池减少 GC 压力
- 使用 pprof 进行性能分析

#### JavaScript 性能
- 避免不必要的重新渲染
- 使用 React.memo 优化组件
- 使用虚拟列表处理大数据集
- 使用代码分割和懒加载

### 5. 安全考虑

#### Go 安全
- 验证所有用户输入
- 使用参数化查询防止 SQL 注入
- 避免硬编码敏感信息

#### JavaScript 安全
- 防止 XSS 攻击
- 验证用户输入
- 使用 HTTPS 通信
- 避免在客户端存储敏感信息
- 使用 Content Security Policy (CSP)

---

## 代码审查检查清单

### 功能性
- [ ] 代码实现了需求的功能
- [ ] 代码处理了所有边界情况
- [ ] 代码处理了错误情况
- [ ] 代码没有明显的 bug

### 可读性
- [ ] 代码易于理解
- [ ] 变量名有意义
- [ ] 函数名清晰
- [ ] 有适当的注释

### 性能
- [ ] 没有明显的性能问题
- [ ] 没有不必要的循环
- [ ] 没有不必要的内存分配
- [ ] 算法复杂度合理

### 安全性
- [ ] 验证了所有输入
- [ ] 没有 SQL 注入漏洞
- [ ] 没有 XSS 漏洞
- [ ] 没有硬编码的敏感信息

### 测试
- [ ] 有单元测试
- [ ] 测试覆盖了主要逻辑
- [ ] 测试覆盖了边界情况
- [ ] 测试通过

### 文档
- [ ] 有适当的注释
- [ ] 有 API 文档
- [ ] 有使用示例
- [ ] 更新了相关文档
