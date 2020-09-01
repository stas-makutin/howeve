package events

import (
	"context"
	"sync"
	"sync/atomic"
)

type eventSubscriberID uint64

type eventListenerFn func(event interface{})

type eventDispatcher struct {
	subscribers      sync.Map
	lastSubscriberID eventSubscriberID
}

func (d *eventDispatcher) subscribe(fn eventListenerFn) (id eventSubscriberID) {
	for {
		id = (eventSubscriberID)(atomic.AddUint64((*uint64)(&d.lastSubscriberID), 1))
		if _, exists := d.subscribers.LoadOrStore(id, fn); !exists {
			break
		}
	}
	return
}

func (d *eventDispatcher) unsubscribe(id eventSubscriberID) bool {
	_, exists := d.subscribers.LoadAndDelete(id)
	return exists
}

func (d *eventDispatcher) send(event interface{}) {
	d.subscribers.Range(func(key, value interface{}) bool {
		value.(eventListenerFn)(event)
		return true
	})
}

func (d *eventDispatcher) receive(ctx context.Context, ch chan<- interface{}, fn func(event interface{}) bool) (id eventSubscriberID) {
	id = d.subscribe(func(event interface{}) {
		if fn == nil || fn(event) {
			select {
			case ch <- event:
			case <-ctx.Done():
			}
		}
	})
	return
}
