# ErrorType 与 ErrorCode 映射关系

本文档说明了系统中定义的 7 种 ErrorType 与其对应的 ErrorCode 映射关系。

## ErrorType 概览

系统中定义了 7 种错误类型：

```go
ErrorTypeNewAPIError     = "new_api_error"     // 内部系统错误
ErrorTypeOpenAIError     = "openai_error"     // OpenAI 上游错误
ErrorTypeClaudeError     = "claude_error"     // Claude 上游错误
ErrorTypeMidjourneyError = "midjourney_error" // Midjourney 错误（未实现）
ErrorTypeGeminiError     = "gemini_error"     // Gemini 错误（未实现）
ErrorTypeRerankError     = "rerank_error"     // Rerank 错误（未实现）
ErrorTypeUpstreamError   = "upstream_error"   // 通用上游错误
```

## ErrorCode 数值范围

| 范围 | 分类 | 说明 |
|------|------|------|
| 1xxx | General Errors | 通用错误 |
| 2xxx | System Errors | 系统错误 |
| 3xxx | Channel Errors | 渠道错误 |
| 4xxx | Client Errors | 客户端错误 |
| 5xxx | Upstream Errors | 上游服务错误 |
| 6xxx | Database Errors | 数据库错误 |
| 7xxx | Quota Errors | 配额错误 |

---

## 1. ErrorTypeNewAPIError (new_api_error)

通过 `types.NewError()` 函数设置，用于系统内部错误。

### 创建函数

```go
func NewError(err error, errorCode ErrorCode, ops ...NewAPIErrorOptions) *NewAPIError
```

### 对应的 ErrorCode

| ErrorCode 常量 | 数值 | 字符串表示 | HTTP 状态码 | 使用场景 |
|---------------|------|-----------|------------|---------|
| `ErrorCodeInvalidRequest` | 1001 | invalid_request | 400 | 通用请求错误 |
| `ErrorCodeSensitiveWordsDetected` | 1002 | sensitive_words_detected | 400 | 敏感词检测 |
| `ErrorCodeCountTokenFailed` | 2001 | count_token_failed | 500 | Token 计数失败 |
| `ErrorCodeModelPriceError` | 2002 | model_price_error | 500 | 模型价格错误 |
| `ErrorCodeInvalidApiType` | 2003 | invalid_api_type | 400 | 无效 API 类型 |
| `ErrorCodeJsonMarshalFailed` | 2004 | json_marshal_failed | 500 | JSON 序列化失败 |
| `ErrorCodeJsonUnmarshalFailed` | 2005 | json_unmarshal_failed | 500 | JSON 反序列化失败 |
| `ErrorCodeDoRequestFailed` | 2006 | do_request_failed | 500 | HTTP 请求失败 |
| `ErrorCodeGetChannelFailed` | 2007 | get_channel_failed | 500 | 获取渠道失败 |
| `ErrorCodeGenRelayInfoFailed` | 2008 | gen_relay_info_failed | 500 | 生成中继信息失败 |
| `ErrorCodeChannelNoAvailableKey` | 3001 | channel_no_available_key | 503 | 渠道无可用密钥 |
| `ErrorCodeChannelParamOverrideInvalid` | 3002 | channel_param_override_invalid | 400 | 渠道参数覆盖无效 |
| `ErrorCodeChannelHeaderOverrideInvalid` | 3003 | channel_header_override_invalid | 400 | 渠道请求头覆盖无效 |
| `ErrorCodeChannelModelMappedError` | 3004 | channel_model_mapped_error | 500 | 渠道模型映射错误 |
| `ErrorCodeChannelAwsClientError` | 3005 | channel_aws_client_error | 500 | AWS 客户端错误 |
| `ErrorCodeChannelInvalidKey` | 3006 | channel_invalid_key | 401 | 渠道密钥无效 |
| `ErrorCodeChannelResponseTimeExceeded` | 3007 | channel_response_time_exceeded | 504 | 渠道响应超时 |
| `ErrorCodeChannelNotAvailable` | 3008 | channel_not_available | 503 | 渠道不可用 |
| `ErrorCodeReadRequestBodyFailed` | 4001 | read_request_body_failed | 400 | 读取请求体失败 |
| `ErrorCodeConvertRequestFailed` | 4002 | convert_request_failed | 400 | 请求转换失败 |
| `ErrorCodeAccessDenied` | 4003 | access_denied | 401 | 访问拒绝 |
| `ErrorCodeBadRequestBody` | 4004 | bad_request_body | 400 | 错误的请求体 |
| `ErrorCodeUnauthorized` | 4005 | unauthorized | 401 | 未授权 |
| `ErrorCodeForbidden` | 4006 | forbidden | 403 | 禁止访问 |
| `ErrorCodeQueryDataError` | 6001 | query_data_error | 500 | 查询数据错误 |
| `ErrorCodeUpdateDataError` | 6002 | update_data_error | 500 | 更新数据错误 |
| `ErrorCodeInsertDataError` | 6003 | insert_data_error | 500 | 插入数据错误 |
| `ErrorCodeDeleteDataError` | 6004 | delete_data_error | 500 | 删除数据错误 |
| `ErrorCodeDatabaseConnectionFailed` | 6005 | database_connection_failed | 500 | 数据库连接失败 |
| `ErrorCodeInsufficientUserQuota` | 7001 | insufficient_user_quota | 402 | 用户配额不足 |
| `ErrorCodePreConsumeTokenQuotaFailed` | 7002 | pre_consume_token_quota_failed | 500 | 预消费 Token 配额失败 |
| `ErrorCodeQuotaExceeded` | 7003 | quota_exceeded | 402 | 配额超限 |

