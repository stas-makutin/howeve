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
	rendered   bool
	loading    bool
	useSockets bool
}

func NewViewProtocols() (r *ViewProtocols) {
	store := actions.GetProtocolViewStore()
	r = &ViewProtocols{
		rendered:   false,
		loading:    store.Loading,
		useSockets: store.UseSocket,
	}
	actions.Subscribe(r)
	return
}

func (ch *ViewProtocols) OnChange(event interface{}) {
	if store, ok := event.(*actions.ProtocolViewStore); ok {
		ch.loading = store.Loading
		ch.useSockets = store.UseSocket
		if ch.rendered {
			vecty.Rerender(ch)
		}
	}
}

func (ch *ViewProtocols) Mount() {
	core.Dispatch(&actions.ProtocolsLoad{Force: false, UseSocket: ch.useSockets})
}

func (ch *ViewProtocols) changeUseSocket(checked, disabled bool) {
	core.Dispatch(actions.ProtocolsUseSocket(checked))
}

func (ch *ViewProtocols) refresh() {
	core.Dispatch(&actions.ProtocolsLoad{Force: true, UseSocket: ch.useSockets})
}

func (ch *ViewProtocols) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *ViewProtocols) Render() vecty.ComponentOrHTML {
	ch.rendered = true
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
				components.NewMdcCheckbox("pt-socket-check", "Use WebSocket", ch.useSockets, false, ch.changeUseSocket),
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
