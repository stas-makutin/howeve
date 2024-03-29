package views

import (
	"strconv"
	"strings"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/stas-makutin/howeve/api"
	"github.com/stas-makutin/howeve/page/actions"
	"github.com/stas-makutin/howeve/page/components"
	"github.com/stas-makutin/howeve/page/core"
)

type ViewProtocols struct {
	vecty.Core
	rendered     bool
	loading      bool
	useSockets   bool
	errorMessage []vecty.MarkupOrChild
	protocols    *api.ProtocolInfoResult
}

func NewViewProtocols() (r *ViewProtocols) {
	store := actions.GetProtocolViewStore()
	r = &ViewProtocols{
		rendered:     false,
		loading:      store.Loading,
		useSockets:   store.UseSocket,
		errorMessage: core.FormatMultilineText(store.DisplayError),
		protocols:    store.Protocols.Value,
	}
	actions.Subscribe(r)
	return
}

func (ch *ViewProtocols) OnChange(event interface{}) {
	if store, ok := event.(*actions.ProtocolViewStore); ok {
		ch.loading = store.Loading
		ch.useSockets = store.UseSocket
		ch.errorMessage = core.FormatMultilineText(store.DisplayError)
		ch.protocols = store.Protocols.Value
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
	return components.NewMdcGrid(
		components.NewMdcGridSingleCellRow(
			components.NewMdcButton("pt-refresh", "Refresh", false, ch.refresh),
			components.NewMdcCheckbox("pt-socket-check", "Use WebSocket", ch.useSockets, false, ch.changeUseSocket),
		),
		core.If(len(ch.errorMessage) > 0, components.NewMdcGridSingleCellRow(
			components.NewMdcBanner("pt-error-banner", "Retry", true, ch.refresh, ch.errorMessage...),
		)),
		&components.SectionTitle{Text: "Protocols"},
		components.NewMdcGridSingleCellRow(
			&protocolsTable{Protocols: ch.protocols},
		),
		core.If(ch.loading, &components.ViewLoading{}),
	)
}

type protocolsTable struct {
	vecty.Core
	Protocols *api.ProtocolInfoResult `vecty:"prop"`
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
				ch.tableBody(),
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

func (ch *protocolsTable) tableBody() vecty.ComponentOrHTML {
	if ch.Protocols == nil || len(ch.Protocols.Protocols) <= 0 {
		return nil
	}

	var content vecty.List
	for _, protocol := range ch.Protocols.Protocols {
		for _, transport := range protocol.Transports {
			content = append(content, ch.tableRow(protocol, transport))
		}
	}

	return elem.TableBody(
		vecty.Markup(
			vecty.Class("mdc-data-table__content"),
		),
		content,
	)
}

func (ch *protocolsTable) tableRow(protocol *api.ProtocolInfoEntry, transport *api.ProtocolTransportInfoEntry) vecty.ComponentOrHTML {
	var dicoverable string
	if transport.Discoverable {
		dicoverable = "Yes"
	} else {
		dicoverable = "No"
	}
	return elem.TableRow(
		vecty.Markup(
			vecty.Class("mdc-data-table__row"),
		),
		ch.tableColumn(vecty.Text(protocol.Name)),
		ch.tableColumn(vecty.Text(strconv.Itoa(int(protocol.ID))), "mdc-data-table__cell--numeric"),
		ch.tableColumn(vecty.Text(transport.Name)),
		ch.tableColumn(vecty.Text(strconv.Itoa(int(transport.ID))), "mdc-data-table__cell--numeric"),
		ch.tableColumn(ch.parametersTable(transport.Params)),
		ch.tableColumn(vecty.Text(dicoverable)),
		ch.tableColumn(ch.parametersTable(transport.Params)),
	)
}

func (ch *protocolsTable) tableColumn(content vecty.MarkupOrChild, classes ...string) vecty.ComponentOrHTML {
	return elem.TableData(
		vecty.Markup(
			vecty.Class(append([]string{"mdc-data-table__cell", "data-table-cell--top"}, classes...)...),
			vecty.Attribute("role", "columnheader"),
			vecty.Attribute("scope", "col"),
		),
		content,
	)
}

func (ch *protocolsTable) parametersTable(params map[string]*api.ParamInfoEntry) vecty.ComponentOrHTML {
	if len(params) <= 0 {
		return vecty.Text("None")
	}

	names := core.ArrangeParams(params)

	return components.NewKeyValueTable(func(builder components.KeyValueTableBuilder) {
		for i, name := range names {
			info := params[name]
			if i > 0 {
				builder.AddDelimiterRow()
			}
			builder.AddKeyValueRow("Name", name)
			if info.Description != "" {
				builder.AddKeyValueRow("Description", info.Description)
			}
			builder.AddKeyValueRow("Type", info.Type)
			if info.DefaultValue != "" {
				builder.AddKeyValueRow("Default Value", info.DefaultValue)
			}
			if info.Type == "enum" {
				builder.AddKeyValueRow("Allowed Values", strings.Join(info.EnumValues, ", "))
			}
		}
	})
}
