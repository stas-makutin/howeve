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

func init() {
	core.AppendStyles(`
#sv-add-service-dialog---content {
	min-height: 15em;
}
#sv-add-service-protocols {
	min-width: 18em;
	margin-top: 0.5em;
	margin-right: 0.5em;
}
#sv-add-service-transports {
	min-width: 18em;
	margin-top: 0.5em;
	margin-right: 0.5em;
}
#sv-add-service-alias {
	min-width: 18em;
	margin-top: 0.5em;
}
#sv-add-service-entry {
	min-width: 16em;
	width: 100%;
	margin-top: 0.5em;
}
#sv-add-service-add-param {
	margin-top: 0.5em;
	margin-bottom: 0.5em;
}
.sv-add-service-param-name {
	min-width: 15%;
	width: 25%;
	margin-top: 0.5em;
	margin-right: 0.5em;
}
.sv-add-service-param-value {
	min-width: 15%;
	width: 55%;
	margin-top: 0.5em;
}
.sv-add-service-param-delete {
	position: relative;
	top: 0.2em;
}
.sv-service-table-entry-cell {
	min-width: 10em;
	word-break: break-word;
	white-space: break-spaces;
}
.sv-service-table-action-cell {
	white-space: break-spaces;
}
.sv-service-table-action {
	white-space: nowrap;
}
.sv-service-table-status-healthy {
	color: green;
}
.sv-service-table-status-unhealthy {
	color: red;
}
.sv-service-table-remove-cell {
	text-align: center;
}
.sv-change-alias-edit-field {
	margin-top: 0.8em;
	width: 100%;
}
.sv-dialogs-center {
	margin: auto;
}
`,
	)
}

const (
	ServicesDialog_None = iota
	ServicesDialog_AddService
	ServicesDialog_RemoveService
	ServicesDialog_ChangeAlias
	ServicesDialog_ViewParameters
	ServicesDialog_ViewStatus
)

type ViewServices struct {
	vecty.Core
	rendered       bool
	renderDialog   int
	loading        bool
	useSockets     bool
	errorMessage   []vecty.MarkupOrChild
	protocols      *core.ProtocolsWrapper
	services       *api.ListServicesResult
	currentService api.ListServicesEntry `vecty:"prop"`
}

func NewViewServices() (r *ViewServices) {
	store := actions.GetServicesViewStore()
	r = &ViewServices{
		rendered:     false,
		renderDialog: ServicesDialog_None,
		loading:      store.Loading > 0,
		useSockets:   store.UseSocket,
		errorMessage: core.FormatMultilineText(store.DisplayError),
		protocols:    core.NewProtocolsWrapper(store.Protocols.Value),
		services:     store.Services.Value,
	}
	actions.Subscribe(r)
	return
}

