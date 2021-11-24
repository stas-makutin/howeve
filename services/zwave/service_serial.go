package zwave

import (
	"context"

	"github.com/stas-makutin/howeve/defs"
)

// ServiceSerial zwave service implementation using serial transport
type ServiceSerial struct {
	ctx context.Context
}

// NewServiceSerial creates new zwave service implementation using serial transport
func NewServiceSerial(ctx context.Context, entry string, params defs.ParamValues) (*defs.Service, error) {
	return nil, nil
}

func (svc *ServiceSerial) Start() error {
	return nil
}

func (svc *ServiceSerial) Stop() {
}

func (svc *ServiceSerial) Send(message defs.Message) error {
	return nil
}
