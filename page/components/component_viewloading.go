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
			elem.Div(
				vecty.Markup(
					vecty.Class("mdc-circular-progress", "mdc-circular-progress--indeterminate"),
					vecty.Style("width", "48px"),
					vecty.Style("height", "48px"),
					vecty.Attribute("role", "progressbar"),
					vecty.Attribute("aria-label", "Page Is Loading"),
					vecty.Attribute("aria-label", "Page Is Loading"),
				),
				elem.Div(
					vecty.Markup(
						vecty.Class("mdc-circular-progress__determinate-container"),
					),
					vecty.Tag(
						"svg",
						vecty.Markup(
							vecty.Namespace("http://www.w3.org/2000/svg"),
							vecty.Class("mdc-circular-progress__determinate-circle-graphic"),
							vecty.Attribute("viewBox", "0 0 48 48"),
						),
						vecty.Tag(
							"circle",
							vecty.Markup(
								vecty.Namespace("http://www.w3.org/2000/svg"),
								vecty.Class("mdc-circular-progress__determinate-track"),
								vecty.Attribute("cx", "24"),
								vecty.Attribute("cy", "24"),
								vecty.Attribute("r", "18"),
								vecty.Attribute("stroke-width", "4"),
							),
						),
						vecty.Tag(
							"circle",
							vecty.Markup(
								vecty.Namespace("http://www.w3.org/2000/svg"),
								vecty.Class("mdc-circular-progress__determinate-circle"),
								vecty.Attribute("cx", "24"),
								vecty.Attribute("cy", "24"),
								vecty.Attribute("r", "18"),
								vecty.Attribute("stroke-width", "4"),
								vecty.Attribute("stroke-dasharray", "113.097"),
								vecty.Attribute("stroke-dashoffset", "113.097"),
							),
						),
					),
				),
				elem.Div(
					vecty.Markup(
						vecty.Class("mdc-circular-progress__indeterminate-container"),
					),
					elem.Div(
						vecty.Markup(
							vecty.Class("mdc-circular-progress__spinner-layer"),
						),
						elem.Div(
							vecty.Markup(
								vecty.Class("mdc-circular-progress__circle-clipper", "mdc-circular-progress__circle-left"),
							),
							vecty.Tag(
								"svg",
								vecty.Markup(
									vecty.Namespace("http://www.w3.org/2000/svg"),
									vecty.Class("mdc-circular-progress__indeterminate-circle-graphic"),
									vecty.Attribute("viewBox", "0 0 48 48"),
								),
								vecty.Tag(
									"circle",
									vecty.Markup(
										vecty.Namespace("http://www.w3.org/2000/svg"),
										vecty.Attribute("cx", "24"),
										vecty.Attribute("cy", "24"),
										vecty.Attribute("r", "18"),
										vecty.Attribute("stroke-width", "4"),
										vecty.Attribute("stroke-dasharray", "113.097"),
										vecty.Attribute("stroke-dashoffset", "56.549"),
									),
								),
							),
						),
						elem.Div(
							vecty.Markup(
								vecty.Class("mdc-circular-progress__gap-patch"),
							),
							vecty.Tag(
								"svg",
								vecty.Markup(
									vecty.Namespace("http://www.w3.org/2000/svg"),
									vecty.Class("mdc-circular-progress__indeterminate-circle-graphic"),
									vecty.Attribute("viewBox", "0 0 48 48"),
								),
								vecty.Tag(
									"circle",
									vecty.Markup(
										vecty.Namespace("http://www.w3.org/2000/svg"),
										vecty.Attribute("cx", "24"),
										vecty.Attribute("cy", "24"),
										vecty.Attribute("r", "18"),
										vecty.Attribute("stroke-width", "4"),
										vecty.Attribute("stroke-dasharray", "113.097"),
										vecty.Attribute("stroke-dashoffset", "56.549"),
									),
								),
							),
						),
						elem.Div(
							vecty.Markup(
								vecty.Class("mdc-circular-progress__circle-clipper", "mdc-circular-progress__circle-right"),
							),
							vecty.Tag(
								"svg",
								vecty.Markup(
									vecty.Namespace("http://www.w3.org/2000/svg"),
									vecty.Class("mdc-circular-progress__indeterminate-circle-graphic"),
									vecty.Attribute("viewBox", "0 0 48 48"),
								),
								vecty.Tag(
									"circle",
									vecty.Markup(
										vecty.Namespace("http://www.w3.org/2000/svg"),
										vecty.Attribute("cx", "24"),
										vecty.Attribute("cy", "24"),
										vecty.Attribute("r", "18"),
										vecty.Attribute("stroke-width", "4"),
										vecty.Attribute("stroke-dasharray", "113.097"),
										vecty.Attribute("stroke-dashoffset", "56.549"),
									),
								),
							),
						),
					),
				),
			),
		),
	)
}
