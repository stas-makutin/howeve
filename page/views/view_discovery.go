package views

import (
	"github.com/google/uuid"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/stas-makutin/howeve/page/actions"
	"github.com/stas-makutin/howeve/page/components"
	"github.com/stas-makutin/howeve/page/core"
)

const (
	DiscoveryDialog_None = iota
	DiscoveryDialog_StartDiscovery
)

type ViewDiscovery struct {
	vecty.Core
	rendered     bool
	loading      bool
	renderDialog int
	protocols    *core.ProtocolsWrapper
	discoveries  map[uuid.UUID]*actions.DiscoveryData
	errorMessage []vecty.MarkupOrChild
}

func NewViewDiscovery() (r *ViewDiscovery) {
	store := actions.GetDiscoveryViewStore()
	r = &ViewDiscovery{
		rendered:     false,
		loading:      store.Loading,
		renderDialog: DiscoveryDialog_None,
		protocols:    core.NewProtocolsWrapper(store.Protocols),
		discoveries:  store.Discoveries,
	}
	actions.Subscribe(r)
	core.Dispatch(actions.DiscoveryLoad{})
	return
}

func (ch *ViewDiscovery) OnChange(event interface{}) {
	if store, ok := event.(*actions.DiscoveryViewStore); ok {
		ch.loading = store.Loading
		ch.protocols = core.NewProtocolsWrapper(store.Protocols)
		if ch.rendered {
			vecty.Rerender(ch)
		}
		core.Dispatch(actions.DiscoveryLoad{})
	}
}

func (ch *ViewDiscovery) load() {
	if ch.protocols == nil && !ch.loading {
		core.Dispatch(actions.DiscoveryLoad{})
	}
}

func (ch *ViewDiscovery) retry() {
}

func (ch *ViewDiscovery) toggleDialog(dialog int) {
	ch.renderDialog = dialog
	vecty.Rerender(ch)
}

func (ch *ViewDiscovery) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *ViewDiscovery) Render() vecty.ComponentOrHTML {
	ch.rendered = true
	return components.NewMdcGrid(
		components.NewMdcGridSingleCellRow(
			components.NewMdcButton("dv-start-discovery", "Start Discovery", ch.protocols == nil,
				func() { ch.toggleDialog(DiscoveryDialog_StartDiscovery) },
			).WithClasses("adjacent-margins"),
		),
		core.If(len(ch.errorMessage) > 0, components.NewMdcGridSingleCellRow(
			components.NewMdcBanner("sv-error-banner", "Retry", true, ch.retry, ch.errorMessage...),
		)),
		&components.SectionTitle{Text: "Discoveries"},
		components.NewMdcGridSingleCellRow(
			&discoveriesTable{Protocols: ch.protocols},
		),
		core.If(ch.loading, &components.ViewLoading{}),
	)
}

type discoveriesTable struct {
	vecty.Core
	Protocols *core.ProtocolsWrapper `vecty:"prop"`
}

func (ch *discoveriesTable) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *discoveriesTable) headerColumn(name string, classes ...string) vecty.ComponentOrHTML {
	return elem.TableHeader(
		vecty.Markup(
			vecty.Class(append([]string{"mdc-data-table__header-cell"}, classes...)...),
			vecty.Attribute("role", "columnheader"),
			vecty.Attribute("scope", "col"),
		),
		vecty.Text(name),
	)
}

func (ch *discoveriesTable) tableBody() vecty.ComponentOrHTML {
	return nil
	// if ch.Services == nil || len(ch.Services.Services) <= 0 {
	// 	return nil
	// }

	// var content vecty.List
	// for _, service := range ch.Services.Services {
	// 	content = append(content, ch.tableRow(&service))
	// }

	// return elem.TableBody(
	// 	vecty.Markup(
	// 		vecty.Class("mdc-data-table__content"),
	// 	),
	// 	content,
	// )
}

func (ch *discoveriesTable) Render() vecty.ComponentOrHTML {
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
					vecty.Attribute("aria-label", "Discoveries"),
				),
				elem.TableHead(
					elem.TableRow(
						vecty.Markup(
							vecty.Class("mdc-data-table__header-row"),
						),
						ch.headerColumn("Discovery ID"),
						ch.headerColumn("Protocol"),
						ch.headerColumn("Transport"),
						ch.headerColumn("Parameters"),
						ch.headerColumn("Status"),
						ch.headerColumn("Stop"),
					),
				),
				ch.tableBody(),
			),
		),
	)
}
