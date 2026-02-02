# Phase 1: 基础设施 - 设计文档

## 架构设计

### 系统架构

```
┌─────────────────────────────────────────┐
│         CLI 工具 (Cobra)                 │
│  ┌──────────────────────────────────┐   │
│  │ init 命令 │ server 命令 │ help   │   │
│  └──────────────────────────────────┘   │
└────────────────┬────────────────────────┘
                 │
┌────────────────▼────────────────────────┐
│      配置管理系统 (Viper)                │
│  ┌──────────────────────────────────┐   │
│  │ 配置加载 │ 环境变量 │ 默认值     │   │
│  └──────────────────────────────────┘   │
└────────────────┬────────────────────────┘
                 │
┌────────────────▼────────────────────────┐
│      HTTP 服务框架 (Gin)                │
│  ┌──────────────────────────────────┐   │
│  │ 路由 │ 中间件 │ 处理器 │ 错误处理│   │
│  └──────────────────────────────────┘   │
└────────────────┬────────────────────────┘
                 │
┌────────────────▼────────────────────────┐
│      日志系统 (Zap)                     │
│  ┌──────────────────────────────────┐   │
│  │ 日志记录 │ 日志级别 │ 输出格式   │   │
│  └──────────────────────────────────┘   │
└─────────────────────────────────────────┘
```

### 模块设计

#### 1. CLI 模块 (cmd/swag-gen/)
- **main.go**: 应用入口，定义根命令
- **init.go**: init 命令实现
- **server.go**: server 命令实现

#### 2. 配置模块 (pkg/config/)
- **config.go**: 配置管理实现

#### 3. 日志模块 (pkg/logger/)
- **logger.go**: 日志系统实现

#### 4. 服务器模块 (pkg/server/)
- **server.go**: HTTP 服务器实现
- **middleware.go**: 中间件实现

### 数据流

#### 初始化流程
```
用户执行 swag-gen init
    ↓
解析命令行参数
    ↓
加载配置文件
    ↓
初始化日志系统
    ↓
执行初始化逻辑
    ↓
生成输出文件
    ↓
返回结果
```

#### 服务器启动流程
```
用户执行 swag-gen server
    ↓
解析命令行参数
    ↓
加载配置文件
    ↓
初始化日志系统
    ↓
创建 HTTP 服务器
    ↓
注册路由和中间件
    ↓
启动服务器
    ↓
监听请求
```

## 详细设计

### 1. 配置管理系统

#### 配置结构
```go
type Config struct {
    Server   ServerConfig   // 服务器配置
    Project  ProjectConfig  // 项目配置
    Parser   ParserConfig   // 解析器配置
    Swagger  SwaggerConfig  // Swagger 配置
    Logger   LoggerConfig   // 日志配置
}
```

#### 配置加载优先级
1. 命令行参数（最高优先级）
2. 环境变量
3. 配置文件
4. 默认值（最低优先级）

### 2. CLI 工具

#### init 命令
```bash
swag-gen init [flags]

Flags:
  -p, --path string        API 源代码路径 (default "./")
  -o, --output string      输出文档路径 (default "./docs")
  -t, --title string       API 标题 (default "API Documentation")
  -v, --version string     API 版本 (default "1.0.0")
  -d, --description string API 描述
```

#### server 命令
```bash
swag-gen server [flags]

Flags:
  -p, --port int    服务器端口 (default 8080)
  -h, --host string 服务器地址 (default "0.0.0.0")
  -d, --docs string 文档路径 (default "./docs")
```

### 3. HTTP 服务框架

#### 路由设计
```
GET  /health              # 健康检查
GET  /swagger             # 获取 Swagger 文档
GET  /swagger/ui          # 获取 Swagger UI
GET  /api/endpoints       # 获取所有端点
POST /api/test            # 执行 API 测试
GET  /api/test/history    # 获取测试历史
GET  /api/test/:testId    # 获取测试详情
DELETE /api/test/history  # 清空测试历史
```

#### 中间件设计
- **LoggerMiddleware**: 记录所有请求
- **CORSMiddleware**: 处理跨域请求
- **ErrorHandlerMiddleware**: 处理错误

### 4. 日志系统

#### 日志级别
- DEBUG: 调试信息
- INFO: 一般信息
- WARN: 警告信息
- ERROR: 错误信息

#### 日志格式
- JSON: 结构化日志
- TEXT: 文本日志

## 接口设计

### 配置接口
```go
// 加载配置
func Load(configPath string) (*Config, error)

// 保存配置
func (c *Config) Save(path string) error
```

### 日志接口
```go
// 初始化日志
func Init(level string, format string) error

// 获取日志记录器
func GetLogger() *zap.Logger

// 日志方法
func Debug(msg string, fields ...zap.Field)
func Info(msg string, fields ...zap.Field)
func Warn(msg string, fields ...zap.Field)
func Error(msg string, fields ...zap.Field)
```

### 服务器接口
```go
// 创建服务器
func New(cfg *config.Config) *Server

// 启动服务器
func (s *Server) Start() error

// 停止服务器
func (s *Server) Stop() error
```

## 错误处理

### 错误分类
- 配置错误: 配置文件不存在、格式错误等
- 命令错误: 命令参数错误、命令不存在等
- 服务器错误: 端口被占用、启动失败等
- 日志错误: 日志初始化失败等

### 错误处理策略
- 记录详细的错误信息
- 提供清晰的错误提示
- 优雅地处理错误

## 性能考虑

### 优化策略
- 使用连接池
- 使用缓存
- 异步处理
- 并发处理

### 性能指标
- 服务器启动时间 < 1 秒
- 单个请求响应时间 < 100ms
- 支持 1000+ 并发连接

## 安全考虑

### 安全措施
- CORS 配置
- 输入验证
- 错误信息不泄露敏感信息
- 日志不记录敏感信息

## 测试策略

### 单元测试
- 配置模块测试
- 日志模块测试
- 服务器模块测试
- CLI 命令测试

### 集成测试
- 配置加载和使用
- 日志记录和输出
- 服务器启动和请求处理

### 性能测试
- 服务器启动时间
- 请求响应时间
- 并发连接处理

## 部署考虑

### 部署方式
- 本地开发
- Docker 容器
- Kubernetes 集群

### 配置管理
- 开发环境配置
- 测试环境配置
- 生产环境配置

## 相关文档

- [开发计划](../../steering/development-plan.md)
- [代码规范](../../steering/code-standards.md)
- [API 设计](../../steering/api-design.md)
- [技术栈](../../steering/tech.md)
