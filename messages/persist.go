package messages

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/stas-makutin/howeve/defs"
)

/*
	message log file format:
	all multibyte numbers use little endian byte order

	<fileHeader>
	<servicesCount> uint16
	<serviceEntries:
		protocol identifier uint8
		transport identifier uint8
		entry string length uint16
		entry string (utf-8)
	>
	<messageEntries:
		service index uint16
		time int64
		uuid [16]byte
		state byte
		payload length uint16
		payload
	>
*/

func (m *messages) load(file string, lengthLimit int) (int, error) {
	f, err := os.Open(file)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, fmt.Errorf("unable to open the message log file: %v", err)
	}
	defer f.Close()

	messages, length, err := readMessages(bufio.NewReader(f), lengthLimit)
	if err != nil {
		return 0, err
	} else {
		m.services = messages.services
		m.entries = messages.entries
	}

	return length, nil
}

func (m *messages) save(file string, dirMode os.FileMode, fileMode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(file), dirMode); err != nil {
		return err
	}
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fileMode)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	return writeMessages(w, m)
}

var fileHeader = []byte("ml01")
var minimalLength int = len(fileHeader) + 2 /* services count */

func serviceEntryLength(entryLength int) int {
	return 1 /* protocol id */ + 1 /* transport id */ + 2 /* entry length */ + entryLength
}

func messageEntryLength(payloadLength int) int {
	return 2 /* service index */ + 8 /* time */ + 16 /* UUID */ + 1 /* state */ + 2 /* payload size field */ + payloadLength
}

func readMessages(r io.Reader, lengthLimit int) (*messages, int, error) {
	if err := readHeader(r); err != nil {
		return nil, 0, err
	}
	messages := newMessages()
	serviceCount, err := readServicesCount(r)
	if err != nil {
		return nil, 0, err
	}
	length := minimalLength
	serviceIndices := make(map[uint16]*defs.ServiceKey)
	var serviceIndex uint16 = 0
	for ; serviceIndex < serviceCount; serviceIndex++ {
		if service, err := readServicesEntry(r, serviceIndex); err != nil {
			return nil, 0, err
		} else {
			messages.services[*service] = 0
			serviceIndices[serviceIndex] = service
			length += serviceEntryLength(len(service.Entry))
		}
	}

	for {
		if message, err := readMessageEntry(r, serviceIndices); err != nil {
			return nil, 0, err
		} else if message == nil {
			break
		} else {
			l := length + messageEntryLength(len(message.Payload))
			if lengthLimit > 0 && l > lengthLimit {
				break
			}
			length = l

			messages.entries = append(messages.entries, message)
			messages.services[*message.ServiceKey] += 1
		}
	}

	for service, count := range messages.services {
		if count <= 0 {
			delete(messages.services, service)
		}
	}

	return messages, length, nil
}

func writeMessages(w io.Writer, messages *messages) error {
	if err := writeHeader(w); err != nil {
		return err
	}

	serviceIndices := make(map[defs.ServiceKey]uint16)
	var serviceIndex uint16 = 0

	if err := writeServicesCount(w, uint16(len(messages.services))); err != nil {
		return err
	}

	for service, serviceMessagesCount := range messages.services {
		if serviceMessagesCount <= 0 {
			continue
		}
		if err := writeServicesEntry(w, service, serviceIndex); err != nil {
			return err
		}
		serviceIndices[service] = serviceIndex
		serviceIndex++
	}

	// messages
	for _, message := range messages.entries {
		if err := writeMessageEntry(w, message, serviceIndices); err != nil {
			return err
		}
	}

	return nil
}

func readError(name string, err error) error {
	if err == io.EOF || err == io.ErrUnexpectedEOF {
		return fmt.Errorf("the message log file corrupted: %s", name)
	}
	return fmt.Errorf("the message log file read failure: %s; %v", name, err)
}

func writeError(name string, err error) error {
	return fmt.Errorf("the message log file write failure: %s; %v", name, err)
}

func readHeader(r io.Reader) error {
	header := make([]byte, len(fileHeader))
	if _, err := io.ReadFull(r, header); err != nil {
		return readError("header", err)
	}
	if !bytes.Equal(header, fileHeader) {
		return fmt.Errorf("the message log header is not valid")
	}
	return nil
}

func writeHeader(w io.Writer) error {
	if _, err := w.Write(fileHeader); err != nil {
		return writeError("header", err)
	}
	return nil
}

func readServicesCount(r io.Reader) (uint16, error) {
	var serviceCount uint16
	if err := binary.Read(r, binary.LittleEndian, &serviceCount); err != nil {
		return 0, readError("services count", err)
	}
	return serviceCount, nil
}

func writeServicesCount(w io.Writer, servicesCount uint16) error {
	if err := binary.Write(w, binary.LittleEndian, uint16(servicesCount)); err != nil {
		return writeError("services count", err)
	}
	return nil
}