func (ch *ViewServices) OnChange(event interface{}) {
	if store, ok := event.(*actions.ServicesViewStore); ok {
		ch.loading = store.Loading > 0
		ch.useSockets = store.UseSocket
		ch.errorMessage = core.FormatMultilineText(store.DisplayError)
		ch.protocols = core.NewProtocolsWrapper(store.Protocols.Value)
		ch.services = store.Services.Value
		/*
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
								"param01":                               "value01",
								"param02":                               "value02",
								"param03":                               "value03",
								"param04":                               "value04",
								"param05":                               "value05",
								"param06":                               "value06",
								"param07":                               "value07",
								"param08":                               "value08",
								"param09":                               "value09",
								"param10":                               "value10",
								"param11":                               "value11",
								"param12":                               "value12",
								"param13":                               "value13",
								"param14":                               "value14",
								"param15":                               "value15",
								"param16":                               "value16",
								"param17":                               "value17",
								"param18":                               "value18",
								"param19":                               "value19",
								"param20":                               "value20",
								"param21":                               "value21",
								"param22":                               "value22",
								"param23":                               "value23",
								"param24":                               "value24",
								"param25":                               "value25",
								"param26":                               "value26",
								"param27":                               "value27",
								"param28":                               "value28 - very-very-very-very-very-very-very-very-very-very-very-very-very-very-very-very-very-very-very-very-very-very-very-very-long",
								"param29":                               "value29",
								"param30":                               "value30",
								"param31":                               "value31",
								"param32":                               "value32",
								"param33-very-very-very-very-very-long": "value33",
								"param34":                               "value34",
								"param35":                               "value35",
								"param36":                               "value36",
								"param37":                               "value37",
								"param38":                               "value38",
								"param39":                               "value39",
								"param40":                               "value40",
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
		*/
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

func (ch *ViewServices) addService(ok bool, service *core.ServiceEntryData) {
	ch.renderDialog = ServicesDialog_None
	if ok {
		core.Console.Log(fmt.Sprintf("add %d: %d: %s", service.Protocol, service.Transport, service.Entry))
	}
	vecty.Rerender(ch)
}

func (ch *ViewServices) changeAlias(ok bool, newAlias string, service *api.ServiceEntry) {
	ch.renderDialog = ServicesDialog_None
	if ok {
		core.Console.Log(fmt.Sprintf("alias %d: %d: %s -> %s", service.Protocol, service.Transport, service.Entry, newAlias))
	}
	vecty.Rerender(ch)
}

func (ch *ViewServices) removeService(ok bool, service *api.ServiceEntry) {
	ch.renderDialog = ServicesDialog_None
	if ok {
		core.Console.Log(fmt.Sprintf("remove %d: %d: %s", service.Protocol, service.Transport, service.Entry))
	}
	vecty.Rerender(ch)
}

func (ch *ViewServices) toggleDialog(dialog int) {
	ch.renderDialog = dialog
	vecty.Rerender(ch)
}

func (ch *ViewServices) openServiceDialog(dialog int, service api.ServiceKey) {
	ch.currentService.ServiceEntry = nil
	if ch.services != nil && len(ch.services.Services) > 0 {
		for _, s := range ch.services.Services {
			if s.Protocol == service.Protocol && s.Transport == service.Transport && s.Entry == service.Entry {
				ch.currentService = s
				ch.toggleDialog(dialog)
			}
		}
	}
}

func (ch *ViewServices) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *ViewServices) Render() vecty.ComponentOrHTML {
	ch.rendered = true
	return components.NewMdcGrid(
		components.NewMdcGridSingleCellRow(
			components.NewMdcButton("sv-add-service", "Add Service", ch.protocols == nil,
				func() { ch.toggleDialog(ServicesDialog_AddService) },
			).WithClasses("adjacent-margins"),
			components.NewMdcButton("sv-refresh", "Refresh", false, ch.refresh),
			components.NewMdcCheckbox("sv-socket-check", "Use WebSocket", ch.useSockets, false, ch.changeUseSocket),
		),
		core.If(len(ch.errorMessage) > 0, components.NewMdcGridSingleCellRow(
			components.NewMdcBanner("sv-error-banner", "Refresh", true, ch.refresh, ch.errorMessage...),
		)),
		&components.SectionTitle{Text: "Services"},
		components.NewMdcGridSingleCellRow(
			&servicesTable{Protocols: ch.protocols, Services: ch.services, OpenDialog: ch.openServiceDialog},
		),
		core.If(ch.renderDialog == ServicesDialog_AddService, newAddServiceDialog(ch.protocols, ch.services, ch.addService)),
		core.If(ch.currentService.ServiceEntry != nil && ch.renderDialog == ServicesDialog_ChangeAlias,
			newChangeServiceAliasDialog(ch.protocols, ch.services, ch.currentService.ServiceEntry, ch.changeAlias),
		),
		core.If(ch.currentService.ServiceEntry != nil && ch.renderDialog == ServicesDialog_ViewParameters,
			&viewServiceParametersDialog{Protocols: ch.protocols, Service: ch.currentService.ServiceEntry, CloseFn: func() { ch.toggleDialog(ServicesDialog_None) }},
		),
		core.If(ch.currentService.ServiceEntry != nil && ch.renderDialog == ServicesDialog_ViewStatus,
			&viewServiceStatusDialog{Protocols: ch.protocols, Service: &ch.currentService, CloseFn: func() { ch.toggleDialog(ServicesDialog_None) }},
		),
		core.If(ch.currentService.ServiceEntry != nil && ch.renderDialog == ServicesDialog_RemoveService,
			&removeServiceDialog{Protocols: ch.protocols, Service: ch.currentService.ServiceEntry, CloseFn: ch.removeService},
		),
		core.If(ch.loading, &components.ViewLoading{}),
	)
}

type servicesTable struct {
	vecty.Core
	Protocols  *core.ProtocolsWrapper  `vecty:"prop"`
	Services   *api.ListServicesResult `vecty:"prop"`
	OpenDialog func(dialog int, service api.ServiceKey)
}

func (ch *servicesTable) Copy() vecty.Component {
	cpy := *ch
	return &cpy
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

	protocolName, transportName := ch.Protocols.ProtocolAndTransportFullNames(service.Protocol, service.Transport)

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
				vecty.Attribute("title", fmt.Sprintf("%d: %s", service.Error.Code, service.Error.Message)),
				event.Click(func(e *vecty.Event) { ch.OpenDialog(ServicesDialog_ViewStatus, serviceKey) }).PreventDefault(),
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
					event.Click(func(e *vecty.Event) { ch.OpenDialog(ServicesDialog_ChangeAlias, serviceKey) }).PreventDefault(),
				),
				vecty.Text("Change Alias"),
			),
			vecty.Text(", "),
			elem.Anchor(
				vecty.Markup(
					vecty.Class("sv-service-table-action"),
					prop.Href("#"),
					event.Click(func(e *vecty.Event) { ch.OpenDialog(ServicesDialog_ViewParameters, serviceKey) }).PreventDefault(),
				),
				vecty.Text("View Parameters"),
			),
		}, "sv-service-table-action-cell"),
		ch.tableColumn(status),
		ch.tableColumn(components.NewMdcIconButton("", "Remove Service", "delete_forever", "delete_forever", false,
			func() { ch.OpenDialog(ServicesDialog_RemoveService, serviceKey) },
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

type changeServiceAliasDialog struct {
	vecty.Core
	Protocols *core.ProtocolsWrapper  `vecty:"prop"`
	Services  *api.ListServicesResult `vecty:"prop"`
	Service   *api.ServiceEntry       `vecty:"prop"`
	NewAlias  string                  `vecty:"prop"`
	CloseFn   func(ok bool, newAlias string, data *api.ServiceEntry)
}

func newChangeServiceAliasDialog(
	protocols *core.ProtocolsWrapper, services *api.ListServicesResult, service *api.ServiceEntry,
	closeFn func(ok bool, newAlias string, data *api.ServiceEntry)) (dialog *changeServiceAliasDialog) {
	dialog = &changeServiceAliasDialog{Protocols: protocols, Services: services, Service: service, CloseFn: closeFn}
	if service != nil {
		dialog.NewAlias = service.Alias
	}
	return
}

func (ch *changeServiceAliasDialog) closeDialog(action string, data interface{}) {
	ch.CloseFn(action == components.MdcDialogActionOK, ch.NewAlias, ch.Service)
}

func (ch *changeServiceAliasDialog) changeAlias(value string) {
	if value != ch.NewAlias {
		ch.NewAlias = value
		vecty.Rerender(ch)
	}
}

func (ch *changeServiceAliasDialog) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *changeServiceAliasDialog) Render() vecty.ComponentOrHTML {
	protocolName, transportName := ch.Protocols.ProtocolAndTransportFullNames(ch.Service.Protocol, ch.Service.Transport)
	aliasMessage := ""
	for _, s := range ch.Services.Services {
		if s.Alias != "" && s.Alias != ch.Service.Alias && s.Alias == ch.NewAlias {
			aliasMessage = fmt.Sprintf("Service with alias '%s' already exists", s.Alias)
			break
		}
	}
	disableChange := ch.Service.Alias == ch.NewAlias || aliasMessage != ""

	return components.NewMdcDialog(
		"sv-remove-service-dialog", "Change Service Alias", false, false, ch.closeDialog, nil,
		[]components.MdcDialogButton{
			{Label: "Cancel", Action: components.MdcDialogActionClose, Default: true},
			{Label: "Change Alias", Action: components.MdcDialogActionOK, Disabled: disableChange},
		},
		components.NewPillBox(
			&components.Capsule{Label: "Protocol", Text: protocolName},
			&components.Capsule{Label: "Transport", Text: transportName},
			vecty.If(ch.Service.Alias != "", &components.Capsule{Label: "Alias", Text: ch.Service.Alias}),
			&components.Capsule{Label: "Entry", Text: ch.Service.Entry},
		),
		components.NewMdcTextField(
			"sv-change-service-alias", "New Alias", ch.NewAlias, false, aliasMessage != "", ch.changeAlias,
			vecty.MarkupIf(aliasMessage != "", vecty.Attribute("title", aliasMessage)),
		).
			WithClasses("sv-change-alias-edit-field").
			WithKey("sv-change-alias-"+ch.NewAlias),
	)
}

type viewServiceParametersDialog struct {
	vecty.Core
	Protocols *core.ProtocolsWrapper `vecty:"prop"`
	Service   *api.ServiceEntry      `vecty:"prop"`
	CloseFn   func()
}

func (ch *viewServiceParametersDialog) closeDialog(action string, data interface{}) {
	ch.CloseFn()
}

func (ch *viewServiceParametersDialog) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *viewServiceParametersDialog) Render() vecty.ComponentOrHTML {
	protocolName, transportName := ch.Protocols.ProtocolAndTransportFullNames(ch.Service.Protocol, ch.Service.Transport)
	names := core.ArrangeParams(ch.Service.Params)

	return components.NewMdcDialog(
		"sv-remove-service-dialog", "Service Parameters", false, false, ch.closeDialog, nil,
		[]components.MdcDialogButton{
			{Label: "Close", Action: components.MdcDialogActionClose, Default: true},
		},
		components.NewPillBox(
			&components.Capsule{Label: "Protocol", Text: protocolName},
			&components.Capsule{Label: "Transport", Text: transportName},
			vecty.If(ch.Service.Alias != "", &components.Capsule{Label: "Alias", Text: ch.Service.Alias}),
			&components.Capsule{Label: "Entry", Text: ch.Service.Entry},
		),
		components.NewKeyValueTable(func(builder components.KeyValueTableBuilder) {
			builder.AddDelimiterRow()
			builder.AddMarkupKeyValueRow(
				elem.Bold(vecty.Text("Name")), elem.Bold(vecty.Text("Value")),
			)
			for _, name := range names {
				builder.AddKeyValueRow(name, ch.Service.Params[name])
			}
			builder.AddDelimiterRow()
		}).WithClasses("sv-dialogs-center"),
	)
}

