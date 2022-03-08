package httpsrv

import (
	"reflect"
	"sync"

	"github.com/stas-makutin/howeve/api"
	"github.com/stas-makutin/howeve/events/handlers"
)

var subscriptionEventTypeMap = map[api.SubscriptionEvent]reflect.Type{
	api.EventDiscoveryStarted:   reflect.TypeOf(&handlers.ProtocolDiscoveryStarted{}),
	api.EventDiscoveryFinished:  reflect.TypeOf(&handlers.ProtocolDiscoveryFinished{}),
	api.EventNewMessage:         reflect.TypeOf(&handlers.NewMessage{}),
	api.EventDropMessage:        reflect.TypeOf(&handlers.DropMessage{}),
	api.EventUpdateMessageState: reflect.TypeOf(&handlers.UpdateMessageState{}),
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

func (s *socketSubscription) subscribe(query *api.Subscription) {
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
