package handlers

import "fmt"

// ErrorCode type
type ErrorCode int

// Error codes
const (
	ErrorUnknownProtocol = ErrorCode(iota + 1)
	ErrorUnknownTransport
	ErrorInvalidProtocolTransport
	ErrorUnknownParameter
	ErrorInvalidParameterValue
	ErrorNoRequiredParameter
	ErrorNoDiscovery
)

// ErrorInfo - error
type ErrorInfo struct {
	Code    ErrorCode     `json:"c"`
	Message string        `json:"m"`
	Params  []interface{} `json:"p,omitempty"`
}

// NewErrorInfo - makes error information structure
func NewErrorInfo(code ErrorCode, args ...interface{}) (e *ErrorInfo) {
	e = &ErrorInfo{Code: code, Params: args}
	switch code {
	case ErrorUnknownProtocol:
		e.Message = fmt.Sprintf("Unknown protocol identifier %d", args...)
	case ErrorUnknownTransport:
		e.Message = fmt.Sprintf("Unknown transport identifier %d", args...)
	case ErrorInvalidProtocolTransport:
		e.Message = fmt.Sprintf("The protocol %s (%d) doesn't support the transport %s (%d)", args...)
	case ErrorUnknownParameter:
		e.Message = fmt.Sprintf("Unknown parameter \"%s\"", args...)
	case ErrorInvalidParameterValue:
		e.Message = fmt.Sprintf("Unknown value \"%s\" of parameter \"%s\"", args...)
	case ErrorNoRequiredParameter:
		e.Message = fmt.Sprintf("Required parameter \"%s\" is missing", args...)
	case ErrorNoDiscovery:
		e.Message = fmt.Sprintf("The discovery is not supported for the protocol %s (%d) and the transport %s (%d)", args...)
	}
	return
}
