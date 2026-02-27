package types

import (
	"errors"
	"net/http"
	"testing"
)

// TestErrorCodeSystem verifies the new error code system works correctly
func TestErrorCodeSystem(t *testing.T) {
	tests := []struct {
		name              string
		errorCode         ErrorCode
		expectedHTTP      int
		expectedLevel     ErrorLevel
		expectedString    string
	}{
		{
			name:           "InvalidRequest",
			errorCode:      ErrorCodeInvalidRequest,
			expectedHTTP:   http.StatusBadRequest,
			expectedLevel:  ErrorLevelWarning,
			expectedString: "invalid_request",
		},
		{
			name:           "ChannelNoAvailableKey",
			errorCode:      ErrorCodeChannelNoAvailableKey,
			expectedHTTP:   http.StatusServiceUnavailable,
			expectedLevel:  ErrorLevelError,
			expectedString: "channel_no_available_key",
		},
		{
			name:           "InsufficientUserQuota",
			errorCode:      ErrorCodeInsufficientUserQuota,
			expectedHTTP:   http.StatusPaymentRequired,
			expectedLevel:  ErrorLevelWarning,
			expectedString: "insufficient_user_quota",
		},
		{
			name:           "DatabaseConnectionFailed",
			errorCode:      ErrorCodeDatabaseConnectionFailed,
			expectedHTTP:   http.StatusInternalServerError,
			expectedLevel:  ErrorLevelCritical,
			expectedString: "database_connection_failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test String() method
			if tt.errorCode.String() != tt.expectedString {
				t.Errorf("ErrorCode.String() = %v, want %v", tt.errorCode.String(), tt.expectedString)
			}

			// Test HTTPStatusCode() method
			if tt.errorCode.HTTPStatusCode() != tt.expectedHTTP {
				t.Errorf("ErrorCode.HTTPStatusCode() = %v, want %v", tt.errorCode.HTTPStatusCode(), tt.expectedHTTP)
			}

			// Test DefaultLevel() method
			if tt.errorCode.DefaultLevel() != tt.expectedLevel {
				t.Errorf("ErrorCode.DefaultLevel() = %v, want %v", tt.errorCode.DefaultLevel(), tt.expectedLevel)
			}
		})
	}
}

// TestNewErrorAutoMapping verifies that NewError auto-maps HTTP status and level
func TestNewErrorAutoMapping(t *testing.T) {
	err := NewError(errors.New("test error"), ErrorCodeInvalidRequest)

	if err.StatusCode != http.StatusBadRequest {
		t.Errorf("NewError() StatusCode = %v, want %v", err.StatusCode, http.StatusBadRequest)
	}

	if err.Level != ErrorLevelWarning {
		t.Errorf("NewError() Level = %v, want %v", err.Level, ErrorLevelWarning)
	}

	if err.errorCode != ErrorCodeInvalidRequest {
		t.Errorf("NewError() errorCode = %v, want %v", err.errorCode, ErrorCodeInvalidRequest)
	}
}

// TestErrorOptions verifies error options work correctly
func TestErrorOptions(t *testing.T) {
	err := NewError(
		errors.New("test error"),
		ErrorCodeGetChannelFailed,
		ErrOptionWithLevel(ErrorLevelCritical),
		ErrOptionWithSkipRetry(),
	)

	if err.Level != ErrorLevelCritical {
		t.Errorf("ErrOptionWithLevel() failed, got %v, want %v", err.Level, ErrorLevelCritical)
	}

	if !err.skipRetry {
		t.Error("ErrOptionWithSkipRetry() failed")
	}
}

// TestErrorCodeFromString verifies string to ErrorCode conversion
func TestErrorCodeFromString(t *testing.T) {
	tests := []struct {
		input       string
		expected    ErrorCode
	}{
		{"invalid_request", ErrorCodeInvalidRequest},
		{"channel_no_available_key", ErrorCodeChannelNoAvailableKey},
		{"insufficient_user_quota", ErrorCodeInsufficientUserQuota},
		{"unknown_code", ErrorCodeInvalidRequest}, // fallback
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ErrorCodeFromString(tt.input)
			if result != tt.expected {
				t.Errorf("ErrorCodeFromString(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestErrorLevelMethods verifies ErrorLevel methods
func TestErrorLevelMethods(t *testing.T) {
	tests := []struct {
		level         ErrorLevel
		expectedStr   string
		expectedValid bool
	}{
		{ErrorLevelInfo, "info", true},
		{ErrorLevelWarning, "warning", true},
		{ErrorLevelError, "error", true},
		{ErrorLevelCritical, "critical", true},
		{ErrorLevel(99), "unknown", false},
	}

	for _, tt := range tests {
		t.Run(tt.expectedStr, func(t *testing.T) {
			if tt.level.String() != tt.expectedStr {
				t.Errorf("ErrorLevel.String() = %v, want %v", tt.level.String(), tt.expectedStr)
			}

			if tt.level.IsValid() != tt.expectedValid {
				t.Errorf("ErrorLevel.IsValid() = %v, want %v", tt.level.IsValid(), tt.expectedValid)
			}
		})
	}
}

// TestIsChannelError verifies channel error detection
func TestIsChannelError(t *testing.T) {
	tests := []struct {
		name     string
		error    *NewAPIError
		expected bool
	}{
		{
			name:     "Channel error",
			error:    NewError(errors.New("test"), ErrorCodeChannelNoAvailableKey),
			expected: true,
		},
		{
			name:     "Non-channel error",
			error:    NewError(errors.New("test"), ErrorCodeInvalidRequest),
			expected: false,
		},
		{
			name:     "Nil error",
			error:    nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsChannelError(tt.error)
			if result != tt.expected {
				t.Errorf("IsChannelError() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestErrorLocalization verifies i18n support
func TestErrorLocalization(t *testing.T) {
	err := NewError(errors.New("test"), ErrorCodeInsufficientUserQuota)

	tests := []struct {
		lang          string
		expectedSubstr string
	}{
		{"en", "quota"},
		{"zh", "配额"},
		{"ja", "クォータ"},
	}

	for _, tt := range tests {
		t.Run(tt.lang, func(t *testing.T) {
			result := err.Localize(tt.lang)
			if !contains(result, tt.expectedSubstr) {
				t.Errorf("Localize(%q) = %v, want substring %v", tt.lang, result, tt.expectedSubstr)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || indexOf(s, substr) >= 0)
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
