package components

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

type SectionTitle struct {
	vecty.Core
	Text string
}

func (ch *SectionTitle) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *SectionTitle) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			vecty.Class("mdc-layout-grid__inner"),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("mdc-layout-grid__cell", "mdc-typography--overline"),
			),
			vecty.Text(ch.Text),
		),
	)
}
