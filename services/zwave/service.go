package zwave

import (
	"context"

	"github.com/stas-makutin/howeve/defs"
)

// Service ZWave service implementation
type Service struct {
	transport defs.Transport
	entry     string
	params    defs.ParamValues

	ctx    context.Context
	cancel context.CancelFunc
	stopCh chan struct{}
}

// NewServiceSerial creates new zwave service implementation using serial transport
func NewService(transport defs.Transport, entry string, params defs.ParamValues) (defs.Service, error) {
	return &Service{
		transport: transport,
		entry:     entry,
		params:    params,
	}, nil
}

func (svc *Service) Start() {
	svc.Stop()

	svc.ctx, svc.cancel = context.WithCancel(context.Background())
	svc.stopCh = make(chan struct{})

	go svc.serviceLoop()
}

func (svc *Service) Stop() {
	if svc.ctx == nil {
		return // already stopped
	}

	svc.cancel()
	<-svc.stopCh

	svc.ctx = nil
	svc.cancel = nil
	svc.stopCh = nil
}

func (svc *Service) Status() defs.ServiceStatus {
	return defs.ServiceStatus{}
}

func (svc *Service) Send(message defs.Message) error {
	return nil
}

func (svc *Service) serviceLoop() {
	defer close(svc.stopCh)
	for {
		<-svc.ctx.Done()
	}
}
