package components

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

type Title struct {
	vecty.Core
}

func (ch *Title) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *Title) Render() vecty.ComponentOrHTML {
	return elem.Section(
		vecty.Markup(
			vecty.Class("app-Title"),
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
