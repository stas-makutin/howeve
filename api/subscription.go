package api

import (
	"encoding/json"
	"fmt"
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

var subscriptionEventTypeMap = map[string]SubscriptionEvent{
	"discoveryStarted":   EventDiscoveryStarted,
	"discoveryFinished":  EventDiscoveryFinished,
	"newMessage":         EventNewMessage,
	"dropMessage":        EventDropMessage,
	"updateMessageState": EventUpdateMessageState,
}
var subscriptionEventNameMap map[SubscriptionEvent]string

func init() {
	subscriptionEventNameMap = make(map[SubscriptionEvent]string)
	for k, v := range subscriptionEventTypeMap {
		subscriptionEventNameMap[v] = k
	}
}

// Subscription query structure
type Subscription struct {
	Subscribe bool                `json:"subscribe"`
	AllEvents bool                `json:"all,omitempty"`
	Events    []SubscriptionEvent `json:"events,omitempty"`
}

func (c SubscriptionEvent) String() string {
	if s, ok := subscriptionEventNameMap[c]; ok {
		return s
	}
	return ""
}

// MarshalJSON func for SubscriptionEvent value
func (c SubscriptionEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

// UnmarshalJSON func for SubscriptionEvent value
func (event *SubscriptionEvent) UnmarshalJSON(data []byte) error {
	var name string
	err := json.Unmarshal(data, &name)
	if err != nil {
		return err
	}
	if ev, ok := subscriptionEventTypeMap[name]; ok {
		*event = ev
		return nil
	}
	return fmt.Errorf("unknown subsctiption event name '%v'", name)
}
