package components

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
	"github.com/hexops/vecty/prop"
)

type MdcIconButton struct {
	vecty.Core
	ID       string `vecty:"prop"`
	Icon     string `vecty:"prop"`
	Text     string `vecty:"prop"`
	Disabled bool   `vecty:"prop"`
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
			prop.Disabled(ch.Disabled),
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
