package views

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

type ViewMessages struct {
	vecty.Core
}

func (ch *ViewMessages) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *ViewMessages) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			vecty.Class("mdc-layout-grid__inner"),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("mdc-layout-grid__cell", "mdc-layout-grid__cell--span-12"),
			),
			vecty.Text("Messages View"),
		),
	)
}