type viewServiceStatusDialog struct {
	vecty.Core
	Protocols *core.ProtocolsWrapper `vecty:"prop"`
	Service   *api.ListServicesEntry `vecty:"prop"`
	CloseFn   func()
}

func (ch *viewServiceStatusDialog) closeDialog(action string, data interface{}) {
	ch.CloseFn()
}

func (ch *viewServiceStatusDialog) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *viewServiceStatusDialog) Render() vecty.ComponentOrHTML {
	protocolName, transportName := ch.Protocols.ProtocolAndTransportFullNames(ch.Service.Protocol, ch.Service.Transport)

	return components.NewMdcDialog(
		"sv-remove-service-dialog", "Service Error", false, false, ch.closeDialog, nil,
		[]components.MdcDialogButton{
			{Label: "Close", Action: components.MdcDialogActionClose, Default: true},
		},
		components.NewPillBox(
			&components.Capsule{Label: "Protocol", Text: protocolName},
			&components.Capsule{Label: "Transport", Text: transportName},
			vecty.If(ch.Service.Alias != "", &components.Capsule{Label: "Alias", Text: ch.Service.Alias}),
			&components.Capsule{Label: "Entry", Text: ch.Service.Entry},
		),
		components.NewKeyValueTable(func(builder components.KeyValueTableBuilder) {
			builder.AddDelimiterRow()
			builder.AddMarkupKeyValueRow(
				elem.Bold(vecty.Text("Code")),
				elem.Span(
					vecty.Markup(
						vecty.Style("white-space", "break-spaces"),
					),
					vecty.Text(
						strconv.FormatInt(int64(ch.Service.Error.Code), 10),
					),
				),
			)
			builder.AddMarkupKeyValueRow(
				elem.Bold(vecty.Text("Message")),
				elem.Span(
					vecty.Markup(
						vecty.Style("white-space", "break-spaces"),
					),
					vecty.Text(
						ch.Service.Error.Message,
					),
				),
			)
			builder.AddDelimiterRow()
		}).WithClasses("sv-dialogs-center"),
	)
}

