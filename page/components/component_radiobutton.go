package components

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/prop"
	"github.com/stas-makutin/howeve/page/core"
)

type MdcRadioButton struct {
	vecty.Core
	core.Classes
	core.Keyable
	ID       string `vecty:"prop"`
	Label    string `vecty:"prop"`
	Name     string `vecty:"prop"`
	Value    string `vecty:"prop"`
	Checked  bool   `vecty:"prop"`
	Disabled bool   `vecty:"prop"`
	changeFn func()
}

func NewMdcRadioButton(id, label, name, value string, checked, disabled bool, changeFn func()) (r *MdcRadioButton) {
	r = &MdcRadioButton{ID: id, Label: label, Name: name, Value: value, Checked: checked, Disabled: disabled, changeFn: changeFn}
	return
}

func (ch *MdcRadioButton) onChange(event *vecty.Event) {
	ch.changeFn()
}

func (ch *MdcRadioButton) WithKey(key interface{}) *MdcRadioButton {
	ch.Keyable.WithKey(key)
	return ch
}

func (ch *MdcRadioButton) WithClasses(classes ...string) *MdcRadioButton {
	ch.Classes.WithClasses(classes...)
	return ch
}

func (ch *MdcRadioButton) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *MdcRadioButton) Render() vecty.ComponentOrHTML {
	idInput := ch.ID + "---input"
	return NewMdcFormField(
		elem.Div(
			vecty.Markup(
				prop.ID(ch.ID),
				vecty.Class("mdc-radio"),
				vecty.MarkupIf(ch.Disabled, vecty.Class("mdc-radio--disabled")),
			),
			elem.Input(
				vecty.Markup(
					vecty.Class("mdc-radio__native-control"),
					prop.ID(idInput),
					prop.Type(prop.TypeRadio),
					prop.Checked(ch.Checked),
					prop.Disabled(ch.Disabled),
					prop.Name(ch.Name),
					prop.Value(ch.Value),
				),
			),
			elem.Div(
				vecty.Markup(
					vecty.Class("mdc-radio__background"),
				),
				elem.Div(
					vecty.Markup(
						vecty.Class("mdc-radio__outer-circle"),
					),
				),
				elem.Div(
					vecty.Markup(
						vecty.Class("mdc-radio__inner-circle"),
					),
				),
			),
			elem.Div(
				vecty.Markup(
					vecty.Class("mdc-radio__ripple"),
				),
			),
			elem.Div(
				vecty.Markup(
					vecty.Class("mdc-radio__focus-ring"),
				),
			),
		),
		elem.Label(
			vecty.Markup(
				prop.For(idInput),
			),
			vecty.Text(ch.Label),
		),
	).WithClasses(ch.Classes.Classes...)
}
