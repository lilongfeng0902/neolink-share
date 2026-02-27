package types

import (
	"fmt"
	"net/http"
)

// ErrorCode is a numeric error code for categorization and fast comparison
// This is the NEW error code system (numeric-based)
type ErrorCode int

// ErrorCodeString is the LEGACY string-based error code type
// Kept for backward compatibility during migration period
type ErrorCodeString string

// String returns the string representation of the error code
func (c ErrorCode) String() string {
	return errorCodeStrings[c]
}

// HTTPStatusCode returns the corresponding HTTP status code for the error
func (c ErrorCode) HTTPStatusCode() int {
	statusCode, ok := errorCodeHTTPStatusMap[c]
	if !ok {
		return http.StatusInternalServerError
	}
	return statusCode
}

// DefaultLevel returns the default error level for the error code
func (c ErrorCode) DefaultLevel() ErrorLevel {
	level, ok := errorCodeLevelMap[c]
	if !ok {
		return ErrorLevelError
	}
	return level
}

// Error code definitions
// Numeric ranges:
//   1xxx - General errors
//   2xxx - System errors
//   3xxx - Channel errors
//   4xxx - Client errors
//   5xxx - Upstream errors
//   6xxx - Database errors
//   7xxx - Quota errors
//   8xxx - Authentication errors (reserved)
//   9xxx - Miscellaneous errors (reserved)

const (
	// General Errors (1xxx)

	ErrorCodeInvalidRequest ErrorCode = 1001
	ErrorCodeSensitiveWordsDetected ErrorCode = 1002
	ErrorCodeViolationFeeGrokCSAM ErrorCode = 1003

	// System Errors (2xxx)

	ErrorCodeCountTokenFailed ErrorCode = 2001
	ErrorCodeModelPriceError ErrorCode = 2002
	ErrorCodeInvalidApiType ErrorCode = 2003
	ErrorCodeJsonMarshalFailed ErrorCode = 2004
	ErrorCodeJsonUnmarshalFailed ErrorCode = 2005
	ErrorCodeDoRequestFailed ErrorCode = 2006
	ErrorCodeGetChannelFailed ErrorCode = 2007
	ErrorCodeGenRelayInfoFailed ErrorCode = 2008

	// Channel Errors (3xxx)

	ErrorCodeChannelNoAvailableKey ErrorCode = 3001
	ErrorCodeChannelParamOverrideInvalid ErrorCode = 3002
	ErrorCodeChannelHeaderOverrideInvalid ErrorCode = 3003
	ErrorCodeChannelModelMappedError ErrorCode = 3004
	ErrorCodeChannelAwsClientError ErrorCode = 3005
	ErrorCodeChannelInvalidKey ErrorCode = 3006
	ErrorCodeChannelResponseTimeExceeded ErrorCode = 3007
	ErrorCodeChannelNotAvailable ErrorCode = 3008

	// Client Errors (4xxx)

	ErrorCodeReadRequestBodyFailed ErrorCode = 4001
	ErrorCodeConvertRequestFailed ErrorCode = 4002
	ErrorCodeAccessDenied ErrorCode = 4003
	ErrorCodeBadRequestBody ErrorCode = 4004
	ErrorCodeUnauthorized ErrorCode = 4005
	ErrorCodeForbidden ErrorCode = 4006

	// Upstream Errors (5xxx)

	ErrorCodeReadResponseBodyFailed ErrorCode = 5001
	ErrorCodeBadResponseStatusCode ErrorCode = 5002
	ErrorCodeBadResponse ErrorCode = 5003
	ErrorCodeBadResponseBody ErrorCode = 5004
	ErrorCodeEmptyResponse ErrorCode = 5005
	ErrorCodeAwsInvokeError ErrorCode = 5006
	ErrorCodeModelNotFound ErrorCode = 5007
	ErrorCodePromptBlocked ErrorCode = 5008
	ErrorCodeRateLimitExceeded ErrorCode = 5009
	ErrorCodeServiceUnavailable ErrorCode = 5010

	// Database Errors (6xxx)

	ErrorCodeQueryDataError ErrorCode = 6001
	ErrorCodeUpdateDataError ErrorCode = 6002
	ErrorCodeInsertDataError ErrorCode = 6003
	ErrorCodeDeleteDataError ErrorCode = 6004
	ErrorCodeDatabaseConnectionFailed ErrorCode = 6005

	// Quota Errors (7xxx)

	ErrorCodeInsufficientUserQuota ErrorCode = 7001
	ErrorCodePreConsumeTokenQuotaFailed ErrorCode = 7002
	ErrorCodeQuotaExceeded ErrorCode = 7003
)

