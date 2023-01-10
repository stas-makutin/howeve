package components

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/stas-makutin/howeve/page/core"
)

func init() {
	core.AppendStyles(`
.cm-pillbox {
	display: flex;
	flex-wrap: wrap;
	gap: 0.2em;
}
`,
	)
}

type PillBox struct {
	vecty.Core
	core.Classes
	core.Keyable
	Content []vecty.MarkupOrChild `vecty:"prop"`
}

func NewPillBox(content ...vecty.MarkupOrChild) *PillBox {
	return &PillBox{Content: content}
}

func (ch *PillBox) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *PillBox) WithKey(key interface{}) *PillBox {
	ch.Keyable.WithKey(key)
	return ch
}

func (ch *PillBox) WithClasses(classes ...string) *PillBox {
	ch.Classes.WithClasses(classes...)
	return ch
}

func (ch *PillBox) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			vecty.Class("cm-pillbox"),
			ch.ApplyClasses(),
		),
		elem.Slot(ch.Content...),
	)
}
