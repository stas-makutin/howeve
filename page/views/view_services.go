package views

import (
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
	services         *api.ListServicesResult
}

func NewViewServices() (r *ViewServices) {
	store := actions.GetServicesViewStore()
	r = &ViewServices{
		rendered:         false,
		addServiceDialog: true,
		loading:          store.Loading,
		useSockets:       store.UseSocket,
		errorMessage:     store.Error,
		services:         store.Services,
	}
	actions.Subscribe(r)
	return
}

func (ch *ViewServices) OnChange(event interface{}) {
	if store, ok := event.(*actions.ServicesViewStore); ok {
		ch.loading = store.Loading
		ch.useSockets = store.UseSocket
		ch.errorMessage = store.Error
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
			components.NewMdcButton("sv-add-service", "Add Service", false, ch.addService).AddClasses("adjacent-margins"),
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
		core.If(ch.addServiceDialog, &addServicesDialog{CloseFn: ch.addServiceAction}),
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
	Transports []string `vecty:"prop"`
	CloseFn    func(action string)
}

func (ch *addServicesDialog) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *addServicesDialog) Render() vecty.ComponentOrHTML {
	if len(ch.Transports) <= 0 {
		core.Console.Log("default set")
		ch.Transports = []string{"Transport 1", "Transport 2"}
	}
	var transports vecty.List
	for i, t := range ch.Transports {
		transports = append(transports, &components.MdcSelectOption{Name: t, Selected: i == 0})
	}

	return components.NewMdcDialog(
		"sv-add-service-dialog", "Add Service", true, true, ch.CloseFn,
		[]components.MdcDialogButton{
			{Label: "Cancel", Action: components.MdcDialogActionClose},
			{Label: "Add Service", Action: components.MdcDialogActionOK, Default: true, Disabled: false},
		},
		components.NewMdcSelect(
			"sv-add-service-protocols", "Protocols", false, func(value string, index int) {
				if index == 0 {
					core.Console.Log("first set")
					ch.Transports = []string{"Transport 1", "Transport 2"}
				} else {
					core.Console.Log("second set")
					ch.Transports = []string{"Transport 3", "Transport 4", "Transport 5"}
				}
				vecty.Rerender(ch)
			},
			&components.MdcSelectOption{Name: "Protocol 1", Selected: true},
			&components.MdcSelectOption{Name: "Protocol 2"},
		).AddClasses("adjacent-margins"),
		components.NewMdcSelect(
			"sv-add-service-transports", "Transports", false, func(value string, index int) {},
			transports,
		).AddClasses("adjacent-margins"),
		components.NewMdcTextField("sv-add-service-alias", "Alias", "", false).AddClasses("adjacent-margins"),
		elem.Break(
			vecty.Markup(
				vecty.Key("br"),
			),
		),
		components.NewMdcTextField("sv-add-service-entry", "Entry", "", false).AddClasses("adjacent-margins"),
	)
}
