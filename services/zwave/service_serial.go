package zwave

import (
	"github.com/stas-makutin/howeve/defs"
)

// ServiceSerial zwave service implementation using serial transport
type ServiceSerial struct {
}

// NewServiceSerial creates new zwave service implementation using serial transport
func NewServiceSerial(entry string, params defs.ParamValues) (defs.Service, error) {
	return &ServiceSerial{}, nil
}

func (svc *ServiceSerial) Start() {
}

func (svc *ServiceSerial) Stop() {
}

func (svc *ServiceSerial) Status() defs.ServiceStatus {
	return defs.ServiceStatus{}
}

func (svc *ServiceSerial) Send(message defs.Message) error {
	return nil
}
