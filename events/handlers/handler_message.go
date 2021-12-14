package handlers

import "github.com/stas-makutin/howeve/defs"

// SendNewMessage sends NewMessage event
func SendNewMessage(service *defs.ServiceKey, message *defs.Message) {
	Dispatcher.SendAsync(&NewMessage{
		Header:     *NewHeader(""),
		ServiceKey: service, Message: message,
	})
}

// SendDropMessage sends DropMessage event
func SendDropMessage(service *defs.ServiceKey, message *defs.Message) {
	Dispatcher.SendAsync(&DropMessage{
		Header:     *NewHeader(""),
		ServiceKey: service, Message: message,
	})
}

// SendUpdateMessageState sends UpdateMessageState event
func SendUpdateMessageState(service *defs.ServiceKey, message *defs.Message, prevState defs.MessageState) {
	Dispatcher.SendAsync(&UpdateMessageState{
		Header:     *NewHeader(""),
		ServiceKey: service, Message: message, PrevState: prevState,
	})
}
