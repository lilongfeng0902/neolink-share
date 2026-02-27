# 错误码分析与改进建议

> 文档创建时间：2026-02-27
> 分析范围：New API 项目错误码体系

## 📋 目录

- [1. 错误码体系现状](#1-错误码体系现状)
- [2. 敏感信息保护机制](#2-敏感信息保护机制)
- [3. 不应暴露给用户的错误信息](#3-不应暴露给用户的错误信息)
- [4. 改进建议](#4-改进建议)
- [5. 立即行动清单](#5-立即行动清单)

---

## 1. 错误码体系现状

### 1.1 错误码定义文件位置

| 文件路径 | 行号 | 用途 |
|---------|------|------|
| `types/error.go` | 38-87 | 核心错误码常量定义（32个错误码） |
| `types/error.go` | 89-183 | NewAPIError 结构体及方法 |
| `types/channel_error.go` | 3-10 | ChannelError 渠道错误结构体 |
| `dto/error.go` | 17-38 | API 响应错误 DTO 结构 |
| `dto/task.go` | - | 任务错误 DTO |
| `dto/claude.go` | - | Claude 错误 DTO |
| `service/error.go` | 59-172 | 错误包装和处理函数 |

### 1.2 错误码分类

项目中共定义了 **32 个错误码**，分为以下几类：

#### 通用错误码
- `invalid_request` - 无效请求
- `sensitive_words_detected` - 检测到敏感词

#### New API 内部错误
- `count_token_failed` - 令牌计数失败
- `model_price_error` - 模型价格错误
- `invalid_api_type` - 无效 API 类型
- `json_marshal_failed` - JSON 序列化失败
- `do_request_failed` - 请求失败
- `get_channel_failed` - 获取渠道失败
- `gen_relay_info_failed` - 生成中继信息失败

#### 渠道错误（`channel:` 前缀）
- `channel:no_available_key` - 渠道无可用密钥
- `channel:param_override_invalid` - 参数覆盖无效
- `channel:header_override_invalid` - 头部覆盖无效
- `channel:model_mapped_error` - 模型映射错误
- `channel:aws_client_error` - AWS 客户端错误
- `channel:invalid_key` - 渠道密钥无效
- `channel:response_time_exceeded` - 响应时间超限

#### 客户端请求错误
- `read_request_body_failed` - 读取请求体失败
- `convert_request_failed` - 转换请求失败
- `access_denied` - 访问被拒绝

#### 请求错误
- `bad_request_body` - 错误的请求体

#### 响应错误
- `read_response_body_failed` - 读取响应体失败
- `bad_response_status_code` - 错误的响应状态码
- `bad_response` - 错误的响应
- `bad_response_body` - 错误的响应体
- `empty_response` - 空响应
- `aws_invoke_error` - AWS 调用错误
- `model_not_found` - 模型未找到
- `prompt_blocked` - 提示词被阻止

#### SQL 错误
- `query_data_error` - 查询数据错误
- `update_data_error` - 更新数据错误

#### 配额错误
- `insufficient_user_quota` - 用户配额不足
- `pre_consume_token_quota_failed` - 预消费令牌配额失败

### 1.3 错误类型定义

```go
type ErrorType string

const (
    ErrorTypeNewAPIError     ErrorType = "new_api_error"
    ErrorTypeOpenAIError     ErrorType = "openai_error"
    ErrorTypeClaudeError     ErrorType = "claude_error"
    ErrorTypeMidjourneyError ErrorType = "midjourney_error"
    ErrorTypeGeminiError     ErrorType = "gemini_error"
    ErrorTypeRerankError     ErrorType = "rerank_error"
    ErrorTypeUpstreamError   ErrorType = "upstream_error"
)
```

---

## 2. 敏感信息保护机制

### 2.1 现有保护措施

项目已在 `common/str.go:176-239` 实现 `MaskSensitiveInfo()` 函数，可自动过滤：

- ✅ **URL 地址** - 域名、路径、查询参数
- ✅ **IP 地址** - IPv4 地址
- ✅ **纯域名** - 无协议前缀的域名
- ✅ **邮箱地址** - 用户邮箱前缀脱敏

**示例效果：**
```go
// 原始信息
"https://api.openai.com/v1/chat/completions?key=sk-abc123"

// 脱敏后
"https://***.com/***/***/?key=***"
```

### 2.2 应用位置

敏感信息过滤已在以下位置应用：

| 位置 | 文件路径 | 行号 | 说明 |
|------|----------|------|------|
| `MaskSensitiveError()` | `types/error.go` | 144 | 错误脱敏方法 |
| `ToOpenAIError()` | `types/error.go` | 176-177 | OpenAI 格式转换时脱敏 |
| `ToClaudeError()` | `types/error.go` | 204-205 | Claude 格式转换时脱敏 |
| `TaskErrorWrapper()` | `service/error.go` | 160 | 任务错误包装时脱敏 |

### 2.3 例外处理

```go
// types/error.go:141-143
if e.errorCode == ErrorCodeCountTokenFailed {
    return errStr  // Token 计数错误不脱敏
}
```

---

## 3. 不应暴露给用户的错误信息

### 🔴 高风险（包含敏感配置/内部信息）

| 错误码 | 风险等级 | 当前问题 | 潜在影响 | 位置 |
|--------|----------|----------|----------|------|
| `channel:invalid_key` | **严重** | 可能暴露密钥格式、前缀、长度等信息 | 上游 API 密钥泄露风险 | `types/error.go:59` |
| `ChannelError.UsingKey` | **严重** | 直接显示使用的密钥明文 | 密钥直接泄露 | `types/channel_error.go:9` |
| `do_request_failed` | **高** | 错误信息可能包含上游 URL、IP、端口 | 暴露上游服务地址，便于攻击者探测 | `types/error.go:49` |
| `get_channel_failed` | **高** | 可能暴露数据库查询细节、表结构 | SQL 注入风险线索 | `types/error.go:50` |
| `aws_invoke_error` | **高** | 可能包含 AWS 配置、ARN、区域信息 | 暴露云基础设施配置 | `types/error.go:76` |
| `channel:aws_client_error` | **高** | AWS 错误详情，包含配置信息 | 暴露 AWS 配置 | `types/error.go:58` |
| `query_data_error` | **高** | 包含 SQL 错误信息、表结构 | SQL 注入攻击线索 | `types/error.go:81` |
| `update_data_error` | **高** | 包含 SQL 错误信息、表结构 | SQL 注入攻击线索 | `types/error.go:82` |

#### 详细说明

**1. ChannelError.UsingKey 字段（最严重）**

```go
// types/channel_error.go:9
type ChannelError struct {
    ChannelId   int    `json:"channel_id"`
    ChannelType int    `json:"channel_type"`
    ChannelName string `json:"channel_name"`
    IsMultiKey  bool   `json:"is_multi_key"`
    AutoBan     bool   `json:"auto_ban"`
    UsingKey    string `json:"using_key"`  // ⚠️ 直接暴露密钥
}
```

**风险：** 该字段会直接将使用的 API 密钥返回给用户，造成严重的安全漏洞。

**2. 网络请求错误**

```go
// service/error.go:157-161
if strings.Contains(lowerText, "post") || strings.Contains(lowerText, "dial") ||
   strings.Contains(lowerText, "http") {
    common.SysLog(fmt.Sprintf("error: %s", text))
    text = common.MaskSensitiveInfo(text)  // ✅ 已脱敏
}
```

虽然已有脱敏处理，但仍需确保所有错误路径都经过脱敏。

### 🟡 中风险（技术实现细节）

| 错误码 | 风险等级 | 当前问题 | 用户体验 | 位置 |
|--------|----------|----------|----------|------|
| `json_marshal_failed` | **中** | 暴露 JSON 序列化技术细节 | 用户无法理解，暴露技术栈 | `types/error.go:48` |
| `count_token_failed` | **中** | Token 计数内部错误机制 | 暴露内部计费逻辑 | `types/error.go:45` |
| `gen_relay_info_failed` | **中** | 中继信息生成失败 | 暴露架构设计信息 | `types/error.go:51` |
| `convert_request_failed` | **中** | 协议转换细节 | 用户无需知晓技术实现 | `types/error.go:64` |
| `read_response_body_failed` | **中** | 网络读取底层错误 | 技术细节过载 | `types/error.go:71` |
| `pre_consume_token_quota_failed` | **中** | 配额预消费业务逻辑 | 暴露内部业务流程 | `types/error.go:86` |
| `channel:no_available_key` | **中** | 渠道密钥管理细节 | 可能暴露多渠道架构 | `types/error.go:54` |
| `channel:model_mapped_error` | **中** | 模型映射逻辑 | 暴露模型映射策略 | `types/error.go:57` |

#### 详细说明

**1. 技术细节暴露问题**

```go
// 示例错误消息
"json_marshal_failed: invalid character 'a' looking for beginning of value"
```

这类错误消息对普通用户没有任何帮助，反而暴露了技术实现细节。

**2. 内部业务逻辑暴露**

```go
// 示例
"pre_consume_token_quota_failed: insufficient quota for user 12345"
```

暴露了用户 ID 和配额预消费的业务逻辑。

### 🟢 合理暴露（用户可理解）

| 错误码 | 评估 | 说明 |
|--------|------|------|
| `invalid_request` | ✅ 合理 | 通用请求错误，用户可理解 |
| `insufficient_user_quota` | ✅ 合理 | 配额不足，用户明确知道需要充值 |
| `model_not_found` | ✅ 合理 | 模型不存在，用户可修正请求 |
| `access_denied` | ✅ 合理 | 权限不足，用户明确知道问题 |
| `sensitive_words_detected` | ✅ 合理 | 敏感词拦截，用户了解原因 |
| `prompt_blocked` | ✅ 合理 | 内容被阻止，用户可调整内容 |
| `empty_response` | ✅ 合理 | 空响应，用户知晓服务异常 |
| `bad_request_body` | ✅ 合理 | 请求体格式错误 |

---

## 4. 改进建议

### 4.1 修复 ChannelError 密钥泄露（高优先级）

**问题：** `types/channel_error.go:9` 的 `UsingKey` 字段会直接暴露 API 密钥

**解决方案：**

#### 方案 1：完全移除（推荐）

```go
type ChannelError struct {
    ChannelId   int    `json:"channel_id"`
    ChannelType int    `json:"channel_type"`
    ChannelName string `json:"channel_name"`
    IsMultiKey  bool   `json:"is_multi_key"`
    AutoBan     bool   `json:"auto_ban"`
    UsingKey    string `json:"-"`  // 不对外序列化
}
```

#### 方案 2：仅脱敏显示

```go
type ChannelError struct {
    ChannelId   int    `json:"channel_id"`
    ChannelType int    `json:"channel_type"`
    ChannelName string `json:"channel_name"`
    IsMultiKey  bool   `json:"is_multi_key"`
    AutoBan     bool   `json:"auto_ban"`
    UsingKey    string `json:"using_key,omitempty"`  // 序列化时手动脱敏
}

// 在序列化前
if len(error.UsingKey) > 10 {
    error.UsingKey = error.UsingKey[:7] + "***...***"
}
```

### 4.2 错误信息分级展示（中优先级）

为用户错误添加 **内部/外部** 标识和友好消息：

```go
// types/error.go
type ErrorDetail struct {
    Code         ErrorCode
    UserMessage  string  // 给用户的友好提示
    InternalOnly bool    // 是否仅内部可见
    StatusCode   int     // HTTP 状态码
}

// 错误消息配置
var errorMessages = map[ErrorCode]ErrorDetail{
    // 用户可见错误
    ErrorCodeInvalidRequest: {
        UserMessage:  "请求参数有误，请检查后重试",
        InternalOnly: false,
        StatusCode:   400,
    },
    ErrorCodeInsufficientUserQuota: {
        UserMessage:  "您的配额不足，请充值后继续使用",
        InternalOnly: false,
        StatusCode:   402,
    },
    ErrorCodeModelNotFound: {
        UserMessage:  "请求的模型不存在或已下线",
        InternalOnly: false,
        StatusCode:   404,
    },

    // 仅内部可见（对外隐藏详情）
    ErrorCodeDoRequestFailed: {
        UserMessage:  "上游服务暂时不可用，请稍后重试",
        InternalOnly: true,
        StatusCode:   503,
    },
    ErrorCodeQueryDataError: {
        UserMessage:  "服务暂时异常，请联系管理员",
        InternalOnly: true,
        StatusCode:   500,
    },
    ErrorCodeGetChannelFailed: {
        UserMessage:  "服务暂时异常，请稍后重试",
        InternalOnly: true,
        StatusCode:   500,
    },
    ErrorCodeChannelInvalidKey: {
        UserMessage:  "当前渠道配置异常，已自动切换",
        InternalOnly: true,
        StatusCode:   500,
    },
}
```

**使用方式：**

```go
func (e *NewAPIError) ToUserResponse() interface{} {
    detail, ok := errorMessages[e.errorCode]
    if !ok {
        detail = ErrorDetail{
            UserMessage:  "服务异常，请稍后重试",
            InternalOnly: true,
            StatusCode:   500,
        }
    }

    // 仅返回用户友好的消息
    return gin.H{
        "error": detail.UserMessage,
        "code":  string(e.errorCode),
    }
}

// 日志中记录完整错误
func (e *NewAPIError) ToLogMessage() string {
    return fmt.Sprintf("[%s] %s: %v", e.errorCode, e.errorType, e.Err)
}
```

### 4.3 扩展使用 ErrOptionWithHideErrMsg

项目已提供 `types/error.go:364-371` 的隐藏错误选项，应在更多场景使用：

**当前实现：**

```go
func ErrOptionWithHideErrMsg(replaceStr string) NewAPIErrorOptions {
    return func(e *NewAPIError) {
        if common.DebugEnabled {
            fmt.Printf("ErrOptionWithHideErrMsg: %s, origin error: %s", replaceStr, e.Err)
        }
        e.Err = errors.New(replaceStr)
    }
}
```

**建议应用场景：**

```go
// 渠道相关错误
return types.NewError(err, types.ErrorCodeChannelInvalidKey,
    types.ErrOptionWithHideErrMsg("当前渠道配置异常，已自动切换"))

// 数据库错误
return types.NewError(err, types.ErrorCodeQueryDataError,
    types.ErrOptionWithHideErrMsg("服务暂时异常，请联系管理员"))

// 网络请求错误
return types.NewError(err, types.ErrorCodeDoRequestFailed,
    types.ErrOptionWithHideErrMsg("上游服务暂时不可用，请稍后重试"))

// AWS 错误
return types.NewError(err, types.ErrorCodeAwsInvokeError,
    types.ErrOptionWithHideErrMsg("云服务调用失败，请稍后重试"))
```

### 4.4 响应体中的敏感信息处理

**问题：** `service/error.go:93-107` 会将完整响应体包含在错误中

```go
buildErrWithBody := func(message string) error {
    if message == "" {
        return fmt.Errorf("bad response status code %d, body: %s",
            resp.StatusCode, string(responseBody))  // ⚠️ 暴露上游响应
    }
    return fmt.Errorf("bad response status code %d, message: %s, body: %s",
        resp.StatusCode, message, string(responseBody))
}
```

**改进建议：**

```go
buildErrWithBody := func(message string) error {
    // 调试模式：记录完整响应
    if common.DebugEnabled {
        logger.LogError(ctx, fmt.Sprintf("bad response status code %d, body: %s",
            resp.StatusCode, string(responseBody)))
    }

    // 用户模式：返回友好消息
    if message == "" {
        return fmt.Errorf("upstream service returned status %d", resp.StatusCode)
    }
    return fmt.Errorf("upstream service error: %s", message)
}
```

### 4.5 错误码重构建议（低优先级）

按领域重新分类错误码，区分用户可见和内部错误：

```go
const (
    // ===== 用户端错误 (4xx) - 对外暴露 =====
    ErrorCodeInvalidRequest         ErrorCode = "invalid_request"
    ErrorCodeInsufficientQuota      ErrorCode = "insufficient_quota"
    ErrorCodeAccessDenied           ErrorCode = "access_denied"
    ErrorCodeModelNotFound          ErrorCode = "model_not_found"
    ErrorCodeRateLimitExceeded      ErrorCode = "rate_limit_exceeded"

    // ===== 服务端错误 (5xx) - 对外隐藏详情 =====
    ErrorCodeInternalServiceError   ErrorCode = "internal_service_error"
    ErrorCodeUpstreamServiceError   ErrorCode = "upstream_service_error"
    ErrorCodeDatabaseError          ErrorCode = "database_error"
    ErrorCodeChannelError           ErrorCode = "channel_error"

    // ===== 内部错误码 - 仅用于日志和监控 =====
    // 格式：internal:category:action
    errorCodeDoRequestFailed        ErrorCode = "internal:network:do_request_failed"
    errorCodeQueryDataError         ErrorCode = "internal:database:query_failed"
    errorCodeUpdateDataError        ErrorCode = "internal:database:update_failed"
    errorCodeGetChannelFailed       ErrorCode = "internal:channel:get_failed"
    errorCodeJsonMarshalFailed      ErrorCode = "internal:serialize:json_failed"
    errorCodeCountTokenFailed       ErrorCode = "internal:token:count_failed"
)
```

**优势：**
1. 用户只能看到通用错误码，无法获取内部实现细节
2. 开发者可通过内部错误码快速定位问题
3. 便于错误监控和告警系统分类统计

---

## 5. 立即行动清单

### 5.1 高优先级（必须修复）

- [ ] **修复 ChannelError.UsingKey 密钥泄露**
  - 位置：`types/channel_error.go:9`
  - 方案：改为 `json:"-"` 或添加脱敏逻辑
  - 影响：严重安全漏洞

- [ ] **检查所有渠道错误是否暴露密钥信息**
  - 搜索关键词：`channel:invalid_key`, `UsingKey`
  - 确保所有返回给用户的响应都已脱敏

### 5.2 中优先级（建议修复）

- [ ] **为技术性错误码添加用户友好消息**
  - 定义错误消息映射表
  - 实现 `ToUserResponse()` 方法
  - 更新错误处理流程

- [ ] **扩展使用 ErrOptionWithHideErrMsg**
  - 在所有 `channel:*` 错误使用
  - 在所有 SQL 错误使用
  - 在所有网络/IO 错误使用

- [ ] **改进响应体错误处理**
  - 修改 `buildErrWithBody` 函数
  - 仅在调试模式输出完整响应

### 5.3 低优先级（优化建议）

- [ ] **重构错误码分类结构**
  - 区分用户可见和内部错误码
  - 统一错误码命名规范

- [ ] **添加错误码文档**
  - 为每个错误码编写说明文档
  - 包含触发原因、解决方法

- [ ] **建立错误码监控**
  - 统计各错误码出现频率
  - 设置告警阈值

---

## 6. 最佳实践建议

### 6.1 错误处理原则

1. **用户视角**：只返回用户能理解和处理的信息
2. **开发者视角**：日志中记录完整的错误上下文
3. **安全原则**：绝不暴露敏感配置、密钥、内部路径
4. **调试友好**：开发环境可输出详细错误，生产环境脱敏

### 6.2 错误响应格式示例

**对外（用户）：**

```json
{
  "error": {
    "message": "上游服务暂时不可用，请稍后重试",
    "code": "upstream_service_error",
    "type": "api_error"
  }
}
```

**内部（日志）：**

```
[ERROR] [internal:network:do_request_failed] 2026-02-27 10:30:45
Context: user_id=12345, model=gpt-4, channel_id=6
Error: Post "https://api.openai.com/v1/chat/completions": dial tcp: lookup api.openai.com: no such host
Channel: OpenAI (Channel #6), Key: sk-***...*** (multi-key mode)
Retry: skipped (non-retryable error)
```

### 6.3 代码审查检查点

在代码审查时，关注以下问题：

- [ ] 错误消息是否包含 URL、IP、密钥等敏感信息？
- [ ] 错误消息是否暴露内部实现细节？
- [ ] 错误消息对用户是否可理解？
- [ ] 是否需要��用 `ErrOptionWithHideErrMsg`？
- [ ] 是否区分了用户响应和日志记录？

---

## 附录 A：错误码完整列表

### A.1 按风险等级分类

| 风险等级 | 错误码 | 需要修复 |
|---------|--------|----------|
| 🔴 严重 | `channel:invalid_key` | ✅ 是 |
| 🔴 严重 | `ChannelError.UsingKey` | ✅ 是 |
| 🔴 高 | `do_request_failed` | ✅ 是 |
| 🔴 高 | `get_channel_failed` | ✅ 是 |
| 🔴 高 | `aws_invoke_error` | ✅ 是 |
| 🔴 高 | `channel:aws_client_error` | ✅ 是 |
| 🔴 高 | `query_data_error` | ✅ 是 |
| 🔴 高 | `update_data_error` | ✅ 是 |
| 🟡 中 | `json_marshal_failed` | ⚠️ 建议 |
| 🟡 中 | `count_token_failed` | ⚠️ 建议 |
| 🟡 中 | `gen_relay_info_failed` | ⚠️ 建议 |
| 🟡 中 | `convert_request_failed` | ⚠️ 建议 |
| 🟡 中 | `read_response_body_failed` | ⚠️ 建议 |
| 🟡 中 | `pre_consume_token_quota_failed` | ⚠️ 建议 |
| 🟢 低 | `invalid_request` | ❌ 否 |
| 🟢 低 | `insufficient_user_quota` | ❌ 否 |
| 🟢 低 | `model_not_found` | ❌ 否 |

### A.2 按功能模块分类

#### 通用模块
- `invalid_request`
- `sensitive_words_detected`

#### 渠道管理
- `channel:no_available_key`
- `channel:param_override_invalid`
- `channel:header_override_invalid`
- `channel:model_mapped_error`
- `channel:aws_client_error`
- `channel:invalid_key`
- `channel:response_time_exceeded`

#### 网络请求
- `do_request_failed`
- `read_request_body_failed`
- `bad_request_body`
- `read_response_body_failed`
- `bad_response_status_code`
- `bad_response`
- `bad_response_body`
- `empty_response`

#### 数据存储
- `query_data_error`
- `update_data_error`

#### 业务逻辑
- `count_token_failed`
- `model_price_error`
- `insufficient_user_quota`
- `pre_consume_token_quota_failed`
- `access_denied`
- `model_not_found`
- `prompt_blocked`

#### 协议转换
- `convert_request_failed`
- `json_marshal_failed`
- `gen_relay_info_failed`
- `invalid_api_type`

#### 云服务
- `aws_invoke_error`

---

## 附录 B：相关文件清单

| 文件路径 | 关键行号 | 说明 |
|---------|----------|------|
| `types/error.go` | 38-87 | 错误码常量定义 |
| `types/error.go` | 89-183 | NewAPIError 结构体及方法 |
| `types/channel_error.go` | 3-22 | ChannelError 结构体（⚠️ 密钥泄露风险） |
| `common/str.go` | 167-239 | MaskSensitiveInfo 函数 |
| `service/error.go` | 84-127 | RelayErrorHandler 函数 |
| `service/error.go` | 148-172 | TaskErrorWrapper 函数 |

---

## 附录 C：参考资料

- [OWASP 错误处理指南](https://cheatsheetseries.owasp.org/cheatsheets/Error_Handling_Cheat_Sheet.html)
- [RESTful API 错误响应最佳实践](https://restfulapi.net/http-status-codes/)
- [Go 错误处理最佳实践](https://go.dev/doc/tutorial/errors)

---

**文档维护：** 请在修改错误码时及时更新本文档
**反馈渠道：** 如发现问题或建议，请联系技术负责人
