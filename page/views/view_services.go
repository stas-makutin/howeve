package views

import (
	"fmt"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/stas-makutin/howeve/api"
	"github.com/stas-makutin/howeve/page/actions"
	"github.com/stas-makutin/howeve/page/components"
	"github.com/stas-makutin/howeve/page/core"
)

type ViewServices struct {
	vecty.Core
	rendered         bool
	addServiceDialog bool
	loading          bool
	useSockets       bool
	errorMessage     string
	protocols        *api.ProtocolInfoResult
	services         *api.ListServicesResult
}

func NewViewServices() (r *ViewServices) {
	store := actions.GetServicesViewStore()
	r = &ViewServices{
		rendered:         false,
		addServiceDialog: false,
		loading:          store.Loading > 0,
		useSockets:       store.UseSocket,
		errorMessage:     store.Error,
		protocols:        store.Protocols,
		services:         store.Services,
	}
	actions.Subscribe(r)
	return
}

func (ch *ViewServices) OnChange(event interface{}) {
	if store, ok := event.(*actions.ServicesViewStore); ok {
		ch.loading = store.Loading > 0
		ch.useSockets = store.UseSocket
		ch.errorMessage = store.Error
		ch.protocols = store.Protocols
		ch.services = store.Services
		if ch.rendered {
			vecty.Rerender(ch)
		}
	}
}

func (ch *ViewServices) Mount() {
	core.Dispatch(&actions.ServicesLoad{Force: false, UseSocket: ch.useSockets})
}

func (ch *ViewServices) changeUseSocket(checked, disabled bool) {
	core.Dispatch(actions.ServicesUseSocket(checked))
}

func (ch *ViewServices) refresh() {
	core.Dispatch(&actions.ServicesLoad{Force: true, UseSocket: ch.useSockets})
}

func (ch *ViewServices) addService() {
	ch.addServiceDialog = true
	vecty.Rerender(ch)
}

func (ch *ViewServices) addServiceAction(action string) {
	ch.addServiceDialog = false
	vecty.Rerender(ch)
}

func (ch *ViewServices) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *ViewServices) Render() vecty.ComponentOrHTML {
	ch.rendered = true
	return components.NewMdcGrid(
		components.NewMdcGridSingleCellRow(
			components.NewMdcButton("sv-add-service", "Add Service", ch.protocols == nil, ch.addService).WithClasses("adjacent-margins"),
			components.NewMdcButton("sv-refresh", "Refresh", false, ch.refresh),
			components.NewMdcCheckbox("sv-socket-check", "Use WebSocket", ch.useSockets, false, ch.changeUseSocket),
		),
		core.If(ch.errorMessage != "", components.NewMdcGridSingleCellRow(
			components.NewMdcBanner("sv-error-banner", ch.errorMessage, "Retry", ch.refresh),
		)),
		&components.SectionTitle{Text: "Services"},
		components.NewMdcGridSingleCellRow(
			&servicesTable{Services: ch.services},
		),
		core.If(ch.addServiceDialog, &addServicesDialog{Protocols: ch.protocols, CloseFn: ch.addServiceAction}),
		core.If(ch.loading, &components.ViewLoading{}),
	)
}

type servicesTable struct {
	vecty.Core
	Services *api.ListServicesResult `vecty:"prop"`
}

func (ch *servicesTable) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *servicesTable) Render() vecty.ComponentOrHTML {
	return nil
}

type addServicesDialog struct {
	vecty.Core
	Protocols *api.ProtocolInfoResult
	CloseFn   func(action string)
	protocol  api.ProtocolIdentifier
	transport api.TransportIdentifier
}

func (ch *addServicesDialog) changeProtocol(value string, index int) {
}

func (ch *addServicesDialog) changeTransport(value string, index int) {
}

func (ch *addServicesDialog) addParameter() {
}

func (ch *addServicesDialog) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *addServicesDialog) Render() vecty.ComponentOrHTML {
	protocolsKey := fmt.Sprintf("sv-add-p-%p", ch.Protocols)

	return components.NewMdcDialog(
		"sv-add-service-dialog", "Add Service", true, true, ch.CloseFn,
		[]components.MdcDialogButton{
			{Label: "Cancel", Action: components.MdcDialogActionClose},
			{Label: "Add Service", Action: components.MdcDialogActionOK, Default: true, Disabled: false},
		},
		ch.RenderProtocols(protocolsKey),
		ch.RenderTransports(protocolsKey),
		components.NewMdcTextField("sv-add-service-alias", "Alias", "", false).WithClasses("adjacent-margins").WithKey("text-alias"),
		elem.Break(
			vecty.Markup(
				vecty.Key("br-1"),
			),
		),
		components.NewMdcTextField("sv-add-service-entry", "Entry", "", false).WithClasses("adjacent-margins").WithKey("text-entry"),
		elem.Div(
			vecty.Markup(
				vecty.Key("param-title"),
				vecty.Class("mdc-typography--overline"),
			),
			vecty.Text("Parameters"),
		),
		components.NewMdcButton("sv-add-service-add-param", "Add Parameter", false, ch.addParameter).WithKey("add-param-btn"),
	)
}

func (ch *addServicesDialog) RenderProtocols(protocolsKey string) vecty.ComponentOrHTML {
	var options vecty.List
	if ch.Protocols != nil {
		notFound := true
		for _, protocol := range ch.Protocols.Protocols {
			selected := protocol.ID == ch.protocol
			options = append(options, &components.MdcSelectOption{
				Name:     fmt.Sprintf("%s (%d)", protocol.Name, protocol.ID),
				Selected: selected,
			})
			if selected {
				notFound = false
			}
		}
		if notFound && len(ch.Protocols.Protocols) > 0 {
			ch.protocol = ch.Protocols.Protocols[0].ID
			options[0].(*components.MdcSelectOption).Selected = true
		}
	}

	return components.NewMdcSelect(
		"sv-add-service-protocols", "Protocols", false, ch.changeProtocol, options,
	).WithKey(protocolsKey)
}

func (ch *addServicesDialog) RenderTransports(protocolsKey string) vecty.ComponentOrHTML {
	transportsKey := fmt.Sprintf("%s-%d", protocolsKey, ch.protocol)

	var protocol *api.ProtocolInfoEntry
	if ch.Protocols != nil {
		for _, p := range ch.Protocols.Protocols {
			if p.ID == ch.protocol {
				protocol = p
				break
			}
		}
	}

	var options vecty.List
	if protocol != nil {
		notFound := true
		for _, transport := range protocol.Transports {
			selected := transport.ID == ch.transport
			options = append(options, &components.MdcSelectOption{
				Name:     fmt.Sprintf("%s (%d)", transport.Name, transport.ID),
				Selected: selected,
			})
			if selected {
				notFound = false
			}
		}
		if notFound && len(protocol.Transports) > 0 {
			ch.transport = protocol.Transports[0].ID
			options[0].(*components.MdcSelectOption).Selected = true
		}
	}

	return components.NewMdcSelect(
		"sv-add-service-transports", "Transports", false, ch.changeTransport, options,
	).WithKey(transportsKey)
}
