# 错误码系统文档

此目录包含 New API 错误码系统的完整文档。

## 📚 文档

### 1. [错误码改进提案](./error-code-improvements_CN.md)
**📋 设计提案和架构概述**

本文档详细说明了错误码系统的改进建议，包括：
- 当前状态分析（优缺点）
- 详细的改进提案
- 实施策略和迁移计划
- 改进优先级矩阵

**适用于**: 了解变更背后的原因

---

### 2. [错误码参考](./ERROR_CODES.md)
**📖 完整的错误码参考**

所有错误码的自动生成参考，包括：
- 所有错误码及其数值
- HTTP 状态码映射
- 错误级别分配
- 使用示例
- 类别组织

**适用于**: 查找特定错误码及其含义

**如何重新生成**:
```bash
go run tools/generate_error_doc.go > docs/ERROR_CODES.md
```

---

### 3. [迁移指南](./error-code-migration-guide_CN.md)
**🔄 分步迁移说明**

从旧的基于字符串的错误码迁移到新的数字系统的实用指南：
- 迁移前后的代码示例
- 迁移步骤
- 常见问题和解决方案
- 测试策略

**适用于**: 更新现有代码的开发者

---

## 🏗️ 架构概述

新的错误码系统由以下几个组件组成：

### 核心文件

| 文件 | 描述 |
|------|-------------|
| `types/error.go` | 核心错误实现，包含 `NewAPIError` 结构体 |
| `types/error_code.go` | 数字错误码定义和 HTTP 映射 |
| `types/error_level.go` | 错误严重级别定义 |
| `types/error_i18n.go` | 错误消息国际化支持 |

### 工具

| 文件 | 描述 |
|------|-------------|
| `tools/generate_error_doc.go` | 自动生成 ERROR_CODES.md 文档 |

---

## 🎯 快速参考

### 错误码类别

| 类别 | 范围 | HTTP 状态 | 示例 |
|----------|-------|-------------|---------|
| 通用错误 | 1xxx | 400 | 无效请求 (1001) |
| 系统错误 | 2xxx | 500 | Token 计数失败 (2001) |
| 渠道错误 | 3xxx | 503 | 渠道不可用 (3008) |
| 客户端错误 | 4xxx | 401 | 访问被拒绝 (4003) |
| 上游错误 | 5xxx | 502 | 错误响应 (5003) |
| 数据库错误 | 6xxx | 500 | 查询数据错误 (6001) |
| 配额错误 | 7xxx | 402 | 配额不足 (7001) |

### 错误级别

| 级别 | 描述 | 颜色 |
|-------|-------------|-------|
| info | 信息性消息 | 青色 |
| warning | 非关键警告 | 黄色 |
| error | 错误事件 | 红色 |
| critical | 严重故障 | 洋红色 |

### 支持的语言

| 代码 | 语言 |
|------|----------|
| en | English |
| zh | 中文 |
| ja | 日本語 |
| fr | Français |
| ru | Русский |
| vi | Tiếng Việt |

---

## 💻 使用示例

### 创建错误

```go
import "github.com/QuantumNous/new-api/types"

// 自动映射 HTTP 状态码和错误级别
err := types.NewError(
    errors.New("渠道不可用"),
    types.ErrorCodeChannelNoAvailableKey,
)

// 访问错误属性
fmt.Println(err.StatusCode)         // 503
fmt.Println(err.Level)              // error
fmt.Println(err.errorCode.String()) // "channel_no_available_key"
```

### 本地化错误消息

```go
// 获取本地化错误消息
lang := types.GetLanguageFromContext("zh-CN")
message := err.Localize(lang)
// 返回: "渠道不可用"
```

### 自定义错误级别

```go
// 覆盖默认错误级别
err := types.NewError(
    errors.New("自定义错误"),
    types.ErrorCodeInvalidRequest,
    types.ErrOptionWithLevel(types.ErrorLevelCritical),
)
```

### 错误选项

```go
err := types.NewError(
    errors.New("出错了"),
    types.ErrorCodeInvalidRequest,
    types.ErrOptionWithSkipRetry(),           // 跳过重试
    types.ErrOptionWithNoRecordErrorLog(),    // 不记录日志
    types.ErrOptionWithHideErrMsg("Hidden"),  // 隐藏消息
)
```

---

## 🔄 迁移摘要

### 主要变更

1. **ErrorCode 类型**: 从 `string` 改为 `int`
2. **HTTP 状态自动映射**: 无需手动设置状态码
3. **错误级别**: 内置严重级别
4. **国际化**: 支持 6 种语言
5. **更好的性能**: 整数比较 vs 字符串比较

### 破坏性变更

- `ErrorCode` 现在是 `int` 而不是 `string`
- 使用 `.String()` 方法获取字符串表示
- 旧常量重命名为 `LegacyErrorCodeString`（向后兼容）

### 优势

✅ 类型安全（编译时检查）
✅ 更快的错误比较
✅ 自动 HTTP 状态映射
✅ 内置错误级别用于监控
✅ 易于国际化
✅ 更好的开发者体验

---

## 📊 统计

- **错误码总数**: 34
- **类别**: 7 (1xxx-7xxx)
- **支持语言**: 6
- **文档文件**: 3
- **代码文件**: 4 (types/ + tools/)

---

## 🛠️ 维护

### 添加新错误码

1. 添加到 `types/error_code.go`:
   ```go
   ErrorCodeMyNewError ErrorCode = 1009
   ```

2. 添加字符串映射:
   ```go
   errorCodeStrings[ErrorCodeMyNewError] = "my_new_error"
   ```

3. 添加 HTTP 状态映射:
   ```go
   errorCodeHTTPStatusMap[ErrorCodeMyNewError] = http.StatusBadRequest
   ```

4. 添加级别映射:
   ```go
   errorCodeLevelMap[ErrorCodeMyNewError] = ErrorLevelWarning
   ```

5. 添加到 `tools/generate_error_doc.go`:
   ```go
   {1009, "ErrorCodeMyNewError", "General Errors (1xxx)", 400, "warning", "My new error description"},
   ```

6. 添加本地化消息到 `types/error_i18n.go`

7. 重新生成文档:
   ```bash
   go run tools/generate_error_doc.go > docs/ERROR_CODES.md
   ```

---

## 📞 支持

如有问题、建议或反馈：

1. 查看 [迁移指南](./error-code-migration-guide_CN.md) 了解常见问题
2. 阅读 [改进提案](./error-code-improvements_CN.md) 了解设计细节
3. 搜索 [错误码参考](./ERROR_CODES.md) 查找特定代码
4. 在 GitHub 上提交 issue 报告错误或请求功能

---

**最后更新**: 2026-02-26
**版本**: 1.0
**维护者**: New API 团队
