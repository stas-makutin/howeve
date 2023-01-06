package views

import (
	"fmt"
	"strconv"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
	"github.com/hexops/vecty/prop"
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
	protocols        *core.ProtocolsWrapper
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
		protocols:        core.NewProtocolsWrapper(store.Protocols),
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
		ch.protocols = core.NewProtocolsWrapper(store.Protocols)
		//ch.services = store.Services
		ch.services = &api.ListServicesResult{
			Services: []api.ListServicesEntry{
				{
					ServiceEntry: &api.ServiceEntry{
						ServiceKey: &api.ServiceKey{
							Protocol:  api.ProtocolZWave,
							Transport: api.TransportSerial,
							Entry:     "",
						},
						Params: map[string]string{
							"param1": "value1",
							"param2": "value2",
						},
					},
					StatusReply: &api.StatusReply{
						Success: true,
					},
				},
				{
					ServiceEntry: &api.ServiceEntry{
						ServiceKey: &api.ServiceKey{
							Protocol:  api.ProtocolZWave,
							Transport: api.TransportSerial,
							Entry:     "COM2",
						},
						Params: map[string]string{
							"paramA": "valueA",
							"paramB": "valueB",
						},
						Alias: "ZC2",
					},
					StatusReply: &api.StatusReply{
						Success: false,
						Error: &api.ErrorInfo{
							Code:    api.ErrorServiceStatusBad,
							Message: "unable to connect to the service",
						},
					},
				},
				{
					ServiceEntry: &api.ServiceEntry{
						ServiceKey: &api.ServiceKey{
							Protocol:  api.ProtocolZWave,
							Transport: api.TransportSerial,
							Entry:     "https://github.com/material-components/material-components-web/blob/8f0a11e32895f998c326ab4a10601a2e4d5e18db/packages/mdc-textfield/README.md",
						},
						Params: map[string]string{
							"paramX": "valueX",
							"paramY": "valueY",
							"paramZ": "valueZ",
						},
					},
					StatusReply: &api.StatusReply{
						Success: true,
					},
				},
			},
		}
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

func (ch *ViewServices) addServiceAction(ok bool, data *core.ServiceEntryData) {
	ch.addServiceDialog = false
	if ok {
		core.Console.Log(fmt.Sprintf("%d: %d: %s", data.Protocol, data.Transport, data.Entry))
	}
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
			&servicesTable{Protocols: ch.protocols, Services: ch.services},
		),
		core.If(ch.addServiceDialog, newAddServiceDialog(ch.protocols, ch.services, ch.addServiceAction)),
		core.If(ch.loading, &components.ViewLoading{}),
	)
}

type servicesTable struct {
	vecty.Core
	Protocols *core.ProtocolsWrapper  `vecty:"prop"`
	Services  *api.ListServicesResult `vecty:"prop"`
}

func (ch *servicesTable) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *servicesTable) changeAlias(service api.ServiceKey) {
	core.Console.Log(fmt.Sprintf("changeAlias: %d %d %s", service.Protocol, service.Transport, service.Entry))
}

func (ch *servicesTable) viewParameters(service api.ServiceKey) {
	core.Console.Log(fmt.Sprintf("viewParameters: %d %d %s", service.Protocol, service.Transport, service.Entry))
}

func (ch *servicesTable) viewStatus(service api.ServiceKey) {
	core.Console.Log(fmt.Sprintf("viewStatus: %d %d %s", service.Protocol, service.Transport, service.Entry))
}

func (ch *servicesTable) removeService(service api.ServiceKey) {
	core.Console.Log(fmt.Sprintf("removeService: %d %d %s", service.Protocol, service.Transport, service.Entry))
}

func (ch *servicesTable) headerColumn(name string, classes ...string) vecty.ComponentOrHTML {
	return elem.TableHeader(
		vecty.Markup(
			vecty.Class(append([]string{"mdc-data-table__header-cell"}, classes...)...),
			vecty.Attribute("role", "columnheader"),
			vecty.Attribute("scope", "col"),
		),
		vecty.Text(name),
	)
}

func (ch *servicesTable) tableBody() vecty.ComponentOrHTML {
	if ch.Services == nil || len(ch.Services.Services) <= 0 {
		return nil
	}

	var content vecty.List
	for _, service := range ch.Services.Services {
		content = append(content, ch.tableRow(&service))
	}

	return elem.TableBody(
		vecty.Markup(
			vecty.Class("mdc-data-table__content"),
		),
		content,
	)
}

