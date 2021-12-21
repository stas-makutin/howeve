package zwave

import (
	"context"
	"encoding/hex"
	"errors"
	"sync"
	"time"

	"github.com/stas-makutin/howeve/defs"
	"github.com/stas-makutin/howeve/log"
	zw "github.com/stas-makutin/howeve/zwave"
)

// log constants
const (
	// operation
	zwOpService = "ZW"

	zwOcTransportOpen  = "O"
	zwOcTransportRead  = "R"
	zwOcTransportWrite = "W"
	zwOcReplyTimeout   = "T"
	zwOcDataFrame      = "D"
	zwOcUnknownFrame   = "U"

	zwOsSuccess       = "0"
	zwOsFailure       = "F"
	zwOsWrongLength   = "L"
	zwOsWrongChecksum = "C"

	zwOfWriteQueue = "Q"
	zwOfWriteReply = "R"
)

// Service ZWave service implementation
type Service struct {
	transport defs.Transport
	key       *defs.ServiceKey
	params    defs.ParamValues

	sendQueue chan *defs.Message

	ctx    context.Context
	cancel context.CancelFunc
	stopWg sync.WaitGroup
}

// NewServiceSerial creates new zwave service implementation using serial transport
func NewService(transport defs.Transport, entry string, params defs.ParamValues) (defs.Service, error) {
	return &Service{
		transport: transport,
		key:       &defs.ServiceKey{Protocol: defs.ProtocolZWave, Transport: transport.ID(), Entry: entry},
		params:    params,
		sendQueue: make(chan *defs.Message, 10),
	}, nil
}

func (svc *Service) Start() {
	svc.Stop()

	svc.ctx, svc.cancel = context.WithCancel(context.Background())

	svc.stopWg.Add(1)
	go svc.serviceLoop()
}

func (svc *Service) Stop() {
	if svc.ctx == nil {
		return // already stopped
	}

	svc.cancel()
	svc.stopWg.Wait()

	svc.ctx, svc.cancel = nil, nil
}

func (svc *Service) Status() defs.ServiceStatus {
	return defs.ServiceStatus{}
}

func (svc *Service) Send(payload []byte) (*defs.Message, error) {
	if len(payload) <= 0 || len(payload) > 255 {
		return nil, defs.ErrBadPayload
	}

	message := defs.Messages.Register(svc.key, payload, defs.OutgoingPending)

	select {
	default:
		defs.Messages.UpdateState(message.ID, defs.OutgoingRejected)
		return message, defs.ErrSendBusy
	case svc.sendQueue <- message:
	}
	return message, nil
}

func (svc *Service) openTimeout() time.Duration {
	openTimeout := time.Millisecond * 5000
	if v, ok := svc.params[defs.ParamNameOpenAttemptsInterval]; ok {
		openTimeout = time.Duration(v.(uint32)) * time.Millisecond
		if openTimeout <= 0 {
			openTimeout = 100 * time.Millisecond
		}
	}
	return openTimeout
}

func (svc *Service) log(op string, fields ...string) {
	args := make([]string, 0, len(fields)+6)
	args = append(args, log.SrcSVC)
	args = append(args, zwOpService)
	args = append(args, op)
	args = append(args, defs.ProtocolName(svc.key.Protocol))
	args = append(args, defs.TransportName(svc.key.Transport))
	args = append(args, svc.key.Entry)
	args = append(args, fields...)
	log.Report(args...)
}

func (svc *Service) serviceLoop() {
	defer svc.transport.Close()
	defer svc.stopWg.Done()

	openTimeout := svc.openTimeout()
	open := true
	expectReply := false
	buffer := make([]byte, 4096)
	rb, re := 0, 0

ServiceLoop:
	for {
		if open {
			open = false
			expectReply = false
			rb = re
			if err := svc.transport.Open(svc.key.Entry, svc.params); err != nil {
				svc.log(zwOcTransportOpen, zwOsFailure, err.Error())
				select {
				case <-svc.ctx.Done():
					break ServiceLoop
				case <-time.After(openTimeout):
					open = true
					continue
				}
			} else {
				svc.log(zwOcTransportOpen, zwOsSuccess)
			}
		}

		if expectReply {
			expectReply = false
			select {
			case <-svc.ctx.Done():
				break ServiceLoop
			case <-time.After(time.Millisecond * 1500):
				svc.log(zwOcReplyTimeout)
				continue
			case <-svc.transport.ReadyToRead():
			}
		} else {
			select {
			case <-svc.ctx.Done():
				break ServiceLoop
			case message := <-svc.sendQueue:
				if n, err := svc.transport.Write(message.Payload); err != nil || n != len(message.Payload) {
					svc.log(zwOcTransportWrite, zwOsFailure, zwOfWriteQueue, err.Error())
					defs.Messages.UpdateState(message.ID, defs.OutgoingFailed)
					open = true
				} else {
					defs.Messages.UpdateState(message.ID, defs.Outgoing)
					if vr, _ := zw.ValidateDataFrame(message.Payload); vr == zw.FrameOK || vr == zw.FrameWrongChecksum {
						expectReply = true
					}
				}
				continue
			case <-svc.transport.ReadyToRead():
			}
		}

		if rb == re {
			rb, re = 0, 0
		}
		n, err := svc.transport.Read(buffer[re:])
		if errors.Is(svc.ctx.Err(), context.Canceled) {
			break ServiceLoop
		}
		if err != nil {
			svc.log(zwOcTransportRead, zwOsFailure, err.Error())
			open = true
		} else if n > 0 {
			re += n

			// read loop
		ReadLoop:
			for rb < re {
				switch buffer[rb] {
				case zw.FrameASK, zw.FrameNAK, zw.FrameCAN:
					defs.Messages.Register(svc.key, buffer[rb:rb+1], defs.Incoming)
					rb++

				case zw.FrameSOF:
					var reply []byte

					switch vr, pos := zw.ValidateDataFrame(buffer[rb:re]); vr {
					case zw.FrameOK:
						reply = []byte{zw.FrameASK}
						defs.Messages.Register(svc.key, buffer[rb:rb+pos], defs.Incoming)
						rb += pos

					case zw.FrameIncomplete:
						// continue to read
						break ReadLoop

					case zw.FrameWrongLength:
						svc.log(zwOcDataFrame, zwOsWrongLength)
						reply = []byte{zw.FrameNAK}
						rb = re // reset reading indexes - ignore content of reading buffer

					case zw.FrameWrongChecksum:
						svc.log(zwOcDataFrame, zwOsWrongChecksum, hex.EncodeToString(buffer[rb:rb+pos]))
						reply = []byte{zw.FrameNAK}
						rb += pos
					}

					if len(reply) > 0 {
						message := defs.Messages.Register(svc.key, reply, defs.OutgoingPending)
						if n, err := svc.transport.Write(reply); err != nil || n != len(reply) {
							svc.log(zwOcTransportWrite, zwOsFailure, zwOfWriteReply, err.Error())
							defs.Messages.UpdateState(message.ID, defs.OutgoingFailed)
							open = true
							break ReadLoop
						} else {
							defs.Messages.UpdateState(message.ID, defs.Outgoing)
						}
					}
				default:
					rz := rb + 20
					if rz > re {
						rz = re
					}
					svc.log(zwOcUnknownFrame, hex.EncodeToString(buffer[rb:rz]))
					rb = re // reset reading indexes - ignore content of reading buffer
					break ReadLoop
				}
			}
		}
	}
}
