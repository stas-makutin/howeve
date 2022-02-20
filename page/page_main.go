package main

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

type pageMain struct {
	vecty.Core
	viewRoute pageRoute
	tabBar    *mdcTabBar
}

func newPageMain() (r *pageMain) {
	r = &pageMain{}
	subscribeGlobal(r)
	return
}

func (ch *pageMain) routeChange(route pageRoute) {
	if ch.viewRoute != route {
		ch.viewRoute = route
		if ch.tabBar != nil && !(ch.tabBar.jsTabBar.IsUndefined() || ch.tabBar.jsTabBar.IsNull()) {
			ch.tabBar.jsTabBar.Call("activateTab", int(route))
		}
	}
}

func (ch *pageMain) tabChange(tabIndex int) {
	route := pageRoute(tabIndex)
	if ch.viewRoute != route {
		ch.viewRoute = route
		toRoute(route)
	}
}

func (ch *pageMain) Render() vecty.ComponentOrHTML {
	ch.viewRoute = getRoute()
	ch.tabBar = newMdcTabBar(
		"top-tab", ch.tabChange,
		newMdcTab("Protocols", ch.viewRoute == ProtocolViewRoute),
		newMdcTab("Services", ch.viewRoute == ServicesViewRoute),
		newMdcTab("Messages", ch.viewRoute == MessagesViewRoute),
		&title{},
	)
	return elem.Body(ch.tabBar, newViewMain())
}

type viewMain struct {
	vecty.Core
}

func newViewMain() (r *viewMain) {
	r = &viewMain{}
	subscribeGlobal(r)
	return
}

func (ch *viewMain) routeChange(route pageRoute) {
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
