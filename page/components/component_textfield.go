package components

import (
	"syscall/js"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
	"github.com/hexops/vecty/prop"
	"github.com/stas-makutin/howeve/page/core"
)

type MdcTextField struct {
	vecty.Core
	core.Classes
	core.Keyable
	ID             string `vecty:"prop"`
	Label          string `vecty:"prop"`
	Value          string `vecty:"prop"`
	Disabled       bool   `vecty:"prop"`
	Invalid        bool   `vecty:"prop"`
	inputAtributes []vecty.Applyer
	changeFn       func(value string) string
	jsObject       js.Value
}

func NewMdcTextField(id, label, value string, disabled bool, invalid bool, changeFn func(value string) string, inputAtributes ...vecty.Applyer) (r *MdcTextField) {
	r = &MdcTextField{ID: id, Label: label, Value: value, Disabled: disabled, Invalid: invalid, inputAtributes: inputAtributes, changeFn: changeFn}
	return
}

func (ch *MdcTextField) Mount() {
	ch.Unmount()
	ch.jsObject = js.Global().Get("mdc").Get("textField").Get("MDCTextField").Call(
		"attachTo", js.Global().Get("document").Call("getElementById", ch.ID),
	)
}

func (ch *MdcTextField) Unmount() {
	core.SafeJSDestroy(&ch.jsObject, func(v *js.Value) { v.Call("destroy") })
}

func (ch *MdcTextField) WithKey(key interface{}) *MdcTextField {
	ch.Keyable.WithKey(key)
	return ch
}

func (ch *MdcTextField) WithClasses(classes ...string) *MdcTextField {
	ch.Classes.WithClasses(classes...)
	return ch
}

func (ch *MdcTextField) change(event *vecty.Event) {
	event.Target.Call("setCustomValidity", ch.changeFn(event.Target.Get("value").String()))
}

func (ch *MdcTextField) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *MdcTextField) Render() vecty.ComponentOrHTML {
	labelID := ch.ID + "---label"
	return elem.Label(
		vecty.Markup(
			prop.ID(ch.ID),
			vecty.Class("mdc-text-field", "mdc-text-field--outlined", "mdc-text-field--no-label"),
			vecty.MarkupIf(ch.Invalid, vecty.Class("mdc-text-field--invalid")),
			prop.Disabled(ch.Disabled),
			ch.ApplyClasses(),
		),
		elem.Span(
			vecty.Markup(
				vecty.Class("mdc-notched-outline"),
			),
			elem.Span(
				vecty.Markup(
					vecty.Class("mdc-notched-outline__leading"),
				),
			),
			vecty.If(ch.Label != "",
				elem.Span(
					vecty.Markup(
						vecty.Class("mdc-notched-outline__notch"),
					),
					elem.Span(
						vecty.Markup(
							prop.ID(labelID),
							vecty.Class("mdc-floating-label"),
						),
						vecty.Text(ch.Label),
					),
				),
			),
			elem.Span(
				vecty.Markup(
					vecty.Class("mdc-notched-outline__trailing"),
				),
			),
		),
		elem.Input(
			vecty.Markup(
				vecty.Class("mdc-text-field__input"),
				vecty.MarkupIf(ch.Label != "", vecty.Attribute("aria-labelledby", labelID)),
				vecty.MarkupIf(ch.Label == "", vecty.Attribute("aria-label", "Label")),
				vecty.MarkupIf(len(ch.inputAtributes) <= 0, prop.Type(prop.TypeText)),
				prop.Value(ch.Value),
				event.Change(ch.change),
			),
			vecty.Markup(
				ch.inputAtributes...,
			),
		),
	)
}
