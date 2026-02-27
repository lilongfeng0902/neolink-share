# 错误码改进提案

> **文档版本**: 1.0
> **创建日期**: 2026-02-26
> **状态**: 提案中

## 📊 当前状态分析

### ✅ 优势

1. **集中管理**: 错误码在 `types/error.go` 中集中定义
2. **统一封装**: `NewAPIError` 结构体提供统一的错误包装
3. **多格式支持**: 支持多种 AI 服务的错误格式转换（OpenAI、Claude、Gemini 等）
4. **敏感数据屏蔽**: 自动屏蔽 URL、IP 等敏感信息
5. **灵活的选项模式**: 支持跳过重试、隐藏错误消息等
6. **清晰的分类**: 按功能模块分类错误（channel、quota、client、response）
7. **请求追踪**: 支持请求 ID 追踪用于调试

### ❌ 发现的问题

#### 1. 缺少数字错误码系统
```go
// 当前: 只有字符串错误码，不利于快速识别和日志搜索
ErrorCode = "invalid_request"

// 应该有: 数字代码用于高效处理
ErrorCode = 1001  // 带字符串表示 "invalid_request"
```

#### 2. HTTP 状态码映射不一致
HTTP 状态码分散在整个代码库中：
```go
// controller/relay.go
newAPIError = types.NewErrorWithStatusCode(err, types.ErrorCodeInvalidRequest, http.StatusRequestEntityTooLarge)
newAPIError = types.NewError(err, types.ErrorCodeInvalidRequest)
```

错误码和 HTTP 状态码之间没有统一的映射规则。

#### 3. 错误码命名不一致
```go
// 有些有前缀
ErrorCodeChannelNoAvailableKey ErrorCode = "channel:no_available_key"
// 有些没有
ErrorCodeInvalidRequest ErrorCode = "invalid_request"
```

#### 4. 缺少错误级别定义
没有区分 Critical/Error/Warning/Info 级别，难以进行监控和告警。

#### 5. 错误处理不一致
```go
// 有些地方返回原始错误
func SomeFunc() error {
    return errors.New("xxx")  // ❌ 应该返回 *NewAPIError
}
// 有些地方返回 *NewAPIError
func OtherFunc() *types.NewAPIError {
    return types.NewError(...)
}
```

#### 6. 缺少错误码文档
没有错误码的集中列表；开发者必须搜索代码才能找到可用的代码。

#### 7. 没有后端国际化
错误消息只有英文；i18n 完全依赖前端。

---

## 💡 改进提案

### 提案 1: 引入数字错误码系统

**文件**: `types/error_code.go`

