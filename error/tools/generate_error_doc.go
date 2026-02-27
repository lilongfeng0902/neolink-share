//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"os"
	"sort"
	"text/template"
)

// ErrorDoc represents documentation for a single error code
type ErrorDoc struct {
	Code       int
	Name       string
	Category   string
	HTTPStatus int
	Level      string
	Description string
}

// Category groups error codes by their numeric range
func getCategory(code int) string {
	switch {
	case code >= 1000 && code < 2000:
		return "General Errors (1xxx)"
	case code >= 2000 && code < 3000:
		return "System Errors (2xxx)"
	case code >= 3000 && code < 4000:
		return "Channel Errors (3xxx)"
	case code >= 4000 && code < 5000:
		return "Client Errors (4xxx)"
	case code >= 5000 && code < 6000:
		return "Upstream Errors (5xxx)"
	case code >= 6000 && code < 7000:
		return "Database Errors (6xxx)"
	case code >= 7000 && code < 8000:
		return "Quota Errors (7xxx)"
	default:
		return "Other"
	}
}

// Error code definitions (must match types/error_code.go)
var errorCodes = []ErrorDoc{
	// General Errors (1xxx)
	{1001, "ErrorCodeInvalidRequest", "General Errors (1xxx)", 400, "warning", "Invalid request parameters"},
	{1002, "ErrorCodeSensitiveWordsDetected", "General Errors (1xxx)", 400, "warning", "Sensitive words detected in content"},

	// System Errors (2xxx)
	{2001, "ErrorCodeCountTokenFailed", "System Errors (2xxx)", 500, "error", "Failed to count tokens"},
	{2002, "ErrorCodeModelPriceError", "System Errors (2xxx)", 500, "error", "Model pricing configuration error"},
	{2003, "ErrorCodeInvalidApiType", "System Errors (2xxx)", 400, "error", "Invalid API type"},
	{2004, "ErrorCodeJsonMarshalFailed", "System Errors (2xxx)", 500, "error", "Failed to marshal JSON"},
	{2005, "ErrorCodeJsonUnmarshalFailed", "System Errors (2xxx)", 500, "error", "Failed to unmarshal JSON"},
	{2006, "ErrorCodeDoRequestFailed", "System Errors (2xxx)", 500, "error", "Failed to make HTTP request"},
	{2007, "ErrorCodeGetChannelFailed", "System Errors (2xxx)", 500, "critical", "Failed to get channel information"},
	{2008, "ErrorCodeGenRelayInfoFailed", "System Errors (2xxx)", 500, "error", "Failed to generate relay information"},

	// Channel Errors (3xxx)
	{3001, "ErrorCodeChannelNoAvailableKey", "Channel Errors (3xxx)", 503, "error", "No available API key in channel"},
	{3002, "ErrorCodeChannelParamOverrideInvalid", "Channel Errors (3xxx)", 400, "warning", "Invalid channel parameter override"},
	{3003, "ErrorCodeChannelHeaderOverrideInvalid", "Channel Errors (3xxx)", 400, "warning", "Invalid channel header override"},
	{3004, "ErrorCodeChannelModelMappedError", "Channel Errors (3xxx)", 500, "error", "Channel model mapping error"},
	{3005, "ErrorCodeChannelAwsClientError", "Channel Errors (3xxx)", 500, "error", "AWS client configuration error"},
	{3006, "ErrorCodeChannelInvalidKey", "Channel Errors (3xxx)", 401, "warning", "Invalid channel API key"},
	{3007, "ErrorCodeChannelResponseTimeExceeded", "Channel Errors (3xxx)", 504, "warning", "Channel response time exceeded"},
	{3008, "ErrorCodeChannelNotAvailable", "Channel Errors (3xxx)", 503, "critical", "Channel is not available"},

	// Client Errors (4xxx)
	{4001, "ErrorCodeReadRequestBodyFailed", "Client Errors (4xxx)", 400, "warning", "Failed to read request body"},
	{4002, "ErrorCodeConvertRequestFailed", "Client Errors (4xxx)", 400, "warning", "Failed to convert request format"},
	{4003, "ErrorCodeAccessDenied", "Client Errors (4xxx)", 401, "warning", "Access denied"},
	{4004, "ErrorCodeBadRequestBody", "Client Errors (4xxx)", 400, "warning", "Invalid request body"},
	{4005, "ErrorCodeUnauthorized", "Client Errors (4xxx)", 401, "warning", "Unauthorized access"},
	{4006, "ErrorCodeForbidden", "Client Errors (4xxx)", 403, "warning", "Forbidden"},

	// Upstream Errors (5xxx)
	{5001, "ErrorCodeReadResponseBodyFailed", "Upstream Errors (5xxx)", 500, "error", "Failed to read response body"},
	{5002, "ErrorCodeBadResponseStatusCode", "Upstream Errors (5xxx)", 502, "error", "Bad response status code from upstream"},
	{5003, "ErrorCodeBadResponse", "Upstream Errors (5xxx)", 502, "error", "Bad response from upstream service"},
	{5004, "ErrorCodeBadResponseBody", "Upstream Errors (5xxx)", 500, "error", "Invalid response body format"},
	{5005, "ErrorCodeEmptyResponse", "Upstream Errors (5xxx)", 500, "error", "Empty response from upstream"},
	{5006, "ErrorCodeAwsInvokeError", "Upstream Errors (5xxx)", 500, "error", "AWS invocation error"},
	{5007, "ErrorCodeModelNotFound", "Upstream Errors (5xxx)", 404, "warning", "Model not found"},
	{5008, "ErrorCodePromptBlocked", "Upstream Errors (5xxx)", 400, "warning", "Prompt blocked by content filter"},
	{5009, "ErrorCodeRateLimitExceeded", "Upstream Errors (5xxx)", 429, "warning", "Rate limit exceeded"},
	{5010, "ErrorCodeServiceUnavailable", "Upstream Errors (5xxx)", 503, "critical", "Service temporarily unavailable"},

	// Database Errors (6xxx)
	{6001, "ErrorCodeQueryDataError", "Database Errors (6xxx)", 500, "critical", "Database query error"},
	{6002, "ErrorCodeUpdateDataError", "Database Errors (6xxx)", 500, "critical", "Database update error"},
	{6003, "ErrorCodeInsertDataError", "Database Errors (6xxx)", 500, "critical", "Database insert error"},
	{6004, "ErrorCodeDeleteDataError", "Database Errors (6xxx)", 500, "critical", "Database delete error"},
	{6005, "ErrorCodeDatabaseConnectionFailed", "Database Errors (6xxx)", 500, "critical", "Database connection failed"},

	// Quota Errors (7xxx)
	{7001, "ErrorCodeInsufficientUserQuota", "Quota Errors (7xxx)", 402, "warning", "Insufficient user quota"},
	{7002, "ErrorCodePreConsumeTokenQuotaFailed", "Quota Errors (7xxx)", 500, "error", "Failed to pre-consume token quota"},
	{7003, "ErrorCodeQuotaExceeded", "Quota Errors (7xxx)", 402, "warning", "User quota exceeded"},
}

