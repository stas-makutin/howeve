package events

import (
	"context"
	"reflect"
)

// TargetedRequest interface
type TargetedRequest interface {
	SetReceiver(receiver SubscriberID)
}

// RequestTarget struct
type RequestTarget struct {
	ReceiverID SubscriberID
}

// SetReceiver func
func (e *RequestTarget) SetReceiver(receiver SubscriberID) {
	e.ReceiverID = receiver
}

// ResponseTarget func
func (e *RequestTarget) ResponseTarget() ResponseTarget {
	return ResponseTarget{ReceiverID: e.ReceiverID}
}

// TargetedResponse interface
type TargetedResponse interface {
	Receiver() SubscriberID
}

// ResponseTarget struct
type ResponseTarget struct {
	ReceiverID SubscriberID
}

// Receiver func
func (e *ResponseTarget) Receiver() SubscriberID {
	return e.ReceiverID
}

// RequestResponse func
func (d *Dispatcher) RequestResponse(ctx context.Context, request TargetedRequest, responseType reflect.Type, receiveFn func(interface{})) bool {
	ch := make(chan interface{}, 1)

	var id SubscriberID
	id = d.Receive(nil, ch, func(event interface{}) bool {
		if te, ok := event.(TargetedResponse); ok && te.Receiver() == id {
			return reflect.TypeOf(event) == responseType
		}
		return false
	})
	defer d.Unsubscribe(id)
	request.SetReceiver(id)
	d.Send(request)
	select {
	case <-ctx.Done():
		// break
	case event := <-ch:
		receiveFn(event)
		return true
	}
	return false
}
