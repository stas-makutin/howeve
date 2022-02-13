package main

import (
	"syscall/js"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/prop"
)

type mdcTab struct {
	vecty.Core
	text   string
	active bool
}

func newMdcTab(text string, active bool) (r *mdcTab) {
	r = &mdcTab{text: text, active: active}
	return
}

func (ch *mdcTab) Render() vecty.ComponentOrHTML {
	return elem.Button(
		vecty.Markup(
			vecty.Class("mdc-tab"),
			vecty.Attribute("role", "tab"),
			vecty.Attribute("tabIndex", "0"),
		),
		vecty.MarkupIf(
			ch.active,
			vecty.Class("mdc-tab--active"),
			vecty.Attribute("aria-selected", "true"),
		),
		elem.Span(
			vecty.Markup(
				vecty.Class("mdc-tab__content"),
			),
			elem.Span(
				vecty.Markup(
					vecty.Class("mdc-tab__text-label"),
				),
				vecty.Text(ch.text),
			),
		),
		elem.Span(
			vecty.Markup(
				vecty.Class("mdc-tab-indicator"),
			),
			vecty.MarkupIf(
				ch.active,
				vecty.Class("mdc-tab-indicator--active"),
			),
			elem.Span(
				vecty.Class("mdc-tab-indicator__content", "mdc-tab-indicator__content--underline"),
			),
		),
		elem.Span(
			vecty.Markup(
				vecty.Class("mdc-tab__ripple"),
			),
		),
	)
}

type mdcTabBar struct {
	vecty.Core
	id   string
	tabs []vecty.ComponentOrHTML
}

func newMdcTabBar(id string, tabs ...vecty.ComponentOrHTML) (r *mdcTabBar) {
	r = &mdcTabBar{id: id, tabs: tabs}
	subscribeGlobal(r)
	return
}

func (ch *mdcTabBar) mdcInitialized() {
	js.Global().Get("mdc").Get("tabBar").Get("MDCTabBar").Call(
		"attachTo", js.Global().Get("document").Call("getElementById", ch.id),
	)
}

func (ch *mdcTabBar) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			prop.ID(ch.id),
			vecty.Class("mdc-tab-bar"),
			vecty.Attribute("role", "tablist"),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("mdc-tab-scroller"),
			),
			elem.Div(
				vecty.Markup(
					vecty.Class("mdc-tab-scroller__scroll-area"),
				),
				elem.Div(
					vecty.Markup(
						vecty.Class("mdc-tab-scroller__scroll-content"),
					),
					ch.tabs,
				),
			),
		),
	)
}
