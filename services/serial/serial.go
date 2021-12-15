package serial

import (
	"sync"
	"time"

	serial "github.com/albenik/go-serial/v2"
	"github.com/stas-makutin/howeve/defs"
)

const readyToReadPoolInterval = 200 * time.Millisecond

// Transport struct - serial Transport implementation
type Transport struct {
	port   *serial.Port
	lock   sync.RWMutex
	wg     sync.WaitGroup
	stopCh chan struct{}
}

func (t *Transport) ID() defs.TransportIdentifier {
	return defs.TransportSerial
}

// Open func
func (t *Transport) Open(entry string, params defs.ParamValues) (err error) {
	options := []serial.Option{}
	if v, ok := params[ParamNameBaudRate]; ok {
		options = append(options, serial.WithBaudrate(v.(int)))
	}
	if v, ok := params[ParamNameDataBits]; ok {
		switch v {
		case "5":
			options = append(options, serial.WithDataBits(5))
		case "6":
			options = append(options, serial.WithDataBits(6))
		case "7":
			options = append(options, serial.WithDataBits(7))
		case "8":
			options = append(options, serial.WithDataBits(8))
		}
	}
	if v, ok := params[ParamNameParity]; ok {
		switch v {
		case "none":
			options = append(options, serial.WithParity(serial.NoParity))
		case "odd":
			options = append(options, serial.WithParity(serial.OddParity))
		case "even":
			options = append(options, serial.WithParity(serial.EvenParity))
		case "mark":
			options = append(options, serial.WithParity(serial.MarkParity))
		case "space":
			options = append(options, serial.WithParity(serial.SpaceParity))
		}
	}
	if v, ok := params[ParamNameStopBits]; ok {
		switch v {
		case "1":
			options = append(options, serial.WithStopBits(serial.OneStopBit))
		case "1.5":
			options = append(options, serial.WithStopBits(serial.OnePointFiveStopBits))
		case "2":
			options = append(options, serial.WithStopBits(serial.TwoStopBits))
		}
	}
	if v, ok := params[ParamNameReadTimeout]; ok {
		options = append(options, serial.WithReadTimeout(int(v.(uint32))))
	}
	if v, ok := params[ParamNameWriteTimeout]; ok {
		options = append(options, serial.WithWriteTimeout(int(v.(uint32))))
	}

	t.lock.Lock()
	defer t.lock.Unlock()
	t.close()
	t.port, err = serial.Open(entry, options...)
	if t.port != nil {
		t.stopCh = make(chan struct{}, 1)
	}
	return
}

// Close func
func (t *Transport) Close() error {
	t.lock.Lock()
	defer t.lock.Unlock()
	return t.close()
}

func (t *Transport) close() (err error) {
	defer t.wg.Wait()
	if t.port != nil {
		err = t.port.Close()
		t.port = nil
		close(t.stopCh)
		t.stopCh = nil
	}
	return
}

// ReadyToRead function, singal in the channel if something could be read from the port or port state has changed
func (t *Transport) ReadyToRead() <-chan struct{} {
	// stop previous if any
	t.lock.RLock()
	if t.stopCh != nil {
		select {
		case t.stopCh <- struct{}{}:
		default:
		}
	}
	t.lock.RUnlock()
	t.wg.Wait()

	rc := make(chan struct{})
	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		for {
			stopCh := func() <-chan struct{} {
				t.lock.RLock()
				defer t.lock.RUnlock()

				if t.port == nil {
					return nil
				}

				if n, err := t.port.ReadyToRead(); n > 0 || err != nil {
					return nil
				}

				return t.stopCh
			}()

			if stopCh == nil {
				close(rc)
				break
			}

			select {
			case <-rc:
				return
			case <-stopCh:
				close(rc)
				return
			case <-time.After(readyToReadPoolInterval):
			}
		}
	}()
	return rc
}

// Read func
func (t *Transport) Read(p []byte) (int, error) {
	t.lock.RLock()
	defer t.lock.RUnlock()
	if t.port != nil {
		return t.port.Read(p)
	}
	return 0, defs.ErrNotOpen
}

// Write func
func (t *Transport) Write(p []byte) (int, error) {
	t.lock.RLock()
	defer t.lock.RUnlock()
	if t.port != nil {
		return t.port.Write(p)
	}
	return 0, defs.ErrNotOpen
}