// markdownTemplate is the template for generating the documentation
const markdownTemplate = `# Error Codes Reference

> **Auto-generated**: Do not edit manually
> **Generated by**: tools/generate_error_doc.go
> **Last updated**: {{.Timestamp}}

## Overview

This document provides a comprehensive reference of all error codes used in the New API system.

### Error Code Categories

| Category | Range | Description |
|----------|-------|-------------|
| General Errors | 1xxx | General request and validation errors |
| System Errors | 2xxx | Internal system errors |
| Channel Errors | 3xxx | Channel and provider-related errors |
| Client Errors | 4xxx | Client request errors |
| Upstream Errors | 5xxx | Upstream provider errors |
| Database Errors | 6xxx | Database operation errors |
| Quota Errors | 7xxx | Quota and billing errors |

### Error Levels

| Level | Description | Color |
|-------|-------------|-------|
| info | Informational messages | Cyan |
| warning | Warning messages that don't prevent operation | Yellow |
| error | Error events that might allow continuation | Red |
| critical | Critical errors that may cause termination | Magenta |

---

{{range .Categories}}
## {{.Category}}

| Code | Constant Name | HTTP Status | Level | Description |
|------|---------------|-------------|-------|-------------|
{{range .Errors}}| {{.Code}} | ` + "`" + `{{.Name}}` + "`" + ` | {{.HTTPStatus}} | {{.Level}} | {{.Description}} |
{{end}}

---

{{end}}

## Usage Examples

### Creating an Error

` + "```" + `go
import "github.com/QuantumNous/new-api/types"

// Create error with automatic status code and level
err := types.NewError(
    errors.New("channel not available"),
    types.ErrorCodeChannelNoAvailableKey,
)

// Access error properties
fmt.Println(err.StatusCode)        // 503
fmt.Println(err.Level)             // error
fmt.Println(err.errorCode.String()) // "channel_no_available_key"
` + "```" + `

### Localized Error Messages

` + "```" + `go
// Get localized error message
lang := types.GetLanguageFromContext("zh-CN")
message := err.Localize(lang)
// Returns: "渠道不可用"
` + "```" + `

### Custom Error Level

` + "```" + `go
// Override default error level
err := types.NewError(
    errors.New("custom error"),
    types.ErrorCodeInvalidRequest,
    types.ErrOptionWithLevel(types.ErrorLevelCritical),
)
` + "```" + `

### Error Handling in Controllers

` + "```" + `go
func Relay(c *gin.Context) {
    // ... business logic ...
    if err != nil {
        newApiErr := types.NewError(err, types.ErrorCodeInvalidRequest)
        c.JSON(newApiErr.StatusCode, gin.H{
            "error": newApiErr.ToOpenAIError(),
        })
        return
    }
}
` + "```" + `

## HTTP Status Code Mapping

All error codes automatically map to appropriate HTTP status codes:

- **400 Bad Request**: Invalid request parameters, malformed request body
- **401 Unauthorized**: Invalid or missing authentication
- **403 Forbidden**: Access denied
- **404 Not Found**: Model or resource not found
- **402 Payment Required**: Insufficient quota
- **429 Too Many Requests**: Rate limit exceeded
- **500 Internal Server Error**: Internal system errors
- **502 Bad Gateway**: Invalid response from upstream
- **503 Service Unavailable**: No available channels
- **504 Gateway Timeout**: Response time exceeded

## Migration from Legacy Error Codes

If you're migrating from the old string-based error codes:

` + "```" + `go
// Old (deprecated)
types.ErrorCodeChannelNoAvailableKey  // "channel:no_available_key"

// New (numeric)
types.ErrorCodeChannelNoAvailableKey  // 3001
` + "```" + `

The constant names remain the same, but the type changed from string to int.
This provides better type safety and performance.

## See Also

- [Error Handling Improvements Proposal](./error-code-improvements.md)
- [types/error.go](../types/error.go) - Error implementation
- [types/error_code.go](../types/error_code.go) - Error code definitions
- [types/error_level.go](../types/error_level.go) - Error level definitions
- [types/error_i18n.go](../types/error_i18n.go) - Internationalization support

---

**Total Error Codes**: {{.TotalCount}}
**Last Modified**: {{.Timestamp}}
`

