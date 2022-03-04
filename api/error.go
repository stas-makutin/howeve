package api

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
	ErrorServiceStatusBad
	ErrorServiceBadPayload
	ErrorServiceSendBusy
	ErrorOtherError
)

// ErrorInfo - error
type ErrorInfo struct {
	Code    ErrorCode     `json:"c"`
	Message string        `json:"m"`
	Params  []interface{} `json:"p,omitempty"`
	Err     error         `json:"-"`
}

func (e *ErrorInfo) Error() string {
	return e.Message
}

func (e *ErrorInfo) Unwrap() error {
	return e.Err
}
