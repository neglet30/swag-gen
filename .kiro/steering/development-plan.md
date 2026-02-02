# swag-gen 项目开发计划与技术设计

## 项目概述

swag-gen 是一个开源的 API 文档生成工具库，功能与 swag 相同，但提供更强大的 Web UI 界面。可以轻松集成到任何 Go 项目中，通过简单的命令行工具自动解析项目中的 API 注释，生成 Swagger/OpenAPI 规范文档，并提供现代化的 Web UI 界面用于文档查看和 API 测试。

## 核心使用流程
1. **安装**: `go get github.com/your-org/swag-gen`
2. **初始化**: `swag-gen init [参数]` - 扫描项目并生成文档
3. **集成路由**: 在项目中集成 swag-gen 提供的路由
4. **访问 UI**: 通过浏览器访问 API 文档和测试界面

---

## 第一阶段：核心功能开发

### 1.1 代码解析与 Swagger 生成

**功能描述**：
- 扫描 Go 项目源代码
- 解析代码注释中的 Swagger 标签
- 生成标准的 OpenAPI 3.0 规范文档

**技术实现**：
- 使用 `go/ast` 包进行 AST 解析
- 使用 `go/parser` 包解析 Go 源文件
- 支持自定义注释标签格式（如 `@Router`、`@Param`、`@Success` 等）
- 输出 JSON 格式的 OpenAPI 规范

**关键模块**：
```
pkg/parser/
├── ast_parser.go       # AST 解析器
├── comment_parser.go   # 注释解析器
├── swagger_builder.go  # Swagger 文档构建器
└── models.go           # 数据模型
```

### 1.2 CLI 命令行工具

**功能描述**：
- 提供 `swag-gen init` 命令
- 支持自定义参数配置
- 自动扫描项目并生成文档

**技术实现**：
- 使用 `cobra` 库构建 CLI
- 支持配置文件和命令行参数
- 生成 `docs` 目录和相关文件

**命令示例**：
```bash
swag-gen init -p ./api -o ./docs -t "My API"
```

### 1.3 项目集成与配置

**功能描述**：
- 支持项目级配置文件
- 支持多个 API 分组
- 支持自定义输出路径

**技术实现**：
- 使用 YAML/JSON 配置文件
- 配置文件示例：`swag-gen.yaml`
- 支持环境变量覆盖

**配置示例**：
```yaml
project:
  name: "My API"
  version: "1.0.0"
  description: "API Documentation"

output:
  path: "./docs"
  format: "json"

apis:
  - name: "User API"
    path: "./api/user"
    basePath: "/api/v1"
```

---

## 第二阶段：Web UI 开发

### 2.1 后端 API 服务

**功能描述**：
- 提供 REST API 服务
- 动态加载和解析项目
- 实时生成 Swagger 文档
- 支持 API 代理调用

**技术栈**：
- **框架**: Gin Web Framework
- **HTTP 服务**: 标准库 `net/http`
- **JSON 处理**: `encoding/json`
- **CORS 支持**: 中间件处理跨域请求

**核心 API 端点**：
```
GET  /swagger              # 获取 Swagger 文档
GET  /swagger/ui           # 获取 UI 界面
POST /api/test             # 测试 API 调用
GET  /api/endpoints        # 获取所有端点
```

**关键模块**：
```
pkg/server/
├── router.go            # 路由定义
├── handlers.go          # 处理器
├── middleware.go        # 中间件
└── response.go          # 响应格式
```

### 2.2 前端 UI 界面

**功能描述**：
- Swagger 文档展示
- API 测试工具
- 实时代码预览

**技术栈**：
- **框架**: React 18+ 或 Vue 3
- **UI 组件库**: Material-UI 或 Ant Design
- **HTTP 客户端**: Axios
- **编辑器**: Monaco Editor（代码编辑）
- **API 文档展示**: Swagger UI 或 ReDoc

**页面结构**：
```
web/
├── public/
│   └── index.html
├── src/
│   ├── components/
│   │   ├── SwaggerViewer.jsx    # Swagger 文档查看器
│   │   ├── APITester.jsx        # API 测试工具
│   │   └── CodeEditor.jsx       # 代码编辑器
│   ├── pages/
│   │   ├── Dashboard.jsx        # 仪表板
│   │   └── APITest.jsx          # API 测试页面
│   ├── services/
│   │   └── api.js               # API 调用服务
│   └── App.jsx
└── package.json
```

---

## 第三阶段：高级功能

### 3.1 API 代理与测试

**功能描述**：
- 在 Web UI 中直接调用 API
- 支持请求头、参数、Body 自定义
- 显示响应结果和性能指标

**技术实现**：
- 后端代理请求到目标 API
- 记录请求/响应日志
- 支持请求历史记录

### 3.2 文档导出

**功能描述**：
- 导出为 HTML 文档
- 导出为 PDF 文档
- 导出为 Markdown 文档

**技术栈**：
- **HTML 生成**: 模板引擎（`html/template`）
- **PDF 生成**: `github.com/go-echarts/go-echarts` 或 `wkhtmltopdf`
- **Markdown 生成**: 自定义生成器

### 3.3 版本管理与历史

**功能描述**：
- 保存文档版本历史
- 支持版本对比
- 支持版本回滚

**技术实现**：
- 使用文件系统存储版本信息
- Git 集成（可选）

---

## 完整技术栈总结

### 后端技术栈

