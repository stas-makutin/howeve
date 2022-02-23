package core

// subscriberID identifier type for the subscriber
type SubscriberID int

// subscriberFn function signature for the subscriber
type SubscriberFn func(event interface{})

//
var (
	subscribers      = make(map[SubscriberID]SubscriberFn)
	lastSubscriberID SubscriberID
)

// subscribe func
func DispatcherSubscribe(fn SubscriberFn) (id SubscriberID) {
	for {
		lastSubscriberID++
		id = lastSubscriberID
		if _, exists := subscribers[id]; !exists {
			subscribers[id] = fn
			break
		}
	}
	return
}

// unsubscribe func
func DispatcherUnsubscribe(id SubscriberID) {
	delete(subscribers, id)
}

// dispatch func
func Dispatch(event interface{}, receivers ...SubscriberID) {
	if len(receivers) > 0 {
		for _, id := range receivers {
			if fn, ok := subscribers[id]; ok {
				fn(event)
			}
		}
	} else {
		for _, fn := range subscribers {
			fn(event)
		}
	}
}
