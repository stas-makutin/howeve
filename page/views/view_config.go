package views

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

type ViewConfig struct {
	vecty.Core
}

func (ch *ViewConfig) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *ViewConfig) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			vecty.Class("mdc-layout-grid__inner"),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("mdc-layout-grid__cell", "mdc-layout-grid__cell--span-12"),
			),
			vecty.Text("Config View"),
		),
	)
}