// errorCodeStrings maps error codes to their string representations
var errorCodeStrings = map[ErrorCode]string{
	// General Errors (1xxx)
	ErrorCodeInvalidRequest:         "invalid_request",
	ErrorCodeSensitiveWordsDetected: "sensitive_words_detected",
	ErrorCodeViolationFeeGrokCSAM:   "violation_fee.grok_csam",

	// System Errors (2xxx)
	ErrorCodeCountTokenFailed:   "count_token_failed",
	ErrorCodeModelPriceError:    "model_price_error",
	ErrorCodeInvalidApiType:     "invalid_api_type",
	ErrorCodeJsonMarshalFailed:  "json_marshal_failed",
	ErrorCodeJsonUnmarshalFailed: "json_unmarshal_failed",
	ErrorCodeDoRequestFailed:    "do_request_failed",
	ErrorCodeGetChannelFailed:   "get_channel_failed",
	ErrorCodeGenRelayInfoFailed: "gen_relay_info_failed",

	// Channel Errors (3xxx)
	ErrorCodeChannelNoAvailableKey:        "channel_no_available_key",
	ErrorCodeChannelParamOverrideInvalid:  "channel_param_override_invalid",
	ErrorCodeChannelHeaderOverrideInvalid: "channel_header_override_invalid",
	ErrorCodeChannelModelMappedError:      "channel_model_mapped_error",
	ErrorCodeChannelAwsClientError:        "channel_aws_client_error",
	ErrorCodeChannelInvalidKey:            "channel_invalid_key",
	ErrorCodeChannelResponseTimeExceeded:  "channel_response_time_exceeded",
	ErrorCodeChannelNotAvailable:          "channel_not_available",

	// Client Errors (4xxx)
	ErrorCodeReadRequestBodyFailed: "read_request_body_failed",
	ErrorCodeConvertRequestFailed:  "convert_request_failed",
	ErrorCodeAccessDenied:          "access_denied",
	ErrorCodeBadRequestBody:        "bad_request_body",
	ErrorCodeUnauthorized:          "unauthorized",
	ErrorCodeForbidden:             "forbidden",

	// Upstream Errors (5xxx)
	ErrorCodeReadResponseBodyFailed: "read_response_body_failed",
	ErrorCodeBadResponseStatusCode:  "bad_response_status_code",
	ErrorCodeBadResponse:            "bad_response",
	ErrorCodeBadResponseBody:        "bad_response_body",
	ErrorCodeEmptyResponse:          "empty_response",
	ErrorCodeAwsInvokeError:         "aws_invoke_error",
	ErrorCodeModelNotFound:          "model_not_found",
	ErrorCodePromptBlocked:          "prompt_blocked",
	ErrorCodeRateLimitExceeded:      "rate_limit_exceeded",
	ErrorCodeServiceUnavailable:     "service_unavailable",

	// Database Errors (6xxx)
	ErrorCodeQueryDataError:         "query_data_error",
	ErrorCodeUpdateDataError:        "update_data_error",
	ErrorCodeInsertDataError:        "insert_data_error",
	ErrorCodeDeleteDataError:        "delete_data_error",
	ErrorCodeDatabaseConnectionFailed: "database_connection_failed",

	// Quota Errors (7xxx)
	ErrorCodeInsufficientUserQuota:      "insufficient_user_quota",
	ErrorCodePreConsumeTokenQuotaFailed: "pre_consume_token_quota_failed",
	ErrorCodeQuotaExceeded:              "quota_exceeded",
}