```go
type ErrorCode int

const (
    // 通用错误 (1xxx)
    ErrorCodeInvalidRequest ErrorCode = 1001
    ErrorCodeSensitiveWordsDetected ErrorCode = 1002

    // 系统错误 (2xxx)
    ErrorCodeCountTokenFailed ErrorCode = 2001
    ErrorCodeModelPriceError ErrorCode = 2002
    ErrorCodeInvalidApiType ErrorCode = 2003
    ErrorCodeJsonMarshalFailed ErrorCode = 2004
    ErrorCodeDoRequestFailed ErrorCode = 2005
    ErrorCodeGetChannelFailed ErrorCode = 2006
    ErrorCodeGenRelayInfoFailed ErrorCode = 2007

    // 渠道错误 (3xxx)
    ErrorCodeChannelNoAvailableKey ErrorCode = 3001
    ErrorCodeChannelParamOverrideInvalid ErrorCode = 3002
    ErrorCodeChannelHeaderOverrideInvalid ErrorCode = 3003
    ErrorCodeChannelModelMappedError ErrorCode = 3004
    ErrorCodeChannelAwsClientError ErrorCode = 3005
    ErrorCodeChannelInvalidKey ErrorCode = 3006
    ErrorCodeChannelResponseTimeExceeded ErrorCode = 3007

    // 客户端错误 (4xxx)
    ErrorCodeReadRequestBodyFailed ErrorCode = 4001
    ErrorCodeConvertRequestFailed ErrorCode = 4002
    ErrorCodeAccessDenied ErrorCode = 4003
    ErrorCodeBadRequestBody ErrorCode = 4004

    // 上游错误 (5xxx)
    ErrorCodeReadResponseBodyFailed ErrorCode = 5001
    ErrorCodeBadResponseStatusCode ErrorCode = 5002
    ErrorCodeBadResponse ErrorCode = 5003
    ErrorCodeBadResponseBody ErrorCode = 5004
    ErrorCodeEmptyResponse ErrorCode = 5005
    ErrorCodeAwsInvokeError ErrorCode = 5006
    ErrorCodeModelNotFound ErrorCode = 5007
    ErrorCodePromptBlocked ErrorCode = 5008

    // 数据库错误 (6xxx)
    ErrorCodeQueryDataError ErrorCode = 6001
    ErrorCodeUpdateDataError ErrorCode = 6002

    // 配额错误 (7xxx)
    ErrorCodeInsufficientUserQuota ErrorCode = 7001
    ErrorCodePreConsumeTokenQuotaFailed ErrorCode = 7002
)

// String 返回错误码的字符串表示
func (c ErrorCode) String() string {
    // 映射实现
}

// HTTPStatusCode 返回对应的 HTTP 状态码
func (c ErrorCode) HTTPStatusCode() int {
    // 映射实现
}
```

**优势**:
- 快速数字比较用于错误处理逻辑
- 易于记录和搜索
- 通过数字范围清晰分类
- 数字常量的类型安全

---

### 提案 2: 添加错误级别定义

**文件**: `types/error_level.go`

```go
type ErrorLevel int

const (
    ErrorLevelInfo ErrorLevel = iota
    ErrorLevelWarning
    ErrorLevelError
    ErrorLevelCritical
)

func (l ErrorLevel) String() string {
    switch l {
    case ErrorLevelInfo:
        return "info"
    case ErrorLevelWarning:
        return "warning"
    case ErrorLevelError:
        return "error"
    case ErrorLevelCritical:
        return "critical"
    default:
        return "unknown"
    }
}

func (l ErrorLevel) Color() string {
    switch l {
    case ErrorLevelInfo:
        return "\033[36m" // 青色
    case ErrorLevelWarning:
        return "\033[33m" // 黄色
    case ErrorLevelError:
        return "\033[31m" // 红色
    case ErrorLevelCritical:
        return "\033[35m" // 洋红色
    default:
        return "\033[0m" // 重置
    }
}
```

**与 NewAPIError 集成**:
```go
type NewAPIError struct {
    // ... 现有字段
    Level ErrorLevel
}

// 使用
func (e *NewAPIError) Log() {
    level := e.Level
    logger.Logf(level, e.Error())
}
```

**优势**:
- 清晰的严重程度指示
- 易于与监控系统集成（Prometheus、Sentry）
- 支持按严重程度过滤告警

---

### 提案 3: 统一 HTTP 状态码映射

```go
// 在 error_code.go 中
func (c ErrorCode) HTTPStatusCode() int {
    switch c {
    // 客户端错误 (4xx)
    case ErrorCodeInvalidRequest,
         ErrorCodeReadRequestBodyFailed,
         ErrorCodeConvertRequestFailed,
         ErrorCodeBadRequestBody:
        return http.StatusBadRequest

    case ErrorCodeAccessDenied:
        return http.StatusUnauthorized

    case ErrorCodeInsufficientUserQuota:
        return http.StatusPaymentRequired

    case ErrorCodeModelNotFound:
        return http.StatusNotFound

    case ErrorCodeChannelResponseTimeExceeded:
        return http.StatusRequestTimeout

    // 服务器错误 (5xx)
    case ErrorCodeCountTokenFailed,
         ErrorCodeModelPriceError,
         ErrorCodeGetChannelFailed,
         ErrorCodeChannelNoAvailableKey,
         ErrorCodeBadResponseStatusCode,
         ErrorCodeBadResponse:
        return http.StatusInternalServerError

    default:
        return http.StatusInternalServerError
    }
}
```

