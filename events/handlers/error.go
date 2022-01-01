package handlers

import (
	"errors"
	"fmt"

	"github.com/stas-makutin/howeve/defs"
)

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
	ErrorDiscoveryBusy
	ErrorNoDiscoveryID
	ErrorDiscoveryPending
	ErrorDiscoveryFailed
	ErrorServiceNoKey
	ErrorServiceNoID
	ErrorServiceExists
	ErrorServiceAliasExists
	ErrorServiceInitialize
	ErrorServiceKeyNotExists
	ErrorServiceAliasNotExists
)

// ErrorInfo - error
type ErrorInfo struct {
	Code    ErrorCode     `json:"c"`
	Message string        `json:"m"`
	Params  []interface{} `json:"p,omitempty"`
	err     error
}

// newErrorInfo - makes error information structure
func newErrorInfo(code ErrorCode, err error, args ...interface{}) (e *ErrorInfo) {
	e = &ErrorInfo{Code: code, Params: args, err: err}
	switch code {
	case ErrorUnknownProtocol:
		e.Message = fmt.Sprintf("Unknown protocol identifier %d", args...)
	case ErrorUnknownTransport:
		e.Message = fmt.Sprintf("Unknown transport identifier %d", args...)
	case ErrorInvalidProtocolTransport:
		e.Message = fmt.Sprintf(
			"The protocol %s (%d) doesn't support the transport %s (%d)",
			defs.ProtocolName(args[0].(defs.ProtocolIdentifier)), args[0],
			defs.TransportName(args[1].(defs.TransportIdentifier)), args[1],
		)
	case ErrorUnknownParameter:
		e.Message = fmt.Sprintf("Unknown parameter \"%s\"", args...)
	case ErrorInvalidParameterValue:
		e.Message = fmt.Sprintf("Unknown value \"%s\" of parameter \"%s\"", args...)
	case ErrorNoRequiredParameter:
		e.Message = fmt.Sprintf("Required parameter \"%s\" is missing", args...)
	case ErrorNoDiscovery:
		e.Message = fmt.Sprintf(
			"The discovery is not supported for %s (%d) protocol and %s (%d) transport",
			defs.ProtocolName(args[0].(defs.ProtocolIdentifier)), args[0],
			defs.TransportName(args[1].(defs.TransportIdentifier)), args[1],
		)
	case ErrorDiscoveryBusy:
		e.Message = "The discovery is not available - too many discovery quieries are executing at this moment"
	case ErrorNoDiscoveryID:
		e.Message = fmt.Sprintf("The discovery request %s not found", args...)
	case ErrorDiscoveryPending:
		e.Message = fmt.Sprintf("The discovery request %s not completed yet", args...)
	case ErrorDiscoveryFailed:
		e.Message = fmt.Sprintf("The discovery has failed, reason: %s", err.Error())
	case ErrorServiceNoKey:
		e.Message = "The service key fields (protocol, transport, entry) are required"
	case ErrorServiceNoID:
		e.Message = "Either service key fields (protocol, transport, entry) or service alias are required"
	case ErrorServiceExists:
		e.Message = fmt.Sprintf(
			"The service exists already for %s (%d) protocol, %s (%d) transport, and %s entry",
			defs.ProtocolName(args[0].(defs.ProtocolIdentifier)), args[0],
			defs.TransportName(args[1].(defs.TransportIdentifier)), args[1],
			args[2],
		)
	case ErrorServiceAliasExists:
		e.Message = fmt.Sprintf("The service's alias %s exists already", args...)
	case ErrorServiceInitialize:
		e.Message = fmt.Sprintf(
			"The service initialization failed, %s (%d) protocol, %s (%d) transport, and %s entry, reason: %s",
			defs.ProtocolName(args[0].(defs.ProtocolIdentifier)), args[0],
			defs.TransportName(args[1].(defs.TransportIdentifier)), args[1],
			args[2], err.Error(),
		)
	case ErrorServiceKeyNotExists:
		e.Message = fmt.Sprintf(
			"The service not exists for %s (%d) protocol, %s (%d) transport, and %s entry",
			defs.ProtocolName(args[0].(defs.ProtocolIdentifier)), args[0],
			defs.TransportName(args[1].(defs.TransportIdentifier)), args[1],
			args[2],
		)
	case ErrorServiceAliasNotExists:
		e.Message = fmt.Sprintf("The service with alias %s not exists", args...)
	}
	return
}

func (e *ErrorInfo) Error() string {
	return e.Message
}

func (e *ErrorInfo) Unwrap() error {
	return e.err
}

func handleParamsErrors(err error) *ErrorInfo {
	var pe *defs.ParseError
	if errors.As(err, &pe) {
		switch pe.Code {
		case defs.UnknownParamName:
			return newErrorInfo(ErrorUnknownParameter, err, pe.Name)
		case defs.NoRequiredParam:
			return newErrorInfo(ErrorNoRequiredParameter, err, pe.Name)
		}
		return newErrorInfo(ErrorInvalidParameterValue, err, pe.Value, pe.Name)
	}
	return nil
}

func handleProtocolErrors(err error, protocol defs.ProtocolIdentifier, transport defs.TransportIdentifier) *ErrorInfo {
	switch err {
	case defs.ErrTransportNotSupported:
		return newErrorInfo(ErrorUnknownTransport, err, transport)
	case defs.ErrProtocolNotSupported:
		return newErrorInfo(ErrorUnknownProtocol, err, protocol)
	case defs.ErrNoProtocolTransport:
		return newErrorInfo(ErrorInvalidProtocolTransport, err, protocol, transport)
	}
	return nil
}

func handleServiceNotExistsError(key *defs.ServiceKey, alias string) *ErrorInfo {
	if key != nil {
		return newErrorInfo(ErrorServiceKeyNotExists, defs.ErrServiceNotExists, key.Protocol, key.Transport, key.Entry)
	}
	return newErrorInfo(ErrorServiceAliasNotExists, defs.ErrServiceNotExists, alias)
}
