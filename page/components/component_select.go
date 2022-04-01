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
	selectedFlag := "false"
	if ch.Selected {
		selectedFlag = "true"
	}
	value := ch.Value
	if value == "" {
		value = ch.Name
	}
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
			vecty.Attribute("aria-selected", selectedFlag),
			vecty.Attribute("data-value", value),
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
	core.Classes
	core.Keyable
	ID         string
	Label      string
	Disabled   bool
	Options    vecty.List
	changeFn   func(value string, index int)
	jsObject   js.Value
	jsChangeFn js.Func
}

func NewMdcSelect(id, label string, disabled bool, changeFn func(value string, index int), options ...vecty.ComponentOrHTML) (r *MdcSelect) {
	r = &MdcSelect{ID: id, Label: label, Disabled: disabled, changeFn: changeFn, Options: options}
	return
}

func (ch *MdcSelect) Mount() {
	ch.Unmount()
	ch.jsObject = js.Global().Get("mdc").Get("select").Get("MDCSelect").Call(
		"attachTo", js.Global().Get("document").Call("getElementById", ch.ID),
	)
	ch.jsChangeFn = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) > 0 {
			value := args[0].Get("detail").Get("value").String()
			index := args[0].Get("detail").Get("index").Int()
			ch.changeFn(value, index)
		}
		return nil
	})
	ch.jsObject.Call("listen", "MDCSelect:change", ch.jsChangeFn)
}

func (ch *MdcSelect) Unmount() {
	core.SafeJSDestroy(&ch.jsObject, func(v *js.Value) { v.Call("destroy") })
	core.ReleaseJSFunc(&ch.jsChangeFn)
}

func (ch *MdcSelect) WithKey(key interface{}) *MdcSelect {
	ch.Keyable.WithKey(key)
	return ch
}

func (ch *MdcSelect) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *MdcSelect) WithClasses(classes ...string) *MdcSelect {
	ch.Classes.WithClasses(classes...)
	return ch
}

func (ch *MdcSelect) Render() vecty.ComponentOrHTML {
	hasLabel := ch.Label != ""
	labelID := ch.ID + "---label"
	return elem.Div(
		vecty.Markup(
			prop.ID(ch.ID),
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
							vecty.Text(ch.Label),
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
	)
}
