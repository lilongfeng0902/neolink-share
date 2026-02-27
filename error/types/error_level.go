package types

// ErrorLevel represents the severity level of an error
type ErrorLevel int

const (
	// ErrorLevelInfo indicates informational messages that don't require action
	ErrorLevelInfo ErrorLevel = iota

	// ErrorLevelWarning indicates potentially harmful situations or minor issues
	ErrorLevelWarning

	// ErrorLevelError indicates error events that might still allow the application to continue
	ErrorLevelError

	// ErrorLevelCritical indicates critical error events that might cause the application to terminate
	ErrorLevelCritical
)

// String returns the string representation of the error level
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

// Color returns the ANSI color code for the error level
func (l ErrorLevel) Color() string {
	switch l {
	case ErrorLevelInfo:
		return "\033[36m" // Cyan
	case ErrorLevelWarning:
		return "\033[33m" // Yellow
	case ErrorLevelError:
		return "\033[31m" // Red
	case ErrorLevelCritical:
		return "\033[35m" // Magenta
	default:
		return "\033[0m" // Reset
	}
}

// ResetColor returns the ANSI reset code
func (l ErrorLevel) ResetColor() string {
	return "\033[0m"
}

// IsValid checks if the error level is valid
func (l ErrorLevel) IsValid() bool {
	return l >= ErrorLevelInfo && l <= ErrorLevelCritical
}

// MarshalJSON implements json.Marshaler interface
func (l ErrorLevel) MarshalJSON() ([]byte, error) {
	return []byte(`"` + l.String() + `"`), nil
}

// UnmarshalJSON implements json.Unmarshaler interface
func (l *ErrorLevel) UnmarshalJSON(data []byte) error {
	str := string(data)
	switch str {
	case `"info"`:
		*l = ErrorLevelInfo
	case `"warning"`:
		*l = ErrorLevelWarning
	case `"error"`:
		*l = ErrorLevelError
	case `"critical"`:
		*l = ErrorLevelCritical
	default:
		*l = ErrorLevelInfo // default
	}
	return nil
}
