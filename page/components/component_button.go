package components

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/prop"
)

type MdcButton struct {
	vecty.Core
	ID       string
	Text     string
	Disabled bool
}

func NewMdcButton(id string, text string, disabled bool) (r *MdcButton) {
	r = &MdcButton{ID: id, Text: text, Disabled: disabled}
	return
}

func (ch *MdcButton) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *MdcButton) Render() vecty.ComponentOrHTML {
	return elem.Button(
		vecty.Markup(
			prop.ID(ch.ID),
			vecty.Class("mdc-button", "mdc-button--outlined"),
			vecty.MarkupIf(
				ch.Disabled,
				prop.Disabled(true),
			),
		),
		elem.Span(
			vecty.Markup(
				vecty.Class("mdc-button__ripple"),
			),
		),
		elem.Span(
			vecty.Markup(
				vecty.Class("mdc-button__focus-ring"),
			),
		),
		elem.Span(
			vecty.Markup(
				vecty.Class("mdc-button__label"),
			),
			vecty.Text(ch.Text),
		),
	)
}
