# swag-gen API 文档

## 概述

swag-gen 是一个强大的 Go API 文档生成工具，可以自动解析 Go 源代码中的注释，生成标准的 OpenAPI 3.0 规范文档。

## 核心功能

### 1. 代码解析
- 自动扫描 Go 项目源代码
- 支持递归目录扫描
- 支持并发文件解析
- 支持自定义注释标签格式

### 2. Swagger 生成
- 生成标准的 OpenAPI 3.0 规范
- 支持 JSON 和 YAML 格式输出
- 自动生成 Schema 定义
- 支持复杂的数据类型

### 3. CLI 工具
- 提供 `swag-gen init` 命令
- 支持灵活的参数配置
- 提供清晰的进度反馈
- 完善的错误处理

## 支持的注释标签

### 路由标签
```go
// @Router /api/users [GET]
```
定义 API 端点的路由和 HTTP 方法。

### 摘要和描述
```go
// @Summary 获取所有用户
// @Description 从数据库获取所有用户
```
定义 API 的摘要和详细描述。

### 标签
```go
// @Tags User
```
为 API 分类，支持多个标签。

### 参数
```go
// @Param page query int false "页码"
// @Param id path int true "用户 ID"
// @Param body body CreateUserRequest true "用户信息"
```
定义 API 的参数，支持 query、path、header、body 等位置。

### 响应
```go
// @Success 200 {array} User "成功"
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 500 {object} ErrorResponse "服务器错误"
```
定义 API 的响应，包括状态码、数据类型和描述。

### 已弃用
```go
// @Deprecated
```
标记 API 为已弃用。

## 数据类型支持

### 基本类型
- `string` - 字符串
- `int` - 整数
- `float` - 浮点数
- `bool` - 布尔值
- `array` - 数组

### 复杂类型
- `object` - 对象/结构体
- 嵌套结构体
- 指针类型
- 数组类型

## 使用示例

### 简单的 GET 端点
```go
// GetUsers 获取所有用户
// @Router /api/users [GET]
// @Summary 获取所有用户
// @Description 从数据库获取所有用户
// @Tags User
// @Success 200 {array} User
// @Failure 500 {object} ErrorResponse
func GetUsers() {}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
```

### 带参数的 GET 端点
```go
// GetUserByID 根据 ID 获取用户
// @Router /api/users/{id} [GET]
// @Summary 根据 ID 获取用户
// @Tags User
// @Param id path int true "用户 ID"
// @Success 200 {object} User
// @Failure 404 {object} ErrorResponse
func GetUserByID(id int) {}
```

### POST 端点
```go
// CreateUser 创建用户
// @Router /api/users [POST]
// @Summary 创建用户
// @Tags User
// @Param body body CreateUserRequest true "用户信息"
// @Success 201 {object} User
// @Failure 400 {object} ErrorResponse
func CreateUser(req CreateUserRequest) {}

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
```

### PUT 端点
```go
// UpdateUser 更新用户
// @Router /api/users/{id} [PUT]
// @Summary 更新用户
// @Tags User
// @Param id path int true "用户 ID"
// @Param body body UpdateUserRequest true "用户信息"
// @Success 200 {object} User
// @Failure 404 {object} ErrorResponse
func UpdateUser(id int, req UpdateUserRequest) {}
```

### DELETE 端点
```go
// DeleteUser 删除用户
// @Router /api/users/{id} [DELETE]
// @Summary 删除用户
// @Tags User
// @Param id path int true "用户 ID"
// @Success 204
// @Failure 404 {object} ErrorResponse
func DeleteUser(id int) {}
```

## 最佳实践

### 1. 注释规范
- 每个 API 函数都应该有完整的注释
- 使用清晰的摘要和描述
- 为所有参数提供描述
- 为所有可能的响应提供定义

### 2. 参数命名
- 使用有意义的参数名称
- 遵循 RESTful 命名规范
- 使用小写字母和下划线

### 3. 数据模型
- 为所有请求和响应定义结构体
- 使用 JSON 标签定义字段映射
- 为复杂类型提供清晰的定义

### 4. 错误处理
- 为所有可能的错误定义响应
- 使用标准的 HTTP 状态码
- 提供清晰的错误信息

### 5. 版本管理
- 在 API 信息中指定版本
- 使用语义版本控制
- 标记已弃用的 API

## 常见问题

### Q: 如何处理复杂的嵌套结构？
A: swag-gen 自动支持嵌套结构体。只需在注释中引用结构体类型，工具会自动生成相应的 Schema 定义。

### Q: 如何处理数组类型？
A: 使用 `{array} TypeName` 格式定义数组类型。例如：`@Success 200 {array} User`

### Q: 如何处理指针类型？
A: swag-gen 自动处理指针类型，生成的 Schema 会正确反映字段的可选性。

### Q: 如何处理自定义类型？
A: 定义结构体并在注释中引用即可。swag-gen 会自动生成相应的 Schema 定义。

### Q: 如何处理泛型类型？
A: 目前 swag-gen 不支持 Go 泛型。建议使用具体的类型定义。

## 输出文件

### swagger.json
标准的 OpenAPI 3.0 规范文档，JSON 格式。

### swagger.yaml
标准的 OpenAPI 3.0 规范文档，YAML 格式。

### swag-gen.yaml
项目配置文件，包含解析和生成的配置信息。

### README.md
项目说明文件，包含 API 的基本信息。

## 相关资源

- [OpenAPI 3.0 规范](https://spec.openapis.org/oas/v3.0.3)
- [Swagger UI](https://swagger.io/tools/swagger-ui/)
- [Go 官方文档](https://golang.org/doc/)

## 支持和反馈

如有问题或建议，请提交 Issue 或 Pull Request。

---

**版本**: 1.0.0  
**最后更新**: 2026 年 2 月 3 日
