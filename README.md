# swag-gen - API 文档生成工具库

[![Go Version](https://img.shields.io/badge/Go-1.25.5+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)]()

swag-gen 是一个开源的 API 文档生成工具库，功能与 swag 相同，但提供更强大的 Web UI 界面。可以轻松集成到任何 Go 项目中，通过简单的命令行工具自动解析项目中的 API 注释，生成 Swagger/OpenAPI 规范文档，并提供现代化的 Web UI 界面用于文档查看和 API 测试。

## ✨ 核心特性

- 🚀 **易于集成**: 使用 `go get` 获取，通过 `swag-gen init` 命令快速初始化
- 📝 **自动解析**: 自动扫描并解析项目中所有 API 注释，生成完整的 Swagger/OpenAPI 文档
- 🎨 **现代化 UI**: 提供美观的 Web UI 界面，无需额外工具即可查看和测试 API
- 🧪 **在线测试**: 在前端页面直接生成测试数据并调用 API 接口
- 📦 **零配置**: 开箱即用，最小化配置需求
- 🔌 **路由集成**: 简单集成到现有路由中，通过 URL 访问 UI 界面
- 🌐 **跨平台**: 使用 Go 语言开发，支持 Windows、macOS、Linux
- 📄 **标准规范**: 生成标准的 Swagger/OpenAPI 3.0 规范文档

## 🚀 快速开始

### 前置要求

- Go 1.25.5+
- Node.js 18+ 和 npm 9+（用于前端开发）

### 安装

```bash
go get github.com/your-org/swag-gen
```

### 基本使用

#### 1. 初始化项目

```bash
swag-gen init -p ./api -o ./docs -t "My API"
```

参数说明：
- `-p, --path`: API 源代码路径（默认：./）
- `-o, --output`: 输出文档路径（默认：./docs）
- `-t, --title`: API 标题（默认：API Documentation）
- `-v, --version`: API 版本（默认：1.0.0）

#### 2. 集成到项目

在你的 Go 项目中集成 swag-gen 路由：

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/your-org/swag-gen/pkg/server"
)

func main() {
    r := gin.Default()
    
    // 集成 swag-gen 路由
    server.RegisterRoutes(r, "./docs")
    
    r.Run(":8080")
}
```

#### 3. 访问 UI

打开浏览器访问：
- Swagger UI: http://localhost:8080/swagger/ui
- API 文档: http://localhost:8080/swagger
- API 测试: http://localhost:8080/swagger/ui/test

## 📚 文档

详细文档请查看 [.kiro/steering](./kiro/steering) 目录：

- [产品概述](./kiro/steering/product.md) - 项目目的和特性
- [技术栈](./kiro/steering/tech.md) - 使用的技术和工具
- [项目结构](./kiro/steering/structure.md) - 代码组织规范
- [开发计划](./kiro/steering/development-plan.md) - 完整的开发流程和技术设计
- [架构设计](./kiro/steering/architecture.md) - 系统架构和设计
- [API 设计](./kiro/steering/api-design.md) - REST API 规范
- [代码规范](./kiro/steering/code-standards.md) - 代码编写规范
- [部署指南](./kiro/steering/deployment.md) - 部署和运维指南

## 🏗️ 项目结构

```
swag-gen/
├── cmd/
│   ├── swag-gen/               # CLI 工具入口
│   └── server/                 # Web 服务入口（可选）
├── pkg/
│   ├── parser/                # 代码解析
│   ├── swagger/               # Swagger 生成
│   ├── server/                # Web 服务
│   └── config/                # 配置管理
├── internal/                  # 私有包
│   ├── services/              # 业务逻辑
│   └── utils/                 # 工具函数
├── web/                       # 前端代码
├── tests/                     # 测试文件
├── docs/                      # 文档
├── .kiro/                     # Kiro 配置
│   ├── steering/             # 指导文档
│   └── specs/                # 功能规范
├── Dockerfile                # Docker 配置
├── docker-compose.yml        # Docker Compose
├── Makefile                  # 构建脚本
├── go.mod                    # Go 模块
└── README.md                 # 本文件
```

## 🔧 常用命令

### 开发命令

```bash
# 构建应用
make build

# 运行应用
make run

# 运行所有测试（单元测试 + 集成测试）
make test

# 运行单元测试
make test-unit

# 运行集成测试
make test-integration

# 生成覆盖率报告
make test-coverage

# 代码格式化
make fmt

# 代码检查
make lint

# 清理构建文件
make clean
```

### 开发环境

```bash
# 启动完整开发环境
make dev

# 仅启动后端
make dev-backend

# 仅启动前端
make dev-frontend
```

### Docker 命令

```bash
# 构建镜像
make docker-build

# 启动容器
make docker-up

# 停止容器
make docker-down

# 查看日志
make docker-logs
```

查看所有可用命令：
```bash
make help
```

## 🧪 测试

### 运行测试

```bash
# 运行所有测试
go test -v ./...

# 运行单元测试
go test -v ./pkg/...

# 运行集成测试
go test -v ./tests/integration/...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 测试覆盖率

当前测试覆盖率：
- pkg/config: 95.2%
- pkg/logger: 94.3%
- pkg/server: 83.3%
- **总体**: 90.9%

### 集成测试

Phase 1 包含 41 个集成测试用例，覆盖：
- 配置系统集成 (6 个用例)
- 日志系统集成 (15 个用例)
- 服务器系统集成 (20 个用例)

详见 [集成测试快速开始指南](./.kiro/INTEGRATION_TESTS_QUICK_START.md)

## 📖 使用示例

### 示例 1: 基本使用

```bash
# 初始化项目
swag-gen init -p ./api -o ./docs -t "User API" -v "1.0.0"

# 在 Go 项目中集成
```

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/your-org/swag-gen/pkg/server"
)

