package main

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/stas-makutin/howeve/page/actions"
	"github.com/stas-makutin/howeve/page/components"
	"github.com/stas-makutin/howeve/page/core"
	"github.com/stas-makutin/howeve/page/views"
)

type pageMain struct {
	vecty.Core
	viewRoute core.PageRoute
	tabBar    *components.MdcTabBar
}

func newPageMain() (r *pageMain) {
	r = &pageMain{}
	actions.Subscribe(r)
	return
}

func (ch *pageMain) OnRouteChange(route core.PageRoute) {
	if ch.viewRoute != route {
		ch.viewRoute = route
		if ch.tabBar != nil {
			ch.tabBar.ActivateTab(int(route))
		}
	}
}

func (ch *pageMain) tabChange(tabIndex int) {
	route := core.PageRoute(tabIndex)
	if ch.viewRoute != route {
		ch.viewRoute = route
		core.ToRoute(route)
	}
}

func (ch *pageMain) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *pageMain) Render() vecty.ComponentOrHTML {
	ch.viewRoute = core.GetRoute()
	ch.tabBar = components.NewMdcTabBar(
		"top-tab", ch.tabChange,
		components.NewMdcTab("Protocols", ch.viewRoute == core.ProtocolViewRoute),
		components.NewMdcTab("Services", ch.viewRoute == core.ServicesViewRoute),
		components.NewMdcTab("Messages", ch.viewRoute == core.MessagesViewRoute),
		components.NewMdcTab("Discovery", ch.viewRoute == core.DiscoveryViewRoute),
		components.NewMdcTab("Config", ch.viewRoute == core.ConfigViewRoute),
		components.NewMdcTab("Log", ch.viewRoute == core.LogViewRoute),
		&components.Title{},
	)
	return elem.Body(ch.tabBar, newViewMain())
}

type viewMain struct {
	vecty.Core
}

func newViewMain() (r *viewMain) {
	r = &viewMain{}
	actions.Subscribe(r)
	return
}

func (ch *viewMain) OnRouteChange(route core.PageRoute) {
	vecty.Rerender(ch)
}

func (ch *viewMain) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *viewMain) Render() vecty.ComponentOrHTML {
	route := core.GetRoute()
	return elem.Main(
		vecty.Markup(
			vecty.Class("view"),
		),
		vecty.If(route == core.ProtocolViewRoute, views.NewViewProtocols()),
		vecty.If(route == core.ServicesViewRoute, views.NewViewServices()),
		vecty.If(route == core.MessagesViewRoute, &views.ViewMessages{}),
		vecty.If(route == core.DiscoveryViewRoute, &views.ViewDiscovery{}),
		vecty.If(route == core.ConfigViewRoute, views.NewViewConfig()),
		vecty.If(route == core.LogViewRoute, &views.ViewLog{}),
		vecty.If(route == core.NotFoundRoute, &views.ViewNotFound{}),
	)
}