// CategoryData holds errors grouped by category
type CategoryData struct {
	Category string
	Errors   []ErrorDoc
}

// TemplateData holds data for the markdown template
type TemplateData struct {
	Timestamp  string
	Categories []CategoryData
	TotalCount int
}

func main() {
	// Generate current timestamp
	timestamp := "2026-02-26"

	// Group errors by category
	categoryMap := make(map[string][]ErrorDoc)
	for _, err := range errorCodes {
		categoryMap[err.Category] = append(categoryMap[err.Category], err)
	}

	// Convert to sorted categories
	var categories []CategoryData
	for cat, errs := range categoryMap {
		categories = append(categories, CategoryData{
			Category: cat,
			Errors:   errs,
		})
	}

	// Sort categories by name
	sort.Slice(categories, func(i, j int) bool {
		return categories[i].Category < categories[j].Category
	})

	// Sort errors within each category by code
	for i := range categories {
		sort.Slice(categories[i].Errors, func(j, k int) bool {
			return categories[i].Errors[j].Code < categories[i].Errors[k].Code
		})
	}

	// Prepare template data
	data := TemplateData{
		Timestamp:  timestamp,
		Categories: categories,
		TotalCount: len(errorCodes),
	}

	// Parse and execute template
	tmpl, err := template.New("errorcodes").Parse(markdownTemplate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing template: %v\n", err)
		os.Exit(1)
	}

	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing template: %v\n", err)
		os.Exit(1)
	}
}
