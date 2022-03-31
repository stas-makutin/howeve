package components

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
	"github.com/hexops/vecty/prop"
	"github.com/stas-makutin/howeve/page/core"
)

type MdcButton struct {
	vecty.Core
	core.ClassAdder
	ID       string `vecty:"prop"`
	Text     string `vecty:"prop"`
	Disabled bool   `vecty:"prop"`
	clickFn  func()
}

func NewMdcButton(id string, text string, disabled bool, clickFn func()) (r *MdcButton) {
	r = &MdcButton{ID: id, Text: text, Disabled: disabled, clickFn: clickFn}
	return
}

func (ch *MdcButton) onClick(event *vecty.Event) {
	ch.clickFn()
}

func (ch *MdcButton) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *MdcButton) AddClasses(classes ...string) vecty.Component {
	ch.ClassAdder.AddClasses(classes...)
	return ch
}

func (ch *MdcButton) Render() vecty.ComponentOrHTML {
	return elem.Button(
		vecty.Markup(
			prop.ID(ch.ID),
			vecty.Class("mdc-button", "mdc-button--outlined"),
			ch.ApplyClasses(),
			prop.Disabled(ch.Disabled),
			event.Click(ch.onClick),
		),
		elem.Span(
			vecty.Markup(
				vecty.Class("mdc-button__ripple"),
			),
		),
		elem.Span(
			vecty.Markup(
				vecty.Class("mdc-button__focus-ring"),
			),
		),
		elem.Span(
			vecty.Markup(
				vecty.Class("mdc-button__label"),
			),
			vecty.Text(ch.Text),
		),
	)
}
