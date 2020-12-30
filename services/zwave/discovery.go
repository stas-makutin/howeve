package zwave

import (
	"context"

	"github.com/albenik/go-serial/v2/enumerator"
	"github.com/stas-makutin/howeve/services/servicedef"
)

// DiscoverySerial - discover COM ports
func DiscoverySerial(ctx context.Context, params servicedef.ParamValues) ([]*servicedef.ServiceEntryDetails, error) {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		return nil, err
	}
	if len(ports) <= 0 {
		return nil, nil
	}
	se := make([]*servicedef.ServiceEntryDetails, 0, len(ports))
	for _, port := range ports {
		se = append(se, &servicedef.ServiceEntryDetails{
			ServiceEntry: servicedef.ServiceEntry{
				Key: servicedef.ServiceKey{
					Protocol:  servicedef.ProtocolZWave,
					Transport: servicedef.TransportSerial,
					Entry:     port.Name,
				},
			},
			Description: port.Product,
		})
	}
	return se, nil
}