func (ch *servicesTable) tableRow(service *api.ListServicesEntry) vecty.ComponentOrHTML {
	serviceKey := *service.ServiceKey // copy for call enclosures

	pi, pti := ch.Protocols.ProtocolAndTransport(service.Protocol, service.Transport)
	protocolName := strconv.Itoa(int(service.Protocol))
	if pi != nil {
		protocolName = pi.Name + " (" + protocolName + ")"
	}
	transportName := strconv.Itoa(int(service.Transport))
	if pti != nil {
		transportName = pti.Name + " (" + transportName + ")"
	}

	var status *vecty.HTML
	if service.Success {
		status = elem.Span(
			vecty.Markup(
				vecty.Class("sv-service-table-status-healthy"),
			),
			vecty.Text("OK"),
		)
	} else {
		status = elem.Anchor(
			vecty.Markup(
				vecty.Class("sv-service-table-status-unhealthy"),
				prop.Href("#"),
				event.Click(func(e *vecty.Event) { ch.viewStatus(serviceKey) }).PreventDefault(),
			),
			vecty.Text("Unhealthy"),
		)
	}

	return elem.TableRow(
		vecty.Markup(
			vecty.Class("mdc-data-table__row"),
		),
		ch.tableColumn(vecty.Text(service.Alias)),
		ch.tableColumn(vecty.Text(protocolName)),
		ch.tableColumn(vecty.Text(transportName)),
		ch.tableColumn(vecty.Text(service.Entry), "sv-service-table-entry-cell"),
		ch.tableColumn(vecty.List{
			elem.Anchor(
				vecty.Markup(
					vecty.Class("sv-service-table-action"),
					prop.Href("#"),
					event.Click(func(e *vecty.Event) { ch.changeAlias(serviceKey) }).PreventDefault(),
				),
				vecty.Text("Change Alias"),
			),
			vecty.Text(", "),
			elem.Anchor(
				vecty.Markup(
					vecty.Class("sv-service-table-action"),
					prop.Href("#"),
					event.Click(func(e *vecty.Event) { ch.viewParameters(serviceKey) }).PreventDefault(),
				),
				vecty.Text("View Parameters"),
			),
		}, "sv-service-table-action-cell"),
		ch.tableColumn(status),
		ch.tableColumn(components.NewMdcIconButton("", "Remove Service", "delete_forever", "delete_forever", false,
			func() { ch.removeService(serviceKey) },
		), "sv-service-table-remove-cell"),
	)
}

func (ch *servicesTable) tableColumn(content vecty.MarkupOrChild, classes ...string) vecty.ComponentOrHTML {
	return elem.TableData(
		vecty.Markup(
			vecty.Class(append([]string{"mdc-data-table__cell", "data-table-cell--top"}, classes...)...),
			vecty.Attribute("role", "columnheader"),
			vecty.Attribute("scope", "col"),
		),
		content,
	)
}

func (ch *servicesTable) Render() vecty.ComponentOrHTML {
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
						ch.headerColumn("Alias"),
						ch.headerColumn("Protocol"),
						ch.headerColumn("Transport"),
						ch.headerColumn("Entry"),
						ch.headerColumn("Actions"),
						ch.headerColumn("Status"),
						ch.headerColumn("Remove"),
					),
				),
				ch.tableBody(),
			),
		),
	)
}

type addServiceDialog struct {
	vecty.Core
	Protocols *core.ProtocolsWrapper  `vecty:"prop"`
	Services  *api.ListServicesResult `vecty:"prop"`
	CloseFn   func(ok bool, data *core.ServiceEntryData)
	Data      *core.ServiceEntryData

	renderKey int
}

func newAddServiceDialog(protocols *core.ProtocolsWrapper, services *api.ListServicesResult, closeFn func(ok bool, data *core.ServiceEntryData)) *addServiceDialog {
	if services == nil {
		services = &api.ListServicesResult{}
	}
	return &addServiceDialog{
		Protocols: protocols, Services: services, CloseFn: closeFn, Data: core.NewServiceEntryData(protocols.Protocols),
	}
}

func (ch *addServiceDialog) closeDialog(action string, data interface{}) {
	ch.CloseFn(action == components.MdcDialogActionOK, data.(*core.ServiceEntryData))
}

func (ch *addServiceDialog) reRender() {
	ch.renderKey += 1
	vecty.Rerender(ch)
}