func readServicesEntry(r io.Reader, serviceIndex uint16) (*defs.ServiceKey, error) {
	var protocol defs.ProtocolIdentifier
	if err := binary.Read(r, binary.LittleEndian, &protocol); err != nil {
		return nil, readError(fmt.Sprintf("service %d protocol identifier", serviceIndex), err)
	}
	if !protocol.IsValid() {
		return nil, fmt.Errorf("service %d protocol identifier %d is not valid", serviceIndex, protocol)
	}
	var transport defs.TransportIdentifier
	if err := binary.Read(r, binary.LittleEndian, &transport); err != nil {
		return nil, readError(fmt.Sprintf("service %d transport identifier", serviceIndex), err)
	}
	if !transport.IsValid() {
		return nil, fmt.Errorf("service %d transport identifier %d is not valid", serviceIndex, transport)
	}
	var entryLen uint16
	if err := binary.Read(r, binary.LittleEndian, &entryLen); err != nil {
		return nil, readError(fmt.Sprintf("service %d entry length", serviceIndex), err)
	}
	entry := make([]byte, entryLen)
	if entryLen > 0 {
		if _, err := io.ReadFull(r, entry); err != nil {
			return nil, readError(fmt.Sprintf("service %d entry", serviceIndex), err)
		}
	}
	return &defs.ServiceKey{Transport: transport, Protocol: protocol, Entry: string(entry)}, nil
}

func writeServicesEntry(w io.Writer, service defs.ServiceKey, serviceIndex uint16) error {
	if err := binary.Write(w, binary.LittleEndian, service.Protocol); err != nil {
		return writeError(fmt.Sprintf("service %d protocol identifier", serviceIndex), err)
	}
	if err := binary.Write(w, binary.LittleEndian, service.Transport); err != nil {
		return writeError(fmt.Sprintf("service %d transport identifier", serviceIndex), err)
	}
	if err := binary.Write(w, binary.LittleEndian, uint16(len(service.Entry))); err != nil {
		return writeError(fmt.Sprintf("service %d entry length", serviceIndex), err)
	}
	if _, err := w.Write([]byte(service.Entry)); err != nil {
		return writeError(fmt.Sprintf("service %d entry", serviceIndex), err)
	}
	return nil
}

func readMessageEntry(r io.Reader, serviceIndices map[uint16]*defs.ServiceKey) (*message, error) {
	readMessageError := func(name string, err error) (*message, error) {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil, nil // ignore if no more messages or message is not complete
		}
		return nil, readError(name, err)
	}

	var serviceIndex uint16
	if err := binary.Read(r, binary.LittleEndian, &serviceIndex); err != nil {
		return readMessageError("message service index", err)
	}
	service, ok := serviceIndices[serviceIndex]
	if !ok {
		return nil, fmt.Errorf("message service index %d is not valid", serviceIndex)
	}
	var timeNano int64
	if err := binary.Read(r, binary.LittleEndian, &timeNano); err != nil {
		return readMessageError("message time", err)
	}
	var uuid [16]byte
	if err := binary.Read(r, binary.LittleEndian, &uuid); err != nil {
		return readMessageError("message UUID", err)
	}
	var state defs.MessageState
	if err := binary.Read(r, binary.LittleEndian, &state); err != nil {
		return readMessageError("message state", err)
	}
	var payloadLength uint16
	if err := binary.Read(r, binary.LittleEndian, &payloadLength); err != nil {
		return readMessageError("message payload length", err)
	}
	payload := make([]byte, payloadLength)
	if payloadLength > 0 {
		if _, err := io.ReadFull(r, payload); err != nil {
			return readMessageError("message payload", err)
		}
	}
	return &message{
		ServiceKey: service,
		Message: &defs.Message{
			Time:    time.Unix(0, timeNano).UTC(),
			UUID:    uuid,
			State:   state,
			Payload: payload,
		},
	}, nil
}

func writeMessageEntry(w io.Writer, m *message, serviceIndices map[defs.ServiceKey]uint16) error {
	serviceIndex, ok := serviceIndices[*m.ServiceKey]
	if !ok || len(m.Payload) > int(^uint16(0)) {
		return nil // ignore invalid services and large payoads
	}
	if err := binary.Write(w, binary.LittleEndian, serviceIndex); err != nil {
		return writeError("message service index", err)
	}
	if err := binary.Write(w, binary.LittleEndian, m.Time.UTC().UnixNano()); err != nil {
		return writeError("message time", err)
	}
	if _, err := w.Write(m.UUID[:]); err != nil {
		return writeError("message UUID", err)
	}
	if err := binary.Write(w, binary.LittleEndian, m.State); err != nil {
		return writeError("message state", err)
	}
	if err := binary.Write(w, binary.LittleEndian, uint16(len(m.Payload))); err != nil {
		return writeError("message payload length", err)
	}
	if _, err := w.Write(m.Payload); err != nil {
		return writeError("message payload", err)
	}
	return nil
}
