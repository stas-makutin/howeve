package handlers

// ServiceID
type ServiceID struct {
	*ServiceKey
	Alias string `json:"alias,omitempty"`
}

// ServiceEntryWithAlias - service entry with alias
type ServiceEntryWithAlias struct {
	ServiceEntry
	Alias string `json:"alias,omitempty"`
}

// AddService - add new service
type AddService struct {
	RequestHeader
	*ServiceEntryWithAlias
}

// AddServiceReply - add new service reply
type AddServiceReply struct {
	Error   *ErrorInfo `json:"error,omitempty"`
	Success bool       `json:"success,omitempty"`
}

// AddServiceResult - add new service result
type AddServiceResult struct {
	ResponseHeader
	*AddServiceReply
}

// SendToService - send message to service
type SendToService struct {
	RequestHeader
	*ServiceID
}

// SendToServiceResult - send message to service result
type SendToServiceResult struct {
	ResponseHeader
}

// RetriveFromServiceResult - retrieve message(s) from service
type RetriveFromService struct {
	RequestHeader
}

// RetriveFromServiceResult - retrieve message(s) from service result
type RetriveFromServiceResult struct {
	ResponseHeader
}
