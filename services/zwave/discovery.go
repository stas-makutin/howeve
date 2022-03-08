package zwave

import (
	"context"
	"strings"
	"time"

	"github.com/albenik/go-serial/v2/enumerator"
	"github.com/stas-makutin/howeve/api"
	"github.com/stas-makutin/howeve/services/serial"
	zw "github.com/stas-makutin/howeve/zwave"
)

// serial parameters, must match with default service parameters
var discoverSerialParams api.ParamValues = api.ParamValues{
	serial.ParamNameBaudRate:     int32(115200),
	serial.ParamNameDataBits:     "8",
	serial.ParamNameParity:       "none",
	serial.ParamNameStopBits:     "1",
	serial.ParamNameReadTimeout:  uint32(0),
	serial.ParamNameWriteTimeout: uint32(0),
}

// ZW_Version ZWave serial API request data frame
var zwVersionFrame = zw.DataRequest([]byte{zw.ZW_VERSION})

// DiscoverSerial - discover COM ports with ZWave controllers
func DiscoverSerial(ctx context.Context, params api.ParamValues) ([]*api.DiscoveryEntry, error) {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		return nil, err
	}
	if len(ports) <= 0 {
		return nil, nil
	}

	var rv []*api.DiscoveryEntry
	for _, port := range ports {
		select {
		case <-ctx.Done():
			return nil, nil
		default:
		}
		if ok, info := discoverSerialPort(ctx, port.Name, discoverSerialParams); ok {
			if info != "" {
				info = " [" + info + "]"
			}
			rv = append(rv, &api.DiscoveryEntry{
				ServiceKey: api.ServiceKey{
					Protocol:  api.ProtocolZWave,
					Transport: api.TransportSerial,
					Entry:     port.Name,
				},
				Description: port.Product + info,
			})
		}
	}
	return rv, nil
}

func discoverSerialPort(ctx context.Context, port string, params api.ParamValues) (bool, string) {
	t := &serial.Transport{}

	if err := t.Open(port, params); err != nil {
		return false, ""
	}
	defer t.Close()

	select {
	case <-ctx.Done():
		return false, ""
	default:
	}

	if n, err := t.Write(zwVersionFrame); err != nil || n != len(zwVersionFrame) {
		return false, ""
	}

	buffer := make([]byte, 128)
	rb, re := 0, 0
	seq := 0
	for {
		select {
		case <-ctx.Done():
			return false, ""
		case <-time.After(time.Millisecond * 1500):
			return false, ""
		case <-t.ReadyToRead():
		}

		n, err := t.Read(buffer[re:])
		if err != nil || n <= 0 {
			return false, ""
		}
		re += n

	ReadLoop:
		for rb < re {
			switch seq {
			case 0:
				if buffer[rb] != zw.FrameASK {
					return false, ""
				}
				seq = 1
				rb++
			case 1:
				switch vr, _ := zw.ValidateDataFrame(buffer[rb:re]); vr {
				case zw.FrameOK:
					if payload := zw.UnpackResponse(buffer[rb:re]); payload != nil {
						// ZW_Version response
						if payload[0] == zw.ZW_VERSION && len(payload) == 14 {
							return true, strings.TrimRight(string(payload[1:13]), "\x00") // the library version, "Z-Wave x.yy"
						}
					}
					return false, ""
				case zw.FrameIncomplete:
					break ReadLoop
				}
				return false, ""
			}
		}
	}
}
