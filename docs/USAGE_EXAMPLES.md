# swag-gen 使用示例

## 目录
1. [基本用法](#基本用法)
2. [高级用法](#高级用法)
3. [完整示例](#完整示例)
4. [常见场景](#常见场景)

## 基本用法

### 安装
```bash
go get github.com/neglet30/swag-gen
```

### 初始化项目
```bash
swag-gen init -p ./api -o ./docs
```

这个命令会：
1. 扫描 `./api` 目录中的所有 Go 文件
2. 解析其中的 Swagger 注释
3. 生成 Swagger 文档到 `./docs` 目录

### 查看生成的文件
```bash
ls -la ./docs/
# 输出:
# swagger.json      - OpenAPI 3.0 规范文档
# swagger.yaml      - YAML 格式的规范文档
# swag-gen.yaml     - 项目配置文件
# README.md         - 项目说明文件
```

## 高级用法

### 指定 API 信息
```bash
swag-gen init \
  -p ./api \
  -o ./docs \
  -t "My API" \
  -v 2.0.0 \
  -d "My API Description"
```

### 指定输出格式
```bash
# 输出为 JSON 格式
swag-gen init -p ./api -o ./docs -f json

# 输出为 YAML 格式
swag-gen init -p ./api -o ./docs -f yaml
```

### 使用所有参数
```bash
swag-gen init \
  --path ./api \
  --output ./docs \
  --title "User API" \
  --version 1.0.0 \
  --description "User management API" \
  --format json
```

## 完整示例

### 示例 1: 简单的用户 API

**项目结构**:
```
myproject/
├── api/
│   └── user.go
├── main.go
└── go.mod
```

**api/user.go**:
```go
package api

// GetUsers 获取所有用户
// @Router /api/users [GET]
// @Summary 获取所有用户
// @Description 从数据库获取所有用户
// @Tags User
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {array} User "成功"
// @Failure 500 {object} ErrorResponse "服务器错误"
func GetUsers() {}

// GetUserByID 根据 ID 获取用户
// @Router /api/users/{id} [GET]
// @Summary 根据 ID 获取用户
// @Tags User
// @Param id path int true "用户 ID"
// @Success 200 {object} User "成功"
// @Failure 404 {object} ErrorResponse "用户不存在"
func GetUserByID(id int) {}

// CreateUser 创建用户
// @Router /api/users [POST]
// @Summary 创建用户
// @Tags User
// @Param body body CreateUserRequest true "用户信息"
// @Success 201 {object} User "创建成功"
// @Failure 400 {object} ErrorResponse "请求错误"
func CreateUser(req CreateUserRequest) {}

// User 用户模型
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
```

**生成文档**:
```bash
swag-gen init -p ./api -o ./docs -t "User API" -v 1.0.0
```

**生成的 swagger.json 片段**:
```json
{
  "openapi": "3.0.0",
  "info": {
    "title": "User API",
    "version": "1.0.0"
  },
  "paths": {
    "/api/users": {
      "get": {
        "summary": "获取所有用户",
        "description": "从数据库获取所有用户",
        "tags": ["User"],
        "parameters": [
          {
            "name": "page",
            "in": "query",
            "schema": {"type": "integer"}
          }
        ],
        "responses": {
          "200": {
            "description": "成功",
            "content": {
              "application/json": {
                "schema": {
                  "type": "array",
                  "items": {"$ref": "#/components/schemas/User"}
                }
              }
            }
          }
        }
      }
    }
  }
}
```

### 示例 2: 多模块 API

**项目结构**:
```
myproject/
├── api/
│   ├── user/
│   │   └── api.go
│   ├── post/
│   │   └── api.go
│   └── comment/
│       └── api.go
├── main.go
└── go.mod
```

**生成文档**:
```bash
swag-gen init -p ./api -o ./docs -t "Blog API" -v 1.0.0 -d "Blog management API"
```

这会自动扫描所有子目录并生成完整的 API 文档。

### 示例 3: 集成到 Gin 框架

**main.go**:
```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/neglet30/swag-gen/pkg/server"
)

func main() {
	r := gin.Default()

	// 集成 swag-gen 路由
	swagServer := server.NewServer("./docs")
	swagServer.RegisterRoutes(r)

	// 你的其他路由
	r.GET("/api/users", getUsers)

	r.Run(":8080")
}

func getUsers(c *gin.Context) {
	// 实现
}
```

访问 `http://localhost:8080/swagger/ui` 查看 API 文档。

## 常见场景

### 场景 1: 更新现有文档

当你修改了 API 代码后，重新运行命令即可更新文档：

```bash
swag-gen init -p ./api -o ./docs
```

### 场景 2: 生成多个版本的文档

```bash
# 生成 v1 版本
swag-gen init -p ./api/v1 -o ./docs/v1 -v 1.0.0

# 生成 v2 版本
swag-gen init -p ./api/v2 -o ./docs/v2 -v 2.0.0
```

### 场景 3: 在 CI/CD 中使用

**GitHub Actions 示例**:
```yaml
name: Generate API Docs

on: [push]

jobs:
  generate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: 1.25
      - run: go get github.com/neglet30/swag-gen
      - run: swag-gen init -p ./api -o ./docs
      - uses: actions/upload-artifact@v2
        with:
          name: api-docs
          path: ./docs
```

### 场景 4: 自定义输出格式

```bash
# 生成 JSON 格式
swag-gen init -p ./api -o ./docs -f json

# 生成 YAML 格式
swag-gen init -p ./api -o ./docs -f yaml
```

### 场景 5: 处理大型项目

对于大型项目，swag-gen 支持并发处理：

```bash
# 自动使用多个 goroutine 处理文件
swag-gen init -p ./api -o ./docs
```

## 输出示例

### 生成的 swagger.json 结构

```json
{
  "openapi": "3.0.0",
  "info": {
    "title": "API Title",
    "version": "1.0.0",
    "description": "API Description"
  },
  "paths": {
    "/api/endpoint": {
      "get": {
        "summary": "Endpoint summary",
        "description": "Endpoint description",
        "tags": ["Tag"],
        "parameters": [],
        "responses": {
          "200": {
            "description": "Success"
          }
        }
      }
    }
  },
  "components": {
    "schemas": {
      "Model": {
        "type": "object",
        "properties": {
          "id": {"type": "integer"},
          "name": {"type": "string"}
        }
      }
    }
  }
}
```

### 生成的 swag-gen.yaml 结构

```yaml
project:
  name: "API Title"
  version: "1.0.0"
  description: "API Description"

parser:
  path: "./api"
  exclude:
    - "vendor"
    - "node_modules"

output:
  path: "./docs"
  format: "json"

swagger:
  title: "API Title"
  version: "1.0.0"
  description: "API Description"
```

## 故障排除

### 问题 1: 找不到 API 端点

**原因**: 注释格式不正确或文件不在扫描路径中

**解决方案**:
1. 检查注释格式是否正确
2. 确保文件在指定的路径中
3. 检查文件是否是有效的 Go 文件

### 问题 2: 生成的文档不完整

**原因**: 某些注释标签缺失或格式不正确

**解决方案**:
1. 检查所有必需的标签是否都有
2. 验证标签格式是否正确
3. 查看日志输出了解详细信息

### 问题 3: 性能问题

**原因**: 项目过大或文件过多

**解决方案**:
1. 使用 `-p` 参数指定特定的 API 目录
2. 排除不必要的目录
3. 检查系统资源是否充足

---

**版本**: 1.0.0  
**最后更新**: 2026 年 2 月 3 日