func (ch *addServiceDialog) changeProtocol(value string, index int) {
	if ch.Data.ChangeProtocol(ch.Protocols.Protocols, index) {
		ch.reRender()
	}
}

func (ch *addServiceDialog) changeTransport(value string, index int) {
	if ch.Data.ChangeTransport(ch.Protocols, index) {
		ch.reRender()
	}
}

func (ch *addServiceDialog) changeAlias(value string) {
	if value != ch.Data.Alias {
		ch.Data.Alias = value
		ch.reRender()
	}
}

func (ch *addServiceDialog) changeEntry(value string) {
	if value != ch.Data.Entry {
		ch.Data.Entry = value
		ch.reRender()
	}
}

func (ch *addServiceDialog) addParameter() {
	_, transport := ch.Data.ProtocolAndTransport(ch.Protocols)
	if params, ok := ch.Data.Params.AppendAvailable(transport); ok {
		ch.Data.Params = params
		ch.reRender()
	}
}

func (ch *addServiceDialog) changeParameter(paramIndex int, name string) {
	_, transport := ch.Data.ProtocolAndTransport(ch.Protocols)
	if ch.Data.Params.Replace(paramIndex, name, transport) {
		ch.reRender()
	}
}

func (ch *addServiceDialog) changeParameterValue(paramIndex int, value string) {
	if ch.Data.Params.ChangeValue(paramIndex, value) {
		ch.reRender()
	}
}

func (ch *addServiceDialog) removeParameter(paramIndex int) {
	if params, ok := ch.Data.Params.Remove(paramIndex); ok {
		ch.Data.Params = params
		ch.reRender()
	}
}

func (ch *addServiceDialog) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *addServiceDialog) Render() vecty.ComponentOrHTML {
	protocol, transport := ch.Data.ProtocolAndTransport(ch.Protocols)
	availableParams := ch.Data.Params.AvailableNames(transport)
	renderParams, paramsValid := ch.Data.Params.ToRender(transport)
	entryMessage, aliasMessage := ch.Data.ValidateEntryAndAlias(ch.Services.Services, ch.Protocols)
	dataValid := aliasMessage == "" && entryMessage == "" && paramsValid

	return components.NewMdcDialog(
		"sv-add-service-dialog", "Add Service", true, true, ch.closeDialog, ch.Data,
		[]components.MdcDialogButton{
			{Label: "Cancel", Action: components.MdcDialogActionClose},
			{Label: "Add Service", Action: components.MdcDialogActionOK, Default: true, Disabled: !dataValid},
		},
		ch.RenderProtocols(protocol),
		ch.RenderTransports(protocol, transport),
		components.NewMdcTextField(
			"sv-add-service-alias", "Alias", ch.Data.Alias, false, aliasMessage != "", ch.changeAlias,
			vecty.MarkupIf(aliasMessage != "", vecty.Attribute("title", aliasMessage)),
		).
			WithClasses("adjacent-margins").
			WithKey(fmt.Sprintf("text-alias-%d", ch.renderKey)),
		elem.Break(
			vecty.Markup(
				vecty.Key("br-1"),
			),
		),
		components.NewMdcTextField(
			"sv-add-service-entry", "Entry", ch.Data.Entry, false, entryMessage != "", ch.changeEntry,
			vecty.MarkupIf(entryMessage != "", vecty.Attribute("title", entryMessage)),
		).
			WithClasses("adjacent-margins").
			WithKey(fmt.Sprintf("text-entry-%d", ch.renderKey)),
		elem.Div(
			vecty.Markup(
				vecty.Key("param-title"),
				vecty.Class("mdc-typography--overline"),
			),
			vecty.Text("Parameters"),
		),
		components.NewMdcButton("sv-add-service-add-param", "Add Parameter", len(availableParams) <= 0, ch.addParameter).
			WithKey("add-param-btn"),
		ch.RenderParameters(renderParams, availableParams, transport),
	)
}

func (ch *addServiceDialog) RenderProtocols(protocol *api.ProtocolInfoEntry) vecty.ComponentOrHTML {
	var options vecty.List
	for _, p := range ch.Protocols.Protocols {
		options = append(options, &components.MdcSelectOption{
			Name:     fmt.Sprintf("%s (%d)", p.Name, p.ID),
			Selected: protocol.ID == p.ID,
		})
	}
	return components.NewMdcSelect(
		"sv-add-service-protocols", "Protocols", false, ch.changeProtocol, options,
	).WithKey(fmt.Sprintf("sv-add-service-protocols-%d", ch.renderKey))
}

