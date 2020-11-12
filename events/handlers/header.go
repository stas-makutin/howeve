package handlers

import (
	"strconv"
	"sync/atomic"

	"github.com/stas-makutin/howeve/events"
)

// Ordinal is atomically incremented number
type Ordinal uint32

// EventOrdinal is atomically incremented event number
var EventOrdinal Ordinal

// Next function returns next ordinal number
func (o *Ordinal) Next() Ordinal {
	return Ordinal(atomic.AddUint32((*uint32)(o), 1))
}

func (o Ordinal) String() string {
	return strconv.FormatUint(uint64(o), 36)
}

// TraceHeader interface
type TraceHeader interface {
	Ordinal() Ordinal
	TraceID() string
}

// TraceFlow interface
type TraceFlow interface {
	Associate() TraceFlow
}

// Header - event header struct
type Header struct {
	ordinal Ordinal
	traceID string
}

// NewHeader function creates new header and allocates next ordinal
func NewHeader(traceID string) *Header {
	return &Header{ordinal: EventOrdinal.Next(), traceID: traceID}
}

// Ordinal - implementation of TraceHeader
func (h Header) Ordinal() Ordinal {
	return h.ordinal
}

// TraceID - implementation of TraceHeader
func (h Header) TraceID() string {
	return h.traceID
}

// Associate - implementation of TraceFlow
func (h *Header) Associate() Header {
	return Header{ordinal: h.Ordinal(), traceID: h.TraceID()}
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

// Associate - implementation of TraceFlow
func (h *RequestHeader) Associate() ResponseHeader {
	return ResponseHeader{Header: h.Header.Associate(), ResponseTarget: h.ResponseTarget()}
}

// ResponseHeader - combined header of the response in the request-response pair
type ResponseHeader struct {
	Header
	events.ResponseTarget
}
