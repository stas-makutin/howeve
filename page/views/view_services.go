package views

import (
	"fmt"
	"strconv"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
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
							Entry:     "COM3",
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

func (ch *ViewServices) addServiceAction(ok bool, data *addServiceData) {
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
			&servicesTable{Services: ch.services},
		),
		core.If(ch.addServiceDialog, newAddServiceDialog(ch.protocols, ch.services, ch.addServiceAction)),
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

type addServiceData struct {
	api.ServiceKey
	Alias  string
	Params []core.Parameter
}

type addServiceDialog struct {
	vecty.Core
	Protocols *api.ProtocolInfoResult `vecty:"prop"`
	Services  *api.ListServicesResult `vecty:"prop"`
	CloseFn   func(ok bool, data *addServiceData)
	Data      *addServiceData

	renderKey int
}

func newAddServiceDialog(protocols *api.ProtocolInfoResult, services *api.ListServicesResult, closeFn func(ok bool, data *addServiceData)) *addServiceDialog {
	return &addServiceDialog{
		Protocols: protocols, Services: services, CloseFn: closeFn, Data: &addServiceData{},
	}
}

func (ch *addServiceDialog) closeDialog(action string, data interface{}) {
	ch.CloseFn(action == components.MdcDialogActionOK, data.(*addServiceData))
}

func (ch *addServiceDialog) changeProtocol(value string, index int) {
	if ch.Protocols == nil || index < 0 || index >= len(ch.Protocols.Protocols) {
		return
	}

	protocol := ch.Protocols.Protocols[index]
	if ch.Data.Protocol == protocol.ID {
		return
	}

	transportNotSupported := true
	for _, t := range protocol.Transports {
		if t.ID == ch.Data.Transport {
			transportNotSupported = false
			break
		}
	}
	if transportNotSupported && len(protocol.Transports) > 0 {
		ch.Data.Transport = protocol.Transports[0].ID
	}
	vecty.Rerender(ch)
}

func (ch *addServiceDialog) changeTransport(value string, index int) {
	if ch.Protocols == nil || index < 0 || len(ch.Protocols.Protocols) == 0 {
		return
	}

	var transports []*api.ProtocolTransportInfoEntry
	for _, p := range ch.Protocols.Protocols {
		if p.ID == ch.Data.Protocol {
			transports = p.Transports
		}
	}
	if index >= len(transports) {
		return
	}

	transportID := transports[index].ID
	if transportID != ch.Data.Transport {
		ch.Data.Transport = transportID
		ch.renderKey += 1
		vecty.Rerender(ch)
	}
}

func (ch *addServiceDialog) changeAlias(value string) {
	if value != ch.Data.Alias {
		ch.Data.Alias = value
		ch.renderKey += 1
		vecty.Rerender(ch)
	}
}

func (ch *addServiceDialog) changeEntry(value string) {
	if value != ch.Data.Entry {
		ch.Data.Entry = value
		ch.renderKey += 1
		vecty.Rerender(ch)
	}
}

func (ch *addServiceDialog) validateEntryAndAlias() (entryMessage, aliasMessage string) {
	for _, s := range ch.Services.Services {
		if s.Alias != "" && s.Alias == ch.Data.Alias {
			aliasMessage = fmt.Sprintf("Service with alias '%s' already exists", s.Alias)
		}
		if s.Protocol == ch.Data.Protocol && s.Transport == ch.Data.Transport && s.Entry == ch.Data.Entry {
			protocolName, transportName := core.ProtocolAndTransportName(s.Protocol, s.Transport, ch.Protocols)

			entryMessage = fmt.Sprintf(
				"Service with '%s' entry already exists for %s (%v) protocol and %s (%v) transport",
				s.Entry,
				protocolName, s.Protocol,
				transportName, s.Transport,
			)
		}
	}
	return
}

func (ch *addServiceDialog) addParameter() {
	_, transport := ch.protocolAndTransport()
	availableParams := ch.availableParameters(transport)
	if len(availableParams) > 0 {
		name := availableParams[0]
		value, _ := ch.paramDefaultValue(name, transport)
		ch.Data.Params = append(ch.Data.Params, core.Parameter{Name: name, Value: value})
		ch.renderKey += 1
		vecty.Rerender(ch)
	}
}

func (ch *addServiceDialog) changeParameter(paramIndex int, name string) {
	if paramIndex >= 0 && paramIndex < len(ch.Data.Params) {
		p := &(ch.Data.Params[paramIndex])
		if p.Name != name {
			_, transport := ch.protocolAndTransport()
			if value, ok := ch.paramDefaultValue(name, transport); ok {
				p.Name = name
				p.Value = value
				ch.renderKey += 1
				vecty.Rerender(ch)
			}
		}
	}
}

func (ch *addServiceDialog) changeParameterValue(paramIndex int, value string) {
	if paramIndex >= 0 && paramIndex < len(ch.Data.Params) {
		ch.Data.Params[paramIndex].Value = value
		ch.renderKey += 1
		vecty.Rerender(ch)
	}
}

func (ch *addServiceDialog) removeParameter(paramIndex int) {
	if paramIndex >= 0 && paramIndex < len(ch.Data.Params) {
		ch.Data.Params = append(ch.Data.Params[:paramIndex], ch.Data.Params[paramIndex+1:]...)
		ch.renderKey += 1
		vecty.Rerender(ch)
	}
}

func (ch *addServiceDialog) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *addServiceDialog) protocolAndTransport() (protocol *api.ProtocolInfoEntry, transport *api.ProtocolTransportInfoEntry) {
	if ch.Protocols != nil {
		for _, p := range ch.Protocols.Protocols {
			if p.ID == ch.Data.Protocol {
				protocol = p
				break
			}
		}
		if protocol == nil && len(ch.Protocols.Protocols) > 0 {
			protocol = ch.Protocols.Protocols[0]
			ch.Data.Protocol = protocol.ID
		}

		for _, t := range protocol.Transports {
			if t.ID == ch.Data.Transport {
				transport = t
				break
			}
		}
		if transport == nil && len(protocol.Transports) > 0 {
			transport = protocol.Transports[0]
			ch.Data.Transport = transport.ID
		}
	}
	return
}

