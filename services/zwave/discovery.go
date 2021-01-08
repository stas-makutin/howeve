package zwave

import (
	"context"

	"github.com/albenik/go-serial/v2/enumerator"
	"github.com/stas-makutin/howeve/defs"
)

// DiscoverySerial - discover COM ports
func DiscoverySerial(ctx context.Context, params defs.ParamValues) ([]*defs.ServiceEntryDetails, error) {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		return nil, err
	}
	if len(ports) <= 0 {
		return nil, nil
	}
	se := make([]*defs.ServiceEntryDetails, 0, len(ports))
	for _, port := range ports {
		se = append(se, &defs.ServiceEntryDetails{
			ServiceEntry: defs.ServiceEntry{
				Key: defs.ServiceKey{
					Protocol:  defs.ProtocolZWave,
					Transport: defs.TransportSerial,
					Entry:     port.Name,
				},
			},
			Description: port.Product,
		})
	}
	return se, nil
}
