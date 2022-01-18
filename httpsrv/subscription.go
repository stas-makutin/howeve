package httpsrv

import (
	"reflect"
	"sync"

	"github.com/stas-makutin/howeve/events/handlers"
)

// SubscriptionEvent defines types of events to subscribe (web sockets only)
type SubscriptionEvent byte

// SubscriptionEvent types
const (
	EventDiscoveryStarted = SubscriptionEvent(iota)
	EventDiscoveryFinished
	EventNewMessage
	EventDropMessage
	EventUpdateMessageState
)

var subscriptionEventNameMap = map[string]SubscriptionEvent{
	"discoveryStarted":   EventDiscoveryStarted,
	"discoveryFinished":  EventDiscoveryFinished,
	"newMessage":         EventNewMessage,
	"dropMessage":        EventDropMessage,
	"updateMessageState": EventUpdateMessageState,
}

var subscriptionEventTypeMap = map[SubscriptionEvent]reflect.Type{
	EventDiscoveryStarted:   reflect.TypeOf(&handlers.DiscoveryStarted{}),
	EventDiscoveryFinished:  reflect.TypeOf(&handlers.DiscoveryFinished{}),
	EventNewMessage:         reflect.TypeOf(&handlers.NewMessage{}),
	EventDropMessage:        reflect.TypeOf(&handlers.DropMessage{}),
	EventUpdateMessageState: reflect.TypeOf(&handlers.UpdateMessageState{}),
}

// Subscription query structure
type Subscription struct {
	Subscribe bool                `json:"subscribe"`
	AllEvents bool                `json:"all,omitempty"`
	Events    []SubscriptionEvent `json:"events,omitempty"`
}

type socketSubscription struct {
	sync.RWMutex
	data map[reflect.Type]struct{}
}

func newSocketSubscription() *socketSubscription {
	return &socketSubscription{
		data: make(map[reflect.Type]struct{}),
	}
}

func (s *socketSubscription) subscribe(query *Subscription) {
	s.Lock()
	defer s.Unlock()

	var action func(t reflect.Type)
	if query.Subscribe {
		action = func(t reflect.Type) {
			s.data[t] = struct{}{}
		}
	} else {
		action = func(t reflect.Type) {
			delete(s.data, t)
		}
	}

	if query.AllEvents {
		for _, v := range subscriptionEventTypeMap {
			action(v)
		}
	} else {
		for _, se := range query.Events {
			if v, ok := subscriptionEventTypeMap[se]; ok {
				action(v)
			}
		}
	}
}

func (s *socketSubscription) subscribed(v interface{}) bool {
	s.RLock()
	defer s.RUnlock()

	_, ok := s.data[reflect.TypeOf(v)]
	return ok
}
