package components

import (
	"syscall/js"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/prop"
	"github.com/stas-makutin/howeve/page/actions"
)

type MdcTab struct {
	vecty.Core
	Text   string
	Active bool
}

func NewMdcTab(text string, active bool) (r *MdcTab) {
	r = &MdcTab{Text: text, Active: active}
	return
}

func (ch *MdcTab) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *MdcTab) Render() vecty.ComponentOrHTML {
	return elem.Button(
		vecty.Markup(
			vecty.Class("mdc-tab"),
			vecty.Attribute("role", "tab"),
			vecty.Attribute("tabIndex", "0"),
			vecty.MarkupIf(
				ch.Active,
				vecty.Class("mdc-tab--active"),
				vecty.Attribute("aria-selected", "true"),
			),
		),
		elem.Span(
			vecty.Markup(
				vecty.Class("mdc-tab__content"),
			),
			elem.Span(
				vecty.Markup(
					vecty.Class("mdc-tab__text-label"),
				),
				vecty.Text(ch.Text),
			),
		),
		elem.Span(
			vecty.Markup(
				vecty.Class("mdc-tab-indicator"),
				vecty.MarkupIf(
					ch.Active,
					vecty.Class("mdc-tab-indicator--active"),
				),
			),
			elem.Span(
				vecty.Markup(
					vecty.Class("mdc-tab-indicator__content", "mdc-tab-indicator__content--underline"),
				),
			),
		),
		elem.Span(
			vecty.Markup(
				vecty.Class("mdc-tab__ripple"),
			),
		),
	)
}

type MdcTabBar struct {
	vecty.Core
	ID          string
	ActivatedFn func(tabIndex int)
	Tabs        vecty.List
	JsTabBar    js.Value
}

func NewMdcTabBar(id string, activatedFn func(tabIndex int), tabs ...vecty.ComponentOrHTML) (r *MdcTabBar) {
	r = &MdcTabBar{ID: id, ActivatedFn: activatedFn, Tabs: tabs}
	actions.SubscribeGlobal(r)
	return
}

func (ch *MdcTabBar) OnLoad() {
	ch.JsTabBar = js.Global().Get("mdc").Get("tabBar").Get("MDCTabBar").Call(
		"attachTo", js.Global().Get("document").Call("getElementById", ch.ID),
	)
	ch.JsTabBar.Call("listen", "MDCTabBar:activated", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) > 0 {
			ch.ActivatedFn(args[0].Get("detail").Get("index").Int())
		}
		return nil
	}))
}

func (ch *MdcTabBar) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *MdcTabBar) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			prop.ID(ch.ID),
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
					ch.Tabs,
				),
			),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("tab-bar-divider"),
			),
		),
	)
}
