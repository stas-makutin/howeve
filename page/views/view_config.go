package views

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/stas-makutin/howeve/page/actions"
	"github.com/stas-makutin/howeve/page/components"
	"github.com/stas-makutin/howeve/page/core"
)

type ViewConfig struct {
	vecty.Core
	rendered     bool
	loading      bool
	useSockets   bool
	errorMessage []vecty.MarkupOrChild
	config       string
}

func NewViewConfig() (r *ViewConfig) {
	store := actions.GetConfigViewStore()
	r = &ViewConfig{
		rendered:     false,
		loading:      store.Loading,
		useSockets:   store.UseSocket,
		errorMessage: core.FormatMultilineText(store.DisplayError),
		config:       store.ConfigValue(),
	}
	actions.Subscribe(r)
	return
}

func (ch *ViewConfig) OnChange(event interface{}) {
	if store, ok := event.(*actions.ConfigViewStore); ok {
		ch.loading = store.Loading
		ch.useSockets = store.UseSocket
		ch.errorMessage = core.FormatMultilineText(store.DisplayError)
		ch.config = store.ConfigValue()
		if ch.rendered {
			vecty.Rerender(ch)
		}
	}
}

func (ch *ViewConfig) Mount() {
	core.Dispatch(&actions.ConfigLoad{Force: false, UseSocket: ch.useSockets})
}

func (ch *ViewConfig) changeUseSocket(checked, disabled bool) {
	core.Dispatch(actions.ConfigUseSocket(checked))
}

func (ch *ViewConfig) refresh() {
	core.Dispatch(&actions.ConfigLoad{Force: true, UseSocket: ch.useSockets})
}

func (ch *ViewConfig) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *ViewConfig) Render() vecty.ComponentOrHTML {
	ch.rendered = true
	configText := ch.config
	if configText == "" {
		configText = " "
	}
	return components.NewMdcGrid(
		components.NewMdcGridSingleCellRow(
			components.NewMdcButton("cf-refresh", "Refresh", false, ch.refresh),
			components.NewMdcCheckbox("cf-socket-check", "Use WebSocket", ch.useSockets, false, ch.changeUseSocket),
		),
		core.If(len(ch.errorMessage) > 0, components.NewMdcGridSingleCellRow(
			components.NewMdcBanner("cf-error-banner", "Retry", true, ch.refresh, ch.errorMessage...),
		)),
		&components.SectionTitle{Text: "Configuration"},
		components.NewMdcGridSingleCellRow(
			elem.Preformatted(
				vecty.Markup(
					vecty.Class("mdc-elevation--z1"),
					vecty.Style("margin", "0"),
					vecty.Style("padding", "3px"),
					vecty.Style("max-height", "calc(100vh - 150px)"),
					vecty.Style("overflow", "auto"),
				),
				vecty.Text(configText),
			),
		),
		core.If(ch.loading, &components.ViewLoading{}),
	)
}
