# swag-gen API 设计规范

## API 基本信息

- **基础 URL**: `http://localhost:8080`
- **认证**: 暂无（可选后续添加）
- **响应格式**: JSON
- **字符编码**: UTF-8

## 响应格式标准

### 成功响应

```json
{
    "code": 0,
    "message": "success",
    "data": {
        // 实际数据
    }
}
```

### 错误响应

```json
{
    "code": 400,
    "message": "Invalid request",
    "errors": [
        {
            "field": "path",
            "message": "Path is required"
        }
    ]
}
```

---

## Swagger 文档 API

### 1. 获取 Swagger 文档

**请求**:
```
GET /swagger
```

**查询参数**:
- `format` (string, optional): 输出格式，支持 json/yaml，默认 json

**响应**:
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "openapi": "3.0.0",
        "info": {
            "title": "User API",
            "version": "1.0.0",
            "description": "User management API"
        },
        "paths": {
            "/users": {
                "get": {
                    "summary": "Get users",
                    "tags": ["User"],
                    "parameters": [],
                    "responses": {
                        "200": {
                            "description": "Success"
                        }
                    }
                }
            }
        }
    }
}
```

### 2. 获取 Swagger UI

**请求**:
```
GET /swagger/ui
```

**响应**: 返回 HTML 页面，展示 Swagger UI 界面

### 3. 获取所有 API 端点

**请求**:
```
GET /api/endpoints
```

**查询参数**:
- `tag` (string, optional): 按标签过滤
- `method` (string, optional): 按 HTTP 方法过滤

**响应**:
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "total": 25,
        "items": [
            {
                "id": "endpoint-001",
                "method": "GET",
                "path": "/api/users",
                "summary": "Get users",
                "description": "Get all users",
                "tags": ["User"],
                "parameters": [],
                "responses": {
                    "200": {
                        "description": "Success"
                    }
                }
            }
        ]
    }
}
```

---

## API 测试 API

### 1. 执行 API 测试

**请求**:
```
POST /api/test
Content-Type: application/json
```

**请求体**:
```json
{
    "method": "GET",
    "url": "http://localhost:8000/api/v1/users",
    "headers": {
        "Authorization": "Bearer token",
        "Content-Type": "application/json"
    },
    "query": {
        "page": "1",
        "pageSize": "10"
    },
    "body": null
}
```

**响应**:
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "id": "test-001",
        "statusCode": 200,
        "headers": {
            "Content-Type": "application/json"
        },
        "body": {
            "code": 0,
            "message": "success",
            "data": []
        },
        "duration": 123,
        "createdAt": "2024-01-01T00:00:00Z"
    }
}
```

### 2. 获取测试历史

**请求**:
```
GET /api/test/history
```

**查询参数**:
- `page` (int, optional): 页码，默认 1
- `pageSize` (int, optional): 每页数量，默认 10
- `endpoint` (string, optional): 端点过滤

**响应**:
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "total": 50,
        "page": 1,
        "pageSize": 10,
        "items": [
            {
                "id": "test-001",
                "method": "GET",
                "url": "http://localhost:8000/api/v1/users",
                "statusCode": 200,
                "duration": 123,
                "createdAt": "2024-01-01T00:00:00Z"
            }
        ]
    }
}
```

### 3. 获取测试详情

**请求**:
```
GET /api/test/:testId
```

**路径参数**:
- `testId` (string): 测试 ID

**响应**:
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "id": "test-001",
        "method": "GET",
        "url": "http://localhost:8000/api/v1/users",
        "headers": {
            "Authorization": "Bearer token"
        },
        "query": {
            "page": "1"
        },
        "body": null,
        "statusCode": 200,
        "responseHeaders": {
            "Content-Type": "application/json"
        },
        "responseBody": {
            "code": 0,
            "message": "success",
            "data": []
        },
        "duration": 123,
        "createdAt": "2024-01-01T00:00:00Z"
    }
}
```

### 4. 清空测试历史

**请求**:
```
DELETE /api/test/history
```

**响应**:
```json
{
    "code": 0,
    "message": "success",
    "data": null
}
```

---

## 健康检查 API

### 1. 获取服务健康状态

**请求**:
```
GET /health
```

**响应**:
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "status": "healthy",
        "timestamp": "2024-01-01T00:00:00Z",
        "version": "1.0.0"
    }
}
```

---

## 错误代码

| 代码 | 含义 | 说明 |
|------|------|------|
| 0 | 成功 | 请求成功 |
| 400 | 请求错误 | 请求参数无效 |
| 401 | 未授权 | 需要认证 |
| 403 | 禁止访问 | 无权限访问 |
| 404 | 未找到 | 资源不存在 |
| 409 | 冲突 | 资源冲突 |
| 500 | 服务器错误 | 内部服务器错误 |
| 503 | 服务不可用 | 服务暂时不可用 |

---

## 分页规范

所有列表 API 都支持分页，使用以下参数：

- `page` (int): 页码，从 1 开始
- `pageSize` (int): 每页数量，最大 100

分页响应格式：
```json
{
    "total": 100,
    "page": 1,
    "pageSize": 10,
    "items": []
}
```

---

## 时间格式

所有时间戳使用 ISO 8601 格式：`2024-01-01T00:00:00Z`

---

## 速率限制

- 暂无速率限制（可选后续添加）

---

## 版本控制

- 当前 API 版本: v1
- 基础 URL: `/`
- 未来版本将使用 `/v2` 等

---

## 向后兼容性

- 新增字段不会破坏现有客户端
- 删除字段会在新版本中进行
- 字段类型变更会在新版本中进行

---

## CORS 配置

支持跨域请求，允许的源：
- `http://localhost:3000`
- `http://localhost:5173`
- `http://localhost:8080`

允许的方法：
- GET
- POST
- PUT
- DELETE
- OPTIONS

允许的请求头：
- Content-Type
- Authorization

---

## 请求示例

### 使用 curl 获取 Swagger 文档

```bash
curl -X GET http://localhost:8080/swagger
```

### 使用 curl 执行 API 测试

```bash
curl -X POST http://localhost:8080/api/test \
  -H "Content-Type: application/json" \
  -d '{
    "method": "GET",
    "url": "http://localhost:8000/api/v1/users",
    "headers": {
      "Authorization": "Bearer token"
    }
  }'
```

### 使用 curl 获取测试历史

```bash
curl -X GET "http://localhost:8080/api/test/history?page=1&pageSize=10"
```

---

## 响应示例

### 成功响应示例

```json
{
    "code": 0,
    "message": "success",
    "data": {
        "openapi": "3.0.0",
        "info": {
            "title": "My API",
            "version": "1.0.0"
        }
    }
}
```

### 错误响应示例

```json
{
    "code": 400,
    "message": "Invalid request",
    "errors": [
        {
            "field": "url",
            "message": "URL is required"
        }
    ]
}
```

---

## 最佳实践

### 1. 错误处理
- 始终检查响应的 `code` 字段
- 根据错误代码采取相应的处理措施
- 显示用户友好的错误信息

### 2. 请求超时
- 建议设置请求超时为 30 秒
- 对于长时间运行的操作，考虑使用异步处理

### 3. 重试策略
- 对于 5xx 错误，建议重试 3 次
- 使用指数退避策略（1s, 2s, 4s）

### 4. 缓存策略
- 缓存 Swagger 文档（TTL: 1 小时）
- 缓存 API 端点列表（TTL: 30 分钟）
- 不缓存测试结果

### 5. 日志记录
- 记录所有 API 请求和响应
- 记录错误和异常
- 定期清理日志文件
