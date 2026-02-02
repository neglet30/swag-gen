# swag-gen 技术栈确认文档

## 📋 技术栈确认

### ✅ 后端技术栈（已确认）

| 组件 | 技术 | 版本 | 用途 |
|------|------|------|------|
| 语言 | Go | 1.25.5+ | 核心开发语言 |
| Web 框架 | Gin | 最新 | HTTP 服务框架 |
| CLI 框架 | Cobra | 最新 | 命令行工具 |
| 代码解析 | go/ast, go/parser | 标准库 | AST 解析 |
| 配置管理 | Viper | 最新 | 配置文件处理 |
| 日志 | Zap | 最新 | 日志记录 |
| 测试 | Testify | 最新 | 单元测试 |

### ✅ 前端技术栈（已确认）

| 组件 | 技术 | 版本 | 用途 |
|------|------|------|------|
| 框架 | React | 18+ | UI 框架 |
| UI 库 | Material-UI (MUI) | 最新 | 组件库 |
| HTTP 客户端 | Axios | 最新 | API 调用 |
| 状态管理 | Redux Toolkit | 最新 | 状态管理 |
| 路由 | React Router | v6 | 页面路由 |
| 编辑器 | Monaco Editor | 最新 | 代码编辑 |
| API 文档 | Swagger UI | 最新 | 文档展示 |
| 表单处理 | React Hook Form | 最新 | 表单管理 |
| 数据验证 | Zod | 最新 | 数据验证 |
| 构建工具 | Vite | 最新 | 前端构建 |
| 代码格式 | Prettier | 最新 | 代码格式化 |
| 代码检查 | ESLint | 最新 | 代码检查 |

### ✅ 开发工具（已确认）

| 工具 | 用途 |
|------|------|
| Git | 版本控制 |
| Docker | 容器化 |
| Docker Compose | 本地开发编排 |
| GitHub Actions | CI/CD |
| golangci-lint | Go 代码检查 |
| Nginx | Web 服务器 |

---

## 📚 已完成的文档

### 1. 技术栈文档 ✅
**文件**: `.kiro/steering/tech.md`

**内容**:
- 后端技术栈详细说明
- 前端技术栈详细说明（React 确认）
- 开发工具列表
- 常用命令
- 依赖管理
- 代码质量标准

### 2. API 设计文档 ✅
**文件**: `.kiro/steering/api-design.md`

**内容**:
- API 基本信息
- 响应格式标准
- Swagger 文档 API（3 个端点）
- API 测试 API（4 个端点）
- 健康检查 API（1 个端点）
- 错误代码定义
- 分页规范
- 时间格式规范
- CORS 配置
- 请求示例
- 最佳实践

**总计**: 8 个 API 端点

### 3. 代码规范文档 ✅
**文件**: `.kiro/steering/code-standards.md`

**内容**:
- Go 代码规范（8 个方面）
  - 命名规范
  - 代码组织
  - 注释规范
  - 错误处理
  - 并发编程
  - 测试规范
  - 代码格式化
  - 最佳实践

- React/JavaScript 代码规范（7 个方面）
  - 命名规范
  - 代码组织
  - 注释规范
  - 错误处理
  - React 最佳实践
  - 测试规范
  - 代码格式化

- 通用规范（5 个方面）
  - 提交规范
  - 代码审查
  - 文档规范
  - 性能考虑
  - 安全考虑

- 代码审查检查清单

---

## 🎯 技术选择理由

### 为什么选择 React？

1. **生态成熟**: React 拥有最成熟的生态系统
2. **组件丰富**: Material-UI 提供了丰富的企业级组件
3. **社区活跃**: 最大的前端社区，问题解决快
4. **学习资源**: 最多的教程和文档
5. **企业支持**: Facebook 官方维护，企业级应用广泛使用
6. **性能优秀**: 虚拟 DOM 和优化的渲染性能
7. **开发效率**: 丰富的开发工具和调试工具

### 为什么选择 Material-UI？

1. **设计规范**: 遵循 Material Design 规范
2. **组件完整**: 提供了完整的企业级组件库
3. **定制性强**: 支持主题定制和样式覆盖
4. **文档完善**: 详细的文档和示例
5. **社区活跃**: 大量的第三方扩展和插件

### 为什么选择 Redux Toolkit？

1. **状态管理**: 简化了 Redux 的使用
2. **开发效率**: 减少了样板代码
3. **性能优化**: 内置了性能优化
4. **调试工具**: 支持 Redux DevTools
5. **中间件支持**: 支持异步操作和中间件

