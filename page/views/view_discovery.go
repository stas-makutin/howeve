package views

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

type ViewDiscovery struct {
	vecty.Core
}

func (ch *ViewDiscovery) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *ViewDiscovery) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			vecty.Class("mdc-layout-grid__inner"),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("mdc-layout-grid__cell", "mdc-layout-grid__cell--span-12"),
			),
			vecty.Text("Discovery View"),
		),
	)
}
