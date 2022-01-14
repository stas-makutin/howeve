package handlers

import "github.com/stas-makutin/howeve/defs"

// SendNewMessage sends NewMessage event
func SendNewMessage(service *defs.ServiceKey, message *defs.Message) {
	Dispatcher.SendAsync(&NewMessage{
		Header: *NewHeader(""),
		MessageEntry: &MessageEntry{
			ServiceKey: service, Message: message,
		},
	})
}

// SendDropMessage sends DropMessage event
func SendDropMessage(service *defs.ServiceKey, message *defs.Message) {
	Dispatcher.SendAsync(&DropMessage{
		Header: *NewHeader(""),
		MessageEntry: &MessageEntry{
			ServiceKey: service, Message: message,
		},
	})
}

// SendUpdateMessageState sends UpdateMessageState event
func SendUpdateMessageState(service *defs.ServiceKey, message *defs.Message, prevState defs.MessageState) {
	Dispatcher.SendAsync(&UpdateMessageState{
		Header: *NewHeader(""),
		UpdateMessageStateData: &UpdateMessageStateData{
			MessageEntry: &MessageEntry{
				ServiceKey: service, Message: message,
			},
			PrevState: prevState,
		},
	})
}

func handleGetMessage(event *GetMessage) {
	r := &GetMessageResult{ResponseHeader: event.Associate()}
	key, message := defs.Messages.Get(event.ID)
	if key != nil && message != nil {
		r.MessageEntry = &MessageEntry{
			ServiceKey: key,
			Message:    message,
		}
	}
	Dispatcher.Send(r)
}

func handleListMessages(event *ListMessages) {
	r := &ListMessagesResult{ResponseHeader: event.Associate()}

	var fromFn defs.MessageFindFunc
	var itFn defs.MessageFunc

	if event.FromIndex != nil {
		fromFn = defs.Messages.FromIndex(*event.FromIndex, event.FromExclusive)
	} else if event.FromID != nil {
		fromFn = defs.Messages.FromID(*event.FromID, event.FromExclusive)
	} else if event.FromTime != nil {
		fromFn = defs.Messages.FromTime(*event.FromTime, event.FromExclusive)
	} else {
		fromFn = defs.Messages.FromIndex(0, false)
	}

	serviceIndices := make(map[defs.ServiceKey]int)
	itFn = func(index int, key *defs.ServiceKey, message *defs.Message) bool {
		serviceIndex, ok := serviceIndices[*key]
		if !ok {
			serviceIndex = len(r.Services)
			r.Services = append(r.Services, &ServiceID{ServiceKey: key})
			serviceIndices[*key] = serviceIndex
		}
		r.Messages = append(r.Messages, &ListMessage{Message: message, ServiceIndex: serviceIndex})
		return false
	}

	count := event.Count
	if count > 10000 {
		count = 10000
	}

	if event.CountAfterFilter {
		itFn = defs.UntilCounter(count, itFn)
	}
	if len(event.Services) > 0 {
		var serviceKeys []*defs.ServiceKey
		i := 0
		defs.Services.ResolveIDs(
			func(key *defs.ServiceKey, alias string) {
				if key != nil {
					serviceKeys = append(serviceKeys, key)
				}
			},
			func() (*defs.ServiceKey, string, bool) {
				si := event.Services[i]
				i += 1
				return si.ServiceKey, si.Alias, i >= len(event.Services)
			},
		)
		itFn = defs.WithServices(serviceKeys, itFn)
	}
	if len(event.States) > 0 {
		itFn = defs.WithStates(event.States, itFn)
	}
	if len(event.Payloads) > 0 {
		itFn = defs.WithPayload(event.Payloads, itFn)
	}
	if event.UntilIndex != nil {
		itFn = defs.UntilIndex(*event.UntilIndex, event.UntilExclusive, itFn)
	} else if event.UntilID != nil {
		itFn = defs.UntilID(*event.UntilID, event.UntilExclusive, itFn)
	} else if event.UntilTime != nil {
		itFn = defs.UntilTime(*event.UntilTime, event.UntilExclusive, itFn)
	}
	if !event.CountAfterFilter {
		itFn = defs.UntilCounter(count, itFn)
	}

	r.Count = defs.Messages.List(fromFn, itFn)

	if len(r.Services) > 0 {
		i := -1
		defs.Services.ResolveIDs(
			func(key *defs.ServiceKey, alias string) {
				r.Services[i].Alias = alias
			},
			func() (*defs.ServiceKey, string, bool) {
				i += 1
				return r.Services[i].ServiceKey, "", i+1 >= len(r.Services)
			},
		)
	}

	Dispatcher.Send(r)
}