**使用**:
```go
// 之前
newAPIError = types.NewErrorWithStatusCode(err, types.ErrorCodeInvalidRequest, http.StatusBadRequest)

// 之后: 状态码自动确定
newAPIError = types.NewError(err, ErrorCodeInvalidRequest)
statusCode := newAPIError.StatusCode  // 自动映射
```

---

### 提案 4: 标准化错误码命名约定

**按数字范围分类**:
- `1xxx` - 通用错误
- `2xxx` - 系统错误
- `3xxx` - 渠道错误
- `4xxx` - 客户端错误
- `5xxx` - 上游错误
- `6xxx` - 数据库错误
- `7xxx` - 配额错误
- `8xxx` - 认证错误（预留）
- `9xxx` - 其他（预留）

**命名规则**:
1. 常量名使用 PascalCase
2. 不需要前缀（改用数字范围）
3. 描述性名称指示错误原因
4. 有歧义时包含模块上下文

**示例**:
```go
// ✅ 好
ErrorCodeChannelNoAvailableKey ErrorCode = 3001
ErrorCodeInsufficientUserQuota ErrorCode = 7001

// ❌ 差（基于字符串带前缀）
ErrorCodeChannelNoAvailableKey ErrorCode = "channel:no_available_key"
```

---

### 提案 5: 添加错误码文档生成器

**文件**: `tools/generate_error_doc.go`

```go
// +build ignore

package main

import (
    "fmt"
    "os"
    "reflect"
    "strings"
)

type ErrorDoc struct {
    Code       int
    Name       string
    Message    string
    HTTPStatus int
    Level      string
}

func main() {
    // 扫描所有 ErrorCode 常量
    // 生成 ERROR_CODES.md 文档
    // 包含代码、名称、描述、HTTP 状态、级别
}

// 生成 markdown 表格
func generateMarkdownTable(errors []ErrorDoc) string {
    // Markdown 表格生成逻辑
}
```

**使用**:
```bash
go run tools/generate_error_doc.go > docs/ERROR_CODES.md
```

**输出**: `docs/ERROR_CODES.md` 包含完整的错误码参考

---

### 提案 6: 统一错误处理中间件

**文件**: `middleware/error_handler.go`

```go
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()

        // 处理从 handler 返回的错误
        if len(c.Errors) > 0 {
            err := c.Errors.Last().Err

            if newApiErr, ok := err.(*types.NewAPIError); ok {
                // 按级别记录错误
                logger.LogWithLevel(c, newApiErr.Level, newApiErr.Error())

                // 根据中继格式返回响应
                switch c.GetHeader("Relay-Format") {
                case "claude":
                    c.JSON(newApiErr.StatusCode, gin.H{
                        "type":  "error",
                        "error": newApiErr.ToClaudeError(),
                    })
                default:
                    c.JSON(newApiErr.StatusCode, gin.H{
                        "error": newApiErr.ToOpenAIError(),
                    })
                }
            } else {
                // 包装原始错误
                c.JSON(http.StatusInternalServerError, gin.H{
                    "error": gin.H{
                        "message": "内部服务器错误",
                        "type":    "internal_error",
                        "code":    2000, // 未知系统错误
                    },
                })
            }
        }
    }
}
```

**优势**:
- 集中错误处理逻辑
- 一致的错误响应格式
- 自动日志记录和监控

---

### 提案 7: 后端国际化支持

**文件**: `types/error_i18n.go`

