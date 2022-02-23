package components

import (
	"syscall/js"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/prop"
	"github.com/stas-makutin/howeve/page/actions"
)

type MdcCheckbox struct {
	vecty.Core
	ID       string
	Label    string
	Checked  bool
	Disabled bool
}

func NewMdcCheckbox(id string, label string, checked, disabled bool) (r *MdcCheckbox) {
	r = &MdcCheckbox{ID: id, Label: label, Checked: checked, Disabled: disabled}
	actions.SubscribeGlobal(r)
	return
}

func (ch *MdcCheckbox) OnLoad() {
	js.Global().Get("mdc").Get("checkbox").Get("MDCCheckbox").Call(
		"attachTo", js.Global().Get("document").Call("getElementById", ch.ID),
	)
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
		),
		elem.Div(
			vecty.Markup(
				prop.ID(ch.ID),
				vecty.Class("mdc-checkbox"),
				vecty.MarkupIf(
					ch.Disabled,
					vecty.Class("mdc-checkbox--disabled"),
				),
			),
			elem.Input(
				vecty.Markup(
					prop.ID(idInput),
					prop.Type(prop.TypeCheckbox),
					vecty.Class("mdc-checkbox__native-control"),
					vecty.MarkupIf(
						ch.Checked,
						prop.Checked(true),
					),
					vecty.MarkupIf(
						ch.Disabled,
						prop.Disabled(true),
					),
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
