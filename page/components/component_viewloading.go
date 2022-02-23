package components

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

type ViewLoading struct {
	vecty.Core
}

func (ch *ViewLoading) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *ViewLoading) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			vecty.Class("view-loading", "mdc-theme--surface"),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("view-loading__progress"),
			),
		),
	)
}