### 为什么选择 Vite？

1. **构建速度**: 比 Webpack 快 10 倍以上
2. **开发体验**: 即时热更新（HMR）
3. **生产优化**: 自动代码分割和优化
4. **配置简单**: 零配置开箱即用
5. **现代工具**: 基于 ES 模块的现代构建工具

---

## 📊 技术栈对比

### React vs Vue

| 特性 | React | Vue |
|------|-------|-----|
| 学习曲线 | 陡峭 | 平缓 |
| 生态系统 | 最大 | 中等 |
| 社区规模 | 最大 | 中等 |
| 企业应用 | 广泛 | 中等 |
| 性能 | 优秀 | 优秀 |
| 开发效率 | 中等 | 高 |
| 组件库 | 丰富 | 中等 |

**选择 React 的原因**: 企业级应用、生态成熟、社区活跃

### Material-UI vs Ant Design

| 特性 | Material-UI | Ant Design |
|------|-------------|-----------|
| 设计规范 | Material Design | Ant Design |
| 组件数量 | 丰富 | 丰富 |
| 定制性 | 强 | 强 |
| 文档 | 完善 | 完善 |
| 社区 | 活跃 | 活跃 |
| 学习曲线 | 平缓 | 平缓 |

**选择 Material-UI 的原因**: 国际化支持、设计规范、组件完整

---

## 🔧 开发环境要求

### 最低要求
- Go 1.25.5+
- Node.js 18+
- npm 9+ 或 yarn 3+
- Docker 20.10+
- Git 2.30+

### 推荐配置
- Go 1.25.5+
- Node.js 20 LTS
- npm 10+
- Docker 24+
- Git 2.40+

### 开发工具
- VS Code 或 GoLand（后端开发）
- VS Code 或 WebStorm（前端开发）
- Postman 或 Insomnia（API 测试）
- Docker Desktop（容器管理）

---

## 📦 依赖管理

### Go 依赖安装
```bash
go mod download
go mod tidy
```

### npm 依赖安装
```bash
cd web
npm install
```

### 更新依赖
```bash
# Go
go get -u ./...

# npm
npm update
```

---

## 🚀 快速开始

### 1. 克隆项目
```bash
git clone https://github.com/your-org/swag-gen.git
cd swag-gen
```

### 2. 安装依赖
```bash
# 后端依赖
go mod download

# 前端依赖
cd web
npm install
cd ..
```

### 3. 启动开发环境
```bash
# 启动后端
go run ./cmd/swag-gen

# 启动前端（新终端）
cd web
npm run dev
```

### 4. 访问应用
- 前端: http://localhost:5173
- 后端 API: http://localhost:8080

---

## 📋 下一步行动

### 立即开始
1. ✅ 技术栈已确认（React + Material-UI + Redux Toolkit）
2. ✅ API 设计文档已完成
3. ✅ 代码规范文档已完成
4. ⏳ 开始 Phase 1 开发

### Phase 1 任务
1. 创建项目目录结构
2. 初始化 Go 模块和依赖
3. 初始化 React 项目
4. 创建基础 HTTP 服务框架
5. 实现配置管理系统

### 开发工具准备
- [ ] 安装 Go 1.25.5+
- [ ] 安装 Node.js 20 LTS
- [ ] 安装 Docker 24+
- [ ] 安装 VS Code 或 GoLand
- [ ] 安装 golangci-lint
- [ ] 配置 Git

---

## 📞 技术支持

### 文档参考
- [技术栈文档](./steering/tech.md)
- [API 设计文档](./steering/api-design.md)
- [代码规范文档](./steering/code-standards.md)
- [开发计划](./steering/development-plan.md)

### 常见问题
- **如何安装依赖?** → 查看快速开始部分
- **如何启动开发环境?** → 查看快速开始部分
- **代码规范是什么?** → 查看代码规范文档
- **API 设计是什么?** → 查看 API 设计文档

---

## ✅ 确认清单

- [x] 后端技术栈确认
- [x] 前端技术栈确认（React）
- [x] 开发工具确认
- [x] 技术栈文档完成
- [x] API 设计文档完成
- [x] 代码规范文档完成
- [ ] 开始 Phase 1 开发

---

**技术栈确认完成！** 🎉

所有技术选择已确认，详细文档已完成。项目已准备好开始开发。

**确认日期**: 2024 年 1 月  
**前端框架**: React 18+  
**UI 库**: Material-UI  
**状态管理**: Redux Toolkit  
**构建工具**: Vite