```go
type ErrorMessage map[string]string // lang -> message

var errorMessages = map[ErrorCode]ErrorMessage{
    ErrorCodeInvalidRequest: {
        "en": "Invalid request",
        "zh": "请求无效",
        "ja": "無効なリクエスト",
        "fr": "Requête invalide",
    },
    ErrorCodeInsufficientUserQuota: {
        "en": "Insufficient user quota",
        "zh": "用户配额不足",
        "ja": "ユーザークォータが不足しています",
        "fr": "Quota utilisateur insuffisant",
    },
    // ... 更多错误消息
}

// Localize 根据 Accept-Language 请求头返回本地化错误消息
func (e *NewAPIError) Localize(lang string) string {
    if msgs, ok := errorMessages[e.errorCode]; ok {
        if msg, ok := msgs[lang]; ok {
            return msg
        }
        return msgs["en"] // 回退到英文
    }
    return e.Error()
}

// GetLanguageFromContext 从 gin.Context 提取语言
func GetLanguageFromContext(c *gin.Context) string {
    lang := c.GetHeader("Accept-Language")
    if lang == "" {
        return "en" // 默认
    }
    // 解析语言请求头（例如 "zh-CN" -> "zh"）
    parts := strings.Split(lang, "-")
    if len(parts) > 0 {
        return parts[0]
    }
    return "en"
}
```

**在 Controller 中使用**:
```go
func Relay(c *gin.Context) {
    // ...
    lang := types.GetLanguageFromContext(c)
    errorMsg := newAPIError.Localize(lang)
    // ...
}
```

**优势**:
- 所有语言的一致错误消息
- 易于添加新语言
- 后端可以在 API 响应中返回本地化错误

---

### 提案 8: 错误码注册系统

**文件**: `types/error_registry.go`

```go
type ErrorInfo struct {
    Code       ErrorCode
    Name       string
    Message    string
    HTTPStatus int
    Level      ErrorLevel
}

var errorRegistry = map[ErrorCode]ErrorInfo{}
var registryMutex sync.RWMutex

// RegisterError 注册带有元数据的错误码
func RegisterError(info ErrorInfo) {
    registryMutex.Lock()
    defer registryMutex.Unlock()
    errorRegistry[info.Code] = info
}

// GetErrorInfo 通过代码检索错误元数据
func GetErrorInfo(code ErrorCode) (ErrorInfo, bool) {
    registryMutex.RLock()
    defer registryMutex.RUnlock()
    info, ok := errorRegistry[code]
    return info, ok
}

// ListAllErrors 返回所有注册的错误码
func ListAllErrors() []ErrorInfo {
    registryMutex.RLock()
    defer registryMutex.RUnlock()

    errors := make([]ErrorInfo, 0, len(errorRegistry))
    for _, info := range errorRegistry {
        errors = append(errors, info)
    }
    return errors
}

// 初始化
func init() {
    RegisterError(ErrorInfo{
        Code:       ErrorCodeInvalidRequest,
        Name:       "invalid_request",
        Message:    "Invalid request parameters",
        HTTPStatus: http.StatusBadRequest,
        Level:      ErrorLevelWarning,
    })

    RegisterError(ErrorInfo{
        Code:       ErrorCodeInsufficientUserQuota,
        Name:       "insufficient_user_quota",
        Message:    "User quota is insufficient",
        HTTPStatus: http.StatusPaymentRequired,
        Level:      ErrorLevelError,
    })
    // ... 注册所有错误码
}
```

**优势**:
- 集中错误元数据管理
- 易于查询和文档化错误码
- 支持运行时错误码验证

---

## 📋 实施优先级

| 优先级 | 改进项 | 影响 | 工作量 | 依赖项 |
|----------|-------------|--------|--------|--------------|
| **P0** | 统一 HTTP 状态码映射 | 高 | 中 | 无 |
| **P0** | 添加错误级别定义 | 高 | 低 | 无 |
| **P1** | 引入数字错误码 | 中 | 高 | 无 |
| **P1** | 统一错误处理中间件 | 中 | 中 | 错误级别 |
| **P2** | 生成错误码文档 | 低 | 低 | 数字错误码 |
| **P2** | 后端国际化 | 低 | 中 | 无 |
| **P3** | 错误码注册系统 | 低 | 高 | 数字错误码 |

