package main

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

type viewProtocols struct {
	vecty.Core
}

func (ch *viewProtocols) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			vecty.Class("mdc-layout-grid__inner"),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("mdc-layout-grid__cell", "mdc-layout-grid__cell--span-12"),
			),
			vecty.Text("Protocols View"),
		),
	)
}