func (ch *addServiceDialog) availableParameters(transport *api.ProtocolTransportInfoEntry) []string {
	var names []string
ParamLoop:
	for name := range transport.Params {
		for _, p := range ch.Data.Params {
			if name == p.Name {
				continue ParamLoop
			}
		}
		names = append(names, name)
	}
	return names
}

func (ch *addServiceDialog) paramDefaultValue(name string, transport *api.ProtocolTransportInfoEntry) (string, bool) {
	value := ""
	pi, ok := transport.Params[name]
	if ok {
		if pi.DefaultValue != "" {
			value = pi.DefaultValue
		} else if pi.Type == api.ParamTypeEnum {
			value = pi.EnumValues[0]
		} else if pi.Type == api.ParamTypeBool {
			value = "false"
		} else if pi.Type != api.ParamTypeString {
			value = "0"
		}
	}
	return value, ok
}

func (ch *addServiceDialog) Render() vecty.ComponentOrHTML {
	protocol, transport := ch.protocolAndTransport()
	availableParams := ch.availableParameters(transport)
	protocolsKey := fmt.Sprintf("sv-add-p-%p", ch.Protocols)
	transportsKey := fmt.Sprintf("%s-%d", protocolsKey, ch.Data.Protocol)
	entryTooltip, aliasTooltip := ch.validateEntryAndAlias()

	return components.NewMdcDialog(
		"sv-add-service-dialog", "Add Service", true, true, ch.closeDialog, ch.Data,
		[]components.MdcDialogButton{
			{Label: "Cancel", Action: components.MdcDialogActionClose},
			{Label: "Add Service", Action: components.MdcDialogActionOK, Default: true, Disabled: aliasTooltip != "" || entryTooltip != ""},
		},
		ch.RenderProtocols(protocolsKey, protocol),
		ch.RenderTransports(transportsKey, protocol, transport),
		components.NewMdcTextField(
			"sv-add-service-alias", "Alias", ch.Data.Alias, false, aliasTooltip != "", ch.changeAlias,
			vecty.MarkupIf(aliasTooltip != "", vecty.Attribute("title", aliasTooltip)),
		).
			WithClasses("adjacent-margins").
			WithKey(fmt.Sprintf("text-alias-%d", ch.renderKey)),
		elem.Break(
			vecty.Markup(
				vecty.Key("br-1"),
			),
		),
		components.NewMdcTextField(
			"sv-add-service-entry", "Entry", ch.Data.Entry, false, entryTooltip != "", ch.changeEntry,
			vecty.MarkupIf(entryTooltip != "", vecty.Attribute("title", entryTooltip)),
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
		ch.RenderParameters(transportsKey, transport, availableParams),
	)
}

func (ch *addServiceDialog) RenderProtocols(protocolsKey string, protocol *api.ProtocolInfoEntry) vecty.ComponentOrHTML {
	var options vecty.List
	if ch.Protocols != nil {
		for _, p := range ch.Protocols.Protocols {
			options = append(options, &components.MdcSelectOption{
				Name:     fmt.Sprintf("%s (%d)", p.Name, p.ID),
				Selected: protocol.ID == p.ID,
			})
		}
	}

	return components.NewMdcSelect(
		"sv-add-service-protocols", "Protocols", false, ch.changeProtocol, options,
	).WithKey(protocolsKey)
}

func (ch *addServiceDialog) RenderTransports(transportsKey string, protocol *api.ProtocolInfoEntry, transport *api.ProtocolTransportInfoEntry) vecty.ComponentOrHTML {
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
	).WithKey(transportsKey)
}

