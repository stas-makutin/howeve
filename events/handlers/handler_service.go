package handlers

func handleAddService(event *AddService) {
	r := &AddServiceResult{ResponseHeader: event.Associate(), AddServiceReply: &AddServiceReply{}}

	//services.Services.Add()

	Dispatcher.Send(r)
}

func handleSendToService(event *SendToService) {
	r := &SendToServiceResult{ResponseHeader: event.Associate()}
	Dispatcher.Send(r)
}

func handleRetriveFromService(event *RetriveFromService) {
	r := &RetriveFromServiceResult{ResponseHeader: event.Associate()}
	Dispatcher.Send(r)
}
