package events

import (
	"context"
	"sync"
	"sync/atomic"
)

// SubscriberID identifier type for the subscriber
type SubscriberID uint64

// SubscriberFn function signature for the subscriber
type SubscriberFn func(event interface{})

// Dispatcher - event dispatcher struct
type Dispatcher struct {
	subscribers      sync.Map
	lastSubscriberID SubscriberID
}

// Subscribe func
func (d *Dispatcher) Subscribe(fn SubscriberFn) (id SubscriberID) {
	for {
		id = (SubscriberID)(atomic.AddUint64((*uint64)(&d.lastSubscriberID), 1))
		if _, exists := d.subscribers.LoadOrStore(id, fn); !exists {
			break
		}
	}
	return
}

// Unsubscribe func
func (d *Dispatcher) Unsubscribe(id SubscriberID) bool {
	_, exists := d.subscribers.LoadAndDelete(id)
	return exists
}

// Send func
func (d *Dispatcher) Send(event interface{}, receivers ...SubscriberID) {
	if len(receivers) > 0 {
		for _, id := range receivers {
			if fn, ok := d.subscribers.Load(id); ok {
				fn.(SubscriberFn)(event)
			}
		}
	} else {
		d.subscribers.Range(func(key, value interface{}) bool {
			value.(SubscriberFn)(event)
			return true
		})
	}
}

// Receive func
func (d *Dispatcher) Receive(ctx context.Context, ch chan<- interface{}, fn func(event interface{}) bool) (id SubscriberID) {
	id = d.Subscribe(func(event interface{}) {
		if fn == nil || fn(event) {
			if ctx == nil {
				select {
				case ch <- event:
				default:
				}
			} else {
				select {
				case ch <- event:
				case <-ctx.Done():
				}
			}
		}
	})
	return
}
