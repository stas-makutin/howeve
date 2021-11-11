package messages

import (
	"math/rand"
	"os"
	"testing"

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
			t.Errorf("Unable to create temporary message log file: %v", err)
		}
	})

	fileName := file.Name()
	defer os.Remove(fileName)
	file.Close()

	t.Run("Load from empty file", func(t *testing.T) {
		err := m.load(fileName)
		if err == nil {
			t.Error("Empty file must not be valid, the file must contain the header at least")
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
		service := services[rand.Intn(1)]
		payloadLen := 12 + rand.Intn(200)
		payload := make([]byte, payloadLen)
		rand.Read(payload)

		m.push(service, &defs.Message{UUID: uuid.New(), Payload: payload})

		msgCount--
	}

	t.Run("Save log file", func(t *testing.T) {
		if err := m.save(fileName, 0644, 0755); err != nil {
			t.Error(err)
		}
	})

	m2 := newMessages()

	t.Run("Load log file", func(t *testing.T) {
		if err := m2.load(fileName); err != nil {
			t.Error(err)
		}
	})

}
