package main

// subscriberID identifier type for the subscriber
type subscriberID int

// subscriberFn function signature for the subscriber
type subscriberFn func(event interface{})

//
var (
	subscribers      = make(map[subscriberID]subscriberFn)
	lastSubscriberID subscriberID
)

// subscribe func
func dispatcherSubscribe(fn subscriberFn) (id subscriberID) {
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
func dispatcherUnsubscribe(id subscriberID) {
	delete(subscribers, id)
}

// dispatch func
func dispatch(event interface{}, receivers ...subscriberID) {
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
