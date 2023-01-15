package components

import (
	"syscall/js"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
	"github.com/hexops/vecty/prop"
	"github.com/stas-makutin/howeve/page/core"
)

type MdcBanner struct {
	vecty.Core
	core.Classes
	core.Keyable
	ID         string                `vecty:"prop"`
	ButtonText string                `vecty:"prop"`
	Content    []vecty.MarkupOrChild `vecty:"prop"`
	fullWidth  bool
	clickFn    func()
	jsObject   js.Value
}

func NewMdcBanner(id, buttonText string, fullWidth bool, clickFn func(), content ...vecty.MarkupOrChild) (r *MdcBanner) {
	r = &MdcBanner{ID: id, ButtonText: buttonText, fullWidth: fullWidth, clickFn: clickFn, Content: content}
	return
}

func (ch *MdcBanner) Mount() {
	ch.Unmount()
	ch.jsObject = js.Global().Get("mdc").Get("banner").Get("MDCBanner").Call(
		"attachTo", js.Global().Get("document").Call("getElementById", ch.ID),
	)
	ch.jsObject.Call("open")
}

func (ch *MdcBanner) Unmount() {
	core.SafeJSDestroy(&ch.jsObject, func(v *js.Value) { v.Call("destroy") })
}

func (ch *MdcBanner) WithKey(key interface{}) *MdcBanner {
	ch.Keyable.WithKey(key)
	return ch
}

func (ch *MdcBanner) WithClasses(classes ...string) *MdcBanner {
	ch.Classes.WithClasses(classes...)
	return ch
}

func (ch *MdcBanner) onClick(event *vecty.Event) {
	ch.clickFn()
}

func (ch *MdcBanner) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *MdcBanner) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			prop.ID(ch.ID),
			vecty.Class("mdc-banner", "mdc-elevation--z2"),
			vecty.Attribute("role", "banner"),
			ch.ApplyClasses(),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("mdc-banner__content"),
				vecty.Attribute("role", "alertdialog"),
				vecty.Attribute("aria-live", "assertive"),
				vecty.MarkupIf(ch.fullWidth, vecty.Style("max-width", "none")),
			),
			elem.Div(
				vecty.Markup(
					vecty.Class("mdc-banner__graphic-text-wrapper"),
				),
				elem.Div(
					vecty.Markup(
						vecty.Class("mdc-banner__text", "mdc-theme--error"),
					),
					elem.Slot(ch.Content...),
				),
			),
			elem.Div(
				vecty.Markup(
					vecty.Class("mdc-banner__actions"),
				),
				elem.Button(
					vecty.Markup(
						prop.Type(prop.TypeButton),
						vecty.Class("mdc-button", "mdc-banner__primary-action"),
						event.Click(ch.onClick).StopPropagation(),
					),
					elem.Div(
						vecty.Markup(
							vecty.Class("mdc-button__ripple"),
						),
					),
					elem.Div(
						vecty.Markup(
							vecty.Class("mdc-button__label"),
						),
						vecty.Text(ch.ButtonText),
					),
				),
			),
		),
	)
}