// errorCodeHTTPStatusMap maps error codes to their HTTP status codes
var errorCodeHTTPStatusMap = map[ErrorCode]int{
	// General Errors (1xxx)
	ErrorCodeInvalidRequest:         http.StatusBadRequest,
	ErrorCodeSensitiveWordsDetected: http.StatusBadRequest,
	ErrorCodeViolationFeeGrokCSAM:   http.StatusBadRequest,

	// System Errors (2xxx)
	ErrorCodeCountTokenFailed:   http.StatusInternalServerError,
	ErrorCodeModelPriceError:    http.StatusInternalServerError,
	ErrorCodeInvalidApiType:     http.StatusBadRequest,
	ErrorCodeJsonMarshalFailed:  http.StatusInternalServerError,
	ErrorCodeJsonUnmarshalFailed: http.StatusInternalServerError,
	ErrorCodeDoRequestFailed:    http.StatusInternalServerError,
	ErrorCodeGetChannelFailed:   http.StatusInternalServerError,
	ErrorCodeGenRelayInfoFailed: http.StatusInternalServerError,

	// Channel Errors (3xxx)
	ErrorCodeChannelNoAvailableKey:        http.StatusServiceUnavailable,
	ErrorCodeChannelParamOverrideInvalid:  http.StatusBadRequest,
	ErrorCodeChannelHeaderOverrideInvalid: http.StatusBadRequest,
	ErrorCodeChannelModelMappedError:      http.StatusInternalServerError,
	ErrorCodeChannelAwsClientError:        http.StatusInternalServerError,
	ErrorCodeChannelInvalidKey:            http.StatusUnauthorized,
	ErrorCodeChannelResponseTimeExceeded:  http.StatusGatewayTimeout,
	ErrorCodeChannelNotAvailable:          http.StatusServiceUnavailable,

	// Client Errors (4xxx)
	ErrorCodeReadRequestBodyFailed: http.StatusBadRequest,
	ErrorCodeConvertRequestFailed:  http.StatusBadRequest,
	ErrorCodeAccessDenied:          http.StatusUnauthorized,
	ErrorCodeBadRequestBody:        http.StatusBadRequest,
	ErrorCodeUnauthorized:          http.StatusUnauthorized,
	ErrorCodeForbidden:             http.StatusForbidden,

	// Upstream Errors (5xxx)
	ErrorCodeReadResponseBodyFailed: http.StatusInternalServerError,
	ErrorCodeBadResponseStatusCode:  http.StatusBadGateway,
	ErrorCodeBadResponse:            http.StatusBadGateway,
	ErrorCodeBadResponseBody:        http.StatusInternalServerError,
	ErrorCodeEmptyResponse:          http.StatusInternalServerError,
	ErrorCodeAwsInvokeError:         http.StatusInternalServerError,
	ErrorCodeModelNotFound:          http.StatusNotFound,
	ErrorCodePromptBlocked:          http.StatusBadRequest,
	ErrorCodeRateLimitExceeded:      http.StatusTooManyRequests,
	ErrorCodeServiceUnavailable:     http.StatusServiceUnavailable,

	// Database Errors (6xxx)
	ErrorCodeQueryDataError:         http.StatusInternalServerError,
	ErrorCodeUpdateDataError:        http.StatusInternalServerError,
	ErrorCodeInsertDataError:        http.StatusInternalServerError,
	ErrorCodeDeleteDataError:        http.StatusInternalServerError,
	ErrorCodeDatabaseConnectionFailed: http.StatusInternalServerError,

	// Quota Errors (7xxx)
	ErrorCodeInsufficientUserQuota:      http.StatusPaymentRequired,
	ErrorCodePreConsumeTokenQuotaFailed: http.StatusInternalServerError,
	ErrorCodeQuotaExceeded:              http.StatusPaymentRequired,
}

