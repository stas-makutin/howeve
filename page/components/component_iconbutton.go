package components

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
	"github.com/hexops/vecty/prop"
	"github.com/stas-makutin/howeve/page/core"
)

type MdcIconButton struct {
	vecty.Core
	core.Classes
	core.Keyable
	ID       string `vecty:"prop"`
	Icon     string `vecty:"prop"`
	IconOn   string `vecty:"prop"`
	Text     string `vecty:"prop"`
	Disabled bool   `vecty:"prop"`
	clickFn  func()
}

func NewMdcIconButton(id, text, icon, iconOn string, disabled bool, clickFn func()) (r *MdcIconButton) {
	r = &MdcIconButton{ID: id, Icon: icon, IconOn: iconOn, Text: text, Disabled: disabled, clickFn: clickFn}
	return
}

func (ch *MdcIconButton) WithKey(key interface{}) *MdcIconButton {
	ch.Keyable.WithKey(key)
	return ch
}

func (ch *MdcIconButton) WithClasses(classes ...string) *MdcIconButton {
	ch.Classes.WithClasses(classes...)
	return ch
}

func (ch *MdcIconButton) onClick(event *vecty.Event) {
	ch.clickFn()
}

func (ch *MdcIconButton) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *MdcIconButton) Render() vecty.ComponentOrHTML {
	return elem.Button(
		vecty.Markup(
			prop.ID(ch.ID),
			vecty.Class("mdc-icon-button"),
			vecty.Attribute("title", ch.Text),
			vecty.Attribute("aria-label", ch.Text),
			vecty.Attribute("aria-pressed", "false"),
			prop.Disabled(ch.Disabled),
			event.Click(ch.onClick),
			ch.ApplyClasses(),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("mdc-icon-button__ripple"),
			),
		),
		elem.Span(
			vecty.Markup(
				vecty.Class("mdc-icon-button__focus-ring"),
			),
		),
		elem.Italic(
			vecty.Markup(
				vecty.Class("material-icons", "mdc-icon-button__icon", "mdc-icon-button__icon--on"),
			),
			vecty.Text(ch.IconOn),
		),
		elem.Italic(
			vecty.Markup(
				vecty.Class("material-icons", "mdc-icon-button__icon"),
			),
			vecty.Text(ch.Icon),
		),
	)
}
