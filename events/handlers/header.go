package handlers

import (
	"sync/atomic"

	"github.com/stas-makutin/howeve/events"
)

// Ordinal is atomically incremented number
type Ordinal uint32

var eventOrdinal Ordinal

// Next function returns next ordinal number
func (o *Ordinal) Next() Ordinal {
	return Ordinal(atomic.AddUint32((*uint32)(o), 1))
}

// TraceHeader interface
type TraceHeader interface {
	Associate() TraceHeader
}

// Header - event header struct
type Header struct {
	Ordinal Ordinal
	ID      string
}

// NewHeader function creates new header and allocates next ordinal
func NewHeader(ID string) *Header {
	return &Header{Ordinal: eventOrdinal.Next(), ID: ID}
}

// Associate - implementation of TraceHeader
func (h *Header) Associate() Header {
	return Header{Ordinal: h.Ordinal, ID: h.ID}
}

// RequestHeader - combined header of the request in the request-response pair
type RequestHeader struct {
	Header
	events.RequestTarget
}

// NewRequestHeader function creates new request header and allocates next ordinal
func NewRequestHeader(ID string) *RequestHeader {
	return &RequestHeader{Header: *NewHeader(ID)}
}

// Associate - implementation of TraceHeader
func (h *RequestHeader) Associate() ResponseHeader {
	return ResponseHeader{Header: h.Header.Associate(), ResponseTarget: h.ResponseTarget()}
}

// ResponseHeader - combined header of the response in the request-response pair
type ResponseHeader struct {
	Header
	events.ResponseTarget
}
