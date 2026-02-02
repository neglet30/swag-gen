# 技术栈

## 语言与运行时
- **语言**: Go 1.25.5+
- **模块**: `swag-gen`
- **操作系统**: 跨平台支持（Windows、macOS、Linux）

## 后端技术栈

### 核心框架
- **Web 框架**: Gin Web Framework
- **CLI 框架**: Cobra
- **HTTP 服务**: 标准库 `net/http`
- **代码解析**: `go/ast`、`go/parser`

### 配置与日志
- **配置管理**: Viper
- **日志**: Zap 或 Logrus
- **环境变量**: godotenv

### 测试与质量
- **单元测试**: Go testing 包 + Testify
- **代码覆盖**: go test -cover
- **代码检查**: golangci-lint
- **性能测试**: Go benchmark

## 前端技术栈

### 框架与库
- **框架**: React 18+
- **UI 组件库**: Material-UI (MUI)
- **HTTP 客户端**: Axios
- **状态管理**: Redux Toolkit
- **代码编辑器**: Monaco Editor
- **API 文档展示**: Swagger UI
- **路由**: React Router v6
- **表单处理**: React Hook Form
- **数据验证**: Zod

### 构建工具
- **构建工具**: Vite
- **包管理**: npm 或 yarn
- **代码格式**: Prettier
- **代码检查**: ESLint

## 开发工具

### 版本控制与协作
- **版本控制**: Git
- **代码托管**: GitHub / GitLab

### 容器化与部署
- **容器**: Docker
- **编排**: Docker Compose（开发）/ Kubernetes（生产）

### CI/CD
- **自动化**: GitHub Actions 或 GitLab CI
- **测试**: 自动化单元测试和集成测试
- **构建**: 自动化构建和发布

## 构建与开发命令

### Go 相关命令
```bash
# 构建项目
go build -o swag-gen ./cmd/swag-gen

# 运行测试
go test ./...

# 运行测试并显示覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# 格式化代码
go fmt ./...

# 代码检查（需要安装 golangci-lint）
golangci-lint run

# 整理依赖
go mod tidy

# 下载依赖
go mod download

# 运行应用
go run ./cmd/swag-gen

# 性能测试
go test -bench=. ./...
```

### 前端相关命令
```bash
# 进入前端目录
cd web

# 安装依赖
npm install

# 开发服务器
npm run dev

# 构建生产版本
npm run build

# 代码检查
npm run lint

# 代码格式化
npm run format

# 预览生产构建
npm run preview
```

### Docker 命令
```bash
# 构建 Docker 镜像
docker build -t swag-gen:latest .

# 运行容器
docker run -p 8080:8080 swag-gen:latest

# 使用 Docker Compose
docker-compose up -d
docker-compose down
```

## 依赖管理

### 主要 Go 依赖（计划）
```
github.com/gin-gonic/gin              # Web 框架
github.com/spf13/cobra                # CLI 框架
github.com/spf13/viper                # 配置管理
go.uber.org/zap                       # 日志
github.com/stretchr/testify           # 测试工具
```

### 主要 npm 依赖（计划）
```
react                                 # React 框架
@mui/material                         # Material-UI
axios                                 # HTTP 客户端
@reduxjs/toolkit                      # Redux 状态管理
react-redux                           # React Redux 绑定
react-router-dom                      # 路由
monaco-editor                         # 代码编辑器
swagger-ui-react                      # Swagger UI
react-hook-form                       # 表单处理
zod                                   # 数据验证
```

## 代码质量标准

### Go 代码规范
- 遵循 [Effective Go](https://golang.org/doc/effective_go) 规范
- 使用 `go fmt` 进行代码格式化
- 使用 `golangci-lint` 进行代码检查
- 编写清晰易读的代码，添加适当注释
- 与实现一起编写单元测试
- 测试覆盖率目标：>80%

### JavaScript/TypeScript 代码规范
- 遵循 ESLint 规则
- 使用 Prettier 进行代码格式化
- 编写清晰易读的代码，添加 JSDoc 注释
- 与实现一起编写单元测试
- 使用 TypeScript 进行类型检查（可选）

### 通用规范
- 提交前运行所有测试
- 保持代码简洁，避免过度设计
- 定期重构和优化代码
- 编写有意义的提交信息