| 层级 | 技术 | 用途 |
|------|------|------|
| **语言** | Go 1.25.5+ | 核心开发语言 |
| **Web 框架** | Gin | HTTP 服务框架 |
| **CLI 框架** | Cobra | 命令行工具 |
| **代码解析** | go/ast, go/parser | AST 解析 |
| **配置管理** | Viper | 配置文件处理 |
| **日志** | Zap 或 Logrus | 日志记录 |
| **测试** | Testing, Testify | 单元测试 |
| **构建** | Go Build, Makefile | 项目构建 |

### 前端技术栈

| 层级 | 技术 | 用途 |
|------|------|------|
| **框架** | React 18+ 或 Vue 3 | UI 框架 |
| **UI 组件** | Material-UI / Ant Design | 组件库 |
| **HTTP 客户端** | Axios | API 调用 |
| **编辑器** | Monaco Editor | 代码编辑 |
| **API 文档** | Swagger UI / ReDoc | 文档展示 |
| **构建工具** | Vite / Webpack | 前端构建 |
| **包管理** | npm / yarn | 依赖管理 |

### 开发工具

| 工具 | 用途 |
|------|------|
| **版本控制** | Git |
| **代码检查** | golangci-lint |
| **代码格式** | go fmt, prettier |
| **API 测试** | Postman / curl |
| **容器化** | Docker |
| **CI/CD** | GitHub Actions / GitLab CI |

---

## 开发流程

### Phase 1: 基础设施（第 1-2 周）
- [ ] 项目结构搭建
- [ ] 配置管理系统
- [ ] CLI 命令行工具框架
- [ ] 基础 HTTP 服务框架

### Phase 2: 核心功能（第 3-4 周）
- [ ] Go 代码解析器
- [ ] Swagger 文档生成
- [ ] CLI init 命令实现
- [ ] 单元测试

### Phase 3: Web UI（第 5-6 周）
- [ ] 前端项目初始化
- [ ] Swagger 文档展示
- [ ] 基础样式和布局
- [ ] 路由集成

### Phase 4: API 测试工具（第 7-8 周）
- [ ] API 代理功能
- [ ] 测试工具 UI
- [ ] 请求/响应展示
- [ ] 测试历史记录

### Phase 5: 高级功能（第 9-10 周）
- [ ] 文档导出功能
- [ ] 版本管理
- [ ] 性能优化
- [ ] 安全加固

### Phase 6: 测试与部署（第 11-12 周）
- [ ] 集成测试
- [ ] 性能测试
- [ ] Docker 容器化
- [ ] 文档编写
- [ ] 发布准备

---

## 项目目录结构

```
swag-gen/
├── cmd/
│   ├── swag-gen/            # CLI 工具入口
│   │   └── main.go
│   └── server/              # Web 服务入口（可选）
│       └── main.go
├── pkg/
│   ├── parser/              # 代码解析模块
│   │   ├── ast_parser.go
│   │   ├── comment_parser.go
│   │   └── models.go
│   ├── swagger/             # Swagger 生成模块
│   │   ├── builder.go
│   │   └── models.go
│   ├── server/              # Web 服务模块
│   │   ├── router.go
│   │   ├── handlers.go
│   │   ├── middleware.go
│   │   └── response.go
│   └── config/              # 配置管理模块
│       └── config.go
├── internal/
│   ├── services/            # 业务逻辑
│   │   ├── parser_service.go
│   │   └── swagger_service.go
│   └── utils/               # 工具函数
│       └── utils.go
├── web/                     # 前端代码
│   ├── public/
│   ├── src/
│   │   ├── components/
│   │   ├── pages/
│   │   ├── services/
│   │   └── App.jsx
│   └── package.json
├── tests/                   # 测试文件
│   ├── unit/
│   ├── integration/
│   └── fixtures/
├── docs/                    # 文档
│   ├── api.md
│   ├── architecture.md
│   └── deployment.md
├── .kiro/
│   ├── steering/
│   └── specs/
├── Dockerfile              # Docker 配置
├── docker-compose.yml      # Docker Compose 配置
├── Makefile               # 构建脚本
├── go.mod                 # Go 模块定义
├── go.sum                 # Go 依赖锁定
└── README.md              # 项目说明
```

---

## 关键设计决策

### 1. 架构设计
- **工具库模式**: 作为可集成的工具库，而不是独立应用
- **分层架构**: Parser → Swagger Builder → Server → UI
- **模块化设计**: 各功能模块独立，便于维护和扩展

### 2. 代码解析策略
- 使用 Go 标准库 `go/ast` 和 `go/parser`
- 支持自定义注释标签格式
- 支持增量解析和缓存

### 3. 集成方式
- 提供 CLI 工具用于初始化和生成文档
- 提供路由处理器用于集成到现有项目
- 支持嵌入式使用和独立服务两种模式

### 4. 前端框架选择
- React 生态成熟，组件丰富
- Vue 学习曲线平缓，开发效率高
- 建议选择 React + Material-UI 或 Vue 3 + Ant Design

### 5. 部署方案
- Docker 容器化部署
- 支持 Docker Compose 本地开发
- 支持 Kubernetes 生产部署

---

## 下一步行动

1. **确认技术栈**: 确认前端框架选择（React 或 Vue）
2. **创建项目规范**: 编写详细的代码规范和 API 设计文档
3. **开始开发**: 从基础设施和核心功能开始
4. **建立 CI/CD**: 设置自动化测试和部署流程
