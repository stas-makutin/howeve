package services

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/stas-makutin/howeve/defs"
)

type discoveryEntry struct {
	results []*defs.DiscoveryEntry
	cancel  context.CancelFunc
	active  bool
}

type discoveryRegistry struct {
	lock    sync.Mutex
	entries map[uuid.UUID]discoveryEntry

	ctx         context.Context
	cancel      context.CancelFunc
	stopWg      sync.WaitGroup
	activeCount uint32
}

func newDiscoveryRegistry() *discoveryRegistry {
	d := &discoveryRegistry{}
	d.entries = make(map[uuid.UUID]discoveryEntry)
	d.ctx, d.cancel = context.WithCancel(context.Background())
	return d
}

func (d *discoveryRegistry) stop() {
	d.cancel()
	d.stopWg.Wait()
	d.entries = nil
}

func (d *discoveryRegistry) Discover(protocol defs.ProtocolIdentifier, transport defs.TransportIdentifier, params defs.RawParamValues) (uuid.UUID, error) {
	return uuid.Nil, nil
}

func (d *discoveryRegistry) Discovery(id uuid.UUID, stop bool) ([]*defs.DiscoveryEntry, error) {
	return nil, nil
}
