package services

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/google/uuid"
	"github.com/stas-makutin/howeve/defs"
	"github.com/stas-makutin/howeve/events/handlers"
)

// discoveryEntry defines discovery query entry
type discoveryEntry struct {
	id uuid.UUID

	results []*defs.DiscoveryEntry
	err     error

	completed uint32
	ctx       context.Context
	cancel    context.CancelFunc
}

func newDiscoveryEntry(ctx context.Context) *discoveryEntry {
	de := &discoveryEntry{}
	de.id = uuid.New()
	de.ctx, de.cancel = context.WithCancel(ctx)
	return de
}

// discoveryRegistry is the controlling strucuture to uexecute limited number of discovery queries
// maxActive - max number of goroutines allocated for discovery queries
// maxEntries - max number of discovery queries allowed, must be >= maxActive
type discoveryRegistry struct {
	maxEntries int
	maxActive  int

	lock      sync.Mutex
	entries   map[uuid.UUID]*discoveryEntry
	completed []*discoveryEntry

	ctx         context.Context
	cancel      context.CancelFunc
	stopWg      sync.WaitGroup
	activeCount int
}

// newDiscoveryRegistry creates new discovery registry with provided limits
func newDiscoveryRegistry(maxEntries, maxActive int) *discoveryRegistry {
	if maxActive > maxEntries {
		maxEntries = maxActive
	}
	d := &discoveryRegistry{}
	d.maxEntries, d.maxActive = maxEntries, maxActive
	d.entries = make(map[uuid.UUID]*discoveryEntry)
	d.ctx, d.cancel = context.WithCancel(context.Background())
	return d
}

// stop stops all running discovery queries and frees resources
func (d *discoveryRegistry) stop() {
	d.cancel()
	d.stopWg.Wait()
	d.entries = nil
	d.completed = nil
	d.ctx = nil
}

// Discover is the part of defs.ServiceRegistry implementation, this function accept, if possible, new discovery query and returns its id (UUID)
// This id then could be used to get discovery query results using Discovery method
func (d *discoveryRegistry) Discover(protocol defs.ProtocolIdentifier, transport defs.TransportIdentifier, params defs.RawParamValues) (uuid.UUID, error) {
	to, _, err := defs.ResolveProtocolAndTransport(protocol, transport)
	if err != nil {
		return uuid.Nil, err
	}
	if to.DiscoveryFunc == nil {
		return uuid.Nil, defs.ErrNoDiscovery
	}
	pv, err := to.DiscoveryParams.ParseValues(params)
	if err != nil {
		return uuid.Nil, err
	}

	d.lock.Lock()
	defer d.lock.Unlock()

	if d.activeCount >= d.maxActive {
		return uuid.Nil, defs.ErrDiscoveryBusy
	}
	if len(d.entries) > d.maxEntries {
		if len(d.completed) == 0 {
			return uuid.Nil, defs.ErrDiscoveryBusy
		}
		delete(d.entries, d.completed[0].id)
		copy(d.completed[0:], d.completed[1:])
		d.completed = d.completed[1:]
	}

	de := newDiscoveryEntry(d.ctx)
	d.entries[de.id] = de

	handlers.SendDiscoveryStarted(de.id, protocol, transport, params)

	d.activeCount++
	d.stopWg.Add(1)
	go d.discoveryRunner(de, to.DiscoveryFunc, pv)

	return uuid.Nil, nil
}

// Discovery method returns the state or results of discovery query, identifying by its id
func (d *discoveryRegistry) Discovery(id uuid.UUID, stop bool) ([]*defs.DiscoveryEntry, error) {
	d.lock.Lock()
	defer d.lock.Unlock()

	de, ok := d.entries[id]
	if !ok {
		return nil, defs.ErrNoDiscoveryID
	}
	if stop {
		de.cancel()
	}
	if atomic.LoadUint32(&de.completed) != 1 {
		return nil, defs.ErrDiscoveryPending
	}
	return de.results, de.err
}

// discoveryRunner is the goroutine which executes discovery query
func (d *discoveryRegistry) discoveryRunner(de *discoveryEntry, df defs.DiscoveryFunc, params defs.ParamValues) {
	defer atomic.StoreUint32(&de.completed, 1)
	defer d.stopWg.Done()

	de.results, de.err = df(de.ctx, params)

	handlers.SendDiscoveryFinished(de.id, de.results, de.err)

	d.lock.Lock()
	defer d.lock.Unlock()
	d.completed = append(d.completed, de)
	d.activeCount--
}
