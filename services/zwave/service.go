package zwave

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/stas-makutin/howeve/defs"
	zw "github.com/stas-makutin/howeve/zwave"
)

// Service ZWave service implementation
type Service struct {
	transport defs.Transport
	entry     string
	params    defs.ParamValues

	sendQueue chan defs.Message

	ctx    context.Context
	cancel context.CancelFunc
	stopWg sync.WaitGroup
}

// NewServiceSerial creates new zwave service implementation using serial transport
func NewService(transport defs.Transport, entry string, params defs.ParamValues) (defs.Service, error) {
	return &Service{
		transport: transport,
		entry:     entry,
		params:    params,
		sendQueue: make(chan defs.Message, 10),
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

func (svc *Service) Send(message defs.Message) error {
	if len(message.Payload) <= 0 {
		return defs.ErrBadPayload
	}
	select {
	default:
		return defs.ErrSendBusy
	case svc.sendQueue <- message:
	}
	return nil
}

func (svc *Service) serviceLoop() {
	defer svc.transport.Close()
	defer svc.stopWg.Done()

	openTimeout := time.Millisecond * 5000
	if v, ok := svc.params[defs.ParamNameOpenAttemptsInterval]; ok {
		openTimeout = time.Duration(v.(uint32)) * time.Millisecond
		if openTimeout <= 0 {
			openTimeout = 100 * time.Millisecond
		}
	}

	open := true
	expectReply := false
	buffer := make([]byte, 4096)
	rb, re := 0, 0

ServiceLoop:
	for {
		if open {
			open = false
			expectReply = false
			rb, re = 0, 0
			if err := svc.transport.Open(svc.entry, svc.params); err != nil {
				// TODO error logging
				select {
				case <-svc.ctx.Done():
					break ServiceLoop
				case <-time.After(openTimeout):
					open = true
					continue
				}
			} else {
				// TODO log open port event
			}
		}

		if expectReply {
			expectReply = false
			select {
			case <-svc.ctx.Done():
				break ServiceLoop
			case <-time.After(time.Millisecond * 1500):
				// TODO error expected reply timeout event
				// TODO additional actions?
				continue
			case <-svc.transport.ReadyToRead():
			}
		} else {
			select {
			case <-svc.ctx.Done():
				break ServiceLoop
			case msg := <-svc.sendQueue:
				if n, err := svc.transport.Write(msg.Payload); err != nil || n != len(msg.Payload) {
					// TODO error logging
					open = true
				} else {
					// TODO update message status

					if len(msg.Payload) > 0 && msg.Payload[0] == zw.FrameSOF { // TODO - proper data frame validation?
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
			log.Println("Failed to read from port:", err)
			open = true
		} else if n > 0 {
			re += n

			// read loop
		ReadLoop:
			for rb < re {
				switch buffer[rb] {
				case zw.FrameASK, zw.FrameNAK, zw.FrameCAN:
					log.Printf("< % x\n", buffer[rb:rb+1])
					rb++
				case zw.FrameSOF:
					if rb+2 < re {
						l := int(buffer[rb+1])
						var reply []byte
						if l <= 2 || l > 253 {
							// log.Printf("< % x - wrong SOF length\n", buffer[rb:re])
							rb = re
							reply = []byte{zw.FrameNAK}
						} else if l+2 <= re-rb {
							if zw.Checksum(buffer[rb+1:rb+l+1]) == buffer[rb+l+1] {
								// TODO report new message

								rb += l + 2
								reply = []byte{zw.FrameASK}
							} else {
								// log.Printf("< % x - wrong SOF checksum\n", buffer[rb:re])
								rb = re
								reply = []byte{zw.FrameNAK}
							}
						} else {
							break ReadLoop // incomplete SOF frame
						}

						if len(reply) > 0 {
							// TODO report new message

							if n, err := svc.transport.Write(reply); err != nil || n != len(reply) {
								// TODO error logging
								open = true
								break ReadLoop
							} else {
								// TODO update message status
							}
						}
					} else {
						break ReadLoop // incomplete SOF frame
					}
				default:

					rb = re
					break ReadLoop
				}
			}
		}
	}

}
