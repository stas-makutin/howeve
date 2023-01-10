package components

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/stas-makutin/howeve/page/core"
)

func init() {
	core.AppendStyles(`
.cm-capsule-shell {
	display: flex;
	border: 1px solid var(--mdc-theme-primary, #6200ee);
	border-radius: var(--mdc-shape-small, 4px); 
}
.cm-capsule-label {
	padding: 0.15em 0.4em;
	border-radius: var(--mdc-shape-small, 4px) 0 0 var(--mdc-shape-small, 4px);
}
.cm-capsule-text {
	white-space: normal;
	word-break: break-all;
	padding: 0.15em 0.4em;
	border-radius: 0 var(--mdc-shape-small, 4px) var(--mdc-shape-small, 4px) 0;
}
`,
	)
}

type Capsule struct {
	vecty.Core
	core.Keyable
	Label string `vecty:"prop"`
	Text  string `vecty:"prop"`
}

func (ch *Capsule) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *Capsule) WithKey(key interface{}) *Capsule {
	ch.Keyable.WithKey(key)
	return ch
}

func (ch *Capsule) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			vecty.Class("cm-capsule-shell"),
		),
		elem.Bold(
			vecty.Markup(
				vecty.Class("cm-capsule-label", "mdc-theme--primary-bg", "mdc-theme--on-primary"),
			),
			vecty.Text(ch.Label),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("cm-capsule-text"),
			),
			vecty.Text(ch.Text),
		),
	)
}
