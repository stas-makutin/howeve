package components

import (
	"syscall/js"

	"github.com/hexops/vecty"
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
	return nil
}

type MdcSelect struct {
	vecty.Core
	ID       string
	Label    string
	Disabled bool
	Options  []MdcSelectOption
	jsObject js.Value
}

func NewMdcSelect(id, label string, disabled bool, options ...MdcSelectOption) (r *MdcSelect) {
	r = &MdcSelect{ID: id, Label: label, Disabled: disabled, Options: options}
	return
}

func (ch *MdcSelect) Mount() {
	ch.jsObject = js.Global().Get("mdc").Get("textField").Get("MdcSelect").Call(
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

func (ch *MdcSelect) Render() vecty.ComponentOrHTML {
	return nil
}