func (ch *addServiceDialog) RenderParameters(transportsKey string, transport *api.ProtocolTransportInfoEntry, availableParams []string) vecty.KeyedList {
	key := fmt.Sprintf("%s-p%d", transportsKey, ch.renderKey)
	var result vecty.List
	for i, param := range ch.Data.Params {
		pi, ok := transport.Params[param.Name]
		if !ok {
			continue
		}

		paramIndex := i
		paramKey := fmt.Sprintf("%s-%s", key, param.Name)

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
		switch pi.Type {
		case "enum":
			var enumOptions vecty.List
			for _, enumValue := range pi.EnumValues {
				enumOptions = append(enumOptions, &components.MdcSelectOption{Name: enumValue, Selected: enumValue == param.Value})
			}
			result = append(result,
				components.NewMdcSelect(paramValueKey, "Parameter Value", false, func(value string, index int) {
					ch.changeParameterValue(paramIndex, value)
				}, enumOptions).
					WithKey(paramValueKey).
					WithClasses("sv-add-service-param-value"),
			)
		case "bool":
			result = append(result,
				elem.Div(
					vecty.Markup(
						vecty.Class("sv-add-service-param-value"),
						vecty.Style("display", "inline-block"),
						vecty.Key(paramValueKey),
					),
					components.NewMdcRadioButton(
						paramValueKey, "True", paramValueKey, "true", param.Value == "true" || param.Value == "1", false, func() {
							ch.changeParameterValue(paramIndex, "true")
						},
					),
					components.NewMdcRadioButton(
						paramValueKey, "False", paramValueKey, "false", !(param.Value == "true" || param.Value == "1"), false, func() {
							ch.changeParameterValue(paramIndex, "false")
						},
					),
				),
			)
		case "string":
			result = append(result,
				components.NewMdcTextField(paramValueKey, "Parameter Value", param.Value, false, false, func(value string) {
					ch.changeParameterValue(paramIndex, value)
				}).
					WithKey(paramValueKey).
					WithClasses("sv-add-service-param-value"),
			)
		default:
			var min, max, msg string
			switch pi.Type {
			case api.ParamTypeInt8:
				l := int64(api.ParamTypeInt8Min)
				u := int64(api.ParamTypeInt8Max)
				v, err := strconv.ParseInt(param.Value, 10, 8)
				min = strconv.FormatInt(l, 10)
				max = strconv.FormatInt(u, 10)
				if err != nil {
					msg = err.Error()
				} else if v < l {
					msg = fmt.Sprintf("The value must be equal or greater than %d", l)
				} else if v > u {
					msg = fmt.Sprintf("The value must be equal or less than %d", u)
				}
			case api.ParamTypeInt16:
				l := int64(api.ParamTypeInt16Min)
				u := int64(api.ParamTypeInt16Max)
				v, err := strconv.ParseInt(param.Value, 10, 16)
				min = strconv.FormatInt(l, 10)
				max = strconv.FormatInt(u, 10)
				if err != nil {
					msg = err.Error()
				} else if v < l {
					msg = fmt.Sprintf("The value must be equal or greater than %d", l)
				} else if v > u {
					msg = fmt.Sprintf("The value must be equal or less than %d", u)
				}
			case api.ParamTypeInt32:
				l := int64(api.ParamTypeInt32Min)
				u := int64(api.ParamTypeInt32Max)
				v, err := strconv.ParseInt(param.Value, 10, 32)
				min = strconv.FormatInt(l, 10)
				max = strconv.FormatInt(u, 10)
				if err != nil {
					msg = err.Error()
				} else if v < l {
					msg = fmt.Sprintf("The value must be equal or greater than %d", l)
				} else if v > u {
					msg = fmt.Sprintf("The value must be equal or less than %d", u)
				}
			case api.ParamTypeInt64:
				l := int64(api.ParamTypeInt64Min)
				u := int64(api.ParamTypeInt64Max)
				v, err := strconv.ParseInt(param.Value, 10, 64)
				min = strconv.FormatInt(l, 10)
				max = strconv.FormatInt(u, 10)
				if err != nil {
					msg = err.Error()
				} else if v < l {
					msg = fmt.Sprintf("The value must be equal or greater than %d", l)
				} else if v > u {
					msg = fmt.Sprintf("The value must be equal or less than %d", u)
				}
			case api.ParamTypeUint16:
				l := uint64(api.ParamTypeUint16Min)
				u := uint64(api.ParamTypeUint16Max)
				v, err := strconv.ParseUint(param.Value, 10, 16)
				min = strconv.FormatUint(l, 10)
				max = strconv.FormatUint(u, 10)
				if err != nil {
					msg = err.Error()
				} else if v < l {
					msg = fmt.Sprintf("The value must be equal or greater than %d", l)
				} else if v > u {
					msg = fmt.Sprintf("The value must be equal or less than %d", u)
				}
			case api.ParamTypeUint32:
				l := uint64(api.ParamTypeUint32Min)
				u := uint64(api.ParamTypeUint32Max)
				v, err := strconv.ParseUint(param.Value, 10, 32)
				min = strconv.FormatUint(l, 10)
				max = strconv.FormatUint(u, 10)
				if err != nil {
					msg = err.Error()
				} else if v < l {
					msg = fmt.Sprintf("The value must be equal or greater than %d", l)
				} else if v > u {
					msg = fmt.Sprintf("The value must be equal or less than %d", u)
				}
			case api.ParamTypeUint64:
				l := uint64(api.ParamTypeUint64Min)
				u := uint64(api.ParamTypeUint64Max)
				v, err := strconv.ParseUint(param.Value, 10, 64)
				min = strconv.FormatUint(l, 10)
				max = strconv.FormatUint(u, 10)
				if err != nil {
					msg = err.Error()
				} else if v < l {
					msg = fmt.Sprintf("The value must be equal or greater than %d", l)
				} else if v > u {
					msg = fmt.Sprintf("The value must be equal or less than %d", u)
				}
			default:
				l := uint64(api.ParamTypeUint8Min)
				u := uint64(api.ParamTypeUint8Max)
				v, err := strconv.ParseUint(param.Value, 10, 8)
				min = strconv.FormatUint(l, 10)
				max = strconv.FormatUint(u, 10)
				if err != nil {
					msg = err.Error()
				} else if v < l {
					msg = fmt.Sprintf("The value must be equal or greater than %d", l)
				} else if v > u {
					msg = fmt.Sprintf("The value must be equal or less than %d", u)
				}
			}

			result = append(result,
				components.NewMdcTextField(
					paramValueKey, "Parameter Value", param.Value, false, msg != "",
					func(value string) {
						ch.changeParameterValue(paramIndex, value)
					},
					prop.Type(prop.TypeNumber),
					vecty.Attribute("min", min),
					vecty.Attribute("max", max),
					vecty.MarkupIf(msg != "", vecty.Attribute("title", msg)),
				).
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
	return result.WithKey(key + "-params")
}
