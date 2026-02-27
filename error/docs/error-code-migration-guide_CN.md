# 错误码迁移指南

> **版本**: 1.0
> **最后更新**: 2026-02-26

## 概述

本指南帮助您从旧的基于字符串的错误码迁移到新的数字错误码系统。

## 有什么变化？

### 之前（旧系统）

```go
type ErrorCode string

const (
    ErrorCodeInvalidRequest ErrorCode = "invalid_request"
    ErrorCodeChannelNoAvailableKey ErrorCode = "channel:no_available_key"
    // ...
)
```

### 之后（新系统）

```go
type ErrorCode int

const (
    ErrorCodeInvalidRequest ErrorCode = 1001
    ErrorCodeChannelNoAvailableKey ErrorCode = 3001
    // ...
)

// 辅助方法
func (c ErrorCode) String() string
func (c ErrorCode) HTTPStatusCode() int
func (c ErrorCode) DefaultLevel() ErrorLevel
```

## 新系统的优势

1. **类型安全**: 数字代码防止拼写错误并提供编译时检查
2. **性能**: 整数比较比字符串比较更快
3. **自动映射**: HTTP 状态码自动确定
4. **错误级别**: 内置严重级别用于监控
5. **国际化**: 易于添加本地化错误消息

## 破坏性变更

### 1. ErrorCode 类型变更

```go
// 旧
var code types.ErrorCode = "invalid_request"

// 新
var code types.ErrorCode = types.ErrorCodeInvalidRequest  // 1001
```

### 2. 字符串转换

```go
// 旧
fmt.Sprintf("Error: %s", errorCode)  // "invalid_request"

// 新
fmt.Sprintf("Error: %s", errorCode.String())  // "invalid_request"
fmt.Sprintf("Error: %d", errorCode)  // 1001
```

### 3. 旧常量重命名

旧的字符串常量重命名为 `LegacyErrorCodeString`：

```go
// 仍可用于向后兼容
types.LegacyErrorCodeInvalidRequest         // "invalid_request"
types.LegacyErrorCodeChannelNoAvailableKey  // "channel:no_available_key"
```

## 迁移步骤

### 步骤 1: 更新错误创建

**之前:**
```go
newApiErr := types.NewError(err, types.ErrorCodeInvalidRequest)
newApiErr.StatusCode = http.StatusBadRequest  // 手动设置
```

**之后:**
```go
// HTTP 状态码自动映射
newApiErr := types.NewError(err, types.ErrorCodeInvalidRequest)
// newApiErr.StatusCode 自动设置为 400
```

### 步骤 2: 更新错误检查

**之前:**
```go
if err.errorCode == "invalid_request" {
    // 处理无效请求
}
```

**之后:**
```go
if err.errorCode == types.ErrorCodeInvalidRequest {
    // 处理无效请求
}

// 或使用辅助方法
if types.IsChannelError(err) {
    // 处理渠道错误
}
```

### 步骤 3: 更新错误级别使用

**新功能:**
```go
// 访问错误级别
level := err.Level  // ErrorLevelWarning, ErrorLevelError 等

// 按级别记录日志
logger.Logf(err.Level, "Error occurred: %s", err.Error())

// 覆盖默认级别
err := types.NewError(
    errors.New("custom"),
    types.ErrorCodeInvalidRequest,
    types.ErrOptionWithLevel(types.ErrorLevelCritical),
)
```

### 步骤 4: 使用本地化消息

**新功能:**
```go
// 获取本地化错误消息
lang := types.GetLanguageFromContext(c.GetHeader("Accept-Language"))
message := err.Localize(lang)

// 示例:
// "en" -> "Invalid request parameters"
// "zh" -> "请求参数无效"
// "ja" -> "無效なリクエストパラメータ"
```

## 代码示例

### Controller 错误处理

**之前:**
```go
func Relay(c *gin.Context) {
    if err != nil {
        newApiErr := types.NewError(err, types.ErrorCodeInvalidRequest)
        c.JSON(http.StatusBadRequest, gin.H{
            "error": newApiErr.ToOpenAIError(),
        })
        return
    }
}
```

**之后:**
```go
func Relay(c *gin.Context) {
    if err != nil {
        newApiErr := types.NewError(err, types.ErrorCodeInvalidRequest)
        // StatusCode 从错误码自动映射
        c.JSON(newApiErr.StatusCode, gin.H{
            "error": newApiErr.ToOpenAIError(),
        })
        return
    }
}
```

### 错误比较

**之前:**
```go
if newApiErr.errorCode == "channel:no_available_key" {
    // 处理错误
}
```

**之后:**
```go
if newApiErr.errorCode == types.ErrorCodeChannelNoAvailableKey {
    // 处理错误
}

// 或使用辅助函数
if types.IsChannelError(newApiErr) {
    // 处理渠道错误 (3xxx 范围)
}

if types.IsSkipRetryError(newApiErr) {
    // 跳过重试逻辑
}
```

### 自定义错误选项

**新功能:**
```go
err := types.NewError(
    errors.New("出错了"),
    types.ErrorCodeInvalidRequest,
    types.ErrOptionWithSkipRetry(),           // 跳过重试
    types.ErrOptionWithNoRecordErrorLog(),    // 不记录日志
    types.ErrOptionWithLevel(types.ErrorLevelWarning),  // 自定义级别
    types.ErrOptionWithHideErrMsg("Hidden"),  // 隐藏消息
)
```