### 使用示例

```go
return types.NewError(err, types.ErrorCodeGetChannelFailed, types.ErrOptionWithSkipRetry())
```

---

## 2. ErrorTypeOpenAIError (openai_error)

通过 `types.WithOpenAIError()` 或 `types.NewOpenAIError()` 设置，用于 OpenAI 兼容 API 的上游错误。

### 创建函数

```go
func WithOpenAIError(openAIError OpenAIError, statusCode int, ops ...NewAPIErrorOptions) *NewAPIError
func NewOpenAIError(err error, errorCode ErrorCode, statusCode int, ops ...NewAPIErrorOptions) *NewAPIError
```

### ErrorCode 映射机制

ErrorCode 通过 `ErrorCodeFromString(openAIError.Code)` 从字符串动态映射。

### 常见上游错误映射

| 上游 Error.Code 字符串 | 映射的 ErrorCode 常量 | 数值 | HTTP 状态码 |
|----------------------|---------------------|------|------------|
| `invalid_request_error` | `ErrorCodeInvalidRequest` | 1001 | 400 |
| `context_length_exceeded` | `ErrorCodeInvalidRequest` | 1001 | 400 |
| `rate_limit_exceeded` | `ErrorCodeRateLimitExceeded` | 5009 | 429 |
| `insufficient_quota` | `ErrorCodeInsufficientUserQuota` | 7001 | 402 |
| `billing_not_active` | `ErrorCodeInsufficientUserQuota` | 7001 | 402 |
| `model_not_found` | `ErrorCodeModelNotFound` | 5007 | 404 |
| `service_unavailable` | `ErrorCodeServiceUnavailable` | 5010 | 503 |
| `server_error` | `ErrorCodeBadResponse` | 5003 | 502 |
| `api_key_error` | `ErrorCodeChannelInvalidKey` | 3006 | 401 |
| 未匹配的字符串 | `ErrorCodeInvalidRequest` | 1001 | 400 (默认) |

### 使用示例

```go
// 在 relay/channel/openai/relay-openai.go 中
if oaiError := simpleResponse.GetOpenAIError(); oaiError != nil && oaiError.Type != "" {
    return nil, types.WithOpenAIError(*oaiError, resp.StatusCode)
}
```

---

## 3. ErrorTypeClaudeError (claude_error)

通过 `types.WithClaudeError()` 设置，用于 Claude (Anthropic) API 的上游错误。

### 创建函数

```go
func WithClaudeError(claudeError ClaudeError, statusCode int, ops ...NewAPIErrorOptions) *NewAPIError
```

### ErrorCode 映射机制

ErrorCode 通过 `ErrorCodeFromString(claudeError.Type)` 从字符串动态映射。

### 常见上游错误映射

| 上游 Error.Type 字符串 | 映射的 ErrorCode 常量 | 数值 | 说明 |
|----------------------|---------------------|------|------|
| `invalid_request_error` | `ErrorCodeInvalidRequest` | 1001 | 请求无效 |
| `type_error` | `ErrorCodeInvalidRequest` | 1001 | 类型错误 |
| `permission_error` | `ErrorCodeAccessDenied` | 4003 | 权限错误 |
| `authentication_error` | `ErrorCodeUnauthorized` | 4005 | 认证错误 |
| `rate_limit_error` | `ErrorCodeRateLimitExceeded` | 5009 | 速率限制 |
| `overloaded_error` | `ErrorCodeServiceUnavailable` | 5010 | 服务过载 |
| 未匹配的字符串 | `ErrorCodeInvalidRequest` | 1001 | 默认 |

