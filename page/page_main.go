package main

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

type pageMain struct {
	vecty.Core
}

func (ch *pageMain) Render() vecty.ComponentOrHTML {
	route := getRoute()
	view := &viewMain{}

	return elem.Body(
		newMdcTabBar(
			"top-tab",
			func(tabIndex int) {
				view.changeView(pageRoute(tabIndex))
			},
			newMdcTab("Protocols", route == ProtocolViewRoute),
			newMdcTab("Services", route == ServicesViewRoute),
			newMdcTab("Messages", route == MessagesViewRoute),
			&title{},
		),
		view,
	)
}

type viewMain struct {
	vecty.Core
}

func (ch *viewMain) changeView(view pageRoute) {
	toRoute(view)
	vecty.Rerender(ch)
}

func (ch *viewMain) Render() vecty.ComponentOrHTML {
	route := getRoute()
	return elem.Main(
		vecty.If(route == ProtocolViewRoute, &viewProtocols{}),
		vecty.If(route == ServicesViewRoute, &viewServices{}),
		vecty.If(route == MessagesViewRoute, &viewMessages{}),
		vecty.If(route == NotFoundRoute, &viewNotFound{}),
	)
}
