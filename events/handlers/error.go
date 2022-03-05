package handlers

import (
	"errors"
	"fmt"

	"github.com/stas-makutin/howeve/api"
	"github.com/stas-makutin/howeve/defs"
)

// newErrorInfo - makes error information structure
func newErrorInfo(code api.ErrorCode, err error, args ...interface{}) (e *api.ErrorInfo) {
	e = &api.ErrorInfo{Code: code, Params: args, Err: err}
	switch code {
	case api.ErrorUnknownProtocol:
		e.Message = fmt.Sprintf("Unknown protocol identifier %d", args...)
	case api.ErrorUnknownTransport:
		e.Message = fmt.Sprintf("Unknown transport identifier %d", args...)
	case api.ErrorInvalidProtocolTransport:
		e.Message = fmt.Sprintf(
			"The protocol %s (%d) doesn't support the transport %s (%d)",
			defs.ProtocolName(args[0].(api.ProtocolIdentifier)), args[0],
			defs.TransportName(args[1].(api.TransportIdentifier)), args[1],
		)
	case api.ErrorUnknownParameter:
		e.Message = fmt.Sprintf("Unknown parameter \"%s\"", args...)
	case api.ErrorInvalidParameterValue:
		e.Message = fmt.Sprintf("Unknown value \"%s\" of parameter \"%s\"", args...)
	case api.ErrorNoRequiredParameter:
		e.Message = fmt.Sprintf("Required parameter \"%s\" is missing", args...)
	case api.ErrorNoDiscovery:
		e.Message = fmt.Sprintf(
			"The discovery is not supported for %s (%d) protocol and %s (%d) transport",
			defs.ProtocolName(args[0].(api.ProtocolIdentifier)), args[0],
			defs.TransportName(args[1].(api.TransportIdentifier)), args[1],
		)
	case api.ErrorDiscoveryBusy:
		e.Message = "The discovery is not available - too many discovery quieries are executing at this moment"
	case api.ErrorNoDiscoveryID:
		e.Message = fmt.Sprintf("The discovery request %s not found", args...)
	case api.ErrorDiscoveryPending:
		e.Message = fmt.Sprintf("The discovery request %s not completed yet", args...)
	case api.ErrorDiscoveryFailed:
		e.Message = fmt.Sprintf("The discovery has failed, reason: %s", err.Error())
	case api.ErrorServiceNoKey:
		e.Message = "The service key fields (protocol, transport, entry) are required"
	case api.ErrorServiceNoID:
		e.Message = "Either service key fields (protocol, transport, entry) or service alias are required"
	case api.ErrorServiceExists:
		e.Message = fmt.Sprintf(
			"The service exists already for %s (%d) protocol, %s (%d) transport, and %s entry",
			defs.ProtocolName(args[0].(api.ProtocolIdentifier)), args[0],
			defs.TransportName(args[1].(api.TransportIdentifier)), args[1],
			args[2],
		)
	case api.ErrorServiceAliasExists:
		e.Message = fmt.Sprintf("The service's alias %s exists already", args...)
	case api.ErrorServiceInitialize:
		e.Message = fmt.Sprintf(
			"The service initialization failed, %s (%d) protocol, %s (%d) transport, and %s entry, reason: %s",
			defs.ProtocolName(args[0].(api.ProtocolIdentifier)), args[0],
			defs.TransportName(args[1].(api.TransportIdentifier)), args[1],
			args[2], err.Error(),
		)
	case api.ErrorServiceKeyNotExists:
		e.Message = fmt.Sprintf(
			"The service not exists for %s (%d) protocol, %s (%d) transport, and %s entry",
			defs.ProtocolName(args[0].(api.ProtocolIdentifier)), args[0],
			defs.TransportName(args[1].(api.TransportIdentifier)), args[1],
			args[2],
		)
	case api.ErrorServiceAliasNotExists:
		e.Message = fmt.Sprintf("The service with alias %s not exists", args...)
	case api.ErrorServiceStatusBad:
		e.Message = err.Error()
	case api.ErrorServiceBadPayload:
		e.Message = "The payload is not valid for the service"
	case api.ErrorServiceSendBusy:
		e.Message = "The service is too busy and unable to send the message"
	case api.ErrorOtherError:
		e.Message = err.Error()
	}
	return
}

func handleParamsErrors(err error) *api.ErrorInfo {
	var pe *defs.ParseError
	if errors.As(err, &pe) {
		switch pe.Code {
		case defs.UnknownParamName:
			return newErrorInfo(api.ErrorUnknownParameter, err, pe.Name)
		case defs.NoRequiredParam:
			return newErrorInfo(api.ErrorNoRequiredParameter, err, pe.Name)
		}
		return newErrorInfo(api.ErrorInvalidParameterValue, err, pe.Value, pe.Name)
	}
	return nil
}

func handleProtocolErrors(err error, protocol api.ProtocolIdentifier, transport api.TransportIdentifier) *api.ErrorInfo {
	switch err {
	case defs.ErrTransportNotSupported:
		return newErrorInfo(api.ErrorUnknownTransport, err, transport)
	case defs.ErrProtocolNotSupported:
		return newErrorInfo(api.ErrorUnknownProtocol, err, protocol)
	case defs.ErrNoProtocolTransport:
		return newErrorInfo(api.ErrorInvalidProtocolTransport, err, protocol, transport)
	}
	return nil
}

func handleServiceNotExistsError(key *api.ServiceKey, alias string) *api.ErrorInfo {
	if key != nil {
		return newErrorInfo(api.ErrorServiceKeyNotExists, defs.ErrServiceNotExists, key.Protocol, key.Transport, key.Entry)
	}
	return newErrorInfo(api.ErrorServiceAliasNotExists, defs.ErrServiceNotExists, alias)
}
