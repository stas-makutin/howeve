package zwave

// Frame types
const (
	// Data frame Start of Frame (SOF)
	FrameSOF = 0x01
	// Acknowledge frame (data frame is accepted)
	FrameASK = 0x06
	// Not-acknowledge frame (data frame is not valid and not accepted)
	FrameNAK = 0x15
	// Data frame is valid but not accepted (race condition or other reasons)
	FrameCAN = 0x18
)

// Data frame structure:
//   SOF byte
//   Length byte
//   Data frame type byte
//   Serial API command ID
//   ... Serial API command parameters
//	 Checksum (calculates for all frame fields except SOF byte)

// Data frame type
const (
	// Request data frame
	FrameRequest = 0
	// Response data frame
	FrameResponse = 1
)

// Checksum calculates data frame checksum
func Checksum(src []byte) (rv byte) {
	rv = 0xff
	for _, v := range src {
		rv ^= v
	}
	return
}

// DataFrame creates valid data frame for provided frame type (request or response) and body (serial api command + parameters)
func DataFrame(frameType byte, body []byte) (frame []byte) {
	l := byte(len(body)) + 2
	frame = append([]byte{l, frameType}, body...)
	frame = append(frame, Checksum(frame))
	frame = append([]byte{FrameSOF}, frame...)
	return
}

// DataRequest creates request data frame for provided body (serial api command + parameters)
func DataRequest(body []byte) (frame []byte) {
	return DataFrame(FrameRequest, body)
}

// DataResponse creates response data frame for provided body (serial api command + parameters)
func DataResponse(body []byte) (frame []byte) {
	return DataFrame(FrameResponse, body)
}
