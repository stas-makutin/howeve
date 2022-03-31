package components

import (
	"syscall/js"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/prop"
	"github.com/stas-makutin/howeve/page/core"
)

type MdcTextField struct {
	vecty.Core
	core.ClassAdder
	ID       string `vecty:"prop"`
	Label    string `vecty:"prop"`
	Value    string `vecty:"prop"`
	Disabled bool   `vecty:"prop"`
	jsObject js.Value
}

func NewMdcTextField(id, label, value string, disabled bool) (r *MdcTextField) {
	r = &MdcTextField{ID: id, Label: label, Value: value, Disabled: disabled}
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

func (ch *MdcTextField) Key() interface{} {
	return ch
}

func (ch *MdcTextField) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *MdcTextField) AddClasses(classes ...string) vecty.Component {
	ch.ClassAdder.AddClasses(classes...)
	return ch
}

func (ch *MdcTextField) Render() vecty.ComponentOrHTML {
	labelID := ch.ID + "---label"
	return elem.Label(
		vecty.Markup(
			prop.ID(ch.ID),
			vecty.Class("mdc-text-field", "mdc-text-field--outlined", "mdc-text-field--no-label"),
			prop.Disabled(ch.Disabled),
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
				prop.Type(prop.TypeText),
				prop.Value(ch.Value),
			),
		),
	)
}
