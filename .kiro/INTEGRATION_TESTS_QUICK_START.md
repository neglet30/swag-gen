# 集成测试快速开始指南

## 快速运行

### 运行所有集成测试

```bash
go test -v ./tests/integration/... -timeout 30s
```

### 运行特定的集成测试

```bash
# 运行配置集成测试
go test -v ./tests/integration/config_integration_test.go -timeout 30s

# 运行日志集成测试
go test -v ./tests/integration/logger_integration_test.go -timeout 30s

# 运行服务器集成测试
go test -v ./tests/integration/server_integration_test.go -timeout 30s
```

### 运行特定的测试用例

```bash
# 运行单个测试
go test -v ./tests/integration/ -run TestConfigLoadAndServerInit -timeout 30s

# 运行匹配模式的测试
go test -v ./tests/integration/ -run TestConfig -timeout 30s
```

### 生成覆盖率报告

```bash
# 生成覆盖率文件
go test ./... -coverprofile=coverage.out -timeout 30s

# 查看覆盖率报告
go tool cover -html=coverage.out
```

## 测试结构

### 配置集成测试 (6 个用例)

**文件**: `tests/integration/config_integration_test.go`

```
TestConfigLoadAndServerInit          - 配置加载与服务器初始化
TestEnvironmentVariableOverride      - 环境变量覆盖
TestConfigWithLoggerIntegration      - 配置与日志集成
TestConfigWithServerAndLogger        - 完整集成
TestMultipleConfigLoads              - 多次加载
TestConfigDefaultValues              - 默认值
```

### 日志集成测试 (15 个用例)

**文件**: `tests/integration/logger_integration_test.go`

```
TestLoggerWithServerIntegration      - 日志与服务器集成
TestLoggerLevelConfiguration         - 日志级别配置
TestLoggerFormatConfiguration        - 日志格式配置
TestLoggerWithConfigIntegration      - 日志与配置集成
TestLoggerMultipleInitializations    - 多次初始化
TestLoggerWithFields                 - 带字段的日志
TestLoggerSync                       - 日志同步
TestLoggerClose                      - 日志关闭
TestLoggerWithServerRequests         - 服务器请求日志
TestLoggerErrorHandling              - 错误日志
TestLoggerConcurrentAccess           - 并发访问
TestLoggerWithDifferentLevels        - 不同级别
TestLoggerJSONFormat                 - JSON 格式
TestLoggerTextFormat                 - 文本格式
TestLoggerWithComplexFields          - 复杂字段
```

### 服务器集成测试 (20 个用例)

**文件**: `tests/integration/server_integration_test.go`

```
TestCompleteRequestFlow              - 完整请求流程
TestConcurrentRequests               - 并发请求
TestDifferentHTTPMethods             - 不同 HTTP 方法
TestCORSMiddlewareIntegration        - CORS 中间件
TestPreflight                        - 预检请求
TestLoggerMiddlewareIntegration      - 日志中间件
TestMultipleEndpoints                - 多个端点
TestResponseFormat                   - 响应格式
TestErrorResponse                    - 错误响应
TestRequestWithBody                  - 请求体
TestRequestHeaders                   - 请求头
TestResponseTime                     - 响应时间
TestServerWithDifferentEnvironments  - 不同环境
TestRequestWithBody                  - 请求体处理
TestRequestHeaders                   - 请求头处理
TestResponseTime                     - 响应时间
TestServerWithDifferentEnvironments  - 环境配置
TestCompleteRequestFlow              - 完整流程
TestConcurrentRequests               - 并发处理
TestDifferentHTTPMethods             - HTTP 方法
```

## 常见命令

### 查看测试帮助

```bash
go test -h
```

### 运行测试并显示详细输出

```bash
go test -v ./tests/integration/...
```

### 运行测试并显示覆盖率

```bash
go test -v -cover ./tests/integration/...
```

### 运行测试并生成覆盖率文件

```bash
go test -v -coverprofile=coverage.out ./tests/integration/...
```

### 查看覆盖率报告

