package main

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

type viewMessages struct {
	vecty.Core
}

func (ch *viewMessages) Render() vecty.ComponentOrHTML {
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
