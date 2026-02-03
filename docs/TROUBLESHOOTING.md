# swag-gen 故障排除指南

## 目录
1. [常见错误](#常见错误)
2. [解析问题](#解析问题)
3. [生成问题](#生成问题)
4. [性能问题](#性能问题)
5. [调试技巧](#调试技巧)

## 常见错误

### 错误 1: "path not found"

**错误信息**:
```
Error: path not found: ./api
```

**原因**:
- 指定的路径不存在
- 路径拼写错误
- 相对路径不正确

**解决方案**:
```bash
# 检查路径是否存在
ls -la ./api

# 使用绝对路径
swag-gen init -p /absolute/path/to/api -o ./docs

# 检查当前工作目录
pwd
```

### 错误 2: "invalid format"

**错误信息**:
```
Error: invalid format: xyz
```

**原因**:
- 指定的格式不支持
- 格式名称拼写错误

**解决方案**:
```bash
# 使用支持的格式
swag-gen init -p ./api -o ./docs -f json
swag-gen init -p ./api -o ./docs -f yaml
swag-gen init -p ./api -o ./docs -f yml
```

### 错误 3: "permission denied"

**错误信息**:
```
Error: permission denied: ./docs
```

**原因**:
- 没有写入权限
- 目录被锁定

**解决方案**:
```bash
# 检查权限
ls -la ./docs

# 修改权限
chmod 755 ./docs

# 使用有权限的目录
swag-gen init -p ./api -o /tmp/docs
```

### 错误 4: "invalid title"

**错误信息**:
```
Error: invalid title: title cannot be empty
```

**原因**:
- 标题为空
- 标题包含无效字符

**解决方案**:
```bash
# 提供有效的标题
swag-gen init -p ./api -o ./docs -t "My API"

# 使用引号处理特殊字符
swag-gen init -p ./api -o ./docs -t "My API v1.0"
```

## 解析问题

### 问题 1: 找不到 API 端点

**症状**:
- 生成的文档中没有 API 端点
- 端点数量少于预期

**原因**:
- 注释格式不正确
- 文件不在扫描路径中
- 注释标签缺失

**调试步骤**:

1. 检查注释格式:
```go
// 正确的格式
// @Router /api/users [GET]
// @Summary 获取用户

// 错误的格式
//@Router /api/users [GET]  // 缺少空格
// @router /api/users [GET]  // 小写 router
```

2. 检查文件是否被扫描:
```bash
# 查看项目中的所有 Go 文件
find ./api -name "*.go" -type f
```

3. 检查注释标签:
```go
// 必需的标签
// @Router /api/users [GET]
// @Summary 获取用户

// 可选的标签
// @Description 详细描述
// @Tags User
// @Success 200 {array} User
```

### 问题 2: 参数解析错误

**症状**:
- 参数没有被正确识别
- 参数类型不正确

**原因**:
- 参数格式不正确
- 参数类型不支持

**调试步骤**:

1. 检查参数格式:
```go
// 正确的格式
// @Param page query int false "页码"
// @Param id path int true "用户 ID"
// @Param body body CreateUserRequest true "用户信息"

// 错误的格式
// @Param page query int  // 缺少 required 和 description
// @Param id path int true  // 缺少 description
```

2. 检查参数位置:
```go
// 支持的位置
// @Param name query string false "查询参数"
// @Param id path int true "路径参数"
// @Param Authorization header string false "请求头"
// @Param body body RequestBody true "请求体"
```

3. 检查参数类型:
```go
// 支持的类型
// @Param count query int false "整数"
// @Param price query float false "浮点数"
// @Param active query bool false "布尔值"
// @Param name query string false "字符串"
// @Param items query array false "数组"
```

### 问题 3: 响应类型不正确

**症状**:
- 响应类型显示为 object 而不是数组
- 响应类型不匹配

**原因**:
- 响应格式不正确
- 类型定义不完整

**调试步骤**:

1. 检查响应格式:
```go
// 正确的格式
// @Success 200 {array} User "成功"
// @Success 200 {object} User "成功"
// @Failure 400 {object} ErrorResponse "错误"

// 错误的格式
// @Success 200 User "成功"  // 缺少类型标记
// @Success 200 {User} "成功"  // 错误的类型标记
```

2. 检查类型定义:
```go
// 确保类型已定义
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// 使用已定义的类型
// @Success 200 {object} User
```

## 生成问题

### 问题 1: 文件写入失败

**症状**:
- 输出文件未生成
- 错误信息提示写入失败

**原因**:
- 输出目录不存在
- 没有写入权限
- 磁盘空间不足

**解决方案**:

```bash
# 创建输出目录
mkdir -p ./docs

# 检查权限
ls -la ./docs

# 检查磁盘空间
df -h

# 使用有权限的目录
swag-gen init -p ./api -o /tmp/docs
```

### 问题 2: 生成的文档不完整

**症状**:
- 某些端点缺失
- 某些字段缺失

**原因**:
- 解析过程中出错
- 某些注释格式不正确

**调试步骤**:

1. 检查日志输出:
```bash
# 运行命令并查看详细输出
swag-gen init -p ./api -o ./docs 2>&1 | tee output.log
```

2. 检查每个文件:
```bash
# 逐个检查 API 文件
for file in ./api/*.go; do
  echo "Checking $file"
  grep -n "@Router" "$file"
done
```

3. 验证生成的文档:
```bash
# 检查 JSON 格式
cat ./docs/swagger.json | jq .

# 检查 YAML 格式
cat ./docs/swagger.yaml
```

### 问题 3: Schema 定义不正确

**症状**:
- 生成的 Schema 缺少字段
- 字段类型不正确

**原因**:
- 结构体定义不完整
- JSON 标签缺失

**解决方案**:

```go
// 正确的结构体定义
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 错误的结构体定义
type User struct {
	ID   int     // 缺少 JSON 标签
	Name string  // 缺少 JSON 标签
}
```

## 性能问题

### 问题 1: 解析速度慢

**症状**:
- 命令执行时间长
- CPU 使用率高

**原因**:
- 项目文件过多
- 文件过大
- 系统资源不足

**解决方案**:

```bash
# 指定特定的 API 目录
swag-gen init -p ./api -o ./docs

# 排除不必要的目录
# (在代码中配置)

# 检查系统资源
top
free -h
```

### 问题 2: 内存使用过高

**症状**:
- 内存占用不断增加
- 系统变慢

**原因**:
- 项目过大
- 内存泄漏

**解决方案**:

```bash
# 分批处理
# 将大型项目分成多个小项目

# 监控内存使用
watch -n 1 'free -h'

# 使用内存分析工具
go tool pprof
```

### 问题 3: 磁盘空间不足

**症状**:
- 生成失败
- 错误信息提示磁盘空间不足

**解决方案**:

```bash
# 检查磁盘空间
df -h

# 清理不必要的文件
rm -rf ./docs/*

# 使用其他磁盘
swag-gen init -p ./api -o /mnt/other/docs
```

## 调试技巧

### 技巧 1: 启用详细日志

```bash
# 设置日志级别为 DEBUG
export LOG_LEVEL=debug
swag-gen init -p ./api -o ./docs
```

### 技巧 2: 保存输出日志

```bash
# 将输出保存到文件
swag-gen init -p ./api -o ./docs > output.log 2>&1

# 查看日志
cat output.log
```

### 技巧 3: 逐步调试

```bash
# 1. 检查路径
ls -la ./api

# 2. 检查文件
find ./api -name "*.go" -type f

# 3. 检查注释
grep -r "@Router" ./api

# 4. 运行命令
swag-gen init -p ./api -o ./docs

# 5. 检查输出
ls -la ./docs
cat ./docs/swagger.json | jq .
```

### 技巧 4: 使用临时目录

```bash
# 使用临时目录测试
swag-gen init -p ./api -o /tmp/test-docs

# 检查结果
ls -la /tmp/test-docs
cat /tmp/test-docs/swagger.json | jq .
```

### 技巧 5: 验证 JSON/YAML

```bash
# 验证 JSON 格式
cat ./docs/swagger.json | jq . > /dev/null && echo "JSON valid" || echo "JSON invalid"

# 验证 YAML 格式
cat ./docs/swagger.yaml | python3 -c "import yaml; yaml.safe_load(__import__('sys').stdin)" && echo "YAML valid" || echo "YAML invalid"
```

## 获取帮助

### 查看帮助信息
```bash
swag-gen init --help
```

### 查看版本
```bash
swag-gen version
```

### 提交 Issue
如果问题无法解决，请提交 Issue 并包含：
1. 错误信息
2. 命令行参数
3. 项目结构
4. 相关的 Go 代码

---

**版本**: 1.0.0  
**最后更新**: 2026 年 2 月 3 日
