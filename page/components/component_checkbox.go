package components

import (
	"syscall/js"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
	"github.com/hexops/vecty/prop"
	"github.com/stas-makutin/howeve/page/core"
)

type MdcCheckbox struct {
	vecty.Core
	core.Classes
	core.Keyable
	ID       string `vecty:"prop"`
	Label    string `vecty:"prop"`
	Checked  bool   `vecty:"prop"`
	Disabled bool   `vecty:"prop"`
	changeFn func(checked, disabled bool)
	jsObject js.Value
}

func NewMdcCheckbox(id string, label string, checked, disabled bool, changeFn func(checked, disabled bool)) (r *MdcCheckbox) {
	r = &MdcCheckbox{ID: id, Label: label, Checked: checked, Disabled: disabled, changeFn: changeFn}
	return
}

func (ch *MdcCheckbox) Mount() {
	ch.Unmount()
	ch.jsObject = js.Global().Get("mdc").Get("checkbox").Get("MDCCheckbox").Call(
		"attachTo", js.Global().Get("document").Call("getElementById", ch.ID),
	)
}

func (ch *MdcCheckbox) Unmount() {
	core.SafeJSDestroy(&ch.jsObject, func(v *js.Value) { v.Call("destroy") })
}

func (ch *MdcCheckbox) WithKey(key interface{}) *MdcCheckbox {
	ch.Keyable.WithKey(key)
	return ch
}

func (ch *MdcCheckbox) WithClasses(classes ...string) *MdcCheckbox {
	ch.Classes.WithClasses(classes...)
	return ch
}

func (ch *MdcCheckbox) onClick(event *vecty.Event) {
	ch.Checked = ch.jsObject.Get("checked").Bool()
	ch.changeFn(ch.Checked, ch.Disabled)
}

func (ch *MdcCheckbox) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *MdcCheckbox) Render() vecty.ComponentOrHTML {
	idInput := ch.ID + "---input"
	return elem.Div(
		vecty.Markup(
			vecty.Class("mdc-form-field"),
			ch.ApplyClasses(),
		),
		elem.Div(
			vecty.Markup(
				prop.ID(ch.ID),
				vecty.Class("mdc-checkbox"),
				vecty.MarkupIf(
					ch.Disabled,
					vecty.Class("mdc-checkbox--disabled"),
				),
				event.Click(ch.onClick),
			),
			elem.Input(
				vecty.Markup(
					prop.ID(idInput),
					prop.Type(prop.TypeCheckbox),
					vecty.Class("mdc-checkbox__native-control"),
					prop.Checked(ch.Checked),
					prop.Disabled(ch.Disabled),
				),
			),
			elem.Div(
				vecty.Markup(
					vecty.Class("mdc-checkbox__background"),
				),
				vecty.Tag(
					"svg",
					vecty.Markup(
						vecty.Namespace("http://www.w3.org/2000/svg"),
						vecty.Class("mdc-checkbox__checkmark"),
						vecty.Attribute("viewBox", "0 0 24 24"),
					),
					vecty.Tag(
						"path",
						vecty.Markup(
							vecty.Namespace("http://www.w3.org/2000/svg"),
							vecty.Class("mdc-checkbox__checkmark-path"),
							vecty.Attribute("fill", "none"),
							vecty.Attribute("d", "M1.73,12.91 8.1,19.28 22.79,4.59"),
						),
					),
				),
				elem.Div(
					vecty.Markup(
						vecty.Class("mdc-checkbox__mixedmark"),
					),
				),
			),
			elem.Div(
				vecty.Markup(
					vecty.Class("mdc-checkbox__ripple"),
				),
			),
			elem.Div(
				vecty.Markup(
					vecty.Class("mdc-checkbox__focus-ring"),
				),
			),
		),
		elem.Label(
			vecty.Markup(
				prop.For(idInput),
			),
			vecty.Text(ch.Label),
		),
	)
}