type removeServiceDialog struct {
	vecty.Core
	Protocols *core.ProtocolsWrapper `vecty:"prop"`
	Service   *api.ServiceEntry      `vecty:"prop"`
	CloseFn   func(ok bool, service *api.ServiceEntry)
}

func (ch *removeServiceDialog) closeDialog(action string, data interface{}) {
	ch.CloseFn(action == components.MdcDialogActionOK, ch.Service)
}

func (ch *removeServiceDialog) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *removeServiceDialog) Render() vecty.ComponentOrHTML {
	protocolName, transportName := ch.Protocols.ProtocolAndTransportFullNames(ch.Service.Protocol, ch.Service.Transport)

	return components.NewMdcDialog(
		"sv-remove-service-dialog", "Remove the Service?", false, false, ch.closeDialog, nil,
		[]components.MdcDialogButton{
			{Label: "Cancel", Action: components.MdcDialogActionClose, Default: true},
			{Label: "Remove", Action: components.MdcDialogActionOK},
		},
		components.NewPillBox(
			&components.Capsule{Label: "Protocol", Text: protocolName},
			&components.Capsule{Label: "Transport", Text: transportName},
			vecty.If(ch.Service.Alias != "", &components.Capsule{Label: "Alias", Text: ch.Service.Alias}),
			&components.Capsule{Label: "Entry", Text: ch.Service.Entry},
		),
	)
}
