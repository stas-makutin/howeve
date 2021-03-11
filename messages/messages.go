package messages

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/stas-makutin/howeve/defs"
)

var fileHeader = []byte("ml01")

// message struct defines entry in the message log
type message struct {
	time time.Time
	*defs.ServiceKey
	*defs.Message
}

// returns byte length of the message entry
func (m *message) length() int {
	return 8 /* time */ + 2 /* service index */ + len(m.UUID) + len(m.Payload) + 2 /* payload size field */
}

// messages log container
type messages struct {
	sync.Mutex
	services map[defs.ServiceKey]int
	entries  []*message
}

func (m *messages) load(file string) error {
	err := func() error {
		m.services = make(map[defs.ServiceKey]int)
		m.entries = nil

		f, err := os.Open(file)
		if err != nil {
			// TODO
		}
		defer f.Close()

		r := bufio.NewReader(f)

		header := make([]byte, len(fileHeader))
		if _, err := io.ReadFull(r, header); err != nil {
			if err == io.ErrUnexpectedEOF {
				// convert
			}
			return err
		}
		if !bytes.Equal(header, fileHeader) {
			// TODO wrong signature error
			return err
		}
		header = nil

		var serviceCount uint16
		if err := binary.Read(r, binary.LittleEndian, &serviceCount); err != nil {
			return err
		}
		serviceIndices := make(map[uint16]*defs.ServiceKey)
		var serviceIndex uint16 = 0
		for serviceCount > 0 {
			var protocol defs.ProtocolIdentifier
			if err := binary.Read(r, binary.LittleEndian, &protocol); err != nil {
				return err
			}
			var transport defs.TransportIdentifier
			if err := binary.Read(r, binary.LittleEndian, &transport); err != nil {
				return err
			}
			var entryLen uint16
			if err := binary.Read(r, binary.LittleEndian, &entryLen); err != nil {
				return err
			}
			entry := make([]byte, entryLen)
			if entryLen > 0 {
				if _, err := io.ReadFull(r, entry); err != nil {
					if err == io.ErrUnexpectedEOF {
						// convert
					}
					return err
				}
			}

			// TODO validate protocol
			// TODO validate transport

			key := defs.ServiceKey{Transport: transport, Protocol: protocol, Entry: string(entry)}
			m.services[key] = 0
			serviceIndices[serviceIndex] = &key
			serviceIndex++
			serviceCount--
		}

		for {

		}

		return nil
	}()
	if err != nil {

	}
	return err
}

func (m *messages) save(file string, dirMode os.FileMode, fileMode os.FileMode) error {
	err := func() error {
		if err := os.MkdirAll(filepath.Dir(file), dirMode); err != nil {
			return err
		}
		f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fileMode)
		if err != nil {
			return err
		}
		defer f.Close()
		w := bufio.NewWriter(f)

		// header
		if _, err := w.Write(fileHeader); err != nil {
			return err
		}

		// services
		serviceIndices := make(map[defs.ServiceKey]uint16)
		var serviceIndex uint16 = 0

		if err := binary.Write(w, binary.LittleEndian, uint16(len(m.services))); err != nil {
			return err
		}
		for k, v := range m.services {
			if v <= 0 {
				continue
			}
			if err := binary.Write(w, binary.LittleEndian, k.Protocol); err != nil {
				return err
			}
			if err := binary.Write(w, binary.LittleEndian, k.Transport); err != nil {
				return err
			}
			if err := binary.Write(w, binary.LittleEndian, uint16(len(k.Entry))); err != nil {
				return err
			}
			if _, err := w.Write([]byte(k.Entry)); err != nil {
				return err
			}
			serviceIndices[k] = serviceIndex
			serviceIndex++
		}

		// messages
		for _, msg := range m.entries {
			si, ok := serviceIndices[*msg.ServiceKey]
			if !ok {
				continue
			}
			if err := binary.Write(w, binary.LittleEndian, msg.time); err != nil {
				return err
			}
			if err := binary.Write(w, binary.LittleEndian, si); err != nil {
				return err
			}
			if _, err := w.Write(msg.UUID[:]); err != nil {
				return err
			}
			if err := binary.Write(w, binary.LittleEndian, uint16(len(msg.Payload))); err != nil {
				return err
			}
			if _, err := w.Write(msg.Payload); err != nil {
				return err
			}
		}

		w.Flush()

		return nil
	}()
	if err != nil {
		os.Remove(file)

		// todo error conversion
	}
	return err
}