---

## 🔄 迁移策略

### 阶段 1: 基础 (P0)
1. ✅ 添加 `ErrorLevel` 定义
2. ✅ 实现 `HTTPStatusCode()` 方法
3. ✅ 更新 `NewAPIError` 以包含 `Level` 字段

### 阶段 2: 核心重构 (P1)
1. ✅ 迁移到数字错误码
2. ✅ 更新所有错误创建以使用新系统
3. ✅ 实现统一错误处理中间件

### 阶段 3: 增强 (P2)
1. ✅ 实现国际化
2. ✅ 生成错误码文档
3. ✅ 添加错误码注册

### 阶段 4: 高级功能 (P3)
1. ⏳ 带元数据的错误码注册
2. ⏳ 运行时错误验证
3. ⏳ 错误分析集成

---

## 📝 代码示例

### 当前系统之前

```go
// 创建错误
err := types.NewError(
    errors.New("渠道不可用"),
    types.ErrorCodeChannelNoAvailableKey,
)

// 使用错误
if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
        "error": err.ToOpenAIError(),
    })
}
```

### 提案系统之后

```go
// 创建自动状态码和级别的错误
err := types.NewError(
    errors.New("渠道不可用"),
    types.ErrorCodeChannelNoAvailableKey,
    types.ErrOptionWithLevel(types.ErrorLevelError),
)

// 使用错误（状态码自动映射）
if err != nil {
    // 状态码从错误码自动确定
    c.JSON(err.StatusCode, gin.H{
        "error": err.ToOpenAIError(),
    })
}

// 本地化错误消息
lang := types.GetLanguageFromContext(c)
errorMsg := err.Localize(lang)
```

---

## 🎯 预期成果

### 可维护性
- ✅ 清晰的数字错误码用于快速识别
- ✅ 集中错误元数据管理
- ✅ 自动生成文档

### 可靠性
- ✅ 一致的 HTTP 状态码映射
- ✅ 基于错误级别的日志记录和告警
- ✅ 减少错误处理中的人为错误

### 开发者体验
- ✅ 易于查找和使用错误码
- ✅ 清晰的错误分类
- ✅ 更好的 IDE 自动完成
- ✅ 全面的错误码文档

### 用户体验
- ✅ 本地化错误消息
- ✅ 所有 API 端点的一致错误格式
- ✅ 更好的错误消息用于调试

---

## 🔒 向后兼容性

为了在迁移期间保持向后兼容性：

1. **保留字符串错误码类型** 作为别名
   ```go
   type ErrorCodeString string
   // 将旧字符串代码映射到新数字代码
   ```

2. **过渡期的双重支持**
   ```go
   func (c ErrorCode) LegacyString() string {
       // 将数字转换为旧字符串格式
   }
   ```

3. **旧错误创建方法的弃用警告**
   ```go
   // 已弃用: 使用带有数字 ErrorCode 的 NewError
   func NewErrorLegacy(err error, code string) *NewAPIError {
       // 带弃用警告的实现
   }
   ```

4. **渐进式迁移**:
   - 阶段 1: 在旧系统旁边添加新系统
   - 阶段 2: 迁移所有新代码以使用新系统
   - 阶段 3: 增量迁移现有代码
   - 阶段 4: 完全迁移后删除旧系统

---

## 📚 参考

- 当前实现: `types/error.go`
- 错误处理最佳实践: [Go Error Handling](https://go.dev/doc/error-handle)
- HTTP 状态码: [MDN HTTP Status](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status)
- i18n 模式: [Go i18n](https://github.com/nicksnyder/go-i18n)

---

## 📞 联系与讨论

有关此提案的问题或建议：
1. 在仓库中创建 issue
2. 在团队沟通渠道中开始讨论
3. 提交 PR 以改进此提案

---

**最后更新**: 2026-02-26
**下次审查**: 阶段 1 完成后
