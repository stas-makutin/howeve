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
	errorMessage string
	config       string
}

func NewViewConfig() (r *ViewConfig) {
	store := actions.GetConfigViewStore()
	r = &ViewConfig{
		rendered:     false,
		loading:      store.Loading,
		useSockets:   store.UseSocket,
		errorMessage: store.Error,
		config:       store.Config,
	}
	actions.Subscribe(r)
	return
}

func (ch *ViewConfig) OnChange(event interface{}) {
	if store, ok := event.(*actions.ConfigViewStore); ok {
		ch.loading = store.Loading
		ch.useSockets = store.UseSocket
		ch.errorMessage = store.Error
		ch.config = store.Config
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
				components.NewMdcButton("pt-refresh", "Refresh", false, ch.refresh),
				components.NewMdcCheckbox("pt-socket-check", "Use WebSocket", ch.useSockets, false, ch.changeUseSocket),
			),
		),
		vecty.If(ch.errorMessage != "", elem.Div(
			vecty.Markup(
				vecty.Class("mdc-layout-grid__inner"),
			),
			elem.Div(
				vecty.Markup(
					vecty.Class("mdc-layout-grid__cell"),
				),
				components.NewMdcBanner("pt-error-banner", ch.errorMessage, "Retry", ch.refresh),
			),
		)),
		elem.Div(
			vecty.Markup(
				vecty.Class("mdc-layout-grid__inner"),
			),
			elem.Div(
				vecty.Markup(
					vecty.Class("mdc-layout-grid__cell", "mdc-layout-grid__cell--span-12"),
				),
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
		),
		vecty.If(ch.loading, &components.ViewLoading{}),
	)
}