### 使用示例

```go
// 在 relay/channel/claude/relay-claude.go 中
if claudeError := claudeResponse.GetClaudeError(); claudeError != nil && claudeError.Type != "" {
    return types.WithClaudeError(*claudeError, http.StatusInternalServerError)
}
```

---

## 4. ErrorTypeMidjourneyError (midjourney_error)

**状态：已定义但未实现**

该错误类型已定义，但代码中未找到对应的包装函数或使用场景。

### 可能的未来实现

```go
// 预期函数（尚未实现）
func WithMidjourneyError(midjourneyError MidjourneyError, statusCode int, ops ...NewAPIErrorOptions) *NewAPIError
```

---

## 5. ErrorTypeGeminiError (gemini_error)

**状态：已定义但未实现**

该错误类型已定义，但代码中未找到对应的包装函数或使用场景。

Gemini 相关错误目前通过 `ErrorTypeNewAPIError` 处理，使用 `types.NewError()` 函数。

### 当前 Gemini 错误处理示例

```go
// relay/gemini_handler.go
return types.NewError(fmt.Errorf("failed to copy request to GeminiChatRequest: %w", err),
    types.ErrorCodeInvalidRequest, types.ErrOptionWithSkipRetry())
```

---

## 6. ErrorTypeRerankError (rerank_error)

**状态：已定义但未实现**

该错误类型已定义，但代码中未找到对应的包装函数或使用场景。

Rerank 相关错误目前通过 `ErrorTypeNewAPIError` 处理。

### 当前 Rerank 错误处理示例

```go
// relay/rerank_handler.go
return types.NewError(err, types.ErrorCodeConvertRequestFailed, types.ErrOptionWithSkipRetry())
```

---

## 7. ErrorTypeUpstreamError (upstream_error)

**状态：默认上游错误类型**

该错误类型作为上游错误的默认 type 使用。

### 使用场景

当 `WithOpenAIError()` 或 `WithClaudeError()` 接收到空的 Type 字段时，会自动设置为 `upstream_error`。

```go
// types/error.go:302
if openAIError.Type == "" {
    openAIError.Type = "upstream_error"
}

// types/error.go:327
if claudeError.Type == "" {
    claudeError.Type = "upstream_error"
}
```

---

## 关键发现总结

### 实际使用的 ErrorType

| ErrorType | 状态 | 主要使用场景 |
|-----------|------|-------------|
| `ErrorTypeNewAPIError` | ✅ 活跃 | 系统内部错误 |
| `ErrorTypeOpenAIError` | ✅ 活跃 | OpenAI 兼容 API 上游错误 |
| `ErrorTypeClaudeError` | ✅ 活跃 | Claude API 上游错误 |
| `ErrorTypeMidjourneyError` | ⚠️ 预留 | 未实现 |
| `ErrorTypeGeminiError` | ⚠️ 预留 | 未实现 |
| `ErrorTypeRerankError` | ⚠️ 预留 | 未实现 |
| `ErrorTypeUpstreamError` | 🔄 默认 | 作为上游错误的默认 type |

### ErrorCode 映射逻辑

1. **NewAPIError**: 直接传入数值型 ErrorCode
2. **OpenAIError**: 从 `openAIError.Code` 字符串通过 `ErrorCodeFromString()` 转换
3. **ClaudeError**: 从 `claudeError.Type` 字符串通过 `ErrorCodeFromString()` 转换

### ErrorCodeFromString() 转换规则

```go
func ErrorCodeFromString(s string) ErrorCode {
    for code, str := range errorCodeStrings {
        if str == s {
            return code
        }
    }
    return ErrorCodeInvalidRequest // 默认返回 1001
}
```

如果字符串无法匹配到已知的 ErrorCode，会默认返回 `ErrorCodeInvalidRequest` (1001)。

---

## 相关文件

- `types/error.go` - ErrorType 定义和错误结构
- `types/error_code.go` - ErrorCode 定义和映射表
- `service/error.go` - 错误处理辅助函数
- `relay/` - 各渠道错误处理实现

---

## 更新时间

文档生成时间: 2026-02-27