// errorCodeLevelMap maps error codes to their default severity levels
var errorCodeLevelMap = map[ErrorCode]ErrorLevel{
	// General Errors (1xxx)
	ErrorCodeInvalidRequest:         ErrorLevelWarning,
	ErrorCodeSensitiveWordsDetected: ErrorLevelWarning,
	ErrorCodeViolationFeeGrokCSAM:   ErrorLevelWarning,

	// System Errors (2xxx)
	ErrorCodeCountTokenFailed:   ErrorLevelError,
	ErrorCodeModelPriceError:    ErrorLevelError,
	ErrorCodeInvalidApiType:     ErrorLevelError,
	ErrorCodeJsonMarshalFailed:  ErrorLevelError,
	ErrorCodeJsonUnmarshalFailed: ErrorLevelError,
	ErrorCodeDoRequestFailed:    ErrorLevelError,
	ErrorCodeGetChannelFailed:   ErrorLevelCritical,
	ErrorCodeGenRelayInfoFailed: ErrorLevelError,

	// Channel Errors (3xxx)
	ErrorCodeChannelNoAvailableKey:        ErrorLevelError,
	ErrorCodeChannelParamOverrideInvalid:  ErrorLevelWarning,
	ErrorCodeChannelHeaderOverrideInvalid: ErrorLevelWarning,
	ErrorCodeChannelModelMappedError:      ErrorLevelError,
	ErrorCodeChannelAwsClientError:        ErrorLevelError,
	ErrorCodeChannelInvalidKey:            ErrorLevelWarning,
	ErrorCodeChannelResponseTimeExceeded:  ErrorLevelWarning,
	ErrorCodeChannelNotAvailable:          ErrorLevelCritical,

	// Client Errors (4xxx)
	ErrorCodeReadRequestBodyFailed: ErrorLevelWarning,
	ErrorCodeConvertRequestFailed:  ErrorLevelWarning,
	ErrorCodeAccessDenied:          ErrorLevelWarning,
	ErrorCodeBadRequestBody:        ErrorLevelWarning,
	ErrorCodeUnauthorized:          ErrorLevelWarning,
	ErrorCodeForbidden:             ErrorLevelWarning,

	// Upstream Errors (5xxx)
	ErrorCodeReadResponseBodyFailed: ErrorLevelError,
	ErrorCodeBadResponseStatusCode:  ErrorLevelError,
	ErrorCodeBadResponse:            ErrorLevelError,
	ErrorCodeBadResponseBody:        ErrorLevelError,
	ErrorCodeEmptyResponse:          ErrorLevelError,
	ErrorCodeAwsInvokeError:         ErrorLevelError,
	ErrorCodeModelNotFound:          ErrorLevelWarning,
	ErrorCodePromptBlocked:          ErrorLevelWarning,
	ErrorCodeRateLimitExceeded:      ErrorLevelWarning,
	ErrorCodeServiceUnavailable:     ErrorLevelCritical,

	// Database Errors (6xxx)
	ErrorCodeQueryDataError:         ErrorLevelCritical,
	ErrorCodeUpdateDataError:        ErrorLevelCritical,
	ErrorCodeInsertDataError:        ErrorLevelCritical,
	ErrorCodeDeleteDataError:        ErrorLevelCritical,
	ErrorCodeDatabaseConnectionFailed: ErrorLevelCritical,

	// Quota Errors (7xxx)
	ErrorCodeInsufficientUserQuota:      ErrorLevelWarning,
	ErrorCodePreConsumeTokenQuotaFailed: ErrorLevelError,
	ErrorCodeQuotaExceeded:              ErrorLevelWarning,
}

// ErrorCodeFromString converts a string representation to an ErrorCode
// Returns ErrorCodeInvalidRequest if not found
func ErrorCodeFromString(s string) ErrorCode {
	for code, str := range errorCodeStrings {
		if str == s {
			return code
		}
	}
	return ErrorCodeInvalidRequest
}

// IsValid checks if the error code is valid
func (c ErrorCode) IsValid() bool {
	_, ok := errorCodeStrings[c]
	return ok
}

// Format implements fmt.Formatter interface
func (c ErrorCode) Format(f fmt.State, verb rune) {
	switch verb {
	case 's', 'v':
		fmt.Fprint(f, c.String())
	case 'd':
		fmt.Fprint(f, int(c))
	default:
		fmt.Fprint(f, c.String())
	}
}