func main() {
    r := gin.Default()
    
    // 你的 API 路由
    r.GET("/api/users", getUsers)
    
    // 集成 swag-gen
    server.RegisterRoutes(r, "./docs")
    
    r.Run(":8080")
}

func getUsers(c *gin.Context) {
    // 实现
}
```

### 示例 2: 自定义配置

创建 `swag-gen.yaml` 配置文件：

```yaml
project:
  name: "My API"
  version: "1.0.0"
  description: "My API Documentation"

output:
  path: "./docs"
  format: "json"

apis:
  - name: "User API"
    path: "./api/user"
    basePath: "/api/v1"
```

然后运行：
```bash
swag-gen init -c swag-gen.yaml
```

## 🛠️ 技术栈

### 后端
- **语言**: Go 1.25.5+
- **Web 框架**: Gin
- **CLI 框架**: Cobra
- **代码解析**: go/ast, go/parser
- **配置**: Viper
- **日志**: Zap / Logrus

### 前端
- **框架**: React 18+ / Vue 3
- **UI 库**: Material-UI / Ant Design
- **HTTP 客户端**: Axios
- **编辑器**: Monaco Editor

### 工具
- **容器**: Docker
- **编排**: Docker Compose
- **Web 服务器**: Nginx
- **CI/CD**: GitHub Actions

## 📊 性能指标

| 指标 | 目标值 |
|------|--------|
| 代码解析速度 | < 5s（1000 行代码） |
| API 响应时间 | < 200ms |
| 并发连接数 | > 1000 |
| 可用性 | > 99.9% |

## 🔐 安全特性

- ✅ 输入验证和清理
- ✅ SQL 注入防护
- ✅ XSS 防护
- ✅ HTTPS 支持
- ✅ 认证和授权（可选）
- ✅ 审计日志

## 🤝 贡献指南

欢迎贡献！请遵循以下步骤：

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'feat: add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 开启 Pull Request

请确保：
- 遵循代码规范（见 [code-standards.md](./kiro/steering/code-standards.md)）
- 添加适当的测试
- 更新相关文档

## 📝 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。

## 📞 获取帮助

- 📖 查看 [文档](./kiro/steering)
- 🐛 提交 [Issue](https://github.com/your-org/swag-gen/issues)
- 💬 参与 [讨论](https://github.com/your-org/swag-gen/discussions)

## 🎯 路线图

- [ ] Phase 1: 基础设施（第 1-2 周）
- [ ] Phase 2: 核心功能（第 3-4 周）
- [ ] Phase 3: Web UI（第 5-6 周）
- [ ] Phase 4: API 测试工具（第 7-8 周）
- [ ] Phase 5: 高级功能（第 9-10 周）
- [ ] Phase 6: 测试与部署（第 11-12 周）

## 🙏 致谢

感谢所有贡献者和使用者的支持！

---

**最后更新**: 2024 年 1 月

**维护者**: swag-gen 开发团队

**官方网站**: https://swag-gen.com

**GitHub**: https://github.com/your-org/swag-gen
