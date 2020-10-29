package handlers

import (
	"sync/atomic"

	"github.com/stas-makutin/howeve/events"
)

// Ordinal is atomically incremented number
type Ordinal uint32

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

// Associate - implementation of TraceHeader
func (h *Header) Associate() Header {
	return Header{Ordinal: h.Ordinal, ID: h.ID}
}

// RequestHeader - combined header of the request in the request-response pair
type RequestHeader struct {
	Header
	events.RequestTarget
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
