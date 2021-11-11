package messages

import (
	"bytes"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stas-makutin/howeve/defs"
)

func TestMessagesPersistence(t *testing.T) {

	m := newMessages()

	t.Run("Load from non-existing file", func(t *testing.T) {
		if err := m.load("messages-non-existing-file"); err != nil {
			t.Error(err)
		}
	})

	file, err := os.CreateTemp("", "messages_*.log")
	t.Run("Create temporary message log file", func(t *testing.T) {
		if err != nil {
			t.Errorf("unable to create temporary message log file: %v", err)
		}
	})

	fileName := file.Name()
	defer os.Remove(fileName)
	file.Close()

	t.Run("Load from empty file", func(t *testing.T) {
		err := m.load(fileName)
		if err == nil {
			t.Error("empty file must not be valid, the file must contain the header at least")
		} else if err.Error() != "the message log file corrupted: header" {
			t.Error(err)
		}
	})

	services := []*defs.ServiceKey{
		{Protocol: defs.ProtocolZWave, Transport: defs.TransportSerial, Entry: "COM1"},
		{Protocol: defs.ProtocolZWave, Transport: defs.TransportSerial, Entry: "COM3"},
	}
	msgCount := 50 + rand.Intn(100)

	for msgCount > 0 {
		service := services[rand.Intn(100)%2]
		payloadLen := 12 + rand.Intn(200)
		payload := make([]byte, payloadLen)
		rand.Read(payload)

		m.push(service, &defs.Message{UUID: uuid.New(), Payload: payload})

		msgCount--
	}

	if !t.Run("Save log file", func(t *testing.T) {
		if err := m.save(fileName, 0644, 0755); err != nil {
			t.Error(err)
		}
	}) {
		return
	}

	m2 := newMessages()

	if !t.Run("Load log file", func(t *testing.T) {
		if err := m2.load(fileName); err != nil {
			t.Error(err)
			return
		}
	}) {
		return
	}

	t.Run("Verify loaded log file", func(t *testing.T) {
		if len(m.services) != len(m2.services) {
			t.Errorf("number of services is different: %d vs %d", len(m.services), len(m2.services))
		} else {
			for k, c := range m.services {
				if c2, ok := m2.services[k]; !ok {
					t.Errorf("service not found: %v", k)
				} else if c != c2 {
					t.Errorf("service %v usage count is different: %v vs %v", k, c, c2)
				}
			}
		}
		if len(m.entries) != len(m2.entries) {
			t.Errorf("number of entries is different: %d vs %d", len(m.entries), len(m2.entries))
		} else {
			for i := 0; i < len(m.entries); i++ {
				entry := m.entries[i]
				entry2 := m2.entries[i]

				if entry.time != entry2.time {
					t.Errorf("message %d: time is different: %s vs %s", i, entry.time.Format(time.RFC3339), entry2.time.Format(time.RFC3339))
				}
				if entry.UUID != entry2.UUID {
					t.Errorf("message %d: uuid is different: %s vs %s", i, entry.UUID, entry2.UUID)
				}
				if !bytes.Equal(entry.Payload, entry2.Payload) {
					t.Errorf("message %d: payload is different", i)
				}
			}
		}
	})
}
