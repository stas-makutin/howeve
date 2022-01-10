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
	// TODO
	Dispatcher.Send(r)
}
