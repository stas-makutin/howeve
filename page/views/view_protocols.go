package views

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/stas-makutin/howeve/page/actions"
	"github.com/stas-makutin/howeve/page/components"
	"github.com/stas-makutin/howeve/page/core"
)

type ViewProtocols struct {
	vecty.Core
	loading bool
}

func NewViewProtocols() (r *ViewProtocols) {
	r = &ViewProtocols{loading: true}
	actions.SubscribeGlobal(r)
	return
}

func (ch *ViewProtocols) OnLoad() {
	actions.ProtocolViewRegister()
}

func (ch *ViewProtocols) OnChange(event interface{}) {
	if s, ok := event.(*actions.ProtocolViewStore); ok {
		ch.loading = s.Loading
		vecty.Rerender(ch)
	}
}

func (ch *ViewProtocols) refresh() {
	core.Dispatch(actions.ProtocolsLoadFetch)
}

func (ch *ViewProtocols) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *ViewProtocols) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			vecty.Class("mdc-layout-grid"),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("mdc-layout-grid__inner"),
			),

			elem.Div(
				vecty.Markup(
					vecty.Class("mdc-layout-grid__cell"),
				),
				components.NewMdcButton("pt-refresh", "Refresh", false),
				components.NewMdcCheckbox("pt-socket-check", "Use WebSocket", false, false),
			),

			elem.Div(
				vecty.Markup(
					vecty.Class("mdc-layout-grid__cell", "mdc-layout-grid__cell--span-12"),
				),
				&protocolsTable{},
			),
		),
		vecty.If(ch.loading, &components.ViewLoading{}),
	)
}

type protocolsTable struct {
	vecty.Core
}

func (ch *protocolsTable) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *protocolsTable) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			vecty.Class("mdc-data-table"),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("mdc-data-table__table-container"),
			),
			elem.Table(
				vecty.Markup(
					vecty.Class("mdc-data-table__table"),
					vecty.Attribute("aria-label", "Protocols"),
				),
				elem.TableHead(
					elem.TableRow(
						vecty.Markup(
							vecty.Class("mdc-data-table__header-row"),
						),
						ch.headerColumn("Protocol"),
						ch.headerColumn("Protocol ID", "mdc-data-table__header-cell--numeric"),
						ch.headerColumn("Transport"),
						ch.headerColumn("Transport ID", "mdc-data-table__header-cell--numeric"),
						ch.headerColumn("Parameters"),
						ch.headerColumn("Discoverable"),
						ch.headerColumn("Discovery Parameters"),
					),
				),
				elem.TableBody(
					vecty.Markup(
						vecty.Class("mdc-data-table__content"),
					),
				),
			),
		),
	)
}

func (ch *protocolsTable) headerColumn(name string, classes ...string) vecty.ComponentOrHTML {
	return elem.TableHeader(
		vecty.Markup(
			vecty.Class(append([]string{"mdc-data-table__header-cell"}, classes...)...),
			vecty.Attribute("role", "columnheader"),
			vecty.Attribute("scope", "col"),
		),
		vecty.Text(name),
	)
}
