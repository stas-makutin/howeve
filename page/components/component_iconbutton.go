package components

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
	"github.com/hexops/vecty/prop"
)

type MdcIconButton struct {
	vecty.Core
	ID       string
	Icon     string
	Text     string
	Disabled bool
	clickFn  func()
}

func NewMdcIconButton(id, text, icon string, disabled bool, clickFn func()) (r *MdcIconButton) {
	r = &MdcIconButton{ID: id, Icon: icon, Text: text, Disabled: disabled, clickFn: clickFn}
	return
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
			vecty.Class("mdc-icon-button", "material-icons"),
			vecty.Attribute("title", ch.Text),
			vecty.MarkupIf(
				ch.Disabled,
				prop.Disabled(true),
			),
			event.Click(ch.onClick),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("mdc-button__ripple"),
			),
		),
		elem.Span(
			vecty.Markup(
				vecty.Class("mdc-button__focus-ring"),
			),
		),
		vecty.Text(ch.Icon),
	)
}
