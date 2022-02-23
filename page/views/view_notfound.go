package views

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/prop"
)

type ViewNotFound struct {
	vecty.Core
}

func (ch *ViewNotFound) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *ViewNotFound) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			vecty.Class("mdc-layout-grid"),
			vecty.Style("text-align", "center"),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("mdc-layout-grid__inner"),
			),
			elem.Div(
				vecty.Markup(
					vecty.Class("mdc-layout-grid__cell", "mdc-layout-grid__cell--span-12"),
				),
				elem.Button(
					vecty.Markup(
						vecty.Class("mdc-button", "mdc-button--icon-leading", "mdc-theme--text-primary-on-light"),
						prop.Disabled(true),
					),
					elem.Italic(
						vecty.Markup(
							vecty.Class("material-icons", "mdc-button__icon"),
							vecty.Attribute("aria-hidden", "true"),
						),
						vecty.Text("web_asset_off"),
					),
					elem.Span(
						vecty.Markup(
							vecty.Class("mdc-button__label"),
						),
						vecty.Text("Page Not Found"),
					),
				),
			),
		),
	)
}