## 错误码参考

### 数字范围

| 范围 | 类别 | 示例 |
|-------|----------|---------|
| 1xxx | 通用错误 | 1001: 无效请求 |
| 2xxx | 系统错误 | 2001: Token 计数失败 |
| 3xxx | 渠道错误 | 3001: 渠道无可用密钥 |
| 4xxx | 客户端错误 | 4001: 读取请求体失败 |
| 5xxx | 上游错误 | 5001: 读取响应体失败 |
| 6xxx | 数据库错误 | 6001: 查询数据错误 |
| 7xxx | 配额错误 | 7001: 用户配额不足 |

### HTTP 状态映射

所有错误码自动映射到适当的 HTTP 状态码：

- **400**: Bad Request（验证错误）
- **401**: Unauthorized（认证错误）
- **402**: Payment Required（配额错误）
- **403**: Forbidden（访问被拒绝）
- **404**: Not Found（模型未找到）
- **429**: Too Many Requests（速率限制）
- **500**: Internal Server Error（系统错误）
- **502**: Bad Gateway（上游错误）
- **503**: Service Unavailable（无可用渠道）
- **504**: Gateway Timeout（超时）

### 错误级别

| 级别 | 描述 | 用例 |
|-------|-------------|----------|
| **info** | 信息性 | 正常操作 |
| **warning** | 警告 | 非关键问题 |
| **error** | 错误 | 错误事件 |
| **critical** | 严重 | 系统故障 |

## 向后兼容性

### 旧版支持

旧的基于字符串的错误码仍可作为 `LegacyErrorCodeString` 使用：

```go
// 旧常量（已弃用但仍可用）
types.LegacyErrorCodeInvalidRequest
types.LegacyErrorCodeChannelNoAvailableKey
// 等等
```

### 转换函数

在旧系统和新系统之间转换：

```go
// 旧到新
newCode := types.ErrorCodeFromString("invalid_request")  // 1001

// 新到旧
oldCode := types.ErrorCodeInvalidRequest.String()  // "invalid_request"
```

## 测试迁移

### 1. 单元测试

更新测试中的错误断言：

```go
// 之前
assert.Equal(t, "invalid_request", err.GetErrorCode())

// 之后
assert.Equal(t, types.ErrorCodeInvalidRequest, err.GetErrorCode())

// 或
assert.Equal(t, 1001, int(err.GetErrorCode()))
```

### 2. 集成测试

测试 HTTP 状态码：

```go
// 之前
assert.Equal(t, http.StatusBadRequest, w.Code)

// 之后（仍然有效，状态码自动映射）
assert.Equal(t, http.StatusBadRequest, w.Code)
```

### 3. 错误消息

测试错误消息：

```go
// 字符串表示如之前一样工作
assert.Contains(t, err.Error(), "invalid_request")
```

## 常见问题和解决方案

### 问题 1: 类型不匹配

**错误:**
```go
cannot use "invalid_request" (type untyped string) as type ErrorCode
```

**解决方案:**
```go
// 使用常量而不是字符串
types.ErrorCodeInvalidRequest  // ✅

// 或从字符串转换
types.ErrorCodeFromString("invalid_request")  // ✅
```

### 问题 2: HTTP 状态码缺失

**错误:**
```go
NewApiErr.StatusCode undefined
```

**解决方案:**
```go
// StatusCode 现在从 ErrorCode 自动映射
// 无需手动设置
err := types.NewError(e, types.ErrorCodeInvalidRequest)
statusCode := err.StatusCode  // 自动设置为 400
```

### 问题 3: 导入错误

**错误:**
```go
undefined: ErrorCode
```

**解决方案:**
```go
// 确保导入新的 types 包
import "github.com/QuantumNous/new-api/types"

// 使用 ErrorCode 类型
types.ErrorCodeInvalidRequest
```

## 回滚计划

如果需要回滚到旧系统：

1. 将 `types/error.go` 恢复为使用 `string` 类型的 `ErrorCode`
2. 将 `types/error_code.go` 恢复到旧版本
3. 删除 `types/error_level.go`
4. 删除 `types/error_i18n.go`

但是，我们建议保留新系统以获得其优势。

## 获取帮助

- **文档**: 查看 [ERROR_CODES.md](./ERROR_CODES.md) 获取完整的错误码参考
- **提案**: 查看 [error-code-improvements.md](./error-code-improvements_CN.md) 了解设计细节
- **问题**: 在 GitHub Issues 上报告问题

## 检查清单

- [ ] 更新错误创建代码
- [ ] 更新错误比较代码
- [ ] 更新单元测试
- [ ] 更新集成测试
- [ ] 更新错误日志
- [ ] 添加错误级别处理
- [ ] 添加本地化消息（如需要）
- [ ] 更新文档
- [ ] 删除旧代码（迁移完成后）

---

**有问题？** 查看 [ERROR_CODES.md](./ERROR_CODES.md) 参考或 [error-code-improvements.md](./error-code-improvements_CN.md) 提案了解更多详情。