```bash
go tool cover -html=coverage.out
```

### 运行测试并显示性能信息

```bash
go test -v -bench=. ./tests/integration/...
```

### 运行测试并显示内存分配信息

```bash
go test -v -benchmem ./tests/integration/...
```

## 测试输出解读

### 成功的测试输出

```
=== RUN   TestConfigLoadAndServerInit
--- PASS: TestConfigLoadAndServerInit (0.01s)
```

- `RUN`: 测试开始运行
- `PASS`: 测试通过
- `(0.01s)`: 测试执行时间

### 失败的测试输出

```
=== RUN   TestConfigLoadAndServerInit
--- FAIL: TestConfigLoadAndServerInit (0.01s)
    config_integration_test.go:20: assertion failed
```

- `FAIL`: 测试失败
- 显示失败的文件和行号
- 显示失败的原因

### 跳过的测试输出

```
=== RUN   TestConfigLoadAndServerInit
--- SKIP: TestConfigLoadAndServerInit (0.00s)
    config_integration_test.go:20: skipped
```

- `SKIP`: 测试被跳过
- 显示跳过的原因

## 调试技巧

### 添加调试输出

在测试中添加 `t.Logf()` 来输出调试信息：

```go
func TestExample(t *testing.T) {
    t.Logf("Debug message: %v", value)
    // ...
}
```

### 运行单个测试进行调试

```bash
go test -v ./tests/integration/ -run TestConfigLoadAndServerInit
```

### 使用 -race 检测竞态条件

```bash
go test -race ./tests/integration/...
```

### 使用 -timeout 设置超时

```bash
go test -timeout 60s ./tests/integration/...
```

## 性能测试

### 运行性能基准测试

```bash
go test -bench=. ./tests/integration/...
```

### 运行性能基准测试并显示内存分配

```bash
go test -bench=. -benchmem ./tests/integration/...
```

### 运行性能基准测试并保存结果

```bash
go test -bench=. -benchmem ./tests/integration/ > benchmark.txt
```

## 覆盖率分析

### 查看覆盖率统计

```bash
go test -cover ./tests/integration/...
```

### 生成详细的覆盖率报告

```bash
go test -coverprofile=coverage.out ./tests/integration/...
go tool cover -html=coverage.out -o coverage.html
```

### 查看特定文件的覆盖率

```bash
go tool cover -html=coverage.out -o coverage.html
# 然后在浏览器中打开 coverage.html
```

## 持续集成

### GitHub Actions 示例

```yaml
name: Integration Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: 1.25.5
      - run: go test -v ./tests/integration/... -timeout 30s
      - run: go test -coverprofile=coverage.out ./tests/integration/...
      - uses: codecov/codecov-action@v2
        with:
          files: ./coverage.out
```

## 常见问题

### Q: 测试运行很慢怎么办？

A: 可以使用 `-timeout` 参数增加超时时间，或者使用 `-run` 参数只运行特定的测试。

### Q: 如何并行运行测试？

A: 使用 `-parallel` 参数：
```bash
go test -parallel 4 ./tests/integration/...
```

### Q: 如何查看测试的详细输出？

A: 使用 `-v` 参数显示详细输出：
```bash
go test -v ./tests/integration/...
```

### Q: 如何检测竞态条件？

A: 使用 `-race` 参数：
```bash
go test -race ./tests/integration/...
```

### Q: 如何生成覆盖率报告？

A: 使用 `-coverprofile` 参数：
```bash
go test -coverprofile=coverage.out ./tests/integration/...
go tool cover -html=coverage.out
```

## 相关文档

- [集成测试规范](./specs/phase-1/integration-tests.md)
- [集成测试总结](./INTEGRATION_TESTS_SUMMARY.md)
- [集成测试报告](./INTEGRATION_TESTS_REPORT.md)
- [集成测试完成总结](./INTEGRATION_TESTS_COMPLETION.md)

## 更多帮助

```bash
# 查看 go test 的完整帮助
go help test

# 查看 go test 的标志说明
go test -h
```
