package components

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/stas-makutin/howeve/page/core"
)

type MdcFormField struct {
	vecty.Core
	core.Classes
	core.Keyable
	Content vecty.List `vecty:"prop"`
}

func NewMdcFormField(content ...vecty.ComponentOrHTML) (r *MdcFormField) {
	r = &MdcFormField{Content: content}
	return
}

func (ch *MdcFormField) WithKey(key interface{}) *MdcFormField {
	ch.Keyable.WithKey(key)
	return ch
}

func (ch *MdcFormField) WithClasses(classes ...string) *MdcFormField {
	ch.Classes.WithClasses(classes...)
	return ch
}

func (ch *MdcFormField) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *MdcFormField) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			vecty.Class("mdc-form-field"),
			ch.ApplyClasses(),
		),
		ch.Content,
	)
}
