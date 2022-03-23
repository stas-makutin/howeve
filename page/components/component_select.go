package components

import (
	"syscall/js"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/prop"
	"github.com/stas-makutin/howeve/page/core"
)

type MdcSelectOption struct {
	vecty.Core
	Name     string
	Value    string
	Selected bool
	Disabled bool
}

func (ch *MdcSelectOption) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *MdcSelectOption) Render() vecty.ComponentOrHTML {
	return elem.ListItem(
		vecty.Markup(
			vecty.Class("mdc-deprecated-list-item"),
			vecty.MarkupIf(
				ch.Selected,
				vecty.Class("mdc-deprecated-list-item--selected"),
			),
			vecty.MarkupIf(
				ch.Disabled,
				vecty.Class("mdc-deprecated-list-item--disabled"),
			),
			vecty.Attribute("role", "option"),
			vecty.MarkupIf(
				ch.Value != "",
				vecty.Attribute("data-value", ch.Value),
			),
		),
		elem.Span(
			vecty.Markup(
				vecty.Class("mdc-deprecated-list-item__ripple"),
			),
		),
		elem.Span(
			vecty.Markup(
				vecty.Class("mdc-deprecated-list-item__text"),
			),
			vecty.Text(ch.Name),
		),
	)
}

type MdcSelect struct {
	vecty.Core
	core.ClassAdder
	ID       string
	Label    string
	Disabled bool
	Options  vecty.List
	jsObject js.Value
}

func NewMdcSelect(id, label string, disabled bool, options ...vecty.ComponentOrHTML) (r *MdcSelect) {
	r = &MdcSelect{ID: id, Label: label, Disabled: disabled, Options: options}
	return
}

func (ch *MdcSelect) Mount() {
	ch.jsObject = js.Global().Get("mdc").Get("select").Get("MDCSelect").Call(
		"attachTo", js.Global().Get("document").Call("getElementById", ch.ID),
	)
}

func (ch *MdcSelect) Unmount() {
	ch.jsObject.Call("destroy")
}

func (ch *MdcSelect) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *MdcSelect) AddClasses(classes ...string) vecty.Component {
	ch.ClassAdder.AddClasses(classes...)
	return ch
}

func (ch *MdcSelect) Render() vecty.ComponentOrHTML {
	hasLabel := ch.Label != ""
	labelID := ch.ID + "---label"
	return elem.Div(
		vecty.Markup(
			vecty.Class("mdc-select", "mdc-select--outlined"),
			ch.ApplyClasses(),
			vecty.MarkupIf(ch.Disabled, vecty.Class("mdc-select--disabled")),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("mdc-select__anchor"),
				vecty.MarkupIf(ch.Disabled, vecty.Attribute("aria-disabled", "true")),
				vecty.MarkupIf(hasLabel, vecty.Attribute("aria-labelledby", labelID)),
			),
			elem.Span(
				vecty.Markup(
					vecty.Class("mdc-notched-outline"),
				),
				elem.Span(
					vecty.Markup(
						vecty.Class("mdc-notched-outline__leading"),
					),
				),
				vecty.If(hasLabel,
					elem.Span(
						vecty.Markup(
							vecty.Class("mdc-notched-outline__notch"),
						),
						elem.Span(
							vecty.Markup(
								prop.ID(labelID),
								vecty.Class("mdc-floating-label"),
							),
						),
					),
				),
				elem.Span(
					vecty.Markup(
						vecty.Class("mdc-notched-outline__trailing"),
					),
				),
			),
			elem.Span(
				vecty.Markup(
					vecty.Class("mdc-select__selected-text-container"),
				),
				elem.Span(
					vecty.Markup(
						vecty.Class("mdc-select__selected-text"),
					),
					vecty.Text(""),
				),
			),
			elem.Span(
				vecty.Markup(
					vecty.Class("mdc-select__dropdown-icon"),
				),
				vecty.Tag(
					"svg",
					vecty.Markup(
						vecty.Namespace("http://www.w3.org/2000/svg"),
						vecty.Class("mdc-select__dropdown-icon-graphic"),
						vecty.Attribute("viewBox", "7 10 10 5"),
						vecty.Attribute("focusable", "false"),
					),
					vecty.Tag(
						"polygon",
						vecty.Markup(
							vecty.Namespace("http://www.w3.org/2000/svg"),
							vecty.Class("mdc-select__dropdown-icon-inactive"),
							vecty.Attribute("stroke", "none"),
							vecty.Attribute("fill-rule", "evenodd"),
							vecty.Attribute("points", "7 10 12 15 17 10"),
						),
					),
					vecty.Tag(
						"polygon",
						vecty.Markup(
							vecty.Namespace("http://www.w3.org/2000/svg"),
							vecty.Class("mdc-select__dropdown-icon-active"),
							vecty.Attribute("stroke", "none"),
							vecty.Attribute("fill-rule", "evenodd"),
							vecty.Attribute("points", "7 15 12 10 17 15"),
						),
					),
				),
			),
			elem.Div(
				vecty.Markup(
					vecty.Class("mdc-select__menu", "mdc-menu", "mdc-menu-surface", "mdc-menu-surface--fullwidth"),
				),
				elem.UnorderedList(
					vecty.Markup(
						vecty.Class("mdc-deprecated-list"),
						vecty.Attribute("role", "listbox"),
						vecty.MarkupIf(hasLabel, vecty.Attribute("aria-label", ch.Label)),
					),
					ch.Options,
				),
			),
		),
	)
}
