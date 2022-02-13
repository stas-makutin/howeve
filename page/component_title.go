package main

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

type title struct {
	vecty.Core
}

func (ch *title) Render() vecty.ComponentOrHTML {
	return elem.Section(
		vecty.Markup(
			vecty.Class("title"),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("mdc-typography--subtitle1", "mdc-theme--secondary"),
			),
			vecty.Text("HOWEVE"),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("mdc-typography--caption", "mdc-theme--text-icon-on-light"),
			),
			vecty.Text("TEST PAGE"),
		),
	)
}
