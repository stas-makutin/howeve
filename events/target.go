package events

import (
	"context"
	"reflect"
)

// TargetedEvents interface
type TargetedEvents interface {
	Receiver() SubscriberID
	SetReceiver(receiver SubscriberID)
}

// EventWithReceiver struct
type EventWithReceiver struct {
	receiver SubscriberID
}

// Receiver func
func (e *EventWithReceiver) Receiver() SubscriberID {
	return e.receiver
}

// SetReceiver func
func (e *EventWithReceiver) SetReceiver(receiver SubscriberID) {
	e.receiver = receiver
}

// RequestResponse func
func (d *Dispatcher) RequestResponse(ctx context.Context, request TargetedEvents, responseType reflect.Type, receiveFn func(interface{})) bool {
	ch := make(chan interface{}, 1)

	var id SubscriberID
	id = d.Receive(nil, ch, func(event interface{}) bool {
		if event != request {
			if te, ok := event.(TargetedEvents); ok && te.Receiver() == id {
				return reflect.TypeOf(event) == responseType
			}
		}
		return false
	})
	defer d.Unsubscribe(id)
	request.SetReceiver(id)
	d.Send(request)
	select {
	case event := <-ch:
		receiveFn(event)
		return true
	case <-ctx.Done():
	}
	return false
}