func (ch *addServiceDialog) RenderTransports(protocol *api.ProtocolInfoEntry, transport *api.ProtocolTransportInfoEntry) vecty.ComponentOrHTML {
	var options vecty.List
	if protocol != nil {
		for _, t := range protocol.Transports {
			options = append(options, &components.MdcSelectOption{
				Name:     fmt.Sprintf("%s (%d)", t.Name, t.ID),
				Value:    strconv.Itoa(int(t.ID)),
				Selected: transport.ID == t.ID,
			})
		}
	}

	return components.NewMdcSelect(
		"sv-add-service-transports", "Transports", false, ch.changeTransport, options,
	).WithKey(fmt.Sprintf("sv-add-service-transports-%d", ch.renderKey))
}

func (ch *addServiceDialog) RenderParameters(renderParams []*core.RenderParameter, availableParams []string, transport *api.ProtocolTransportInfoEntry) vecty.KeyedList {
	baseKey := fmt.Sprintf("sv-add-service-param-%d", ch.renderKey)
	var result vecty.List
	for _, param := range renderParams {
		paramKey := fmt.Sprintf("%s-%s", baseKey, param.Name)
		paramIndex := param.Index

		result = append(result, elem.Break(
			vecty.Markup(
				vecty.Key("br-"+paramKey),
			),
		))

		var options vecty.List
		options = append(options, &components.MdcSelectOption{Name: param.Name, Selected: true})
		for _, name := range availableParams {
			options = append(options, &components.MdcSelectOption{Name: name})
		}

		paramNameKey := paramKey + "-name"
		result = append(result,
			components.NewMdcSelect(paramNameKey, "Parameter", false, func(value string, index int) {
				ch.changeParameter(paramIndex, value)
			}, options).
				WithKey(paramNameKey).
				WithClasses("sv-add-service-param-name"),
		)

		paramValueKey := paramKey + "-value"

		switch data := param.Data.(type) {
		case *core.ParameterEnumData:
			var enumOptions vecty.List
			for _, option := range data.Options {
				enumOptions = append(enumOptions, &components.MdcSelectOption{Name: option.Name, Selected: option.Selected})
			}
			result = append(result,
				components.NewMdcSelect(paramValueKey, "Parameter Value", false, func(value string, index int) {
					ch.changeParameterValue(paramIndex, value)
				}, enumOptions).
					WithKey(paramValueKey).
					WithClasses("sv-add-service-param-value"),
			)
		case *core.ParameterBoolData:
			result = append(result,
				elem.Div(
					vecty.Markup(
						vecty.Class("sv-add-service-param-value"),
						vecty.Style("display", "inline-block"),
						vecty.Key(paramValueKey),
					),
					components.NewMdcRadioButton(
						paramValueKey, "True", paramValueKey, "true", data.BoolValue, false, func() {
							ch.changeParameterValue(paramIndex, "true")
						},
					),
					components.NewMdcRadioButton(
						paramValueKey, "False", paramValueKey, "false", !data.BoolValue, false, func() {
							ch.changeParameterValue(paramIndex, "false")
						},
					),
				),
			)
		case *core.ParameterIntData:
			result = append(result,
				components.NewMdcTextField(
					paramValueKey, "Parameter Value", param.Value, false, data.Error != "",
					func(value string) {
						ch.changeParameterValue(paramIndex, value)
					},
					prop.Type(prop.TypeNumber),
					vecty.Attribute("min", data.Minimum),
					vecty.Attribute("max", data.Maximum),
					vecty.MarkupIf(data.Error != "", vecty.Attribute("title", data.Error)),
				).
					WithKey(paramValueKey).
					WithClasses("sv-add-service-param-value"),
			)
		default:
			result = append(result,
				components.NewMdcTextField(paramValueKey, "Parameter Value", param.Value, false, false, func(value string) {
					ch.changeParameterValue(paramIndex, value)
				}).
					WithKey(paramValueKey).
					WithClasses("sv-add-service-param-value"),
			)
		}

		paramDeleteKey := paramKey + "-delete"

		result = append(result, components.NewMdcIconButton(paramDeleteKey, "Delete Parameter", "delete_forever", "delete_forever", false,
			func() {
				ch.removeParameter(paramIndex)
			}).
			WithKey(paramDeleteKey).
			WithClasses("sv-add-service-param-delete"),
		)
	}
	return result.WithKey(baseKey + "-params")
}
